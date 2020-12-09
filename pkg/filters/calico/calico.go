package calico

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"k8s.io/client-go/rest"

	"github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/filters/util"
)

const (
	// This number is based on the fact that there are two GlobalNetworkPolicies
	// that are created for calico. Each policy is specified for three protocols
	// TCP , UDP and ICMP. Hence it is expected that `iptables --list | grep <ipset name>`
	// would match for 6 entries which would confirm that the ipset/iptables rules
	// have been updated by Calico and we can proceed with the test.
	expectedNumberOfMatchesInIptables = 6
)

type calicoCNI struct {
	Iface                     string
	Nets                      []string
	GlobalNetworkSetManifests []string
}

func New(nets []net.IPNet, iface string) *calicoCNI {
	netsList := make([]string, len(nets))
	for i, n := range nets {
		netsList[i] = n.String()
	}

	return &calicoCNI{
		Iface:                     iface,
		Nets:                      netsList,
		GlobalNetworkSetManifests: make([]string, 0),
	}
}

// SetUp installs the filter in iface
func (b *calicoCNI) SetUp(nets []net.IPNet, iface string) (int64, error) {
	start := time.Now()

	// Get the list of ipsets before applying the filter.
	ipsetList, err := listIpsets()
	if err != nil {
		return 0, fmt.Errorf("listing existing ipsets: %w", err)
	}

	gnsManifests := map[int][]string{}

	// This code splits the GlobalNetworkSet manifests to accomodate a
	// maximum if `util.RulesPerManifest` entries per manifest.
	current := 0
	maxManifests := (len(b.Nets) / util.RulesPerManifest) + 1

	for i := 0; i < maxManifests; i++ {
		remaining := len(b.Nets) - current
		if remaining >= util.RulesPerManifest {
			gnsManifests[i] = b.Nets[current : current+util.RulesPerManifest]
		} else {
			gnsManifests[i] = b.Nets[current : current+remaining]
		}

		current = current + util.RulesPerManifest

		m := struct {
			Index int
			Nets  []string
		}{
			Index: i,
			Nets:  gnsManifests[i],
		}

		rendered, err := util.RenderTemplate(globalNetworkSetTmpl, m)
		if err != nil {
			return 0, fmt.Errorf("rendering GlobalNetworkSet template: %w", err)
		}

		b.GlobalNetworkSetManifests = append(b.GlobalNetworkSetManifests, rendered)
	}
	// Render GlobalNetworkPolicy
	gnpManifest, err := util.RenderTemplate(globalNetworkPolicyTmpl, b)
	if err != nil {
		return 0, fmt.Errorf("rendering GlobalNetworkPolicy template: %w", err)
	}
	// Render GlobalNetworkPolicy
	gnpWorkloadsManifest, err := util.RenderTemplate(gnpTmplForWorkloads, b)
	if err != nil {
		return 0, fmt.Errorf("rendering GlobalNetworkPolicy workloads template: %w", err)
	}

	// Get in cluster config, to create k8s resources.
	config, err := rest.InClusterConfig()
	if err != nil {
		return 0, err
	}

	for _, gnsManifest := range b.GlobalNetworkSetManifests {
		// Decode and apply the GlobalNetworkSet manifest
		if err := util.DecodeAndApply(config, gnsManifest, "CREATE"); err != nil {
			return 0, err
		}
	}

	// Decode and apply the GlobalNetworkPolicy manifest
	if err := util.DecodeAndApply(config, gnpManifest, "CREATE"); err != nil {
		return 0, err
	}
	// Decode and apply the GlobalNetworkPolicy manifest
	if err := util.DecodeAndApply(config, gnpWorkloadsManifest, "CREATE"); err != nil {
		return 0, err
	}

	ready, err := waitUntilReady(ipsetList, len(b.Nets))
	if err != nil {
		return 0, fmt.Errorf("calico ipsets not yet ready: %w", err)
	}

	if !ready {
		return 0, fmt.Errorf("calico ipsets not yet ready")
	}

	elapsed := time.Since(start)

	return elapsed.Nanoseconds(), nil
}

func waitUntilReady(ipsetListBefore map[string]bool, expectedNumberOfEntries int) (bool, error) {
	for i := 0; i <= util.Timeout; i++ {
		ipsetListAfter, err := listIpsets()
		if err != nil {
			return false, fmt.Errorf("listing existing ipsets: %w", err)
		}

		for entry, exists := range ipsetListBefore {
			if ipsetListAfter[entry] == exists {
				delete(ipsetListAfter, entry)
			}
		}

		s := strconv.Itoa(expectedNumberOfEntries)
		expectedOutput := fmt.Sprintf("Number of entries: %s", s)

		for entry, _ := range ipsetListAfter {
			run := fmt.Sprintf("ipset list %s -terse | grep '^Number of entries'", entry)
			output, err := runCmd(run)
			if err != nil {
				return false, fmt.Errorf("executing command %q: %w", run, err)
			}

			output = strings.Trim(output, "\n")
			if output == expectedOutput {
				// Even if the entries are created in the hashset; Calico needs a little more
				// time to update the iptables rules to match the hashset.
				for j := 0; j <= util.Timeout; j++ {
					run = fmt.Sprintf("iptables --list | grep '%s'", entry)
					iptablesOutput, err := runCmd(run)
					if err != nil {
						return false, fmt.Errorf("executing command %q: %w", run, err)
					}

					lines := strings.Count(iptablesOutput, entry)
					if lines == expectedNumberOfMatchesInIptables {
						return true, nil
					}

					time.Sleep(1 * time.Millisecond)
				}
				return true, nil
			}
		}

		// Wait for one millisecond before checking again.
		time.Sleep(1 * time.Millisecond)
	}

	return false, nil
}

func (b *calicoCNI) CleanUp() {
	gnpManifest, err := util.RenderTemplate(globalNetworkPolicyTmpl, b)
	if err != nil {
		fmt.Printf("rendering GlobalNetworkPolicy template: %w", err)
	}

	gnpWorkloadsManifest, err := util.RenderTemplate(gnpTmplForWorkloads, b)
	if err != nil {
		fmt.Printf("rendering GlobalNetworkPolicy workloads template: %w", err)
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("creating cluster config: %w", err)
	}

	for _, gnsManifest := range b.GlobalNetworkSetManifests {
		// Decode and apply the GlobalNetworkSet manifest
		if err := util.DecodeAndApply(config, gnsManifest, "DELETE"); err != nil {
			fmt.Printf("deleting GlobalNetworkSets: %w", err)
		}
	}
	// Decode and apply the GlobalNetworkPolicy manifest
	if err := util.DecodeAndApply(config, gnpManifest, "DELETE"); err != nil {
		fmt.Printf("deleting GlobalNetworkPolicy: %w", err)
	}
	// Decode and apply the GlobalNetworkPolicy manifest
	if err := util.DecodeAndApply(config, gnpWorkloadsManifest, "DELETE"); err != nil {
		fmt.Printf("deleting GlobalNetworkPolicy: %w", err)
	}
}

func runCmd(run string) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("sh", "-c", run)
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Ignoring the error here has the error returned is of the form
	// `exit status`.

	// The actual error is captured in stderr which is then passed to the
	// calling function as an error
	_ = cmd.Run()

	if len(stderr.String()) > 0 {
		return out.String(), fmt.Errorf(stderr.String())
	}

	return out.String(), nil
}

func listIpsets() (map[string]bool, error) {
	run := "ipset list -n"
	// Ignoring the err here as it returns `exit status 1`
	// Actual message is captured in output.
	output, err := runCmd(run)
	if err != nil {
		return nil, fmt.Errorf("executing command %q: %w", run, err)
	}

	pattern := "cali([a-zA-Z0-9]*)[\\:]([a-zA-Z0-9]*)"
	// Map to store ipset list before applying the manifests
	ipsetList := map[string]bool{}
	for _, name := range strings.Split(output, "\n") {
		match, err := regexp.MatchString(pattern, name)
		if err != nil {
			return nil, fmt.Errorf("pattern match %s: %w", pattern, err)
		}
		if match {
			ipsetList[name] = match
		}
	}

	return ipsetList, nil
}

package cilium

import (
	"fmt"
	"net"
	"time"

	"k8s.io/client-go/rest"

	"github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/filters/util"
)

type ciliumCNI struct {
	Iface                             string
	Nets                              []string
	ClusterWideNetworkPolicyManifests []string
	PingIP                            string
}

func New(nets []net.IPNet, iface, testIP string) *ciliumCNI {
	netsList := make([]string, len(nets))
	for i, n := range nets {
		netsList[i] = n.String()
	}

	return &ciliumCNI{
		Iface:                             iface,
		Nets:                              netsList,
		ClusterWideNetworkPolicyManifests: make([]string, 0),
		PingIP:                            testIP,
	}
}

// SetUp installs the filter in iface
func (b *ciliumCNI) SetUp(nets []net.IPNet, iface string) (int64, error) {
	start := time.Now()

	ccnpManifests := map[int][]string{}

	// This code splits the CiliumClusterWideNetworkPolicy manifests to accomodate a
	// maximum if `util.RulesPerManifest` entries per manifest.
	current := 0
	maxManifests := (len(b.Nets) / util.RulesPerManifest) + 1

	for i := 0; i < maxManifests; i++ {
		remaining := len(b.Nets) - current
		if remaining >= util.RulesPerManifest {
			ccnpManifests[i] = b.Nets[current : current+util.RulesPerManifest]
		} else {
			ccnpManifests[i] = b.Nets[current : current+remaining]
		}

		current = current + util.RulesPerManifest

		m := struct {
			Index int
			Nets  []string
		}{
			Index: i,
			Nets:  ccnpManifests[i],
		}

		rendered, err := util.RenderTemplate(allowAllEgressOnHost, m)
		if err != nil {
			return 0, fmt.Errorf("rendering CiliumClusterWideNetworkPolicy template: %w", err)
		}

		b.ClusterWideNetworkPolicyManifests = append(b.ClusterWideNetworkPolicyManifests, rendered)

		rendered, err = util.RenderTemplate(denyPolicyForHosts, m)
		if err != nil {
			return 0, fmt.Errorf("rendering CiliumClusterWideNetworkPolicy template: %w", err)
		}

		b.ClusterWideNetworkPolicyManifests = append(b.ClusterWideNetworkPolicyManifests, rendered)

		rendered, err = util.RenderTemplate(allowAllEgressOnBenchmarkApp, m)
		if err != nil {
			return 0, fmt.Errorf("rendering CiliumClusterWideNetworkPolicy template: %w", err)
		}

		b.ClusterWideNetworkPolicyManifests = append(b.ClusterWideNetworkPolicyManifests, rendered)

		rendered, err = util.RenderTemplate(denyPolicyForBenchmarkApp, m)
		if err != nil {
			return 0, fmt.Errorf("rendering CiliumClusterWideNetworkPolicy template: %w", err)
		}

		b.ClusterWideNetworkPolicyManifests = append(b.ClusterWideNetworkPolicyManifests, rendered)

	}

	// Get in cluster config, to create k8s resources.
	config, err := rest.InClusterConfig()
	if err != nil {
		return 0, err
	}

	for _, ccnpManifest := range b.ClusterWideNetworkPolicyManifests {
		// Decode and apply the CiliumClusterWideNetworkPolicy manifest
		if err := util.DecodeAndApply(config, ccnpManifest, "CREATE"); err != nil {
			return 0, err
		}
	}

	err = waitUntilPolicyIsReady(b.PingIP)
	if err != nil {
		return 0, err
	}
	elapsed := time.Since(start)

	return elapsed.Nanoseconds(), nil
}

func waitUntilPolicyIsReady(pingIP string) error {
	for i := 0; i <= util.Timeout; i++ {
		err := util.CheckPingFail(pingIP)
		// Ping should fail and throw no error. If that happens, break out of loop
		if err == nil {
			return nil
		}

		time.Sleep(1 * time.Millisecond)
	}

	return fmt.Errorf("network policies created but not yet enforced")
}

func waitUntilPolicyIsDeleted(pingIP string) error {
	for i := 0; i <= util.Timeout; i++ {
		err := util.CheckPingSuccess(pingIP)
		// Ping should succeed and throw no error. If that happens, break out of loop
		if err == nil {
			return nil
		}

		time.Sleep(1 * time.Millisecond)
	}

	return fmt.Errorf("network policies deleted but waiting for test ip %q to be reachable again", pingIP)
}

func (b *ciliumCNI) CleanUp() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("creating cluster config: %v", err)
	}

	for _, ccnpManifest := range b.ClusterWideNetworkPolicyManifests {
		// Decode and delete the CiliumClusterwideNetworkPolicy manifest
		if err := util.DecodeAndApply(config, ccnpManifest, "DELETE"); err != nil {
			return fmt.Errorf("deleting CiliumClusterwideNetworkPolicy: %v", err)
		}
	}

	err = waitUntilPolicyIsDeleted(b.PingIP)
	if err != nil {
		return err
	}

	return nil
}

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"net"

	"github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/filters/calico"
	cgroupbpf "github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/filters/cgroup-bpf"
	"github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/filters/ipset"
	"github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/filters/iptables"
	tcbpf "github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/filters/tc-bpf"
	"github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/filters/util"
	"github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/ipnetsgenerator"
)

var (
	iface       string
	countParam  int
	ipnetsParam string
	seed        int64
	filterType  string
	testIP      string
)

type filter interface {
	SetUp(nets []net.IPNet, iface string) (int64, error)
	CleanUp()
}

func init() {
	flag.StringVar(&iface, "iface", "", "Iface to attach the filter to")
	flag.IntVar(&countParam, "count", 0, "Number of entries to generate")
	flag.StringVar(&ipnetsParam, "ipnets", "", "List of ipnets and their weigth to generate (ex. 24:0.7,16:0.1)")
	flag.Int64Var(&seed, "seed", 0, "Seed to use for the random generator")
	flag.StringVar(&filterType, "filter", "", "Type of filter to use (tc-bpf, cgroup-bpf, iptables, ipset)")
	flag.StringVar(&testIP, "test-ip", util.DefaultPingIP, "IP to perform a ping to test if filters were correctly applied")
}

func main() {
	flag.Parse()

	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	ipnetsReq := ipnetsgenerator.ParseIPNetsParam(countParam-1, ipnetsParam)
	nets := ipnetsgenerator.GenerateIPNets(ipnetsReq, seed)

	_, pingIP, err := net.ParseCIDR(testIP + "/32")
	if err != nil {
		fmt.Printf("error parsing testIP: %s\n", err)
		return
	}
	nets = append(nets, *pingIP)

	var filter filter

	switch filterType {
	case "none":
		filter = nil
	case "tc-bpf":
		filter = tcbpf.New()
	case "cgroup-bpf":
		filter = cgroupbpf.New()
	case "iptables":
		filter = iptables.New()
	case "ipset":
		filter = ipset.New()
	case "calico":
		filter = calico.New(nets, iface)
	default:
		fmt.Printf("%q is not a valid filter type", filterType)
		os.Exit(1)
	}

	var setupTime int64

	if filter != nil {
		// Check that the test ip is reachable before applying the filter.
		if err := util.CheckPingSuccess(testIP); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		// Disabling this check for Calico, as Calico has fail safe network policies
		// that allow DNS communication.
		if filterType != "calico" {
			if err := util.CheckDNSSuccess(testIP); err != nil {
				fmt.Printf("%s\n", err)
				return
			}
		}

		var err error
		setupTime, err = filter.SetUp(nets, iface)
		if err != nil {
			fmt.Printf("error setting up filter %s", err)
			return
		}
		// Check that the test ip is not reachable after applying the filter.
		if err := util.CheckPingFail(testIP); err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		// Disabling this check for Calico, as Calico has fail safe network policies
		// that allow DNS communication.
		if filterType != "calico" {
			if err := util.CheckDNSFail(testIP); err != nil {
				fmt.Printf("%s\n", err)
				return
			}
		}

		defer filter.CleanUp()
	}

	if os.Getenv("BENCHMARK_COMMAND") == "MEASURE_SETUP_TIME" {
		fmt.Println(setupTime)
	} else if os.Getenv("BENCHMARK_COMMAND") != "" {
		cmd1 := exec.Command("sh", "-c", os.Getenv("BENCHMARK_COMMAND"))
		stdoutStderr, err := cmd1.CombinedOutput()
		if err != nil {
			fmt.Printf("Error executing command %q: %s\n", os.Getenv("BENCHMARK_COMMAND"), err)
		}
		fmt.Printf("%s\n", stdoutStderr)
	} else {

		var input string
		fmt.Println("Print enter to exit: ")
		fmt.Scanf("%s", &input)
	}
}

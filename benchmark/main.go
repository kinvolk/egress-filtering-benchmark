package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"net"

	"github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/filters/bpf"
	"github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/filters/ipset"
	"github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/filters/iptables"
	"github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/ipnetsgenerator"

	"github.com/sparrc/go-ping"
)

const (
	testIPDefault = "8.8.8.8" // IP used to test if the filter works.
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
	flag.StringVar(&filterType, "filter", "", "Type of filter to use (bpf, iptables, ipset)")
	flag.StringVar(&testIP, "test-ip", testIPDefault, "IP to perform a ping to test if filters were correctly applied")
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
	case "bpf":
		filter = bpf.New()
	case "iptables":
		filter = iptables.New()
	case "ipset":
		filter = ipset.New()
	default:
		fmt.Printf("%q is not a valid filter type", filterType)
		os.Exit(1)
	}

	var setupTime int64

	if filter != nil {
		// Check that the test ip is reachable before applying the filter
		if err := checkPingSuccess(testIP); err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		var err error

		// Check that the test ip is reachable before applying the filter
		if err := checkPingSuccess(testIP); err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		setupTime, err = filter.SetUp(nets, iface)
		if err != nil {
			fmt.Printf("error setting up filter %s", err)
			return
		}

		// Check that the test ip is not reachable after applying the filter
		if err := checkPingFail(testIP); err != nil {
			fmt.Printf("%s\n", err)
			return
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

func createPinger() (*ping.Pinger, error) {
	pinger, err := ping.NewPinger(testIP)
	if err != nil {
		return nil, err
	}

	pinger.Count = 1
	pinger.SetPrivileged(true)
	pinger.Timeout = time.Second
	return pinger, err
}

func checkPingSuccess(ip string) error {
	pinger, err := createPinger()
	if err != nil {
		return err
	}

	pinger.Run()
	if pinger.PacketsRecv != 1 {
		return fmt.Errorf("imposible to ping %q", ip)
	}

	return nil
}

func checkPingFail(ip string) error {
	pinger, err := createPinger()
	if err != nil {
		return err
	}

	pinger.Run()
	if pinger.PacketsRecv != 0 {
		return fmt.Errorf("ping to %q should have fail", ip)
	}

	return nil
}

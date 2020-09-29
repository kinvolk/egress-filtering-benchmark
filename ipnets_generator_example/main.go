package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/kinvolk/k8s-egress-filtering-benchmark/pkg/ipnetsgenerator"
)

var (
	countParam  int
	ipnetsParam string
	seed        int64
)

func init() {
	flag.IntVar(&countParam, "count", 0, "Number of entries to generate")
	flag.StringVar(&ipnetsParam, "ipnets", "", "List of ipnets and their weigth to generate (ex. 24:0.7,16:0.1)")
	flag.Int64Var(&seed, "seed", 0, "Seed to use for the random generator")
}

func main() {
	flag.Parse()

	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	ipnetsReq := ipnetsgenerator.ParseIPNetsParam(countParam, ipnetsParam)
	nets := ipnetsgenerator.GenerateIPNets(ipnetsReq, seed)

	for _, i := range nets {
		fmt.Println(i.String())
	}
}

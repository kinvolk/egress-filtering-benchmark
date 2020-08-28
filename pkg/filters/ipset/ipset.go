package ipset

import (
	"bytes"
	"fmt"
	"net"
	"os/exec"
	"time"
)

const (
	setName    = "myset"
	setType    = "hash:net"
	setSize    = 2000000
	ruleFormat = "add %s %s\n"
)

type ipsetFilter struct {
	iface string
}

func New() *ipsetFilter {
	return &ipsetFilter{}
}

// SetUp installs the filter in iface
func (b *ipsetFilter) SetUp(nets []net.IPNet, iface string) (int64, error) {
	if len(nets) > setSize {
		return 0, fmt.Errorf("Imposible to add %d rules. The maximum allowed is %d",
			len(nets), setSize)
	}

	b.iface = iface

	start := time.Now()

	// create and fill ipset
	var buf bytes.Buffer
	var err error
	fmt.Fprintf(&buf, "create %s %s maxelem %d\n", setName, setType, setSize)
	for _, n := range nets {
		_, err := fmt.Fprintf(&buf, ruleFormat, setName, n.String())
		if err != nil {
			return 0, err
		}
	}
	fmt.Fprintln(&buf, "quit")

	err = execIpSet(buf.Bytes())
	if err != nil {
		return 0, err
	}

	elapsed := time.Since(start)

	// associate iptable to it
	err = execIpTables("-A", "OUTPUT", "-o", iface, "-m", "set", "--match-set",
		setName, "dst", "-j", "DROP", "-m", "comment", "--comment", "benchmark")
	if err != nil {
		return 0, err
	}

	return elapsed.Nanoseconds(), nil
}

func execIpTables(args ...string) error {
	cmd := exec.Command("iptables", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Error executing iptables %w: %s", err, out)
	}

	return nil
}

func execIpSet(cmds []byte) error {
	cmd := exec.Command("ipset", "-")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		stdin.Write(cmds)
	}()

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Error adding iptables rules %w: %s", err, out)
	}

	return nil
}

func (b *ipsetFilter) CleanUp() {
	var buf bytes.Buffer
	// remove iptable associated to it
	execIpTables("-D", "OUTPUT", "-o", b.iface, "-m", "set", "--match-set",
		setName, "dst", "-j", "DROP", "-m", "comment", "--comment", "benchmark")

	fmt.Fprintf(&buf, "destroy %s\n", setName)
	fmt.Fprintln(&buf, "quit")

	execIpSet(buf.Bytes())
}

package iptables

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os/exec"
)

const (
	ruleFormat = "-A OUTPUT -d %s -o %s -j DROP\n"
	rulesPath  = "/tmp/iptables-save.txt"
)

type iptablesFilter struct {
	iface string
	nets  []net.IPNet
}

func New() *iptablesFilter {
	return &iptablesFilter{}
}

// SetUp installs the filter in iface
func (f *iptablesFilter) SetUp(nets []net.IPNet, iface string) error {
	f.iface = iface
	f.nets = nets

	rulesToSave, err := iptablesSave()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(rulesPath, rulesToSave, 0644); err != nil {
		return err
	}

	var b bytes.Buffer
	fmt.Fprintln(&b, "*filter")
	for _, n := range nets {
		_, err := fmt.Fprintf(&b, ruleFormat, n.String(), iface)
		if err != nil {
			return err
		}
	}
	fmt.Fprintln(&b, "COMMIT")

	// add the new rules without removing the existing ones
	return iptablesRestore(b.Bytes(), true)
}

// CleanUp removes the filter
func (f *iptablesFilter) CleanUp() {
	rules, err := ioutil.ReadFile(rulesPath)
	if err != nil {
		return
	}
	iptablesRestore(rules, false)
}

func iptablesRestore(rules []byte, noflush bool) error {
	// restore the previous rules (remove existing ones)
	cmd := exec.Command("iptables-restore")
	if noflush {
		cmd.Args = append(cmd.Args, "--noflush")
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		defer stdin.Close()
		stdin.Write(rules)
	}()

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Error adding iptables rules %w: %s", err, out)
	}

	return nil
}

func iptablesSave() ([]byte, error) {
	cmd := exec.Command("iptables-save")
	return cmd.Output()
}

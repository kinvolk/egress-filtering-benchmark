package bpf

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"net"
	"unsafe"

	"github.com/cilium/ebpf"
)

type bpf struct {
	iface string
}

const (
	mapPath     = "/sys/fs/bpf/tc/globals/lpm_filter"
	bpfCodePath = "/tmp/bpf.o"
)

func New() *bpf {
	return &bpf{}
}

// SetUp installs the filter in iface
func (b *bpf) SetUp(nets []net.IPNet, iface string) error {
	asset, err := Asset("datapath/bpf.o")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(bpfCodePath, asset, 0644)
	if err != nil {
		return err
	}
	defer os.Remove(bpfCodePath)

	b.iface = iface
	b.CleanUp()
	cmd1 := exec.Command("tc", "qdisc", "add", "dev", iface, "clsact")
	if out, err := cmd1.CombinedOutput(); err != nil {
		return fmt.Errorf("Error adding clsact %w: %s", err, out)
	}

	cmd2 := exec.Command("tc", "filter", "add", "dev", iface,
		"egress", "bpf", "da", "obj", bpfCodePath, "sec", "filter_egress")
	if out, err := cmd2.CombinedOutput(); err != nil {
		return fmt.Errorf("Error adding egress filter %w: %s", err, out)
	}

	if err := updateMap(nets); err != nil {
		return err
	}

	return nil
}

// CleanUp removes the filter
func (b *bpf) CleanUp() {
	cmd1 := exec.Command("tc", "filter", "del", "dev", b.iface, "egress")
	cmd1.Run()

	cmd2 := exec.Command("tc", "qdisc", "del", "dev", b.iface, "clsact")
	cmd2.Run()

	cmd3 := exec.Command("rm", mapPath)
	cmd3.Run()
}

func updateMap(nets []net.IPNet) error {
	filterMap, err := ebpf.LoadPinnedMap(mapPath)
	if err != nil {
		return err
	}

	value := uint32(0)

	for _, n := range nets {
		siz, _ := n.Mask.Size()
		IPBigEndian := unsafe.Pointer(&n.IP[0])
		key := []uint32{uint32(siz), *(*uint32)(IPBigEndian)}

		err2 := filterMap.Put(unsafe.Pointer(&key[0]), unsafe.Pointer(&value))
		if err2 != nil {
			return err2
		}
	}

	return nil
}

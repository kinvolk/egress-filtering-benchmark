package cgroupbpf

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"
	"unsafe"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"golang.org/x/sys/unix"
)

type bpf struct {
	prog *ebpf.Program
}

const (
	rootCgroup     = "/sys/fs/cgroup/unified"
	mapName        = "lpm_filter"
	bpfCodePath    = "/tmp/benchmark-cgroup-bpf.o"
	bpfProgramName = "filter_egress"
)

func New() *bpf {
	return &bpf{}
}

// SetUp installs the filter in iface
func (b *bpf) SetUp(nets []net.IPNet, iface string) (int64, error) {
	asset, err := Asset("datapath/bpf.o")
	if err != nil {
		return 0, err
	}
	err = ioutil.WriteFile(bpfCodePath, asset, 0644)
	if err != nil {
		return 0, err
	}
	defer os.Remove(bpfCodePath)

	b.CleanUp()

	unix.Setrlimit(unix.RLIMIT_MEMLOCK, &unix.Rlimit{
		Cur: unix.RLIM_INFINITY,
		Max: unix.RLIM_INFINITY,
	})

	collec, err := ebpf.LoadCollection(bpfCodePath)
	if err != nil {
		return 0, err
	}

	if _, ok := collec.Programs[bpfProgramName]; !ok {
		return 0, fmt.Errorf("Object file doesn't contain program '%s'", bpfProgramName)
	}
	b.prog = collec.Programs[bpfProgramName]

	cgroup, err := os.Open(rootCgroup)
	if err != nil {
		return 0, err
	}
	defer cgroup.Close()

	_, err = link.AttachCgroup(link.CgroupOptions{
		Path:    cgroup.Name(),
		Attach:  ebpf.AttachCGroupInetEgress,
		Program: collec.Programs["filter_egress"],
	})
	if err != nil {
		return 0, err
	}

	start := time.Now()
	if err := updateMap(collec.Maps, nets); err != nil {
		return 0, err
	}
	elapsed := time.Since(start)

	return elapsed.Nanoseconds(), nil
}

// CleanUp removes the filter
func (b *bpf) CleanUp() {
	if b.prog != nil {
		cgroup, err := os.Open(rootCgroup)
		if err != nil {
			return
		}
		defer cgroup.Close()

		b.prog.Detach(int(cgroup.Fd()), ebpf.AttachCGroupInetEgress, 0)
	}
}


func updateMap(maps map[string]*ebpf.Map, nets []net.IPNet) error {
	filterMap, ok := maps[mapName]
	if !ok {
		return fmt.Errorf("Map %s not found in object file.", mapName)
	}

	value := uint8(0)

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

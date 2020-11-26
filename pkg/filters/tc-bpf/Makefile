SHELL=/bin/bash -o pipefail
EXTRA_CFLAGS += -fno-stack-protector
bpf-assets.go: datapath/bpf.c
	clang $(EXTRA_CFLAGS) -O2 -emit-llvm -c datapath/bpf.c -o - | llc -march=bpf -filetype=obj -o datapath/bpf.o
	go-bindata -pkg tcbpf -modtime 1 -o bpf-assets.go datapath/bpf.o
	rm -f datapath/bpf.o

clean:
	rm -f bpf-assets.go

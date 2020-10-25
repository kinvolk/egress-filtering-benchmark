.PHONY: benchmark
benchmark:
	make -C ../pkg/filters/tc-bpf
	make -C ../pkg/filters/cgroup-bpf
	CGO_ENABLED=0 GO111MODULE=on go build \
		-o benchmark \
		github.com/kinvolk/k8s-egress-filtering-benchmark/benchmark

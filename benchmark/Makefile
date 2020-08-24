.PHONY: benchmark
benchmark:
	make -C ../pkg/filters/bpf
	GO111MODULE=on go build \
		-o benchmark \
		github.com/kinvolk/k8s-egress-filtering-benchmark/benchmark

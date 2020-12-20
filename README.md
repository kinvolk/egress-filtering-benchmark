# Egress Filtering Benchmark

This repository contains a set of tools to measure the egress filtering performance using BPF, iptables, ipsets and calico.

## How to Use

1. Setup two computers to run the test. You need to have Docker, iptables and ipset installed and you should be able to connect to those computers with SSH without requiring a password.

2. Download the latest Lokomotive [release](https://github.com/kinvolk/lokomotive/releases)

    Unpack and move to a desired locatin:
    ```bash
    tar xvf lokoctl_$VERSION_linux_amd64.tar.gz
    mv lokoctl_$VERSION_linux_amd64/lokoctl ~/.local/bin/lokoctl
    ```

3. Create a Kubernetes cluster using [Lokomotive](https://github.com/kinvolk/lokomotive) with at least one worker node.

    A minimal working configuration that can be deployed on [Packet](https://www.packet.com/) (acquired by [Equinix Metal](https://metal.equinix.com/))

    Update the variables in `lokocfg.vars` and execute:

    ```bash
    git clone git@github.com:kinvolk/egress-filtering-benchmark.git
    cd lokomotive/calico-benchmark-cluster
    lokoctl cluster apply
    ```

     Set location of kubeconfig in the environment variable **KUBECONFIG_CALICO**:

     ```
     cd egress-filtering-benchmark/lokomotive
     export KUBECONFIG_CALICO=$PWD/assets/cluster-assets/auth/kubeconfig
     ```

     Label the worker node as follows:
     ```
     kubectl label node calico-benchmark-pool-1-worker-0 nodetype=worker-benchmark
     ```
4. Create a Kubernetes cluster with Cilium using [Lokomotive](https://github.com/kinvolk/lokomotive) with at least one worker node.

    A minimal working configuration that can be deployed on [Packet](https://www.packet.com/) (acquired by [Equinix Metal](https://metal.equinix.com/))

    Since Lokomotive doesn't ship with Cilium in the lokoctl binary, we will checkout the project and
    build the binary based on a different branch.
    ```bash
    git clone git@github.com:kinvolk/lokomotive.git
    cd lokomotive
    git checkout imran/cilium-instead-of-calico
    make
    mv lokoctl ~/.local/bin/lokoctl_cilium
    ```
    Update the variables in `lokocfg.vars` and execute:

    ```bash
    git clone git@github.com:kinvolk/egress-filtering-benchmark.git
    cd lokomotive/cilium-benchmark-cluster
    lokoctl_cilium cluster apply
    ```

     Set location of kubeconfig using the environment variable **KUBECONFIG_CILIUM**:

     ```
     cd egress-filtering-benchmark/lokomotive
     export KUBECONFIG_CILIUM=$PWD/assets/cluster-assets/auth/kubeconfig
     ```

     Label the worker node as follows:
     ```
     kubectl label node cilium-benchmark-pool-1-worker-0 nodetype=worker-benchmark
     ```

5. Configure the parameters of the test in the [parameters.py](benchmark/parameters.py) file.

6. Install the required libraries in the client to run the Python script

```
pip install -r requirements.txt
```

7. Execute the tests:

```
$ cd benchmark
$ make
$ python run_tests.py --mode udp --username USERNAME --client CLIENTADDR --server SERVERADDR --kubeconfig-cilium-cluster $KUBECONFIG_CILIUM --kubeconfig-calico-cluster $KUBECONFIG_CALICO
```

This will create some csv files with the information about the test.
You can plot them by your self or follow the next step.

7. Plot the data by running

```
$ python plot_data.py
```

This will create some svg files with the graphs.

8. Cleanup the Lokomotive clusters created in the benchmark process:
  ```bash
  cd lokomotive/calico-benchmark-cluster
  lokoctl cluster destroy --confirm
  cd lokomotive/cilium-benchmark-cluster
  lokoctl_cilium cluster destroy --confirm
  ```
## Credits

The BPF filter is inspired by the [tc-bpf](https://man7.org/linux/man-pages/man8/tc-bpf.8.html) man page and the [Cilium documentation](https://docs.cilium.io/en/latest/bpf/#tc-traffic-control).

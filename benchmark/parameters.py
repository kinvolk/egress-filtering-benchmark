# change these values to setup the test

# List of filters to use. Supported values are none, bpf, iptables and ipset.
filters = ['none','tc-bpf','cgroup-bpf','iptables','ipset','calico','cilium']
# Number of rules to test
number_rules = ['10','100','1000','10000','100000','1000000']
# Number of iterations to run each test.
iterations = 5
# Network interface to bound the filters to.
iface = "bond0"
# target bandwidth used for UDP tests
bandwidth = "1G"
# Seed used for the random generator. 0 takes one from based on the time.
seed = 1
# IP blocks structure. (block:weight,block:weight)
ipnets = "24:0.1,16:0.01"
# Ping interval in seconds for the latency test.
ping_interval=0.001
# Number of pings to perform for the latency test.
ping_count=1000

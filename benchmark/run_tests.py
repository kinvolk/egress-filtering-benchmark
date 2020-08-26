import json
import sys
import csv
import argparse
import time
import subprocess

import pingparsing

# defined in parameters.py
from parameters import  (
    filters,
    number_rules,
    iterations,
    iface,
    bandwidth,
    seed,
    ipnets,
    ping_interval,
    ping_count,
)

parser = argparse.ArgumentParser("Test egress performance")
parser.add_argument("--username")
parser.add_argument("--client")
parser.add_argument("--server")
parser.add_argument("--mode")
args = parser.parse_args()

benchmark_cmd_format = "sudo -E ./benchmark -count {count} -iface {iface} -seed {seed} -ipnets {ipnets} -filter {filter}"
def run_test(filter, nrules, cmd):
    benchmark_cmd = benchmark_cmd_format.format(count=nrules, iface=iface, seed=seed, ipnets=ipnets, filter=filter)
    cmd = "export BENCHMARK_COMMAND='{}' ; {}".format(cmd, benchmark_cmd)
    return run_in_client(cmd)

iperf_cmd_format = "docker run --name iperfclient --net=host networkstatic/iperf3 -c {server} {mode_flags} -O 2 -t 10 -A 2 -J"
def run_iperf_test(filter, nrules, mode):
    flags = ""
    if mode == "udp":
        flags = "-u -l 1470 -b {}".format(bandwidth)

    iperf_cmd = iperf_cmd_format.format(server=args.server, mode_flags=flags)

    run_in_client("docker rm --force iperfclient || true")
    out = run_test(filter, nrules, iperf_cmd)
    if not out:
        return None

    #print("out is: " + out)
    index = out.find("{")
    if index == -1:
        return None
    j = json.loads(out[index:])

    if mode == "udp":
        key = "sum"
    elif mode == "tcp":
        key = "sum_received"

    throughput = float((j["end"][key]["bits_per_second"]))/(10**9)
    cpu = float(j["end"]["cpu_utilization_percent"]["host_total"])

    return (throughput, cpu)

ping_cmd_format = "ping -i {interval} -c {count} {dest}"
def run_ping_test(filter, nrules, mode):
    cmd = ping_cmd_format.format(interval=ping_interval, count=ping_count, dest=args.server)
    result = run_test(filter, nrules, cmd)
    if not result:
        return None
    parser = pingparsing.PingParsing()
    stats = parser.parse(result)
    return stats.rtt_avg

def start_iperf_server():
    run_in_server("docker rm --force iperfserver || true")
    run_in_server("docker run --name iperfserver -d --net=host networkstatic/iperf3 -s -A 2")
    time.sleep(2)

def copy_benchmark_to_client():
    cmd = "scp benchmark ${}@${}:".format(args.username, args.client)
    subprocess.run(cmd, stdout=subprocess.PIPE, shell=True)

def run_in_client(cmd):
    return run_over_ssh(args.client, cmd)

def run_in_server(cmd):
    return run_over_ssh(args.server, cmd)

def run_over_ssh(host, cmd):
    cmd_to_run = 'ssh {}@{} "{}"'.format(args.username, host, cmd)
    result = subprocess.run(cmd_to_run, stdout=subprocess.PIPE, shell=True)
    if result.returncode != 0:
      return None
    return result.stdout.decode("utf-8")

def write_csv(filename, data):
    with open(filename, "w", newline="") as csvfile:
        writer = csv.writer(csvfile, delimiter="\t",
                                quotechar=";", quoting=csv.QUOTE_MINIMAL)
        writer.writerow(["Filter", "Rules"] + ["r"+str(i) for i in range(1, iterations+1)])
        for (filter_name, filter_data) in data.items():
            for (rules_number, results) in filter_data.items():
                writer.writerow([filter_name, rules_number] + results)

# [filter][rules_number][run]
data_throughput = {}
data_cpu = {}
data_ping = {}

#copy_benchmark_to_client()
start_iperf_server()

print("%\tfilter\tnrules\titeration\tthroughput\tcpu\tping\t")

number_of_tests = len(filters)*len(number_rules)*iterations
number_of_tests_executed = 0

# run the tests and collect all the data
for (filter_index, filter) in enumerate(filters):
    data_throughput[filter] = {}
    data_cpu[filter] = {}
    data_ping[filter] = {}

    for (rules_index, n) in enumerate(number_rules):
        data_throughput[filter][n] = []
        data_cpu[filter][n] = []
        data_ping[filter][n] = []

        for i in range(iterations):
            percentage = 100.0*float(number_of_tests_executed)/number_of_tests
            print("{:1.0f}\t{}\t{}\t{}\t".format(percentage, filter, n, i), end="")

            out = run_iperf_test(filter, n, args.mode)
            if not out:
                print("Testing for {}:{}:{} failed".format(filter, n, i))
                continue

            out_ping = run_ping_test(filter, n, args.mode)
            if not out_ping:
                print("Testing for {}:{}:{} failed".format(filter, n, i))
                continue
            out_ping = 1000*out_ping

            print("{}\t{}\t{}".format(out[0], out[1], out_ping))
            number_of_tests_executed += 1

            data_throughput[filter][n].append(out[0])
            data_cpu[filter][n].append(out[1])
            data_ping[filter][n].append(out_ping)

write_csv("throughput.csv", data_throughput)
write_csv("cpu.csv", data_cpu)
write_csv("latency.csv", data_ping)

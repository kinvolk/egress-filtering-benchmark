#!/bin/bash

set -ex

usage()
{
    echo "usage: remote-test.sh [--username USERNAME] [--client SERVERNAME] [--server SERVERNAME]"
}

while [ "$1" != "" ]; do
    case $1 in
        -u | --username )       shift
                                username=$1
                                ;;
        -c | --client )         shift
                                client=$1
                                ;;
        -s | --server )         shift
                                server=$1
                                ;;
        -h | --help )           usage
                                exit
                                ;;
        * )                     usage
                                exit 1
    esac
    shift
done

if [ -z "$client" -o -z "$server" -o -z "$username" ]; then
    usage
    exit 1
fi

make

scp benchmark ${username}@${client}:

ssh ${username}@${server} "docker rm --force iperfserver || true"
ssh ${username}@${server} docker run --name iperfserver -d --net=host networkstatic/iperf3 -s
sleep 2
ssh ${username}@${server} docker ps


COUNT_LIST="10 100 1000 10000 100000"

for COUNT in $COUNT_LIST; do
  for FILTER in none bpf iptables ipset ; do
    ssh ${username}@${client} "docker rm --force iperfclient || true"
    ssh ${username}@${client} "export BENCHMARK_COMMAND=\"docker run --name iperfclient --net=host networkstatic/iperf3 -c $server --json\" ; sudo -E ./benchmark -count $COUNT -iface bond0 -seed 1 -ipnets 24:0.7,16:0.1 -filter $FILTER" > result-$FILTER-$COUNT.json
    cat result-$FILTER-$COUNT.json | jq '.end.sum_sent.bits_per_second' || true
  done
done

for i in result-*.json ; do echo -n "$i: " ; cat $i | jq '.end.sum_sent.bits_per_second' ; done

echo -n "filter;" > result-all.csv
for COUNT in $COUNT_LIST; do
  echo -n "$COUNT;" >> result-all.csv
done >> result-all.csv
echo >> result-all.csv
for FILTER in none bpf iptables ipset ; do
  echo -n "$FILTER;"
  for COUNT in $COUNT_LIST; do
    echo -n "$(cat result-$FILTER-$COUNT.json | jq '.end.sum_sent.bits_per_second');"
  done
  echo
done >> result-all.csv

#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/pkt_cls.h>

#include <iproute2/bpf_elf.h>

#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

#ifndef __section
# define __section(NAME)                  \
   __attribute__((section(NAME), used))
#endif

#define MAP_SIZE 100000

struct bpf_elf_map lpm_filter __section("maps") = {
    .type         = BPF_MAP_TYPE_LPM_TRIE,
    .size_key     = 8, // int + IPv4
    .size_value   = 4,
    .max_elem     = MAP_SIZE,
    .flags        = BPF_F_NO_PREALLOC,
    .pinning      = PIN_GLOBAL_NS,
};

__section("filter_egress")
int tc_egress(struct __sk_buff *skb) {
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;
    struct ethhdr *eth = data;

    if (data + sizeof(*eth) > data_end)
        return TC_ACT_OK;

    if (eth->h_proto != bpf_htons(ETH_P_IP))
        return TC_ACT_OK;

    struct iphdr *ip = data + sizeof(*eth);
    if (data + sizeof(*eth) + sizeof(*ip) > data_end)
        return TC_ACT_OK;

    __u32 lpm_key[2] = {32, ip->daddr};

    if (bpf_map_lookup_elem(&lpm_filter, lpm_key)) {
        // TODO: which one should we use, TC_ACT_STOLEN or TC_ACT_SHOT?
        // According to the documentation https://docs.cilium.io/en/latest/bpf/#tc-traffic-control
        // TC_ACT_SHOT notifies a failure to the upper layer while TC_ACT_STOLEN
        // doesn't
        return TC_ACT_STOLEN;
    }

    return TC_ACT_OK;
}

char __license[] __section("license") = "GPL";

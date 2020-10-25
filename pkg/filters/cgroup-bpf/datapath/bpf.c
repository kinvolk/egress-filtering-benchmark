#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/pkt_cls.h>

#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

#ifndef __section
# define __section(NAME)                  \
   __attribute__((section(NAME), used))
#endif

#define MAP_SIZE 2000000

struct bpf_map_def lpm_filter __section("maps") = {
    .type         = BPF_MAP_TYPE_LPM_TRIE,
    .key_size     = 8, // int + IPv4
    .value_size   = 4,
    .max_entries  = MAP_SIZE,
    .map_flags    = BPF_F_NO_PREALLOC,
};

__section("cgroup/skb1")
int filter_egress(struct __sk_buff *skb) {
    struct iphdr iph;
    __u32 dstip = skb->remote_ip4;

    if (!dstip) {
        bpf_skb_load_bytes(skb, 0, &iph, sizeof(struct iphdr));
        if (iph.version == 4)
            dstip = iph.daddr;
    }

    __u32 lpm_key[2] = {32, dstip};
    if (bpf_map_lookup_elem(&lpm_filter, lpm_key))
        return 0;

    return 1;
}

char __license[] __section("license") = "GPL";

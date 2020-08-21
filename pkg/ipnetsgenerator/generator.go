package ipnetsgenerator

import (
	"math/rand"
	"net"
	"strings"
	"strconv"
)

type IPNetRequest struct {
	Count     int // Number of elements of this type to generate
	PrefixLen int // Length of the prefix (24, 32)...
}

// ParseIPNetsParam parses a string containinig a list of prefixes' lengths and
// their weight (0.0-1.0) and returns a list of IPNetRequest to be used with
// GenerateIPNets.
// The ipnets string is formated as "prefix1Length:prefix1Weigth,prefix2Length:prefix2Weigth,..."
// If the sum of all weights isn't 1.0 the remaining are assigned to the /32 prefix.
func ParseIPNetsParam(count int, ipnets string) []IPNetRequest {
	if ipnets == "" {
		return []IPNetRequest{IPNetRequest{count, 32}}
	}

	ipnetReq := []IPNetRequest{}

	processedCount := 0
	reqs := strings.Split(ipnets, ",")

	for _, req := range reqs {
		pieces := strings.Split(req, ":")
		if len(pieces)!= 2 {
			continue
		}

		length, err := strconv.Atoi(pieces[0])
		if err != nil {
			continue
		}
		weight, err := strconv.ParseFloat(pieces[1], 64)
		if err != nil {
			continue
		}

		c := int(weight*float64(count))
		processedCount += c

		ipnetReq = append(ipnetReq, IPNetRequest{c, length})
	}

	// fill the remaining with /32
	remaininig := count - processedCount
	if remaininig > 0 {
		ipnetReq = append(ipnetReq, IPNetRequest{remaininig, 32})
	}

	return ipnetReq
}

// GenerateIPNets creates an array of random IPv4 net.IPNet objects.
// The ipnets parameter contains a list with th number of objects for each
// prefix length to generate.
func GenerateIPNets(ipnets []IPNetRequest, seed int64) []net.IPNet {
	r := rand.New(rand.NewSource(seed))

	ipnetsGenerated := []net.IPNet{}

	for _, ipnet := range ipnets {
		nets := generateIPNet(r, ipnet.Count, ipnet.PrefixLen)
		ipnetsGenerated = append(ipnetsGenerated, nets...)
	}

	return ipnetsGenerated
}

func generateIPNet(r *rand.Rand, count int, prefixLen int) []net.IPNet {
	// TODO: Check that it is actually possible to generate "count" ranges
	// with "prefixLen". For instance, it's not possible to have 300 combinations
	// of /8 subnets.
	// TODO: How to handle the case when the numer of entries to generate is
	// close to the number of possible combinations?

	// map to save the generated IPNets
	ipNets := make([]net.IPNet, count)

	// map to check if the IP is duplicated
	// use an uint32 as key because net.IPNet nor net.IP can be used as keys
	ipNetsMap := make(map[uint32]bool, count)
	mask := net.CIDRMask(prefixLen, 32)

	for i := 0; i < count; i++ {
		ip := toNetIP(r.Uint32())
		if !ip.IsGlobalUnicast() {
			i -= 1
			continue
		}
		ip = ip.Mask(mask)
		ipNet := net.IPNet{ip, mask}
		// check that the IPNet is unique
		if ok, _ := ipNetsMap[fromNetIP(ip)]; ok {
			i -= 1
			continue
		}
		ipNetsMap[fromNetIP(ip)] = true
		ipNets[i] = ipNet
	}

	return ipNets
}

// The following functions were taken from https://github.com/signalsciences/ipv4
// ToNetIP converts a uint32 to a net.IP (net.IPv4 actually)
func toNetIP(val uint32) net.IP {
	return net.IPv4(byte(val>>24), byte(val>>16&0xFF),
		byte(val>>8)&0xFF, byte(val&0xFF))
}

// FromNetIP converts a IPv4 net.IP to uint32
func fromNetIP(ip net.IP) uint32 {
	return uint32(ip[3]) | uint32(ip[2])<<8 | uint32(ip[1])<<16 | uint32(ip[0])<<24
}

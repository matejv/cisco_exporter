package nat64

import (
	"errors"

	"github.com/lwlcom/cisco_exporter/rpc"
	"github.com/lwlcom/cisco_exporter/util"
)

// ParseNat64 parses cli output and returns values
func (c *nat64Collector) ParseNat64(ostype string, output string) (Nat64Stats, error) {
	if ostype != rpc.IOSXE {
		return Nat64Stats{}, errors.New("'show nat64 statistics global' is not implemented for " + ostype)
	}

	results, err := util.ParseTextfsm(templ_nat64, output)
	if err != nil {
		return Nat64Stats{}, errors.New("Error parsing via templ_nat64: " + err.Error())
	}

	result := results[0]
	stats := Nat64Stats{
		translationsActive:    util.Str2float64(result["nat64_total_active_translations"].(string)),
		translationsExpired:   util.Str2float64(result["nat64_expired_translations"].(string)),
		sessionsFound:         util.Str2float64(result["nat64_sessions_found"].(string)),
		sessionsCreated:       util.Str2float64(result["nat64_sessions_created"].(string)),
		packetsTranslated4to6: util.Str2float64(result["nat64_ipv4_ipv6_translated_packets"].(string)),
		packetsTranslated6to4: util.Str2float64(result["nat64_ipv6_ipv4_translated_packets"].(string)),
	}
	return stats, nil
}

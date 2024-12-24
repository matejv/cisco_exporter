package nat64

import (
	"log"

	"github.com/lwlcom/cisco_exporter/rpc"

	"github.com/lwlcom/cisco_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix string = "cisco_nat64_"

var (
	translationsActiveDesc    *prometheus.Desc
	translationsExpiredDesc   *prometheus.Desc
	sessionsFoundDesc         *prometheus.Desc
	sessionsCreatedDesc       *prometheus.Desc
	packetsTranslated4to6Desc *prometheus.Desc
	packetsTranslated6to4Desc *prometheus.Desc
)

func init() {
	l := []string{"target"}
	translationsActiveDesc = prometheus.NewDesc(prefix+"translations_active", "Currently active NAT64 translations", l, nil)
	translationsExpiredDesc = prometheus.NewDesc(prefix+"translations_expired", "Total number of NAT64 translations removed from session table", l, nil)
	sessionsFoundDesc = prometheus.NewDesc(prefix+"sessions_found", "Count of packets that matched existing session in NAT64 session table", l, nil)
	sessionsCreatedDesc = prometheus.NewDesc(prefix+"sessions_created", "Count of new sessions created in NAT64 session table", l, nil)
	packetsTranslated4to6Desc = prometheus.NewDesc(prefix+"packets_translated_4to6", "Count of packets translated from IPv4 to IPv6", l, nil)
	packetsTranslated6to4Desc = prometheus.NewDesc(prefix+"packets_translated_6to4", "Count of packets translated from IPv6 to IPv4", l, nil)
}

type nat64Collector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &nat64Collector{}
}

// Name returns the name of the collector
func (*nat64Collector) Name() string {
	return "NAT64"
}

// Describe describes the metrics
func (*nat64Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- translationsActiveDesc
	ch <- translationsExpiredDesc
	ch <- sessionsFoundDesc
	ch <- sessionsCreatedDesc
	ch <- packetsTranslated4to6Desc
	ch <- packetsTranslated6to4Desc
}

// Collect collects metrics from Cisco
func (c *nat64Collector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	out, err := client.RunCommand("show nat64 statistics global")
	if err != nil {
		return err
	}
	stats, err := c.ParseNat64(client.OSType, out)
	if err != nil {
		if client.Debug {
			log.Printf("ParseNat64 for %s: %s\n", labelValues[0], err)
		}
		return err
	}

	ch <- prometheus.MustNewConstMetric(translationsActiveDesc, prometheus.GaugeValue, float64(stats.translationsActive), labelValues...)
	ch <- prometheus.MustNewConstMetric(translationsExpiredDesc, prometheus.CounterValue, float64(stats.translationsExpired), labelValues...)
	ch <- prometheus.MustNewConstMetric(sessionsFoundDesc, prometheus.CounterValue, float64(stats.sessionsFound), labelValues...)
	ch <- prometheus.MustNewConstMetric(sessionsCreatedDesc, prometheus.CounterValue, float64(stats.sessionsCreated), labelValues...)
	ch <- prometheus.MustNewConstMetric(packetsTranslated4to6Desc, prometheus.CounterValue, float64(stats.packetsTranslated4to6), labelValues...)
	ch <- prometheus.MustNewConstMetric(packetsTranslated6to4Desc, prometheus.CounterValue, float64(stats.packetsTranslated6to4), labelValues...)

	return nil
}

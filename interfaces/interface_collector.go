package interfaces

import (
	"log"
	"regexp"

	"github.com/lwlcom/cisco_exporter/dynamiclabels"
	"github.com/lwlcom/cisco_exporter/rpc"

	"github.com/lwlcom/cisco_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
)

const prefix string = "cisco_interface_"

type description struct {
	receiveBytesDesc     *prometheus.Desc
	receiveErrorsDesc    *prometheus.Desc
	receiveDropsDesc     *prometheus.Desc
	receiveBroadcastDesc *prometheus.Desc
	receiveMulticastDesc *prometheus.Desc
	transmitBytesDesc    *prometheus.Desc
	transmitErrorsDesc   *prometheus.Desc
	transmitDropsDesc    *prometheus.Desc
	adminStatusDesc      *prometheus.Desc
	operStatusDesc       *prometheus.Desc
	errorStatusDesc      *prometheus.Desc
	speedDesc            *prometheus.Desc
}

func newDescriptions(dynLabels dynamiclabels.Labels) *description {
	d := &description{}
	l := []string{"target", "name", "description", "mac"}
	l = append(l, dynLabels.Keys()...)

	d.receiveBytesDesc = prometheus.NewDesc(prefix+"receive_bytes", "Received data in bytes", l, nil)
	d.receiveErrorsDesc = prometheus.NewDesc(prefix+"receive_errors", "Number of errors caused by incoming packets", l, nil)
	d.receiveDropsDesc = prometheus.NewDesc(prefix+"receive_drops", "Number of dropped incoming packets", l, nil)
	d.receiveBroadcastDesc = prometheus.NewDesc(prefix+"receive_broadcast", "Received broadcast packets", l, nil)
	d.receiveMulticastDesc = prometheus.NewDesc(prefix+"receive_multicast", "Received multicast packets", l, nil)
	d.transmitBytesDesc = prometheus.NewDesc(prefix+"transmit_bytes", "Transmitted data in bytes", l, nil)
	d.transmitErrorsDesc = prometheus.NewDesc(prefix+"transmit_errors", "Number of errors caused by outgoing packets", l, nil)
	d.transmitDropsDesc = prometheus.NewDesc(prefix+"transmit_drops", "Number of dropped outgoing packets", l, nil)
	d.adminStatusDesc = prometheus.NewDesc(prefix+"admin_up", "Admin operational status", l, nil)
	d.operStatusDesc = prometheus.NewDesc(prefix+"up", "Interface operational status", l, nil)
	d.errorStatusDesc = prometheus.NewDesc(prefix+"error_status", "Admin and operational status differ", l, nil)
	d.speedDesc = prometheus.NewDesc(prefix+"speed", "Interface speed in in bps", l, nil)
	return d
}

type interfaceCollector struct {
	descriptionRe *regexp.Regexp
}

// NewCollector creates a new collector
func NewCollector(descRe *regexp.Regexp) collector.RPCCollector {
	return &interfaceCollector{descriptionRe: descRe}
}

// Name returns the name of the collector
func (*interfaceCollector) Name() string {
	return "Interfaces"
}

// Describe describes the metrics
func (*interfaceCollector) Describe(ch chan<- *prometheus.Desc) {
	d := newDescriptions(nil)
	ch <- d.receiveBytesDesc
	ch <- d.receiveErrorsDesc
	ch <- d.receiveDropsDesc
	ch <- d.receiveBroadcastDesc
	ch <- d.receiveMulticastDesc
	ch <- d.transmitBytesDesc
	ch <- d.transmitDropsDesc
	ch <- d.transmitErrorsDesc
	ch <- d.adminStatusDesc
	ch <- d.operStatusDesc
	ch <- d.errorStatusDesc
	ch <- d.speedDesc
}

// Collect collects metrics from Cisco
func (c *interfaceCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {
	out, err := client.RunCommand("show interface")
	if err != nil {
		return err
	}
	items, err := c.Parse(client.OSType, out)
	if err != nil {
		if client.Debug {
			log.Printf("Parse interfaces for %s: %s\n", labelValues[0], err.Error())
		}
		return nil
	}
	if client.OSType == rpc.IOSXE {
		out, err := client.RunCommand("show vlans")
		if err != nil {
			return err
		}
		vlans, err := c.ParseVlans(client.OSType, out)
		if err != nil {
			if client.Debug {
				log.Printf("Parse vlans for %s: %s\n", labelValues[0], err.Error())
			}
			return nil
		}
		for _, vlan := range vlans {
			for i, item := range items {
				if item.Name == vlan.Name {
					items[i].InputBytes = vlan.InputBytes
					items[i].OutputBytes = vlan.OutputBytes
					break
				}
			}
		}
	}

	for _, item := range items {
		c.collectForInterface(item, ch, labelValues)
	}

	return nil
}

func (c *interfaceCollector) collectForInterface(item Interface, ch chan<- prometheus.Metric, labelValues []string) {
	l := append(labelValues, item.Name, item.Description, item.MacAddress)
	dynLabels := dynamiclabels.ParseDescription(item.Description, c.descriptionRe)
	l = append(l, dynLabels.Values()...)
	d := newDescriptions(dynLabels)

	errorStatus := 0
	if item.AdminStatus != item.OperStatus {
		errorStatus = 1
	}
	adminStatus := 0
	if item.AdminStatus == "up" {
		adminStatus = 1
	}
	operStatus := 0
	if item.OperStatus == "up" {
		operStatus = 1
	}
	ch <- prometheus.MustNewConstMetric(d.receiveBytesDesc, prometheus.GaugeValue, item.InputBytes, l...)
	ch <- prometheus.MustNewConstMetric(d.receiveErrorsDesc, prometheus.GaugeValue, item.InputErrors, l...)
	ch <- prometheus.MustNewConstMetric(d.receiveDropsDesc, prometheus.GaugeValue, item.InputDrops, l...)
	ch <- prometheus.MustNewConstMetric(d.transmitBytesDesc, prometheus.GaugeValue, item.OutputBytes, l...)
	ch <- prometheus.MustNewConstMetric(d.transmitErrorsDesc, prometheus.GaugeValue, item.OutputErrors, l...)
	ch <- prometheus.MustNewConstMetric(d.transmitDropsDesc, prometheus.GaugeValue, item.OutputDrops, l...)
	ch <- prometheus.MustNewConstMetric(d.receiveBroadcastDesc, prometheus.GaugeValue, item.InputBroadcast, l...)
	ch <- prometheus.MustNewConstMetric(d.receiveMulticastDesc, prometheus.GaugeValue, item.InputMulticast, l...)
	ch <- prometheus.MustNewConstMetric(d.adminStatusDesc, prometheus.GaugeValue, float64(adminStatus), l...)
	ch <- prometheus.MustNewConstMetric(d.operStatusDesc, prometheus.GaugeValue, float64(operStatus), l...)
	ch <- prometheus.MustNewConstMetric(d.errorStatusDesc, prometheus.GaugeValue, float64(errorStatus), l...)
	ch <- prometheus.MustNewConstMetric(d.speedDesc, prometheus.GaugeValue, item.Speed, l...)
}

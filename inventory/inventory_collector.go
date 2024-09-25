package inventory

import (
	"errors"
	"log"

	"github.com/lwlcom/cisco_exporter/collector"
	"github.com/lwlcom/cisco_exporter/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

const name_inventory string = "cisco_inventory_item"
const name_transceiver string = "cisco_interface_transceiver"

var (
	inventoryItemDesc   *prometheus.Desc
	transceiverItemDesc *prometheus.Desc
)

func init() {
	l_inv := []string{"target", "name", "description", "part_number", "serial_number"}
	l_transc := []string{"target", "name", "description", "vendor_name", "vendor_part_number", "serial_number"}
	inventoryItemDesc = prometheus.NewDesc(name_inventory, "Hardware inventory info", l_inv, nil)
	transceiverItemDesc = prometheus.NewDesc(name_transceiver, "Transceiver inventory info", l_transc, nil)
}

type inventoryCollector struct {
}

// NewCollector creates a new collector
func NewCollector() collector.RPCCollector {
	return &inventoryCollector{}
}

// Name returns the name of the collector
func (*inventoryCollector) Name() string {
	return "Inventory"
}

// Describe describes the metrics
func (*inventoryCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- inventoryItemDesc
	ch <- transceiverItemDesc
}

// Collect collects metrics from Cisco
func (c *inventoryCollector) Collect(client *rpc.Client, ch chan<- prometheus.Metric, labelValues []string) error {

	switch client.OSType {
	case rpc.IOS, rpc.IOSXE:
		interfaces, err := client.GetInterfaceNames(false)
		if err != nil {
			if client.Debug {
				log.Printf("Get interfaces command on %s: %s\n", labelValues[0], err.Error())
			}
			return nil
		}
		out, err := client.RunCommand("show inventory fru")
		if err != nil {
			if client.Debug {
				log.Printf("show inventory command on %s: %s\n", labelValues[0], err.Error())
			}
			return nil
		}
		inventory_items, transceiver_items, err := c.ParseInventory(client.OSType, out, interfaces)
		if err != nil {
			if client.Debug {
				log.Printf("show inventory parsing on %s: %s\n", labelValues[0], err.Error())
			}
			return nil
		}

		for _, transceiver_item := range transceiver_items {
			out, err = client.RunCommand("show idprom interface " + transceiver_item.Name)
			if err != nil {
				if client.Debug {
					log.Printf("show idprom command on %s %s: %s\n", labelValues[0], transceiver_item.Name, err.Error())
				}
				return nil
			}
			transceiver, err := c.ParseIdprom(client.OSType, transceiver_item.Name, out)
			if err != nil {
				if client.Debug {
					log.Printf("show idprom parsing on %s %s: %s\n", labelValues[0], transceiver_item.Name, err.Error())
				}
				return nil
			}
			l := append(labelValues, transceiver.Name)
			l = append(l, transceiver.Description)
			l = append(l, transceiver.Vendor)
			l = append(l, transceiver.PartNumber)
			l = append(l, transceiver.SerialNumber)
			ch <- prometheus.MustNewConstMetric(transceiverItemDesc, prometheus.GaugeValue, float64(1), l...)
		}

		for _, inventory_item := range inventory_items {
			l := append(labelValues, inventory_item.Name)
			l = append(l, inventory_item.Description)
			l = append(l, inventory_item.PartNumber)
			l = append(l, inventory_item.SerialNumber)
			ch <- prometheus.MustNewConstMetric(inventoryItemDesc, prometheus.GaugeValue, float64(1), l...)

		}

	case rpc.NXOS:
		return errors.New("inventory collector for NXOS not implemented")
	}

	return nil
}

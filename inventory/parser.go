package inventory

import (
	"errors"
	"slices"

	"github.com/lwlcom/cisco_exporter/rpc"
	"github.com/lwlcom/cisco_exporter/util"
)

/* ParseInventory parses cli output and returns two lists:
 * list of FRU items
 * list of transceivers
 */
func (c *inventoryCollector) ParseInventory(ostype string, output string, interface_names []string) ([]InventoryItem, []InventoryItem, error) {
	if ostype != rpc.IOSXE && ostype != rpc.IOS {
		return nil, nil, errors.New("'show inventory' is not implemented for " + ostype)
	}
	items := []InventoryItem{}
	transceivers := []InventoryItem{}

	results_inventory, err := util.ParseTextfsm(templ_inventory, output)
	if err != nil {
		return nil, nil, errors.New("Error parsing via templ_inventory: " + err.Error())
	}
	for _, result := range results_inventory {
		x := InventoryItem{
			Name:        result["NAME"].(string),
			Description: result["DESCR"].(string),
			PartNumber:  result["PID"].(string),
			// VID does not contain interesting info for us
			SerialNumber: result["SN"].(string),
		}
		// some models return short interface names in inventory
		long_name, long_err := util.InterfaceShortToLong(x.Name)
		if slices.Contains(interface_names, x.Name) {
			transceivers = append(transceivers, x)
		} else if long_err == nil && slices.Contains(interface_names, long_name) {
			x.Name = long_name
			transceivers = append(transceivers, x)
		} else {
			items = append(items, x)
		}
	}
	return items, transceivers, nil
}

// ParseIdprom parses cli output and extracts inventory info
func (c *inventoryCollector) ParseIdprom(ostype string, ifname string, output string) (TransceiverItem, error) {
	if ostype != rpc.IOSXE && ostype != rpc.IOS {
		return TransceiverItem{}, errors.New("Idprom data is not implemented for " + ostype)
	}

	results_idprom, err := util.ParseTextfsm(templ_idprom, output)
	if err != nil {
		return TransceiverItem{}, errors.New("Error parsing via templ_idprom: " + err.Error())
	}
	for _, result := range results_idprom {
		x := TransceiverItem{
			Name:         ifname,
			Description:  result["TYPE"].(string),
			Vendor:       result["VENDOR"].(string),
			PartNumber:   result["PID"].(string),
			SerialNumber: result["SN"].(string),
		}
		return x, nil
	}
	return TransceiverItem{}, nil
}

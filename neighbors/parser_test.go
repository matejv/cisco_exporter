package neighbors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseIOSXEIpv4Interfaces(t *testing.T) {
	ostype := "IOSXE"
	out := `#show ip interface brief
Interface              IP-Address      OK? Method Status                Protocol
Vlan1                  100.64.6.1      YES manual up                    up      
TwentyFiveGigE1/0/1    unassigned      YES unset  up                    up       
`

	c := neighborsCollector{}
	interfaces, err := c.ParseInterfacesIPv4(ostype, out)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Vlan1", interfaces[0], "Interface with IPv4")
	assert.NotContains(t, interfaces, "TwentyFiveGigE1/0/1", "Interface without IPv4")
}

func TestParseIOSXEIpv6Interfaces(t *testing.T) {
	ostype := "IOSXE"
	out := `#show ipv6 interface brief
Vlan1                  [up/up]
    unassigned
Vlan2                  [up/down]
    FE80::AA4F:B1FF:FE58:949F
    2001:DB8:F53F:CC::1
`

	c := neighborsCollector{}
	interfaces, err := c.ParseInterfacesIPv6(ostype, out)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Vlan2", interfaces[0], "Interface with IPv6")
	assert.NotContains(t, interfaces, "Vlan1", "Interface without IPv6")
}

func TestParseIOSXEIpv4Neighbors(t *testing.T) {
	ostype := "iosxe"
	out := `#show arp detail | include via
Interface, via Vlan2, last updated 71582 minutes ago.
Interface, via Vlan11, last updated 71582 minutes ago.
Incomplete, via Vlan11, last updated 0 minute ago.
Incomplete, via Vlan11, last updated 0 minute ago.
Incomplete, via Vlan11, last updated 0 minute ago.
Incomplete, via Vlan11, last updated 0 minute ago.
Incomplete, via Vlan11, last updated 0 minute ago.
`

	c := neighborsCollector{}
	var interfaces_data = make(map[string]*InterfaceNeighors)
	interfaces_data["Vlan11"] = &InterfaceNeighors{}
	interfaces_data["Vlan2"] = &InterfaceNeighors{}

	err := c.ParseIPv4Neighbors(ostype, out, interfaces_data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 5.0, interfaces_data["Vlan11"].Incomplete, "Incomplete IPv4")
	assert.Equal(t, 0.0, interfaces_data["Vlan2"].Reachable, "Reachable IPv4")
	assert.NotContains(t, interfaces_data, "TwentyFiveGigE1/0/48", "Unknown interface IPv4")
}

func TestParseIOSXEIpv6Neighbors(t *testing.T) {
	ostype := "iosxe"
	out := `#show ipv6 neighbors 
IPv6 Address                              Age Link-layer Addr State Interface
2001:DB8:EC39:0:5060:3BE5:FD4F:8311         0 00AA.c21e.30cd  REACH Vl2
2001:DB8:EC39:0:5CC5:1F6B:A000:BD61         0 00aa.4da6.0a6a  STALE Vl2
2001:DB8:EC39:0:7196:69BD:9C3:710E          0 00aa.f063.7713  STALE Vl2
2001:DB8:F70E:AA:ECA:FBFF:FE28:75B5         0 00aa.fb28.75b5  STALE Vl4
2001:DB8:F70E:AA:1D17:6BCF:3B7A:8C45        0 00aa.6235.2e60  REACH Vl4
2001:DB8:F70E:AA:38F1:E51F:DF31:C659        0 -               INCMP Vl4
2001:DB8:F70E:AA:4CBC:52C2:AD18:3C82        0 -               INCMP Vl4
2001:DB8:F70E:AA:5944:CD47:32B6:2FF5        0 00AA.5d40.a801  STALE Vl4
2001:DB8:F70E:AA:74BC:C4B9:4C24:4680        0 -               INCMP Vl4
2001:DB8:F70E:AA:81A7:E629:3C31:7CD5        0 00aa.c648.495d  REACH Vl4
2001:DB8:F70E:AA:B051:88B3:DECD:A058        0 00AA.24d2.698e  STALE Vl4
2001:DB8:F70E:AA:C533:89DF:1BDF:852F        0 00aa.2346.e561  STALE Vl4
2001:DB8:F70E:AA:C59C:B68E:B848:A0B3        0 00aa.994b.cca0  STALE Vl4
2001:DB8:F70E:AA:E089:21CD:A1E6:CEB5        0 00aa.880e.5a12  STALE Vl4
`

	c := neighborsCollector{}
	var interfaces_data = make(map[string]*InterfaceNeighors)
	interfaces_data["Vlan2"] = &InterfaceNeighors{}
	interfaces_data["Vlan4"] = &InterfaceNeighors{}

	err := c.ParseIPv6Neighbors(ostype, out, interfaces_data)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2.0, interfaces_data["Vlan2"].Stale, "Stale IPv6 Vlan2")
	assert.Equal(t, 6.0, interfaces_data["Vlan4"].Stale, "Stale IPv6 Vlan4")
	assert.Equal(t, 2.0, interfaces_data["Vlan4"].Reachable, "Reachable IPv6")
	assert.Equal(t, 3.0, interfaces_data["Vlan4"].Incomplete, "Incomplete IPv6")
	assert.NotContains(t, interfaces_data, "Vl4", "Unknown interface IPv6")
}

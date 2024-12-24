package nat64

/*
 * # ASR1000
 * # show nat64 statistics global
 * NAT64 Statistics
 *
 * Total active translations: 238814 (0 static, 238814 dynamic; 238814 extended)
 * Sessions found: 761031179766
 * Sessions created: 4144454051
 * Expired translations: 4149896490
 * Global Stats:
 *
 * 	Packets translated (IPv4 -> IPv6)
 * 	   Stateless: 0
 * 	   Stateful: 536010887327
 * 	   MAP-T: 0
 * 	Packets translated (IPv6 -> IPv4)
 * 	   Stateless: 0
 * 	   Stateful: 229167550920
 * 	   MAP-T: 0
 */
var templ_nat64 = `# show nat64 statistics global
Value nat64_total_active_translations (\d+)
Value nat64_sessions_found (\d+)
Value nat64_sessions_created (\d+)
Value nat64_expired_translations (\d+)
Value nat64_ipv4_ipv6_translated_packets (\d+)
Value nat64_ipv6_ipv4_translated_packets (\d+)

Start
  ^Total\s+active\s+translations:\s+${nat64_total_active_translations}
  ^Sessions\s+found:\s+${nat64_sessions_found}
  ^Sessions\s+created:\s+${nat64_sessions_created}
  ^Expired\s+translations:\s+${nat64_expired_translations}
  ^\s+Packets translated \(IPv4 \-\> IPv6\) -> IP46
  ^\s+Packets translated \(IPv6 \-\> IPv4\) -> IP64

IP46
  ^\s+Stateful:\s+${nat64_ipv4_ipv6_translated_packets} -> Start

IP64
  ^\s+Stateful:\s+${nat64_ipv6_ipv4_translated_packets} -> Record
`

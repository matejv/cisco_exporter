package inventory

/*
 * # C9500
 * NAME: "Slot 1 Supervisor", DESCR: "Cisco Catalyst 9500 Series Router"
 * PID: C9500-48Y4C       , VID: V04  , SN: 456
 *
 * NAME: "TwentyFiveGigE1/0/1", DESCR: "GE T"
 * PID: XSUEG1-M1RN-GC    , VID:      , SN: 123
 */
var templ_inventory = `# show inventory
Value NAME (.*)
Value DESCR (.*)
Value PID (([\S+]+|.*))
Value VID (.*)
Value SN ([\w+\d+]+)

Start
 ^NAME:\s+"${NAME}",\s+DESCR:\s+"${DESCR}"
 ^PID:\s+${PID}.*,.*VID:\s+${VID},.*SN:\s+${SN} -> Record
 ^PID:\s+,.*VID:\s+${VID},.*SN: -> Record
 ^PID:\s+${PID}.*,.*VID:\s+${VID},.*SN: -> Record
 ^PID:\s+,.*VID:\s+${VID},.*SN:\s+${SN} -> Record
 ^PID:\s+${PID}.*,.*VID:\s+${VID}.*
 ^PID:\s+,.*VID:\s+${VID}.*
 ^.*SN:\s+${SN} -> Record
 ^.*SN: -> Record
`

/*
 * # C9500
 * IDPROM for transceiver HundredGigE1/0/49:
 *   Description                               = QSFP28 optics (type 134)
 *   Transceiver Type:                         = QSFP 100GE DWDM2 (462)
 *   Product Identifier (PID)                  =
 *   Vendor Revision                           = 10
 *   Serial Number (SN)                        = 1234
 *   Vendor Name                               = INPHI CORP
 *   Vendor OUI (IEEE company ID)              = 00.21.B8 (8632)
 *   CLEI code                                 =
 *   Cisco part number                         =
 *   Device State                              = Enabled.
 *   Date code (yy/mm/dd)                      = 20/11/06
 *   Connector type                            = LC
 *   Encoding                                  =
 *   Nominal bitrate per channel               = 25GE (25500 Mbits/s)
 *
 * # C9200
 * General SFP Information
 * -----------------------------------------------
 * Identifier            :   SFP/SFP+
 * Ext.Identifier        :   SFP function is defined by two-wire interface ID only
 * Connector             :   LC connector
 * Transceiver
 *  10/40GE Comp code       :   10G BASE-LR
 *  SONET Comp code      :   Unknown
 *  GE Comp code         :   Unknown
 *  Link length          :   Unknown
 *  Technology           :   Unknown
 *  Media                :   Single Mode
 *  Speed                :   Unknown
 * Encoding              :   64B/66B
 * BR_Nominal            :   10300 Mbps
 * Length(9um)-km        :   10 km
 * Length(9um)           :   10000 m
 * Length(50um)          :   GBIC does not support 50 micron multi mode OM2 fibre
 * Length(62.5um)        :   GBIC does not support 62.5 micron multi mode OM1 fibre
 * Length(Copper)        :   GBIC does not support 50 micron multi mode OM4 fibre
 * Vendor Name           :   XenOpt
 * Vendor Part Number    :   XTS31A-10LY-TC
 * Vendor Revision       :   0x56 0x30 0x32 0x20
 * Vendor Serial Number  :   1234
 * Wavelength            :   1310.00 nm
 * CC_BASE               :   0x7E
 * -----------------------------------------------
 *
 */
var templ_idprom = `# show idprom interface ...
Value NAME (\S+)
Value DESCRIPTION (([^\(]+[^\s\(])?)
Value TYPE (([^\(]+[^\s\(])?)
Value PID (([^\(]+[^\s\(])?)
Value SN (([^\(]+[^\s\(])?)
Value VENDOR (([^\(]+[^\s\(])?)

Start
 ^IDPROM for transceiver ${NAME}:
 ^\s*Description\s+[=:]\s*${DESCRIPTION}(?: \(.+)?
 ^\s*Transceiver Type.*\s+[=:]\s*${TYPE}(?: \(.+)?
 ^\s*Product Identifier.*\s+[=:]\s*${PID}(?: \(.+)?
 ^\s*(?:Vendor )?Serial Number.*\s+[=:]\s*${SN}\s*
 ^\s*Vendor Name.*\s+[=:]\s*${VENDOR}\s*
 ^\s*Vendor Part Number.*\s+[=:]\s*${PID}(?: \(.+)?
 # skip Comp unknown codes
 ^\s*.* Comp [Cc]ode\s+[=:]\s*Unknown -> Next
 ^\s*.* Comp [Cc]ode\s+[=:]\s*${TYPE}
`

// Copyright (c) 2018 Zededa, Inc.
// SPDX-License-Identifier: Apache-2.0

// Also blocks the VNC ports (5900...) if ssh is blocked
// Always blocks 4822
// Also always blocks port 8080

package iptables

import (
	"fmt"

	"github.com/lf-edge/eve/pkg/pillar/base"
)

// ControlProtocolMarkingIDMap : Map describing the control flow
// marking values that we intend to use.
// XXX only used by nim hence no concurrency. But LockedStringMap would be better
var ControlProtocolMarkingIDMap = map[string]string{
	// INPUT flows for HTTP, SSH & GUACAMOLE
	"in_http_ssh_guacamole": "1",
	// INPUT flows for VNC
	"in_vnc": "2",
	// There was some feature here that used marking values "3" & "4".
	// Marking values "3" & "4" are unused as of now.

	// OUTPUT flows for all types
	"out_all": "5",
	// App initiated UDP flows towards dom0 for DHCP
	"app_dhcp": "6",
	// App initiated TCP/UDP flows towards dom0 for DNS
	"app_dns": "7",
	// VPN control packets
	"in_vpn_control": "8",
	// ICMP and ICMPv6
	"in_icmp": "9",
	// DHCP packets originating from outside
	// (e.g. DHCP multicast requests from other devices on the same network)
	"in_dhcp": "10",
}

func UpdateSshAccess(log *base.LogObject, enable bool, first bool) {

	log.Functionf("updateSshAccess(enable %v first %v)\n",
		enable, first)

	if first {
		// Always blocked
		dropPortRange(log, 8080, 8080)
		allowLocalPortRange(log, 4822, 4822)
		allowLocalPortRange(log, 5900, 5999)
		markControlFlows(log)
	}
	if enable {
		allowPortRange(log, 22, 22)
	} else {
		dropPortRange(log, 22, 22)
	}
}

func UpdateVncAccess(log *base.LogObject, enable bool) {

	log.Functionf("updateVncAccess(enable %v\n", enable)

	if enable {
		allowPortRange(log, 5900, 5999)
	} else {
		dropPortRange(log, 5900, 5999)
	}
}

func allowPortRange(log *base.LogObject, startPort int, endPort int) {
	log.Functionf("allowPortRange(%d, %d)\n", startPort, endPort)
	// Delete these rules
	// iptables -D INPUT -p tcp --dport 22 -j REJECT --reject-with tcp-reset
	// ip6tables -D INPUT -p tcp --dport 22 -j REJECT --reject-with tcp-reset
	var portStr string
	if startPort == endPort {
		portStr = fmt.Sprintf("%d", startPort)
	} else {
		portStr = fmt.Sprintf("%d:%d", startPort, endPort)
	}
	IptableCmd(log, "-D", "INPUT", "-p", "tcp", "--dport", portStr, "-j", "REJECT", "--reject-with", "tcp-reset")
	Ip6tableCmd(log, "-D", "INPUT", "-p", "tcp", "--dport", portStr, "-j", "REJECT", "--reject-with", "tcp-reset")
}

// Like above but allow for 127.0.0.1 to 127.0.0.1 and block for other IPs
func allowLocalPortRange(log *base.LogObject, startPort int, endPort int) {
	log.Functionf("allowLocalPortRange(%d, %d)\n", startPort, endPort)
	// Add these rules
	// iptables -A INPUT -p tcp -s 127.0.0.1 -d 127.0.0.1 --dport 22 -j ACCEPT
	// iptables -A INPUT -p tcp --dport 22 -j REJECT --reject-with tcp-reset
	// iptables -A INPUT -p tcp -s ::1 -d ::1 --dport 22 -j ACCEPT
	// ip6tables -A INPUT -p tcp --dport 22 -j REJECT --reject-with tcp-reset
	var portStr string
	if startPort == endPort {
		portStr = fmt.Sprintf("%d", startPort)
	} else {
		portStr = fmt.Sprintf("%d:%d", startPort, endPort)
	}
	IptableCmd(log, "-A", "INPUT", "-p", "tcp", "--dport", portStr,
		"-s", "127.0.0.1", "-d", "127.0.0.1", "-j", "ACCEPT")
	IptableCmd(log, "-A", "INPUT", "-p", "tcp", "--dport", portStr,
		"-j", "REJECT", "--reject-with", "tcp-reset")
	Ip6tableCmd(log, "-A", "INPUT", "-p", "tcp", "--dport", portStr,
		"-s", "::1", "-d", "::1", "-j", "ACCEPT")
	Ip6tableCmd(log, "-A", "INPUT", "-p", "tcp", "--dport", portStr,
		"-j", "REJECT", "--reject-with", "tcp-reset")
}

func dropPortRange(log *base.LogObject, startPort int, endPort int) {
	log.Functionf("dropPortRange(%d, %d)\n", startPort, endPort)
	// Add these rules
	// iptables -A INPUT -p tcp --dport 22 -j REJECT --reject-with tcp-reset
	// ip6tables -A INPUT -p tcp --dport 22 -j REJECT --reject-with tcp-reset
	var portStr string
	if startPort == endPort {
		portStr = fmt.Sprintf("%d", startPort)
	} else {
		portStr = fmt.Sprintf("%d:%d", startPort, endPort)
	}
	IptableCmd(log, "-A", "INPUT", "-p", "tcp", "--dport", portStr, "-j", "REJECT", "--reject-with", "tcp-reset")
	Ip6tableCmd(log, "-A", "INPUT", "-p", "tcp", "--dport", portStr, "-j", "REJECT", "--reject-with", "tcp-reset")
}

// With flow monitoring happening, any unmarked connections:
// 1) Not matching any of the INPUT ACL in PREROUTING chain
// 2) Not initiated by applications
// will be dropped (sent out of dummy interface). But, we still
// want some control protocols running on dom0 to run. We mark such
// connections with markings from reserved space and let the ACLs
// in INPUT chain make the ACCEPT/DROP/REJECT decisions.
func markControlFlows(log *base.LogObject) {
	// Mark HTTP, ssh and guacamole packets
	// Pick flow marking values 1, 2, 3 from the reserved space.
	portStr := "22,4822"
	IptableCmd(log, "-t", "mangle", "-I", "PREROUTING", "1", "-p", "tcp",
		"--match", "multiport", "--dports", portStr,
		"-j", "CONNMARK", "--set-mark", ControlProtocolMarkingIDMap["in_http_ssh_guacamole"])

	Ip6tableCmd(log, "-t", "mangle", "-I", "PREROUTING", "1", "-p", "tcp",
		"--match", "multiport", "--dports", portStr,
		"-j", "CONNMARK", "--set-mark", ControlProtocolMarkingIDMap["in_http_ssh_guacamole"])

	// Mark VNC packets
	portStr = "5900:5999"
	IptableCmd(log, "-t", "mangle", "-I", "PREROUTING", "2", "-p", "tcp",
		"--match", "multiport", "--dports", portStr,
		"-j", "CONNMARK", "--set-mark", ControlProtocolMarkingIDMap["in_vnc"])

	Ip6tableCmd(log, "-t", "mangle", "-I", "PREROUTING", "2", "-p", "tcp",
		"--match", "multiport", "--dports", portStr,
		"-j", "CONNMARK", "--set-mark", ControlProtocolMarkingIDMap["in_vnc"])

	// Mark strongswan VPN control packets
	portStr = "4500,500"
	IptableCmd(log, "-t", "mangle", "-I", "PREROUTING", "3", "-p", "udp",
		"--match", "multiport", "--dports", portStr,
		"-j", "CONNMARK", "--set-mark", ControlProtocolMarkingIDMap["in_vpn_control"])
	// Allow esp protocol packets
	IptableCmd(log, "-t", "mangle", "-I", "PREROUTING", "4", "-p", "esp",
		"-j", "CONNMARK", "--set-mark", ControlProtocolMarkingIDMap["in_vpn_control"])

	// Allow all ICMP and ICMPv6
	IptableCmd(log, "-t", "mangle", "-I", "PREROUTING", "5", "-p", "icmp",
		"-j", "CONNMARK", "--set-mark", ControlProtocolMarkingIDMap["in_icmp"])
	Ip6tableCmd(log, "-t", "mangle", "-I", "PREROUTING", "3", "-p", "icmpv6",
		"-j", "CONNMARK", "--set-mark", ControlProtocolMarkingIDMap["in_icmp"])

	// Allow all incoming DHCP traffic
	IptableCmd(log, "-t", "mangle", "-I", "PREROUTING", "6", "-p", "udp",
		"--dport", "bootps:bootpc",
		"-j", "CONNMARK", "--set-mark", ControlProtocolMarkingIDMap["in_dhcp"])

	// Mark all un-marked local traffic generated by local services.
	IptableCmd(log, "-t", "mangle", "-I", "OUTPUT",
		"-j", "CONNMARK", "--restore-mark")
	IptableCmd(log, "-t", "mangle", "-A", "OUTPUT", "-m", "mark", "!", "--mark", "0",
		"-j", "ACCEPT")
	IptableCmd(log, "-t", "mangle", "-A", "OUTPUT",
		"-j", "MARK", "--set-mark", ControlProtocolMarkingIDMap["out_all"])
	IptableCmd(log, "-t", "mangle", "-A", "OUTPUT",
		"-j", "CONNMARK", "--save-mark")
	//IptableCmd(log, "-t", "mangle", "-A", "OUTPUT",
	//	"-j", "CONNMARK", "--set-mark", "5")

	// XXX Later when we support Lisp we should have the above marking
	// checks for IPv6 also.
	Ip6tableCmd(log, "-t", "mangle", "-I", "OUTPUT", "1",
		"-j", "CONNMARK", "--set-mark", ControlProtocolMarkingIDMap["out_all"])

}

// Copyright (c) 2016 Intel Corporation
// Copyright (c) 2019 Huawei Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package persistapi

import (
	"github.com/vishvananda/netlink"
)

// ============= sandbox level resources =============

type NetworkInterface struct {
	Name     string
	HardAddr string
	Addrs    []netlink.Addr
}

// TapInterface defines a tap interface
type TapInterface struct {
	ID       string
	Name     string
	TAPIface NetworkInterface
	// remove VMFds and VhostFds
}

// NetworkInterfacePair defines a pair between VM and virtual network interfaces.
type NetworkInterfacePair struct {
	TapInterface
	VirtIface            NetworkInterface
	NetInterworkingModel int
}

type PhysicalEndpoint struct {
	BDF            string
	Driver         string
	VendorDeviceID string
}

type MacvtapEndpoint struct {
	// This is for showing information.
	// Remove this field won't impact anything.
	PCIAddr string
}

type TapEndpoint struct {
	TapInterface TapInterface
}

type BridgedMacvlanEndpoint struct {
	NetPair NetworkInterfacePair
}

type VethEndpoint struct {
	NetPair NetworkInterfacePair
}

type IPVlanEndpoint struct {
	NetPair NetworkInterfacePair
}

type VhostUserEndpoint struct {
	// This is for showing information.
	// Remove these fields won't impact anything.
	IfaceName string
	PCIAddr   string
}

// NetworkEndpoint contains network interface information
type NetworkEndpoint struct {
	Type string

	// One and only one of these below are not nil according to Type.
	Physical       *PhysicalEndpoint       `json:",omitempty"`
	Veth           *VethEndpoint           `json:",omitempty"`
	VhostUser      *VhostUserEndpoint      `json:",omitempty"`
	BridgedMacvlan *BridgedMacvlanEndpoint `json:",omitempty"`
	Macvtap        *MacvtapEndpoint        `json:",omitempty"`
	Tap            *TapEndpoint            `json:",omitempty"`
	IPVlan         *IPVlanEndpoint         `json:",omitempty"`
}

// NetworkInfo contains network information of sandbox
type NetworkInfo struct {
	NetNsPath    string
	NetmonPID    int
	NetNsCreated bool
	Endpoints    []NetworkEndpoint
}

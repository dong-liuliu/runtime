// Copyright (c) 2017-2018 Intel Corporation
// Copyright (c) 2018-2019 Huawei Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package drivers

import (
	//"encoding/hex"

	"github.com/kata-containers/runtime/virtcontainers/device/api"
	"github.com/kata-containers/runtime/virtcontainers/device/config"
	persistapi "github.com/kata-containers/runtime/virtcontainers/persist/api"
	"github.com/kata-containers/runtime/virtcontainers/utils"
)

// VhostUserBlkDevice is a block vhost-user based device
type VhostUserBlkDevice struct {
	*GenericDevice
	VhostUserDeviceAttrs *config.VhostUserDeviceAttrs
}

// NewVhostUserBlkDevice creates a new vhost-user block device based on DeviceInfo
func NewVhostUserBlkDevice(devInfo *config.DeviceInfo) *VhostUserBlkDevice {
	return &VhostUserBlkDevice{
		GenericDevice: &GenericDevice{
			ID:         devInfo.ID,
			DeviceInfo: devInfo,
		},
	}
}

//
// VhostUserBlkDevice's implementation of the device interface:
//

// Attach is standard interface of api.Device, it's used to add device to some
// DeviceReceiver
func (device *VhostUserBlkDevice) Attach(devReceiver api.DeviceReceiver) (err error) {
	skip, err := device.bumpAttachCount(true)
	if err != nil {
		return err
	}
	if skip {
		return nil
	}

	// Increment the block index for the sandbox. This is used to determine the name
	// for the block device in the case where the block device is used as container
	// rootfs and the predicted block device name needs to be provided to the agent.
	_, err = devReceiver.GetAndSetSandboxBlockIndex()

	defer func() {
		if err != nil {
			devReceiver.DecrementSandboxBlockIndex()
			device.bumpAttachCount(false)
		}
	}()

	if err != nil {
		return err
	}

	vAttrs := &config.VhostUserDeviceAttrs{
		DevID:      utils.MakeNameID("blk", device.DeviceInfo.ID, maxDevIDSize),
		SocketPath: device.DeviceInfo.HostPath,
		//Type:       VhostUserBlk,
	}

	//return devReceiver.AppendDevice(device)
	deviceLogger().WithField("device", device.DeviceInfo.HostPath).WithField("SocketPath", vAttrs.SocketPath).Infof("Attaching %s device", config.VhostUserBlk)
	device.VhostUserDeviceAttrs = vAttrs
	if err = devReceiver.HotplugAddDevice(device, config.VhostUserBlk); err != nil {
		return err
	}

	return nil
}

// Detach is standard interface of api.Device, it's used to remove device from some
// DeviceReceiver
func (device *VhostUserBlkDevice) Detach(devReceiver api.DeviceReceiver) error {
	_, err := device.bumpAttachCount(false)
	return err
}

// DeviceType is standard interface of api.Device, it returns device type
func (device *VhostUserBlkDevice) DeviceType() config.DeviceType {
	return config.VhostUserBlk
}

// GetDeviceInfo returns device information used for creating
func (device *VhostUserBlkDevice) GetDeviceInfo() interface{} {
	return device.VhostUserDeviceAttrs
}

// Save converts Device to DeviceState
func (device *VhostUserBlkDevice) Save() persistapi.DeviceState {
	ds := device.GenericDevice.Save()
	ds.Type = string(device.DeviceType())

	vAttr := device.VhostUserDeviceAttrs
	if vAttr != nil {
		ds.VhostUserDev = &persistapi.VhostUserDeviceAttrs{
			DevID:      vAttr.DevID,
			SocketPath: vAttr.SocketPath,
			Type:       string(vAttr.Type),
		}
	}
	return ds
}

// Load loads DeviceState and converts it to specific device
func (device *VhostUserBlkDevice) Load(ds persistapi.DeviceState) {
	device.GenericDevice = &GenericDevice{}
	device.GenericDevice.Load(ds)

	dev := ds.VhostUserDev
	if dev == nil {
		return
	}

	device.VhostUserDeviceAttrs = &config.VhostUserDeviceAttrs{
		DevID:      dev.DevID,
		SocketPath: dev.SocketPath,
		Type:       config.DeviceType(dev.Type),
	}
}

// It should implement GetAttachCount() and DeviceID() as api.Device implementation
// here it shares function from *GenericDevice so we don't need duplicate codes

/*
Copyright (c) 2023 Fred78290, Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Based on govc sources https://github.com/vmware/govmomi/govc
*/

package object

import (
	"context"
	"fmt"

	"github.com/Fred78290/govmrest/vim25"
	"github.com/Fred78290/govmrest/vim25/methods"
	"github.com/Fred78290/vmrest-go-client/client/model"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type VirtualMachine struct {
	Common
}

func NewVirtualMachine(c *vim25.Client, ref types.ManagedObjectReference) *VirtualMachine {
	return &VirtualMachine{
		Common: NewCommon(c, ref),
	}
}

// removeKey is a helper function for removing a specific file key from a list
// of keys associated with disks attached to a virtual machine.
func removeKey(l *[]int, key int) {
	for i, k := range *l {
		if k == key {
			*l = append((*l)[:i], (*l)[i+1:]...)
			break
		}
	}
}

func (v VirtualMachine) PowerState(ctx context.Context) (types.VirtualMachinePowerState, error) {

	return types.VirtualMachinePowerStatePoweredOn, nil
}

func (v VirtualMachine) PowerOn(ctx context.Context) error {
	return nil
}

func (v VirtualMachine) PowerOff(ctx context.Context) error {
	return nil
}

func (v VirtualMachine) Reset(ctx context.Context) error {
	return nil
}

func (v VirtualMachine) Suspend(ctx context.Context) error {
	return nil
}

func (v VirtualMachine) ShutdownGuest(ctx context.Context) error {
	return nil
}

func (v VirtualMachine) RebootGuest(ctx context.Context) error {
	return nil
}

func (v VirtualMachine) Destroy(ctx context.Context) error {
	return nil
}

func (v VirtualMachine) Clone(ctx context.Context, folder *Folder, name string, config types.VirtualMachineCloneSpec) error {
	return nil
}

func (v VirtualMachine) Customize(ctx context.Context, spec types.CustomizationSpec) error {
	return nil
}

func (v VirtualMachine) Reconfigure(ctx context.Context, config types.VirtualMachineConfigSpec) error {
	return nil
}

// WaitForIP waits for the VM guest.ipAddress property to report an IP address.
// Waits for an IPv4 address if the v4 param is true.
func (v VirtualMachine) WaitForIP(ctx context.Context, v4 ...bool) (string, error) {
	var ip string

	return ip, nil
}

// Device returns the VirtualMachine's config.hardware.device property.
func (v VirtualMachine) Device(ctx context.Context) (VirtualDeviceList, error) {
	var o mo.VirtualMachine

	err := v.Properties(ctx, v.Reference(), []string{"config.hardware.device", "summary.runtime.connectionState"}, &o)
	if err != nil {
		return nil, err
	}

	// Quoting the SDK doc:
	//   The virtual machine configuration is not guaranteed to be available.
	//   For example, the configuration information would be unavailable if the server
	//   is unable to access the virtual machine files on disk, and is often also unavailable
	//   during the initial phases of virtual machine creation.
	if o.Config == nil {
		return nil, fmt.Errorf("%s Config is not available, connectionState=%s",
			v.Reference(), o.Summary.Runtime.ConnectionState)
	}

	return VirtualDeviceList(o.Config.Hardware.Device), nil
}

func diskFileOperation(op types.VirtualDeviceConfigSpecOperation, fop types.VirtualDeviceConfigSpecFileOperation, device types.BaseVirtualDevice) types.VirtualDeviceConfigSpecFileOperation {
	if disk, ok := device.(*types.VirtualDisk); ok {
		// Special case to attach an existing disk
		if op == types.VirtualDeviceConfigSpecOperationAdd && disk.CapacityInKB == 0 && disk.CapacityInBytes == 0 {
			childDisk := false
			if b, ok := disk.Backing.(*types.VirtualDiskFlatVer2BackingInfo); ok {
				childDisk = b.Parent != nil
			}

			if !childDisk {
				fop = "" // existing disk
			}
		}
		return fop
	}

	return ""
}

func (v VirtualMachine) configureDevice(ctx context.Context, op types.VirtualDeviceConfigSpecOperation, fop types.VirtualDeviceConfigSpecFileOperation, devices ...types.BaseVirtualDevice) error {
	spec := types.VirtualMachineConfigSpec{}

	for _, device := range devices {
		config := &types.VirtualDeviceConfigSpec{
			Device:        device,
			Operation:     op,
			FileOperation: diskFileOperation(op, fop, device),
		}

		spec.DeviceChange = append(spec.DeviceChange, config)
	}

	task, err := v.Reconfigure(ctx, spec)
	if err != nil {
		return err
	}

	return task.Wait(ctx)
}

// AddDevice adds the given devices to the VirtualMachine
func (v VirtualMachine) AddDevice(ctx context.Context, device ...types.BaseVirtualDevice) error {
	return v.configureDevice(ctx, types.VirtualDeviceConfigSpecOperationAdd, types.VirtualDeviceConfigSpecFileOperationCreate, device...)
}

// EditDevice edits the given (existing) devices on the VirtualMachine
func (v VirtualMachine) EditDevice(ctx context.Context, device ...types.BaseVirtualDevice) error {
	return v.configureDevice(ctx, types.VirtualDeviceConfigSpecOperationEdit, types.VirtualDeviceConfigSpecFileOperationReplace, device...)
}

// RemoveDevice removes the given devices on the VirtualMachine
func (v VirtualMachine) RemoveDevice(ctx context.Context, keepFiles bool, device ...types.BaseVirtualDevice) error {
	fop := types.VirtualDeviceConfigSpecFileOperationDestroy
	if keepFiles {
		fop = ""
	}
	return v.configureDevice(ctx, types.VirtualDeviceConfigSpecOperationRemove, fop, device...)
}

// AttachDisk attaches the given disk to the VirtualMachine
func (v VirtualMachine) AttachDisk(ctx context.Context, id string, datastore *Datastore, controllerKey int32, unitNumber int32) error {
	req := types.AttachDisk_Task{
		This:          v.Reference(),
		DiskId:        types.ID{Id: id},
		Datastore:     datastore.Reference(),
		ControllerKey: controllerKey,
		UnitNumber:    &unitNumber,
	}

	res, err := methods.AttachDisk_Task(ctx, v.c, &req)
	if err != nil {
		return err
	}

	task := NewTask(v.c, res.Returnval)
	return task.Wait(ctx)
}

// DetachDisk detaches the given disk from the VirtualMachine
func (v VirtualMachine) DetachDisk(ctx context.Context, id string) error {
	req := types.DetachDisk_Task{
		This:   v.Reference(),
		DiskId: types.ID{Id: id},
	}

	res, err := methods.DetachDisk_Task(ctx, v.c, &req)
	if err != nil {
		return err
	}

	task := NewTask(v.c, res.Returnval)
	return task.Wait(ctx)
}

// BootOptions returns the VirtualMachine's config.bootOptions property.
func (v VirtualMachine) BootOptions(ctx context.Context) (*types.VirtualMachineBootOptions, error) {
	var o mo.VirtualMachine

	return o.Config.BootOptions, nil
}

// SetBootOptions reconfigures the VirtualMachine with the given options.
func (v VirtualMachine) SetBootOptions(ctx context.Context, options *types.VirtualMachineBootOptions) error {
	spec := types.VirtualMachineConfigSpec{}

	spec.BootOptions = options

	return v.Reconfigure(ctx, spec)
}

// IsToolsRunning returns true if VMware Tools is currently running in the guest OS, and false otherwise.
func (v VirtualMachine) IsToolsRunning(ctx context.Context) (bool, error) {
	var o mo.VirtualMachine

	return o.Guest.ToolsRunningStatus == string(types.VirtualMachineToolsRunningStatusGuestToolsRunning), nil
}

// Wait for the VirtualMachine to change to the desired power state.
func (v VirtualMachine) WaitForPowerState(ctx context.Context, state types.VirtualMachinePowerState) error {

	return nil
}

func (v VirtualMachine) Unregister(ctx context.Context) error {
	return nil
}

// QueryEnvironmentBrowser is a helper to get the environmentBrowser property.
func (v VirtualMachine) QueryConfigTarget(ctx context.Context) (*model.VmRestrictionsInformation, error) {
	return nil, nil
}

func (v VirtualMachine) UpgradeVM(ctx context.Context, version string) error {
	return nil
}

// UUID is a helper to get the UUID of the VirtualMachine managed object.
// This method returns an empty string if an error occurs when retrieving UUID from the VirtualMachine object.
func (v VirtualMachine) UUID(ctx context.Context) string {
	return ""
}

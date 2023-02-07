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

package vm

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/vmware/govmomi/govc/cli"
	"github.com/vmware/govmomi/govc/flags"
	"github.com/vmware/govmomi/vim25/types"
)

type extraConfig []types.BaseOptionValue

func (e *extraConfig) String() string {
	return fmt.Sprintf("%v", *e)
}

func (e *extraConfig) Set(v string) error {
	r := strings.SplitN(v, "=", 2)
	if len(r) < 2 {
		return fmt.Errorf("failed to parse extraConfig: %s", v)
	}
	*e = append(*e, &types.OptionValue{Key: r[0], Value: r[1]})
	return nil
}

type extraConfigFile []types.BaseOptionValue

func (e *extraConfigFile) String() string {
	return fmt.Sprintf("%v", *e)
}

func (e *extraConfigFile) Set(v string) error {
	r := strings.SplitN(v, "=", 2)
	if len(r) < 2 {
		return fmt.Errorf("failed to parse extraConfigFile: %s", v)
	}

	var fileContents = ""
	if len(r[1]) > 0 {
		contents, err := os.ReadFile(r[1])
		if err != nil {
			return fmt.Errorf("failed to parse extraConfigFile '%s': %w", v, err)
		}
		fileContents = string(contents)
	}

	*e = append(*e, &types.OptionValue{Key: r[0], Value: fileContents})
	return nil
}

type change struct {
	*flags.VirtualMachineFlag
	*flags.ResourceAllocationFlag

	types.VirtualMachineConfigSpec
	extraConfig     extraConfig
	extraConfigFile extraConfigFile
	Latency         string
	hwUpgradePolicy string
}

func init() {
	cli.Register("vm.change", &change{})
}

var latencyLevels = []string{
	string(types.LatencySensitivitySensitivityLevelLow),
	string(types.LatencySensitivitySensitivityLevelNormal),
	string(types.LatencySensitivitySensitivityLevelHigh),
}

var hwUpgradePolicies = []string{
	string(types.ScheduledHardwareUpgradeInfoHardwareUpgradePolicyOnSoftPowerOff),
	string(types.ScheduledHardwareUpgradeInfoHardwareUpgradePolicyNever),
	string(types.ScheduledHardwareUpgradeInfoHardwareUpgradePolicyAlways),
}

func (cmd *change) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.VirtualMachineFlag, ctx = flags.NewVirtualMachineFlag(ctx)
	cmd.VirtualMachineFlag.Register(ctx, f)

	cmd.CpuAllocation = &types.ResourceAllocationInfo{Shares: new(types.SharesInfo)}
	cmd.MemoryAllocation = &types.ResourceAllocationInfo{Shares: new(types.SharesInfo)}
	cmd.ResourceAllocationFlag = flags.NewResourceAllocationFlag(cmd.CpuAllocation, cmd.MemoryAllocation)
	cmd.ResourceAllocationFlag.ExpandableReservation = false
	cmd.ResourceAllocationFlag.Register(ctx, f)

	f.Int64Var(&cmd.MemoryMB, "m", 0, "Size in MB of memory")
	f.Var(flags.NewInt32(&cmd.NumCPUs), "c", "Number of CPUs")
	f.StringVar(&cmd.GuestId, "g", "", "Guest OS")
	f.StringVar(&cmd.Name, "name", "", "Display name")
	f.StringVar(&cmd.Latency, "latency", "", fmt.Sprintf("Latency sensitivity (%s)", strings.Join(latencyLevels, "|")))
	f.StringVar(&cmd.Annotation, "annotation", "", "VM description")
	f.StringVar(&cmd.Uuid, "uuid", "", "BIOS UUID")
	f.Var(&cmd.extraConfig, "e", "ExtraConfig. <key>=<value>")
	f.Var(&cmd.extraConfigFile, "f", "ExtraConfig. <key>=<absolute path to file>")

	f.Var(flags.NewOptionalBool(&cmd.NestedHVEnabled), "nested-hv-enabled", "Enable nested hardware-assisted virtualization")
	cmd.Tools = &types.ToolsConfigInfo{}
	f.Var(flags.NewOptionalBool(&cmd.Tools.SyncTimeWithHost), "sync-time-with-host", "Enable SyncTimeWithHost")
	f.Var(flags.NewOptionalBool(&cmd.VPMCEnabled), "vpmc-enabled", "Enable CPU performance counters")
	f.Var(flags.NewOptionalBool(&cmd.MemoryHotAddEnabled), "memory-hot-add-enabled", "Enable memory hot add")
	f.Var(flags.NewOptionalBool(&cmd.MemoryReservationLockedToMax), "memory-pin", "Reserve all guest memory")
	f.Var(flags.NewOptionalBool(&cmd.CpuHotAddEnabled), "cpu-hot-add-enabled", "Enable CPU hot add")

	f.StringVar(&cmd.hwUpgradePolicy, "scheduled-hw-upgrade-policy", "", fmt.Sprintf("Schedule hardware upgrade policy (%s)", strings.Join(hwUpgradePolicies, "|")))
}

func (cmd *change) Description() string {
	return `Change VM configuration.

To add ExtraConfig variables that can read within the guest, use the 'guestinfo.' prefix.

Examples:
  govc vm.change -vm $vm -mem.reservation 2048
  govc vm.change -vm $vm -e smc.present=TRUE -e ich7m.present=TRUE
  # Enable both cpu and memory hotplug on a guest:
  govc vm.change -vm $vm -cpu-hot-add-enabled -memory-hot-add-enabled
  govc vm.change -vm $vm -e guestinfo.vmname $vm
  # Read the contents of a file and use them as ExtraConfig value
  govc vm.change -vm $vm -f guestinfo.data="$(realpath .)/vmdata.config"
  # Read the variable set above inside the guest:
  vmware-rpctool "info-get guestinfo.vmname"
  govc vm.change -vm $vm -latency high
  govc vm.change -vm $vm -latency normal
  govc vm.change -vm $vm -uuid 4139c345-7186-4924-a842-36b69a24159b
  govc vm.change -vm $vm -scheduled-hw-upgrade-policy always`
}

func (cmd *change) Process(ctx context.Context) error {
	if err := cmd.VirtualMachineFlag.Process(ctx); err != nil {
		return err
	}
	return nil
}

func (cmd *change) Run(ctx context.Context, f *flag.FlagSet) error {
	return nil
}

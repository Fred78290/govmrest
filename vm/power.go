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

	"github.com/Fred78290/govmrest/flags"
	"github.com/Fred78290/govmrest/object"
	"github.com/vmware/govmomi/govc/cli"
	"github.com/vmware/govmomi/vim25/types"
)

type power struct {
	*flags.OutputFlag
	*flags.ClientFlag
	*flags.SearchFlag

	On       bool
	Off      bool
	Reset    bool
	Reboot   bool
	Shutdown bool
	Suspend  bool
	Force    bool
	Multi    bool
	Wait     bool
}

func init() {
	cli.Register("vm.power", &power{})
}

func (cmd *power) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.OutputFlag, ctx = flags.NewOutputFlag(ctx)
	cmd.OutputFlag.Register(ctx, f)

	cmd.ClientFlag, ctx = flags.NewClientFlag(ctx)
	cmd.ClientFlag.Register(ctx, f)

	cmd.SearchFlag, ctx = flags.NewSearchFlag(ctx, flags.SearchVirtualMachines)
	cmd.SearchFlag.Register(ctx, f)

	f.BoolVar(&cmd.On, "on", false, "Power on")
	f.BoolVar(&cmd.Off, "off", false, "Power off")
	f.BoolVar(&cmd.Reset, "reset", false, "Power reset")
	f.BoolVar(&cmd.Suspend, "suspend", false, "Power suspend")
	f.BoolVar(&cmd.Reboot, "r", false, "Reboot guest")
	f.BoolVar(&cmd.Shutdown, "s", false, "Shutdown guest")
	f.BoolVar(&cmd.Force, "force", false, "Force (ignore state error and hard shutdown/reboot if tools unavailable)")
	f.BoolVar(&cmd.Wait, "wait", true, "Wait for the operation to complete")
}

func (cmd *power) Usage() string {
	return "NAME..."
}

func (cmd *power) Description() string {
	return `Invoke VM power operations.

Examples:
  govc vm.power -on VM1 VM2 VM3
  govc vm.power -on -M VM1 VM2 VM3
  govc vm.power -off -force VM1`
}

func (cmd *power) Process(ctx context.Context) error {
	if err := cmd.ClientFlag.Process(ctx); err != nil {
		return err
	}
	if err := cmd.SearchFlag.Process(ctx); err != nil {
		return err
	}
	opts := []bool{cmd.On, cmd.Off, cmd.Reset, cmd.Suspend, cmd.Reboot, cmd.Shutdown}
	selected := false

	for _, opt := range opts {
		if opt {
			if selected {
				return flag.ErrHelp
			}
			selected = opt
		}
	}

	if !selected {
		return flag.ErrHelp
	}

	return nil
}

// this is annoying, but the likely use cases for Datacenter.PowerOnVM outside of this command would
// use []types.ManagedObjectReference via ContainerView or field such as ResourcePool.Vm rather than the Finder.
func vmReferences(vms []*object.VirtualMachine) []types.ManagedObjectReference {
	refs := make([]types.ManagedObjectReference, len(vms))
	for i, vm := range vms {
		refs[i] = vm.Reference()
	}
	return refs
}

func (cmd *power) Run(ctx context.Context, f *flag.FlagSet) error {
	vms, err := cmd.VirtualMachines(f.Args())
	if err != nil {
		return err
	}

	for _, vm := range vms {
		var task *object.Task

		switch {
		case cmd.On:
			fmt.Fprintf(cmd, "Powering on %s... ", vm.Reference())
			task, err = vm.PowerOn(ctx)
		case cmd.Off:
			fmt.Fprintf(cmd, "Powering off %s... ", vm.Reference())
			task, err = vm.PowerOff(ctx)
		case cmd.Reset:
			fmt.Fprintf(cmd, "Reset %s... ", vm.Reference())
			task, err = vm.Reset(ctx)
		case cmd.Suspend:
			fmt.Fprintf(cmd, "Suspend %s... ", vm.Reference())
			task, err = vm.Suspend(ctx)
		case cmd.Reboot:
			fmt.Fprintf(cmd, "Reboot guest %s... ", vm.Reference())
			err = vm.RebootGuest(ctx)

			if err != nil && cmd.Force && isToolsUnavailable(err) {
				task, err = vm.Reset(ctx)
			}
		case cmd.Shutdown:
			fmt.Fprintf(cmd, "Shutdown guest %s... ", vm.Reference())
			err = vm.ShutdownGuest(ctx)

			if err != nil && cmd.Force && isToolsUnavailable(err) {
				task, err = vm.PowerOff(ctx)
			}
		}

		if err != nil {
			return err
		}

		if cmd.Wait && task != nil {
			err = task.Wait(ctx)
		}
		if err == nil {
			fmt.Fprintf(cmd, "OK\n")
			continue
		}

		if cmd.Force {
			fmt.Fprintf(cmd, "Error: %s\n", err)
			continue
		}

		return err
	}

	return nil
}

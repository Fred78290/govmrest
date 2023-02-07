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

	"github.com/vmware/govmomi/govc/cli"
	"github.com/vmware/govmomi/govc/flags"
)

type info struct {
	*flags.ClientFlag
	*flags.OutputFlag
	*flags.SearchFlag

	General         bool
	ExtraConfig     bool
	Resources       bool
	ToolsConfigInfo bool
}

func init() {
	cli.Register("vm.info", &info{})
}

func (cmd *info) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.ClientFlag, ctx = flags.NewClientFlag(ctx)
	cmd.ClientFlag.Register(ctx, f)

	cmd.OutputFlag, ctx = flags.NewOutputFlag(ctx)
	cmd.OutputFlag.Register(ctx, f)

	cmd.SearchFlag, ctx = flags.NewSearchFlag(ctx, flags.SearchVirtualMachines)
	cmd.SearchFlag.Register(ctx, f)

	f.BoolVar(&cmd.General, "g", true, "Show general summary")
	f.BoolVar(&cmd.ExtraConfig, "e", false, "Show ExtraConfig")
	f.BoolVar(&cmd.Resources, "r", false, "Show resource summary")
	f.BoolVar(&cmd.ToolsConfigInfo, "t", false, "Show ToolsConfigInfo")
}

func (cmd *info) Process(ctx context.Context) error {
	if err := cmd.ClientFlag.Process(ctx); err != nil {
		return err
	}
	if err := cmd.OutputFlag.Process(ctx); err != nil {
		return err
	}
	if err := cmd.SearchFlag.Process(ctx); err != nil {
		return err
	}
	return nil
}

func (cmd *info) Usage() string {
	return `VM...`
}

func (cmd *info) Description() string {
	return `Display info for VM.

The '-r' flag displays additional info for CPU, memory and storage usage,
along with the VM's Datastores, Networks and PortGroups.

Examples:
  govc vm.info $vm
  govc vm.info -r $vm | grep Network:
  govc vm.info -json $vm
  govc find . -type m -runtime.powerState poweredOn | xargs govc vm.info`
}

func (cmd *info) Run(ctx context.Context, f *flag.FlagSet) error {
	return nil
}

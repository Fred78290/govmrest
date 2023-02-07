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
)

type clone struct {
	name          string
	memory        int
	cpus          int
	on            bool
	force         bool
	customization string
	waitForIP     bool
}

func init() {
	cli.Register("vm.clone", &clone{})
}

func (cmd *clone) Register(ctx context.Context, f *flag.FlagSet) {
	f.IntVar(&cmd.memory, "m", 0, "Size in MB of memory")
	f.IntVar(&cmd.cpus, "c", 0, "Number of CPUs")
	f.BoolVar(&cmd.on, "on", true, "Power on VM")
	f.BoolVar(&cmd.force, "force", false, "Create VM if vmx already exists")
	f.StringVar(&cmd.customization, "customization", "", "Customization Specification Name")
	f.BoolVar(&cmd.waitForIP, "waitip", false, "Wait for VM to acquire IP address")
}

func (cmd *clone) Usage() string {
	return "NAME"
}

func (cmd *clone) Description() string {
	return `Clone VM to NAME.

Examples:
  govc vm.clone -vm template-vm new-vm`
}

func (cmd *clone) Process(ctx context.Context) error {
	return nil
}

func (cmd *clone) Run(ctx context.Context, f *flag.FlagSet) error {
	if len(f.Args()) != 1 {
		return flag.ErrHelp
	}

	cmd.name = f.Arg(0)
	if cmd.name == "" {
		return flag.ErrHelp
	}

	return nil
}

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

type register struct {
	name string
}

func init() {
	cli.Register("vm.register", &register{})
}

func (cmd *register) Register(ctx context.Context, f *flag.FlagSet) {

	f.StringVar(&cmd.name, "name", "", "Name of the VM")
}

func (cmd *register) Process(ctx context.Context) error {
	return nil
}

func (cmd *register) Usage() string {
	return "VMX"
}

func (cmd *register) Description() string {
	return `Add an existing VM to the inventory.

VMX is an absolute path to the vm config file.

Examples:
  govc vm.register /path/name.vmx`
}

func (cmd *register) Run(ctx context.Context, f *flag.FlagSet) error {
	return nil
}

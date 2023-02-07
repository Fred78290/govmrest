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

package find

import (
	"context"

	"github.com/Fred78290/govmrest/object"
	"github.com/Fred78290/govmrest/vim25"
)

type Finder struct {
	client *vim25.Client
}

func NewFinder(client *vim25.Client, all ...bool) *Finder {

	f := &Finder{
		client: client,
	}

	return f
}

func (f *Finder) VirtualMachineList(ctx context.Context, path string) ([]*object.VirtualMachine, error) {
	return nil, nil
}

func (f *Finder) VirtualMachine(ctx context.Context, path string) (*object.VirtualMachine, error) {
	vms, err := f.VirtualMachineList(ctx, path)
	if err != nil {
		return nil, err
	}

	if len(vms) > 1 {
		return nil, &MultipleFoundError{"vm", path}
	}

	return vms[0], nil
}

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
	"github.com/Fred78290/govmrest/vim25"
	"github.com/vmware/govmomi/vim25/types"
)

type Reference interface {
	Reference() types.ManagedObjectReference
}

func NewReference(c *vim25.Client, e types.ManagedObjectReference) Reference {
	switch e.Type {
	case "Folder":
		return NewFolder(c, e)
	case "VirtualMachine":
		return NewVirtualMachine(c, e)
	case "Network":
		return NewNetwork(c, e)
	default:
		panic("Unknown managed entity: " + e.Type)
	}
}

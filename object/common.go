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
	"errors"
	"fmt"
	"path"

	"github.com/Fred78290/govmrest/vim25"
	"github.com/vmware/govmomi/vim25/types"
)

var (
	ErrNotSupported = errors.New("product/version specific feature not supported by target")
)

// Common contains the fields and functions common to all objects.
type Common struct {
	InventoryPath string

	c *vim25.Client
	r types.ManagedObjectReference
}

func (c Common) String() string {
	ref := fmt.Sprintf("%v", c.Reference())

	if c.InventoryPath == "" {
		return ref
	}

	return fmt.Sprintf("%s @ %s", ref, c.InventoryPath)
}

func NewCommon(c *vim25.Client, r types.ManagedObjectReference) Common {
	return Common{c: c, r: r}
}

func (c Common) Reference() types.ManagedObjectReference {
	return c.r
}

func (c Common) Client() *vim25.Client {
	return c.c
}

// Name returns the base name of the InventoryPath field
func (c Common) Name() string {
	if c.InventoryPath == "" {
		return ""
	}
	return path.Base(c.InventoryPath)
}

func (c *Common) SetInventoryPath(p string) {
	c.InventoryPath = p
}

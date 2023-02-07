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
package flags

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/Fred78290/govmrest/object"
	"github.com/Fred78290/govmrest/vim25"

	"github.com/Fred78290/govmrest/find"
)

const (
	SearchVirtualMachines = iota + 1
)

type SearchFlag struct {
	common

	*ClientFlag

	t               int
	entity          string
	finder          *find.Finder
	byDatastorePath string
	byDNSName       string
	byIP            string
	byUUID          string

	isset bool
}

func NewSearchFlag(ctx context.Context, t int) (*SearchFlag, context.Context) {
	searchFlagKey := flagKey(fmt.Sprintf("search%d", t))

	if v := ctx.Value(searchFlagKey); v != nil {
		return v.(*SearchFlag), ctx
	}

	v := &SearchFlag{
		t: t,
	}

	v.ClientFlag, ctx = NewClientFlag(ctx)

	ctx = context.WithValue(ctx, searchFlagKey, v)
	return v, ctx
}

func (flag *SearchFlag) Register(ctx context.Context, fs *flag.FlagSet) {
	flag.RegisterOnce(func() {
		flag.ClientFlag.Register(ctx, fs)

		register := func(v *string, f string, d string) {
			f = fmt.Sprintf("%s.%s", strings.ToLower(flag.entity), f)
			d = fmt.Sprintf(d, flag.entity)
			fs.StringVar(v, f, "", d)
		}

		register(&flag.byDatastorePath, "path", "Find %s by path to .vmx file")
		register(&flag.byDNSName, "dns", "Find %s by FQDN")
		register(&flag.byIP, "ip", "Find %s by IP address")
		register(&flag.byUUID, "uuid", "Find %s by UUID")
	})
}

func (flag *SearchFlag) Process(ctx context.Context) error {
	return flag.ProcessOnce(func() error {
		if err := flag.ClientFlag.Process(ctx); err != nil {
			return err
		}

		flags := []string{
			flag.byDatastorePath,
			flag.byDNSName,
			flag.byIP,
			flag.byUUID,
		}

		flag.isset = false
		for _, f := range flags {
			if f != "" {
				if flag.isset {
					return errors.New("cannot use more than one search flag")
				}
				flag.isset = true
			}
		}

		return nil
	})
}

func (flag *SearchFlag) IsSet() bool {
	return flag.isset
}

func (flag *SearchFlag) searchIndex(c *vim25.Client) *object.SearchIndex {
	return object.NewSearchIndex(c)
}

func (flag *SearchFlag) searchByDatastorePath(c *vim25.Client) (object.Reference, error) {
	ctx := context.TODO()

	return flag.searchIndex(c).FindByDatastorePath(ctx, flag.byDatastorePath)
}

func (flag *SearchFlag) searchByDNSName(c *vim25.Client) (object.Reference, error) {
	ctx := context.TODO()

	return flag.searchIndex(c).FindByDnsName(ctx, flag.byDNSName)
}

func (flag *SearchFlag) searchByIP(c *vim25.Client) (object.Reference, error) {
	ctx := context.TODO()

	return flag.searchIndex(c).FindByIp(ctx, flag.byIP)
}

func (flag *SearchFlag) searchByUUID(c *vim25.Client) (object.Reference, error) {
	ctx := context.TODO()

	return flag.searchIndex(c).FindByUuid(ctx, flag.byUUID)
}

func (flag *SearchFlag) search() (object.Reference, error) {
	var ref object.Reference
	var err error

	c, err := flag.Client()
	if err != nil {
		return nil, err
	}

	switch {
	case flag.byDatastorePath != "":
		ref, err = flag.searchByDatastorePath(c)
	case flag.byDNSName != "":
		ref, err = flag.searchByDNSName(c)
	case flag.byIP != "":
		ref, err = flag.searchByIP(c)
	case flag.byUUID != "":
		ref, err = flag.searchByUUID(c)
	default:
		err = errors.New("no search flag specified")
	}

	if err != nil {
		return nil, err
	}

	if ref == nil {
		return nil, fmt.Errorf("no such %s", flag.entity)
	}

	return ref, nil
}

func (flag *SearchFlag) Finder(all ...bool) (*find.Finder, error) {
	if flag.finder != nil {
		return flag.finder, nil
	}

	c, err := flag.Client()
	if err != nil {
		return nil, err
	}

	allFlag := false
	if len(all) == 1 {
		allFlag = all[0]
	}
	finder := find.NewFinder(c, allFlag)

	flag.finder = finder

	return flag.finder, nil
}

func (flag *SearchFlag) VirtualMachine() (*object.VirtualMachine, error) {
	ref, err := flag.search()
	if err != nil {
		return nil, err
	}

	vm, ok := ref.(*object.VirtualMachine)
	if !ok {
		return nil, fmt.Errorf("expected VirtualMachine entity, got %s", ref.Reference().Type)
	}

	return vm, nil
}

func (flag *SearchFlag) VirtualMachines(args []string) ([]*object.VirtualMachine, error) {
	ctx := context.TODO()
	var out []*object.VirtualMachine

	if flag.IsSet() {
		vm, err := flag.VirtualMachine()
		if err != nil {
			return nil, err
		}

		out = append(out, vm)
		return out, nil
	}

	// List virtual machines
	if len(args) == 0 {
		return nil, errors.New("no argument")
	}

	finder, err := flag.Finder()
	if err != nil {
		return nil, err
	}

	var nfe error

	// List virtual machines for every argument
	for _, arg := range args {
		vms, err := finder.VirtualMachineList(ctx, arg)
		if err != nil {
			if _, ok := err.(*find.NotFoundError); ok {
				// Let caller decide how to handle NotFoundError
				nfe = err
				continue
			}
			return nil, err
		}

		out = append(out, vms...)
	}

	return out, nfe
}

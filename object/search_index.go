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

	"github.com/Fred78290/govmrest/vim25"
	"github.com/Fred78290/govmrest/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
)

type SearchIndex struct {
	Common
}

func NewSearchIndex(c *vim25.Client) *SearchIndex {
	s := SearchIndex{
		Common: NewCommon(c, *c.ServiceContent.SearchIndex),
	}

	return &s
}

// FindByDatastorePath finds a virtual machine by its location on a datastore.
func (s SearchIndex) FindByDatastorePath(ctx context.Context, path string) (Reference, error) {
	req := types.FindByDatastorePath{
		This: s.Reference(),
		Path: path,
	}

	res, err := methods.FindByDatastorePath(ctx, s.c, &req)
	if err != nil {
		return nil, err
	}

	if res.Returnval == nil {
		return nil, nil
	}
	return NewReference(s.c, *res.Returnval), nil
}

// FindByDnsName finds a virtual machine by DNS name.
func (s SearchIndex) FindByDnsName(ctx context.Context, dnsName string) (Reference, error) {
	req := types.FindByDnsName{
		This:     s.Reference(),
		DnsName:  dnsName,
		VmSearch: true,
	}

	res, err := methods.FindByDnsName(ctx, s.c, &req)
	if err != nil {
		return nil, err
	}

	if res.Returnval == nil {
		return nil, nil
	}
	return NewReference(s.c, *res.Returnval), nil
}

// FindByIp finds a virtual machine by IP address.
func (s SearchIndex) FindByIp(ctx context.Context, ip string) (Reference, error) {
	req := types.FindByIp{
		This:     s.Reference(),
		Ip:       ip,
		VmSearch: true,
	}

	res, err := methods.FindByIp(ctx, s.c, &req)
	if err != nil {
		return nil, err
	}

	if res.Returnval == nil {
		return nil, nil
	}
	return NewReference(s.c, *res.Returnval), nil
}

// FindByUuid finds a virtual machine or host by UUID.
func (s SearchIndex) FindByUuid(ctx context.Context, uuid string) (Reference, error) {
	var instanceUuid bool

	req := types.FindByUuid{
		This:         s.Reference(),
		Uuid:         uuid,
		VmSearch:     true,
		InstanceUuid: &instanceUuid,
	}

	res, err := methods.FindByUuid(ctx, s.c, &req)
	if err != nil {
		return nil, err
	}

	if res.Returnval == nil {
		return nil, nil
	}
	return NewReference(s.c, *res.Returnval), nil
}

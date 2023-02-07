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

package methods

import (
	"context"

	"github.com/Fred78290/govmrest/vim25"
	"github.com/vmware/govmomi/vim25/types"
)

func FindByDatastorePath(ctx context.Context, c *vim25.Client, req *types.FindByDatastorePath) (*types.FindByDatastorePathResponse, error) {
	return nil, nil
}

func FindByDnsName(ctx context.Context, c *vim25.Client, req *types.FindByDnsName) (*types.FindByDnsNameResponse, error) {
	return nil, nil
}

func FindByIp(ctx context.Context, c *vim25.Client, req *types.FindByIp) (*types.FindByIpResponse, error) {
	return nil, nil
}

func FindByUuid(ctx context.Context, c *vim25.Client, req *types.FindByUuid) (*types.FindByUuidResponse, error) {
	return nil, nil
}

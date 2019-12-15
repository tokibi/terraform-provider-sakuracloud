// Copyright 2016-2019 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func resourceSakuraCloudPacketFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceSakuraCloudPacketFilterCreate,
		Read:   resourceSakuraCloudPacketFilterRead,
		Update: resourceSakuraCloudPacketFilterUpdate,
		Delete: resourceSakuraCloudPacketFilterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expression": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice(types.PacketFilterProtocolsStrings(), false),
						},
						"source_network": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"source_port": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"destination_port": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
						"allow": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "",
						},
					},
				},
			},
			"zone": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				Description:  "target SakuraCloud zone",
				ValidateFunc: validateZone([]string{"is1a", "is1b", "tk1a", "tk1v"}),
			},
		},
	}
}

func resourceSakuraCloudPacketFilterCreate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	pfOp := sacloud.NewPacketFilterOp(client)

	pf, err := pfOp.Create(ctx, zone, expandPacketFilterCreateRequest(d))
	if err != nil {
		return fmt.Errorf("creating SakuraCloud PacketFilter is failed: %s", err)
	}

	d.SetId(pf.ID.String())
	return resourceSakuraCloudPacketFilterRead(d, meta)
}

func resourceSakuraCloudPacketFilterRead(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	pfOp := sacloud.NewPacketFilterOp(client)

	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", d.Id(), err)
	}

	return setPacketFilterResourceData(ctx, d, client, pf)
}

func resourceSakuraCloudPacketFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	pfOp := sacloud.NewPacketFilterOp(client)

	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		return fmt.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", d.Id(), err)
	}

	_, err = pfOp.Update(ctx, zone, pf.ID, expandPacketFilterUpdateRequest(d, pf))
	if err != nil {
		return fmt.Errorf("updating SakuraCloud PacketFilter[%s] is failed: %s", d.Id(), err)
	}

	return resourceSakuraCloudPacketFilterRead(d, meta)
}

func resourceSakuraCloudPacketFilterDelete(d *schema.ResourceData, meta interface{}) error {
	client, ctx, zone := getSacloudClient(d, meta)
	pfOp := sacloud.NewPacketFilterOp(client)

	pf, err := pfOp.Read(ctx, zone, sakuraCloudID(d.Id()))
	if err != nil {
		if sacloud.IsNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("could not read SakuraCloud PacketFilter[%s]: %s", d.Id(), err)
	}

	if err := waitForDeletionByPacketFilterID(ctx, client, zone, pf.ID); err != nil {
		return fmt.Errorf("waiting deletion is failed: PacketFilter[%s] still used by Server: %s", pf.ID, err)
	}

	if err := pfOp.Delete(ctx, zone, pf.ID); err != nil {
		return fmt.Errorf("deleting SakuraCloud PacketFilter[%s] is failed: %s", d.Id(), err)
	}
	return nil
}

func setPacketFilterResourceData(ctx context.Context, d *schema.ResourceData, client *APIClient, data *sacloud.PacketFilter) error {
	d.Set("name", data.Name)
	d.Set("description", data.Description)
	d.Set("zone", getZone(d, client))
	return d.Set("expression", flattenPacketFilterExpressions(data))
}

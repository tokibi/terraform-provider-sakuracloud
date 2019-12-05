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
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const filterAttrName = "filters"

type filterSchemaOption struct {
	excludeTags bool
}

func filterSchema(opt *filterSchemaOption) *schema.Schema {
	if opt == nil {
		opt = &filterSchemaOption{}
	}
	s := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"names": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"tags": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"conditions": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:     schema.TypeString,
						Required: true,
					},

					"values": {
						Type:     schema.TypeList,
						Required: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
		},
	}
	if opt.excludeTags {
		delete(s, "tags")
	}
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: s,
		},
	}
}

var filterNoResultMessage = "Your query returned no results. Please change your filters or selectors and try again"

func filterNoResultErr() error {
	if os.Getenv(resource.TestEnvVar) != "" {
		return nil
	}
	return fmt.Errorf(filterNoResultMessage)
}

type filterFunc func(target interface{}, cond []string) bool

type nameFilterable interface {
	GetName() string
}

func hasNames(target interface{}, cond []string) bool {
	t, ok := target.(nameFilterable)
	if !ok {
		return false
	}
	name := t.GetName()
	for _, c := range cond {
		if !strings.Contains(name, c) {
			return false
		}
	}
	return true
}

type tagFilterable interface {
	HasTag(string) bool
}

func hasTags(target interface{}, cond []string) bool {
	t, ok := target.(tagFilterable)
	if !ok {
		return false
	}
	for _, c := range cond {
		if !t.HasTag(c) {
			return false
		}
	}
	return true

}

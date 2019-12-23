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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSakuraCloudDataSourceContainerRegistry_basic(t *testing.T) {
	resourceName := "data.sakuracloud_container_registry.foobar"
	rand := randomName()
	prefix := acctest.RandStringFromCharSet(60, acctest.CharSetAlpha)
	password := randomPassword()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: buildConfigWithArgs(testAccSakuraCloudDataSourceContainerRegistry_basic, rand, prefix, password),
				Check: resource.ComposeTestCheckFunc(
					testCheckSakuraCloudDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rand),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
				),
			},
		},
	})
}

var testAccSakuraCloudDataSourceContainerRegistry_basic = `
resource "sakuracloud_container_registry" "foobar" {
  name        = "{{ .arg0 }}"
  prefix      = "{{ .arg1 }}"
  visibility  = "readwrite"

  description = "description"
  tags        = ["tag1", "tag2"]

  user {
    name     = "user1"
    password = "{{ .arg2 }}"
  }
  user {
    name     = "user2"
    password = "{{ .arg2 }}"
  }
}

data "sakuracloud_container_registry" "foobar" {
  filter {
    names = [sakuracloud_container_registry.foobar.name]
  }
}`

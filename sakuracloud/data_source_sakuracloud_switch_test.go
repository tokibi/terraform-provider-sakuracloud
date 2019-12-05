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
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSakuraCloudDataSourceSwitch_Basic(t *testing.T) {
	randString1 := acctest.RandStringFromCharSet(5, acctest.CharSetAlpha)
	randString2 := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	name := fmt.Sprintf("%s_%s", randString1, randString2)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		CheckDestroy:              testAccCheckSakuraCloudSwitchDestroy,

		Steps: []resource.TestStep{
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchBase(name),
				Check:  testAccCheckSakuraCloudDataSourceExists("sakuracloud_switch.foobar"),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_switch.foobar"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "name", name),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "description", "description_test"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "tags.#", "3"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "tags.0", "tag1"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "tags.1", "tag2"),
					resource.TestCheckResourceAttr("data.sakuracloud_switch.foobar", "tags.2", "tag3"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchConfig_With_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceExists("data.sakuracloud_switch.foobar"),
				),
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchConfig_NotExists(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_switch.foobar"),
				),
				Destroy: true,
			},
			{
				Config: testAccCheckSakuraCloudDataSourceSwitchConfig_With_NotExists_Tag(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSakuraCloudDataSourceNotExists("data.sakuracloud_switch.foobar"),
				),
				Destroy: true,
			},
		},
	})
}

func testAccCheckSakuraCloudDataSourceSwitchBase(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
`, name)
}

func testAccCheckSakuraCloudDataSourceSwitchConfig(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
  filters {
	names = ["%s"]
  }
}`, name, name)
}

func testAccCheckSakuraCloudDataSourceSwitchConfig_With_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
  filters {
	tags = ["tag1","tag3"]
  }
}`, name)
}

func testAccCheckSakuraCloudDataSourceSwitchConfig_With_NotExists_Tag(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
  filters {
	tags = ["tag1-xxxxxxx","tag3-xxxxxxxx"]
  }
}`, name)
}

func testAccCheckSakuraCloudDataSourceSwitchConfig_NotExists(name string) string {
	return fmt.Sprintf(`
resource "sakuracloud_switch" "foobar" {
  name = "%s"
  description = "description_test"
  tags = ["tag1","tag2","tag3"]
}
data "sakuracloud_switch" "foobar" {
  filters {
	names = ["xxxxxxxxxxxxxxxxxx"]
  }
}`, name)
}

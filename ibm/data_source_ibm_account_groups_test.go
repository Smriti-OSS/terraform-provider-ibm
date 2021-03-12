/**
 * (C) Copyright IBM Corp. 2021.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIbmAccountGroupsDataSourceBasic(t *testing.T) {
	accountGroupParent := fmt.Sprintf("parent_%d", acctest.RandIntRange(10, 100))
	accountGroupName := fmt.Sprintf("name_%d", acctest.RandIntRange(10, 100))
	accountGroupPrimaryContactIamID := fmt.Sprintf("primary_contact_iam_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmAccountGroupsDataSourceConfigBasic(accountGroupParent, accountGroupName, accountGroupPrimaryContactIamID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_account_groups.account_groups", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_account_groups.account_groups", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_account_groups.account_groups", "resources.#"),
					resource.TestCheckResourceAttr("data.ibm_account_groups.account_groups", "resources.0.parent", accountGroupParent),
					resource.TestCheckResourceAttr("data.ibm_account_groups.account_groups", "resources.0.name", accountGroupName),
					resource.TestCheckResourceAttr("data.ibm_account_groups.account_groups", "resources.0.primary_contact_iam_id", accountGroupPrimaryContactIamID),
				),
			},
		},
	})
}

func testAccCheckIbmAccountGroupsDataSourceConfigBasic(accountGroupParent string, accountGroupName string, accountGroupPrimaryContactIamID string) string {
	return fmt.Sprintf(`
		resource "ibm_enterprise_account_group" "enterprise_account_group" {
			parent = "%s"
			name = "%s"
			primary_contact_iam_id = "%s"
		}

		data "ibm_account_groups" "account_groups" {
			name = ibm_enterprise_account_group.enterprise_account_group.name
		}
	`, accountGroupParent, accountGroupName, accountGroupPrimaryContactIamID)
}

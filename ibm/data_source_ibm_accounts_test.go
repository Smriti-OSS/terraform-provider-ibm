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

func TestAccIbmAccountsDataSourceBasic(t *testing.T) {
	accountParent := fmt.Sprintf("parent_%d", acctest.RandIntRange(10, 100))
	accountName := fmt.Sprintf("name_%d", acctest.RandIntRange(10, 100))
	accountOwnerIamID := fmt.Sprintf("owner_iam_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmAccountsDataSourceConfigBasic(accountParent, accountName, accountOwnerIamID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_accounts.accounts", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_accounts.accounts", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_accounts.accounts", "rows_count"),
					resource.TestCheckResourceAttrSet("data.ibm_accounts.accounts", "next_url"),
					resource.TestCheckResourceAttrSet("data.ibm_accounts.accounts", "resources.#"),
					resource.TestCheckResourceAttr("data.ibm_accounts.accounts", "resources.0.parent", accountParent),
					resource.TestCheckResourceAttr("data.ibm_accounts.accounts", "resources.0.name", accountName),
					resource.TestCheckResourceAttr("data.ibm_accounts.accounts", "resources.0.owner_iam_id", accountOwnerIamID),
				),
			},
		},
	})
}

func testAccCheckIbmAccountsDataSourceConfigBasic(accountParent string, accountName string, accountOwnerIamID string) string {
	return fmt.Sprintf(`
		resource "ibm_enterprise_account" "enterprise_account" {
			parent = "%s"
			name = "%s"
			owner_iam_id = "%s"
		}

		data "ibm_accounts" "accounts" {
			name = ibm_enterprise_account.enterprise_account.name
		}
	`, accountParent, accountName, accountOwnerIamID)
}

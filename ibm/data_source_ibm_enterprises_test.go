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

func TestAccIbmEnterprisesDataSourceBasic(t *testing.T) {
	//enterpriseSourceAccountID := fmt.Sprintf("source_account_id_%d", acctest.RandIntRange(10, 100))
	enterpriseName := fmt.Sprintf("enterprise_name_%d", acctest.RandIntRange(10, 100))
	//enterprisePrimaryContactIamID := fmt.Sprintf("primary_contact_iam_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnterprise(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmEnterprisesDataSourceConfigBasic(enterpriseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.#"),
					resource.TestCheckResourceAttr("data.ibm_enterprises.enterprises", "resources.0.name", enterpriseName),
				),
			},
		},
	})
}

func TestAccIbmEnterprisesDataSourceAllArgs(t *testing.T) {
	//enterpriseSourceAccountID := fmt.Sprintf("source_account_id_%d", acctest.RandIntRange(10, 100))
	enterpriseName := fmt.Sprintf("enterprise_name_%d", acctest.RandIntRange(10, 100))
	//enterprisePrimaryContactIamID := fmt.Sprintf("primary_contact_iam_id_%d", acctest.RandIntRange(10, 100))
	enterpriseDomain := fmt.Sprintf("enterprise_domain_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnterprise(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmEnterprisesDataSourceConfig(enterpriseName, enterpriseDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.#"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.url"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.enterprise_account_id"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.crn"),
					resource.TestCheckResourceAttr("data.ibm_enterprises.enterprises", "resources.0.name", enterpriseName),
					resource.TestCheckResourceAttr("data.ibm_enterprises.enterprises", "resources.0.domain", enterpriseDomain),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.state"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.primary_contact_iam_id"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.primary_contact_email"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.created_by"),
				),
			},
		},
	})
}

func testAccCheckIbmEnterprisesDataSourceConfigBasic(enterpriseName string) string {
	return fmt.Sprintf(`
		data "ibm_iam_users" "current_account_users"{
		}
		resource "ibm_enterprise" "enterprise" {
			source_account_id = data.ibm_iam_users.current_account_users.users[0].account_id
			name = "%s"
			primary_contact_iam_id = data.ibm_iam_users.current_account_users.users[0].iam_id
		}

		data "ibm_enterprises" "enterprises" {
			name = ibm_enterprise.enterprise.name
		}
	`, enterpriseName)
}

func testAccCheckIbmEnterprisesDataSourceConfig(enterpriseName string, enterpriseDomain string) string {
	return fmt.Sprintf(`
		data "ibm_iam_users" "current_account_users"{
		}
		resource "ibm_enterprise" "enterprise" {
			source_account_id = data.ibm_iam_users.current_account_users.users[0].account_id
			name = "%s"
			primary_contact_iam_id = data.ibm_iam_users.current_account_users.users[0].iam_id
			domain = "%s"
		}
		data "ibm_enterprises" "enterprises" {
			name = ibm_enterprise.enterprise.name
		}
	`, enterpriseName, enterpriseDomain)
}

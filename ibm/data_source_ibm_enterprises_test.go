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
	enterpriseSourceAccountID := fmt.Sprintf("source_account_id_%d", acctest.RandIntRange(10, 100))
	enterpriseName := fmt.Sprintf("name_%d", acctest.RandIntRange(10, 100))
	enterprisePrimaryContactIamID := fmt.Sprintf("primary_contact_iam_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmEnterprisesDataSourceConfigBasic(enterpriseSourceAccountID, enterpriseName, enterprisePrimaryContactIamID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "rows_count"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "next_url"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.#"),
					resource.TestCheckResourceAttr("data.ibm_enterprises.enterprises", "resources.0.name", enterpriseName),
					resource.TestCheckResourceAttr("data.ibm_enterprises.enterprises", "resources.0.primary_contact_iam_id", enterprisePrimaryContactIamID),
				),
			},
		},
	})
}

func TestAccIbmEnterprisesDataSourceAllArgs(t *testing.T) {
	enterpriseSourceAccountID := fmt.Sprintf("source_account_id_%d", acctest.RandIntRange(10, 100))
	enterpriseName := fmt.Sprintf("name_%d", acctest.RandIntRange(10, 100))
	enterprisePrimaryContactIamID := fmt.Sprintf("primary_contact_iam_id_%d", acctest.RandIntRange(10, 100))
	enterpriseDomain := fmt.Sprintf("domain_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmEnterprisesDataSourceConfig(enterpriseSourceAccountID, enterpriseName, enterprisePrimaryContactIamID, enterpriseDomain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "rows_count"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "next_url"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.#"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.url"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.enterprise_account_id"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.crn"),
					resource.TestCheckResourceAttr("data.ibm_enterprises.enterprises", "resources.0.name", enterpriseName),
					resource.TestCheckResourceAttr("data.ibm_enterprises.enterprises", "resources.0.domain", enterpriseDomain),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.state"),
					resource.TestCheckResourceAttr("data.ibm_enterprises.enterprises", "resources.0.primary_contact_iam_id", enterprisePrimaryContactIamID),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.primary_contact_email"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.created_by"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.updated_at"),
					resource.TestCheckResourceAttrSet("data.ibm_enterprises.enterprises", "resources.0.updated_by"),
				),
			},
		},
	})
}

func testAccCheckIbmEnterprisesDataSourceConfigBasic(enterpriseSourceAccountID string, enterpriseName string, enterprisePrimaryContactIamID string) string {
	return fmt.Sprintf(`
		resource "ibm_enterprise" "enterprise" {
			source_account_id = "%s"
			name = "%s"
			primary_contact_iam_id = "%s"
		}

		data "ibm_enterprises" "enterprises" {
			name = ibm_enterprise.enterprise.name
		}
	`, enterpriseSourceAccountID, enterpriseName, enterprisePrimaryContactIamID)
}

func testAccCheckIbmEnterprisesDataSourceConfig(enterpriseSourceAccountID string, enterpriseName string, enterprisePrimaryContactIamID string, enterpriseDomain string) string {
	return fmt.Sprintf(`
		resource "ibm_enterprise" "enterprise" {
			source_account_id = "%s"
			name = "%s"
			primary_contact_iam_id = "%s"
			domain = "%s"
		}

		data "ibm_enterprises" "enterprises" {
			name = ibm_enterprise.enterprise.name
		}
	`, enterpriseSourceAccountID, enterpriseName, enterprisePrimaryContactIamID, enterpriseDomain)
}

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/IBM/platform-services-go-sdk/enterprisemanagementv1"
)

func TestAccIbmEnterpriseAccountBasic(t *testing.T) {
	var conf enterprisemanagementv1.Account
	//parent := fmt.Sprintf("parent_%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf-gen-account-name_%d", acctest.RandIntRange(10, 100))
	//ownerIamID := fmt.Sprintf("owner_iam_id_%d", acctest.RandIntRange(10, 100))
	//parentUpdate := fmt.Sprintf("parent_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckEnterprise(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmEnterpriseAccountConfigBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmEnterpriseAccountExists("ibm_enterprise_account.enterprise_account", conf),
					resource.TestCheckResourceAttrSet("ibm_enterprise_account.enterprise_account", "parent"),
					resource.TestCheckResourceAttr("ibm_enterprise_account.enterprise_account", "name", name),
					resource.TestCheckResourceAttrSet("ibm_enterprise_account.enterprise_account", "owner_iam_id"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIbmEnterpriseAccountConfigUpdateBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("ibm_enterprise_account.enterprise_account", "parent"),
					resource.TestCheckResourceAttrSet("ibm_enterprise_account.enterprise_account", "name"),
					resource.TestCheckResourceAttrSet("ibm_enterprise_account.enterprise_account", "owner_iam_id"),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_enterprise_account.enterprise_account",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIbmEnterpriseAccountConfigBasic(name string) string {
	return fmt.Sprintf(`
		data "ibm_enterprises" "enterprises_instance" {
		}
		resource "ibm_enterprise_account" "enterprise_account" {
			parent = data.ibm_enterprises.enterprises_instance.resources[0].crn
			name = "%s"
			owner_iam_id = data.ibm_enterprises.enterprises_instance.resources[0].primary_contact_iam_id
		}
	`, name)
}

func testAccCheckIbmEnterpriseAccountConfigUpdateBasic(name string) string {
	return fmt.Sprintf(`
		data "ibm_enterprises" "enterprises_instance" {
		}
		data "ibm_account_groups" "account_groups_instance" {
		}
		resource "ibm_enterprise_account" "enterprise_account" {
			parent = data.ibm_account_groups.account_groups_instance.resources[0].crn
			name = "%s"
			owner_iam_id = data.ibm_enterprises.enterprises_instance.resources[0].primary_contact_iam_id
		}
	`, name)
}
func testAccCheckIbmEnterpriseAccountExists(n string, obj enterprisemanagementv1.Account) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		enterpriseManagementClient, err := testAccProvider.Meta().(ClientSession).EnterpriseManagementV1()
		if err != nil {
			return err
		}

		getAccountOptions := &enterprisemanagementv1.GetAccountOptions{}

		getAccountOptions.SetAccountID(rs.Primary.ID)

		account, _, err := enterpriseManagementClient.GetAccount(getAccountOptions)
		if err != nil {
			return err
		}

		obj = *account
		return nil
	}
}

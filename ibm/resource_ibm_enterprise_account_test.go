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
	parent := fmt.Sprintf("parent_%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("name_%d", acctest.RandIntRange(10, 100))
	ownerIamID := fmt.Sprintf("owner_iam_id_%d", acctest.RandIntRange(10, 100))
	parentUpdate := fmt.Sprintf("parent_%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("name_%d", acctest.RandIntRange(10, 100))
	ownerIamIDUpdate := fmt.Sprintf("owner_iam_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIbmEnterpriseAccountDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmEnterpriseAccountConfigBasic(parent, name, ownerIamID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmEnterpriseAccountExists("ibm_enterprise_account.enterprise_account", conf),
					resource.TestCheckResourceAttr("ibm_enterprise_account.enterprise_account", "parent", parent),
					resource.TestCheckResourceAttr("ibm_enterprise_account.enterprise_account", "name", name),
					resource.TestCheckResourceAttr("ibm_enterprise_account.enterprise_account", "owner_iam_id", ownerIamID),
				),
			},
			resource.TestStep{
				Config: testAccCheckIbmEnterpriseAccountConfigBasic(parentUpdate, nameUpdate, ownerIamIDUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_enterprise_account.enterprise_account", "parent", parentUpdate),
					resource.TestCheckResourceAttr("ibm_enterprise_account.enterprise_account", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_enterprise_account.enterprise_account", "owner_iam_id", ownerIamIDUpdate),
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

func testAccCheckIbmEnterpriseAccountConfigBasic(parent string, name string, ownerIamID string) string {
	return fmt.Sprintf(`

		resource "ibm_enterprise_account" "enterprise_account" {
			parent = "%s"
			name = "%s"
			owner_iam_id = "%s"
		}
	`, parent, name, ownerIamID)
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

func testAccCheckIbmEnterpriseAccountDestroy(s *terraform.State) error {
	enterpriseManagementClient, err := testAccProvider.Meta().(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_enterprise_account" {
			continue
		}

		getAccountOptions := &enterprisemanagementv1.GetAccountOptions{}

		getAccountOptions.SetAccountID(rs.Primary.ID)

		// Try to find the key
		_, response, err := enterpriseManagementClient.GetAccount(getAccountOptions)

		if err == nil {
			return fmt.Errorf("enterprise_account still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for enterprise_account (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

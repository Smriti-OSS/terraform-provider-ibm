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

func TestAccIbmEnterpriseAccountGroupBasic(t *testing.T) {
	var conf enterprisemanagementv1.AccountGroup
	parent := fmt.Sprintf("parent_%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("name_%d", acctest.RandIntRange(10, 100))
	primaryContactIamID := fmt.Sprintf("primary_contact_iam_id_%d", acctest.RandIntRange(10, 100))
	parentUpdate := fmt.Sprintf("parent_%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("name_%d", acctest.RandIntRange(10, 100))
	primaryContactIamIDUpdate := fmt.Sprintf("primary_contact_iam_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIbmEnterpriseAccountGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmEnterpriseAccountGroupConfigBasic(parent, name, primaryContactIamID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmEnterpriseAccountGroupExists("ibm_enterprise_account_group.enterprise_account_group", conf),
					resource.TestCheckResourceAttr("ibm_enterprise_account_group.enterprise_account_group", "parent", parent),
					resource.TestCheckResourceAttr("ibm_enterprise_account_group.enterprise_account_group", "name", name),
					resource.TestCheckResourceAttr("ibm_enterprise_account_group.enterprise_account_group", "primary_contact_iam_id", primaryContactIamID),
				),
			},
			resource.TestStep{
				Config: testAccCheckIbmEnterpriseAccountGroupConfigBasic(parentUpdate, nameUpdate, primaryContactIamIDUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_enterprise_account_group.enterprise_account_group", "parent", parentUpdate),
					resource.TestCheckResourceAttr("ibm_enterprise_account_group.enterprise_account_group", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_enterprise_account_group.enterprise_account_group", "primary_contact_iam_id", primaryContactIamIDUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_enterprise_account_group.enterprise_account_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIbmEnterpriseAccountGroupConfigBasic(parent string, name string, primaryContactIamID string) string {
	return fmt.Sprintf(`

		resource "ibm_enterprise_account_group" "enterprise_account_group" {
			parent = "%s"
			name = "%s"
			primary_contact_iam_id = "%s"
		}
	`, parent, name, primaryContactIamID)
}

func testAccCheckIbmEnterpriseAccountGroupExists(n string, obj enterprisemanagementv1.AccountGroup) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		enterpriseManagementClient, err := testAccProvider.Meta().(ClientSession).EnterpriseManagementV1()
		if err != nil {
			return err
		}

		getAccountGroupOptions := &enterprisemanagementv1.GetAccountGroupOptions{}

		getAccountGroupOptions.SetAccountGroupID(rs.Primary.ID)

		accountGroup, _, err := enterpriseManagementClient.GetAccountGroup(getAccountGroupOptions)
		if err != nil {
			return err
		}

		obj = *accountGroup
		return nil
	}
}

func testAccCheckIbmEnterpriseAccountGroupDestroy(s *terraform.State) error {
	enterpriseManagementClient, err := testAccProvider.Meta().(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_enterprise_account_group" {
			continue
		}

		getAccountGroupOptions := &enterprisemanagementv1.GetAccountGroupOptions{}

		getAccountGroupOptions.SetAccountGroupID(rs.Primary.ID)

		// Try to find the key
		_, response, err := enterpriseManagementClient.GetAccountGroup(getAccountGroupOptions)

		if err == nil {
			return fmt.Errorf("enterprise_account_group still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for enterprise_account_group (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

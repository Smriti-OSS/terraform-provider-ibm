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

func TestAccIbmEnterpriseBasic(t *testing.T) {
	var conf enterprisemanagementv1.Enterprise
	sourceAccountID := fmt.Sprintf("source_account_id_%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("name_%d", acctest.RandIntRange(10, 100))
	primaryContactIamID := fmt.Sprintf("primary_contact_iam_id_%d", acctest.RandIntRange(10, 100))
	sourceAccountIDUpdate := fmt.Sprintf("source_account_id_%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("name_%d", acctest.RandIntRange(10, 100))
	primaryContactIamIDUpdate := fmt.Sprintf("primary_contact_iam_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIbmEnterpriseDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmEnterpriseConfigBasic(sourceAccountID, name, primaryContactIamID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmEnterpriseExists("ibm_enterprise.enterprise", conf),
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "source_account_id", sourceAccountID),
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "name", name),
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "primary_contact_iam_id", primaryContactIamID),
				),
			},
			resource.TestStep{
				Config: testAccCheckIbmEnterpriseConfigBasic(sourceAccountIDUpdate, nameUpdate, primaryContactIamIDUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "source_account_id", sourceAccountIDUpdate),
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "primary_contact_iam_id", primaryContactIamIDUpdate),
				),
			},
		},
	})
}

func TestAccIbmEnterpriseAllArgs(t *testing.T) {
	var conf enterprisemanagementv1.Enterprise
	sourceAccountID := fmt.Sprintf("source_account_id_%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("name_%d", acctest.RandIntRange(10, 100))
	primaryContactIamID := fmt.Sprintf("primary_contact_iam_id_%d", acctest.RandIntRange(10, 100))
	domain := fmt.Sprintf("domain_%d", acctest.RandIntRange(10, 100))
	sourceAccountIDUpdate := fmt.Sprintf("source_account_id_%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("name_%d", acctest.RandIntRange(10, 100))
	primaryContactIamIDUpdate := fmt.Sprintf("primary_contact_iam_id_%d", acctest.RandIntRange(10, 100))
	domainUpdate := fmt.Sprintf("domain_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIbmEnterpriseDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmEnterpriseConfig(sourceAccountID, name, primaryContactIamID, domain),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmEnterpriseExists("ibm_enterprise.enterprise", conf),
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "source_account_id", sourceAccountID),
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "name", name),
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "primary_contact_iam_id", primaryContactIamID),
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "domain", domain),
				),
			},
			resource.TestStep{
				Config: testAccCheckIbmEnterpriseConfig(sourceAccountIDUpdate, nameUpdate, primaryContactIamIDUpdate, domainUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "source_account_id", sourceAccountIDUpdate),
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "primary_contact_iam_id", primaryContactIamIDUpdate),
					resource.TestCheckResourceAttr("ibm_enterprise.enterprise", "domain", domainUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_enterprise.enterprise",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIbmEnterpriseConfigBasic(sourceAccountID string, name string, primaryContactIamID string) string {
	return fmt.Sprintf(`

		resource "ibm_enterprise" "enterprise" {
			source_account_id = "%s"
			name = "%s"
			primary_contact_iam_id = "%s"
		}
	`, sourceAccountID, name, primaryContactIamID)
}

func testAccCheckIbmEnterpriseConfig(sourceAccountID string, name string, primaryContactIamID string, domain string) string {
	return fmt.Sprintf(`

		resource "ibm_enterprise" "enterprise" {
			source_account_id = "%s"
			name = "%s"
			primary_contact_iam_id = "%s"
			domain = "%s"
		}
	`, sourceAccountID, name, primaryContactIamID, domain)
}

func testAccCheckIbmEnterpriseExists(n string, obj enterprisemanagementv1.Enterprise) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		enterpriseManagementClient, err := testAccProvider.Meta().(ClientSession).EnterpriseManagementV1()
		if err != nil {
			return err
		}

		getEnterpriseOptions := &enterprisemanagementv1.GetEnterpriseOptions{}

		getEnterpriseOptions.SetEnterpriseID(rs.Primary.ID)

		enterprise, _, err := enterpriseManagementClient.GetEnterprise(getEnterpriseOptions)
		if err != nil {
			return err
		}

		obj = *enterprise
		return nil
	}
}

func testAccCheckIbmEnterpriseDestroy(s *terraform.State) error {
	enterpriseManagementClient, err := testAccProvider.Meta().(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_enterprise" {
			continue
		}

		getEnterpriseOptions := &enterprisemanagementv1.GetEnterpriseOptions{}

		getEnterpriseOptions.SetEnterpriseID(rs.Primary.ID)

		// Try to find the key
		_, response, err := enterpriseManagementClient.GetEnterprise(getEnterpriseOptions)

		if err == nil {
			return fmt.Errorf("enterprise still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for enterprise (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

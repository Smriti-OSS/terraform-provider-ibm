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
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"time"

	"github.com/IBM/platform-services-go-sdk/enterprisemanagementv1"
)

func resourceIbmEnterprise() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmEnterpriseCreate,
		ReadContext:   resourceIbmEnterpriseRead,
		UpdateContext: resourceIbmEnterpriseUpdate,
		DeleteContext: resourceIbmEnterpriseDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"source_account_id": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				Description:      "The ID of the account that is used to create the enterprise.",
				DiffSuppressFunc: applyOnce,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the enterprise. This field must have 3 - 60 characters.",
			},
			"primary_contact_iam_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IAM ID of the enterprise primary contact, such as `IBMid-0123ABC`. The IAM ID must already exist.",
			},
			"domain": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A domain or subdomain for the enterprise, such as `example.com` or `my.example.com`.",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the enterprise.",
			},
			"enterprise_account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The enterprise account ID.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Cloud Resource Name (CRN) of the enterprise.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the enterprise.",
			},
			"primary_contact_email": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email of the primary contact of the enterprise.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time stamp at which the enterprise was created.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IAM ID of the user or service that created the enterprise.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time stamp at which the enterprise was last updated.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IAM ID of the user or service that updated the enterprise.",
			},
		},
	}
}

func resourceIbmEnterpriseCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}
	createEnterpriseOptions := &enterprisemanagementv1.CreateEnterpriseOptions{}
	userDetails, err := meta.(ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(err)
	}
	accountID := userDetails.userAccount
	userManagement, err := meta.(ClientSession).UserManagementAPI()
	if err != nil {
		return diag.FromErr(err)
	}
	client := userManagement.UserInvite()
	res, err := client.ListUsers(accountID)
	if err != nil {
		return diag.FromErr(err)
	}
	var iamID string
	for _, userInfo := range res {
		if userInfo.State == "ACTIVE" && userInfo.AccountID == accountID {
			iamID = userInfo.IamID
		}
	}
	d.Set("source_account_id", accountID)
	d.Set("primary_contact_iam_id", iamID)
	createEnterpriseOptions.SetSourceAccountID(d.Get("source_account_id").(string))
	createEnterpriseOptions.SetName(d.Get("name").(string))
	createEnterpriseOptions.SetPrimaryContactIamID(d.Get("primary_contact_iam_id").(string))
	if _, ok := d.GetOk("domain"); ok {
		createEnterpriseOptions.SetDomain(d.Get("domain").(string))
	}
	createEnterpriseResponse, response, err := enterpriseManagementClient.CreateEnterpriseWithContext(context, createEnterpriseOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateEnterpriseWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	d.SetId(*createEnterpriseResponse.EnterpriseID)
	return resourceIbmEnterpriseRead(context, d, meta)
}

func resourceIbmEnterpriseRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getEnterpriseOptions := &enterprisemanagementv1.GetEnterpriseOptions{}

	getEnterpriseOptions.SetEnterpriseID(d.Id())

	enterprise, response, err := enterpriseManagementClient.GetEnterpriseWithContext(context, getEnterpriseOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetEnterpriseWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	//if err = d.Set("source_account_id", enterprise.); err != nil {
	//	return diag.FromErr(fmt.Errorf("Error setting source_account_id: %s", err))
	//}
	if err = d.Set("name", enterprise.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("primary_contact_iam_id", enterprise.PrimaryContactIamID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting primary_contact_iam_id: %s", err))
	}
	if err = d.Set("domain", enterprise.Domain); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting domain: %s", err))
	}
	if err = d.Set("url", enterprise.URL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting url: %s", err))
	}
	if err = d.Set("enterprise_account_id", enterprise.EnterpriseAccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting enterprise_account_id: %s", err))
	}
	if err = d.Set("crn", enterprise.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("state", enterprise.State); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}
	if err = d.Set("primary_contact_email", enterprise.PrimaryContactEmail); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting primary_contact_email: %s", err))
	}
	if err = d.Set("created_at", enterprise.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("created_by", enterprise.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("updated_at", enterprise.UpdatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("updated_by", enterprise.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
	}

	return nil
}

func resourceIbmEnterpriseUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateEnterpriseOptions := &enterprisemanagementv1.UpdateEnterpriseOptions{}

	updateEnterpriseOptions.SetEnterpriseID(d.Id())

	hasChange := false

	//if d.HasChange("source_account_id") {
	//
	//	updateEnterpriseOptions.SetSourceAccountID(d.Get("source_account_id").(string))
	//	hasChange = true
	//}
	if d.HasChange("name") {
		updateEnterpriseOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("primary_contact_iam_id") {
		updateEnterpriseOptions.SetPrimaryContactIamID(d.Get("primary_contact_iam_id").(string))
		hasChange = true
	}
	if d.HasChange("domain") {
		updateEnterpriseOptions.SetDomain(d.Get("domain").(string))
		hasChange = true
	}

	if hasChange {
		response, err := enterpriseManagementClient.UpdateEnterpriseWithContext(context, updateEnterpriseOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateEnterpriseWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
	}

	return resourceIbmEnterpriseRead(context, d, meta)
}

func resourceIbmEnterpriseDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	return nil
}

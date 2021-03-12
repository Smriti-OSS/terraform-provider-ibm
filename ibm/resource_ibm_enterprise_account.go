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
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/enterprisemanagementv1"
)

func resourceIbmEnterpriseAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmEnterpriseAccountCreate,
		ReadContext:   resourceIbmEnterpriseAccountRead,
		UpdateContext: resourceIbmEnterpriseAccountUpdate,
		DeleteContext: resourceIbmEnterpriseAccountDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"parent": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The CRN of the parent under which the account will be created. The parent can be an existing account group or the enterprise itself.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the account. This field must have 3 - 60 characters.",
			},
			"owner_iam_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The IAM ID of the account owner, such as `IBMid-0123ABC`. The IAM ID must already exist.",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the account.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Cloud Resource Name (CRN) of the account.",
			},
			"enterprise_account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The enterprise account ID.",
			},
			"enterprise_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The enterprise ID that the account is a part of.",
			},
			"enterprise_path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path from the enterprise to this particular account.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the account.",
			},
			"paid": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The type of account - whether it is free or paid.",
			},
			"owner_email": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The email address of the owner of the account.",
			},
			"is_enterprise_account": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag to indicate whether the account is an enterprise account or not.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time stamp at which the account was created.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IAM ID of the user or service that created the account.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time stamp at which the account was last updated.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IAM ID of the user or service that updated the account.",
			},
		},
	}
}

func resourceIbmEnterpriseAccountCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createAccountOptions := &enterprisemanagementv1.CreateAccountOptions{}

	createAccountOptions.SetParent(d.Get("parent").(string))
	createAccountOptions.SetName(d.Get("name").(string))
	createAccountOptions.SetOwnerIamID(d.Get("owner_iam_id").(string))

	createAccountResponse, response, err := enterpriseManagementClient.CreateAccountWithContext(context, createAccountOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateAccountWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(*createAccountResponse.AccountID)

	return resourceIbmEnterpriseAccountRead(context, d, meta)
}

func resourceIbmEnterpriseAccountRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountOptions := &enterprisemanagementv1.GetAccountOptions{}

	getAccountOptions.SetAccountID(d.Id())

	account, response, err := enterpriseManagementClient.GetAccountWithContext(context, getAccountOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetAccountWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	if err = d.Set("parent", account.Parent); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting parent: %s", err))
	}
	if err = d.Set("name", account.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("owner_iam_id", account.OwnerIamID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting owner_iam_id: %s", err))
	}
	if err = d.Set("url", account.URL); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting url: %s", err))
	}

	if err = d.Set("crn", account.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("enterprise_account_id", account.EnterpriseAccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting enterprise_account_id: %s", err))
	}
	if err = d.Set("enterprise_id", account.EnterpriseID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting enterprise_id: %s", err))
	}
	if err = d.Set("enterprise_path", account.EnterprisePath); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting enterprise_path: %s", err))
	}
	if err = d.Set("state", account.State); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}
	if err = d.Set("paid", account.Paid); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting paid: %s", err))
	}
	if err = d.Set("owner_email", account.OwnerEmail); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting owner_email: %s", err))
	}
	if err = d.Set("is_enterprise_account", account.IsEnterpriseAccount); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting is_enterprise_account: %s", err))
	}
	if err = d.Set("created_at", account.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("created_by", account.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("updated_at", account.UpdatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("updated_by", account.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
	}

	return nil
}

func resourceIbmEnterpriseAccountUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enterpriseManagementClient, err := meta.(ClientSession).EnterpriseManagementV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateAccountOptions := &enterprisemanagementv1.UpdateAccountOptions{}

	updateAccountOptions.SetAccountID(d.Id())

	hasChange := false

	if d.HasChange("parent") {
		updateAccountOptions.SetParent(d.Get("parent").(string))
		hasChange = true
	}
	/** Removed as update call requires only parent **/
	//if d.HasChange("name") {
	//
	//	updateAccountOptions.SetName(d.Get("name").(string))
	//	hasChange = true
	//}
	//if d.HasChange("owner_iam_id") {
	//	updateAccountOptions.SetOwnerIamID(d.Get("owner_iam_id").(string))
	//	hasChange = true
	//}

	if hasChange {
		response, err := enterpriseManagementClient.UpdateAccountWithContext(context, updateAccountOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateAccountWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
	}

	return resourceIbmEnterpriseAccountRead(context, d, meta)
}

func resourceIbmEnterpriseAccountDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	d.SetId("")

	return nil
}

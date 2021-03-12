provider "ibm" {
  ibmcloud_api_key = var.ibmcloud_api_key
}

// Provision enterprise resource instance
resource "ibm_enterprise" "enterprise_instance" {
  source_account_id = var.enterprise_source_account_id
  name = var.enterprise_name
  primary_contact_iam_id = var.enterprise_primary_contact_iam_id
  domain = var.enterprise_domain
}

// Provision enterprise_account_group resource instance
resource "ibm_enterprise_account_group" "enterprise_account_group_instance" {
  parent = var.enterprise_account_group_parent
  name = var.enterprise_account_group_name
  primary_contact_iam_id = var.enterprise_account_group_primary_contact_iam_id
}

// Provision enterprise_account resource instance
resource "ibm_enterprise_account" "enterprise_account_instance" {
  parent = var.enterprise_account_parent
  name = var.enterprise_account_name
  owner_iam_id = var.enterprise_account_owner_iam_id
}

// Create enterprises data source
data "ibm_enterprises" "enterprises_instance" {
  name = var.enterprises_name
}

// Create account_groups data source
data "ibm_account_groups" "account_groups_instance" {
  name = var.account_groups_name
}

// Create accounts data source
data "ibm_accounts" "accounts_instance" {
  name = var.accounts_name
}

---
layout: "ibm"
page_title: "IBM : enterprise_account"
sidebar_current: "docs-ibm-resource-enterprise-account"
description: |-
  Manages enterprise_account.
---

# ibm\_enterprise_account

Provides a resource for enterprise_account. This allows enterprise_account to be created, updated and deleted.

## Example Usage

```hcl
resource "enterprise_account" "enterprise_account" {
  parent = "parent"
  name = "name"
  owner_iam_id = "owner_iam_id"
}
```

## Argument Reference

The following arguments are supported:

* `parent` - (Required, string) The CRN of the parent under which the account will be created. The parent can be an existing account group or the enterprise itself.
* `name` - (Required, string) The name of the account. This field must have 3 - 60 characters.
* `owner_iam_id` - (Required, string) The IAM ID of the account owner, such as `IBMid-0123ABC`. The IAM ID must already exist.

## Attribute Reference

The following attributes are exported:

* `id` - The unique identifier of the enterprise_account.
* `url` - The URL of the account.
* `crn` - The Cloud Resource Name (CRN) of the account.
* `enterprise_account_id` - The enterprise account ID.
* `enterprise_id` - The enterprise ID that the account is a part of.
* `enterprise_path` - The path from the enterprise to this particular account.
* `state` - The state of the account.
* `paid` - The type of account - whether it is free or paid.
* `owner_email` - The email address of the owner of the account.
* `is_enterprise_account` - The flag to indicate whether the account is an enterprise account or not.
* `created_at` - The time stamp at which the account was created.
* `created_by` - The IAM ID of the user or service that created the account.
* `updated_at` - The time stamp at which the account was last updated.
* `updated_by` - The IAM ID of the user or service that updated the account.

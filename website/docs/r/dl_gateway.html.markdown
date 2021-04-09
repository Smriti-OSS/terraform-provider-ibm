---

subcategory: "Direct Link Gateway"
layout: "ibm"
page_title: "IBM : dl_gateway"
description: |-
  Manages IBM Direct Link Gateway.
---

# ibm\_dl_gateway

Provides a direct link gateway resource. This allows direct link gateway to be created, and updated and deleted.

## Example Usage
In the following example, you can create Direct link of Dedicated type:

```hcl
data "ibm_dl_routers" "test_dl_routers" {
		offering_type = "dedicated"
		location_name = "dal10"
	}

resource ibm_dl_gateway test_dl_gateway {
  bgp_asn =  64999
  global = true 
  metered = false
  name = "Gateway1"
  resource_group = "bf823d4f45b64ceaa4671bee0479346e"
  speed_mbps = 1000 
  type =  "dedicated" 
  cross_connect_router = data.ibm_dl_routers.test_dl_routers.cross_connect_routers[0].router_name
  location_name = data.ibm_dl_routers.test_dl_routers.location_name
  customer_name = "Customer1" 
  carrier_name = "Carrier1"

}   
```
In the following example, you can create Direct link of Connect type:
```
data "ibm_dl_ports" "test_ds_dl_ports" {
 
 }
resource "ibm_dl_gateway" "test_dl_connect" {
  bgp_asn =  64999
  global = true
  metered = false
  name = "dl-connect-gw-1"
  speed_mbps = 1000
  type =  "connect"
  port =  data.ibm_dl_ports.test_ds_dl_ports.ports[0].port_id
}

```

## Argument Reference

The following arguments are supported:

* `bgp_asn` - (Required, Forces new resource, integer) The BGP ASN of the Gateway to be created. Example: 64999
* `global` - (Required, boolean) Gateways with global routing (true) can connect to networks outside their associated region.
* `metered` -  (Required, boolean) Metered billing option. When true gateway usage is billed per gigabyte. When false there is no per gigabyte usage charge, instead a flat rate is charged for the gateway.
* `name` - (Required, boolean) The unique user-defined name for this gateway. Example: myGateway
* `speed_mbps` - (Required, integer) Gateway speed in megabits per second. Example: 10.254.30.78/30
* `type` - (Required, Forces new resource, string) Gateway type. Allowable values: [dedicated,connect]. 
* `bgp_base_cidr` - (Optional, string) (DEPRECATED) BGP base CIDR. Field is deprecated. See bgp_ibm_cidr and bgp_cer_cidr for details on how to create a gateway using either automatic or explicit IP assignment. Any bgp_base_cidr value set will be ignored. 
* `bgp_cer_cidr` - (Optional, Forces new resource, string) BGP customer edge router CIDR. For auto IP assignment, omit bgp_cer_cidr and bgp_ibm_cidr. IBM will automatically select values for bgp_cer_cidr and bgp_ibm_cidr.
* `bgp_ibm_cidr` - (Optional, Forces new resource, string) BGP IBM CIDR. For auto IP assignment, omit bgp_cer_cidr and bgp_ibm_cidr. IBM will automatically select values for bgp_cer_cidr and bgp_ibm_cidr.
* `resource_group` - (Optional, Forces new resource, string) Resource group for this resource. If unspecified, the account's default resource group is used. 
* `carrier_name` - (Required for 'dedicated' type, Forces new resource, string) Carrier name. Constraints: 1 ≤ length ≤ 128, Value must match regular expression ^[a-z][A-Z][0-9][ -_]$. Example: myCarrierName
* `cross_connect_router` - (Required for 'dedicated' type,  Forces new resource, string) Cross connect router. Example: xcr01.dal03
* `customer_name` - (Required for 'dedicated' type, Forces new resource, string) Customer name. Constraints: 1 ≤ length ≤ 128, Value must match regular expression ^[a-z][A-Z][0-9][ -_]$. Example: newCustomerName
* `location_name` - (Required for 'dedicated' type, Forces new resource, string) Gateway location. Example: dal03
* `port` - (Required for Direct link Connect type, Forces new resource, string) gateway port for type=connect gateways



## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of this gateway. 
* `name` - The unique user-defined name for this gateway. 
* `crn` - The CRN (Cloud Resource Name) of this gateway. 
* `created_at` - The date and time resource was created.
* `location_display_name` - Gateway location long name. 
* `resource_group` - Resource group reference
* `bgp_asn` - IBM BGP ASN.
* `bgp_status` - Gateway BGP status.
* `completion_notice_reject_reason` - Reason for completion notice rejection. 
* `link_status` - Gateway link status. Only included on type=dedicated gateways. Example: down, up.
* `port` - gateway port for type=connect gateways
* `vlan` - VLAN allocated for this gateway. Only set for type=connect gateways created directly through the IBM portal. 
* `provider_api_managed` - Indicates whether gateway changes must be made via a provider portal.
* `operational_status` - Gateway operational status. For gateways pending LOA approval, patch operational_status to the appropriate value to approve or reject its LOA. Example: loa_accepted

**NOTE:** `operational_status`(Gateway operational status) and `loa_reject_reason`(LOA reject reason) cannot be updated using terraform as the status and reason keeps changing with different workflow actions.   

## Import

ibm_dl_gateway can be imported using gateway id, eg

```
$ terraform import ibm_dl_gateway.example 5ffda12064634723b079acdb018ef308
```

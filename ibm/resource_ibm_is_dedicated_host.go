// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const (
	isDedicatedHostStable     = "stable"
	isDedicatedHostDeleting   = "deleting"
	isDedicatedHostDeleteDone = "done"
	isDedicatedHostFailed     = "failed"

	isDedicatedHostUpdating             = "updating"
	isDedicatedHostProvisioningDone     = "done"
	isDedicatedHostWaiting              = "waiting"
	isDedicatedHostSuspended            = "suspended"
	isDedicatedHostActionStatusStopping = "stopping"
	isDedicatedHostActionStatusStopped  = "stopped"
	isDedicatedHostStatusPending        = "pending"
	isDedicatedHostStatusRunning        = "running"
	isDedicatedHostStatusFailed         = "failed"
)

func resourceIbmIsDedicatedHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsDedicatedHostCreate,
		ReadContext:   resourceIbmIsDedicatedHostRead,
		UpdateContext: resourceIbmIsDedicatedHostUpdate,
		DeleteContext: resourceIbmIsDedicatedHostDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance_placement_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If set to true, instances can be placed on this dedicated host.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: InvokeValidator("ibm_is_dedicated_host", "name"),
				Description:  "The unique user-defined name for this dedicated host. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"profile": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Globally unique name of the dedicated host profile to use for this dedicated host.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The unique identifier for the resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.",
			},
			"host_group": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of the dedicated host group for this dedicated host.",
			},
			"available_memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of memory in gibibytes that is currently available for instances.",
			},
			"available_vcpu": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The available VCPU for the dedicated host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VCPU architecture.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of VCPUs assigned.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the dedicated host was created.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this dedicated host.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this dedicated host.",
			},
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of instances that are allocated to this dedicated host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this virtual server instance.",
						},
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this virtual server instance.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this virtual server instance.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this virtual server instance (and default system hostname).",
						},
					},
				},
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the dedicated host resource.",
			},
			"memory": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total amount of memory in gibibytes for this host.",
			},
			"provisionable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this dedicated host is available for instance creation.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
			"socket_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of sockets for this host.",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The administrative state of the dedicated host.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the dedicated host on which the unexpected property value was encountered.",
			},
			"supported_instance_profiles": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of instance profiles that can be used by instances placed on this dedicated host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this virtual server instance profile.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this virtual server instance profile.",
						},
					},
				},
			},
			"vcpu": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The total VCPU of the dedicated host.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"architecture": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VCPU architecture.",
						},
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of VCPUs assigned.",
						},
					},
				},
			},
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique name of the zone this dedicated host resides in.",
			},
		},
	}
}

func resourceIbmIsDedicatedHostValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		})

	resourceValidator := ResourceValidator{ResourceName: "ibm_is_dedicated_host", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsDedicatedHostCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	createDedicatedHostOptions := &vpcv1.CreateDedicatedHostOptions{}
	dedicatedHostPrototype := vpcv1.DedicatedHostPrototype{}

	if dhname, ok := d.GetOk("name"); ok {

		namestr := dhname.(string)
		dedicatedHostPrototype.Name = &namestr
	}
	if insplacementenabled, ok := d.GetOk("instance_placement_enabled"); ok {
		insplacementenabledbool := insplacementenabled.(bool)
		dedicatedHostPrototype.InstancePlacementEnabled = &insplacementenabledbool
	}

	if dhprofile, ok := d.GetOk("profile"); ok {
		dhprofilename := dhprofile.(string)
		dedicatedHostProfileIdentity := vpcv1.DedicatedHostProfileIdentity{
			Name: &dhprofilename,
		}
		dedicatedHostPrototype.Profile = &dedicatedHostProfileIdentity
	}

	if dhgroup, ok := d.GetOk("host_group"); ok {
		dhgroupid := dhgroup.(string)
		dedicatedHostGroupIdentity := vpcv1.DedicatedHostGroupIdentity{
			ID: &dhgroupid,
		}
		dedicatedHostPrototype.Group = &dedicatedHostGroupIdentity
	}

	if resgroup, ok := d.GetOk("resource_group"); ok {
		resgroupid := resgroup.(string)
		resourceGroupIdentity := vpcv1.ResourceGroupIdentity{
			ID: &resgroupid,
		}
		dedicatedHostPrototype.ResourceGroup = &resourceGroupIdentity
	}

	createDedicatedHostOptions.SetDedicatedHostPrototype(&dedicatedHostPrototype)

	dedicatedHost, response, err := vpcClient.CreateDedicatedHostWithContext(context, createDedicatedHostOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateDedicatedHostWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(*dedicatedHost.ID)

	_, err = isWaitForDedicatedHostAvailable(vpcClient, d.Id(), d.Timeout(schema.TimeoutCreate), d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceIbmIsDedicatedHostRead(context, d, meta)
}

func resourceIbmIsDedicatedHostRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getDedicatedHostOptions := &vpcv1.GetDedicatedHostOptions{}

	getDedicatedHostOptions.SetID(d.Id())

	dedicatedHost, response, err := vpcClient.GetDedicatedHostWithContext(context, getDedicatedHostOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetDedicatedHostWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	if err = d.Set("available_memory", intValue(dedicatedHost.AvailableMemory)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting available_memory: %s", err))
	}
	availableVcpuMap := resourceIbmIsDedicatedHostVCPUToMap(*dedicatedHost.AvailableVcpu)
	if err = d.Set("available_vcpu", []map[string]interface{}{availableVcpuMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting available_vcpu: %s", err))
	}
	if err = d.Set("created_at", dedicatedHost.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("crn", dedicatedHost.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	//groupMap := resourceIbmIsDedicatedHostDedicatedHostGroupReferenceToMap(*dedicatedHost.Group)

	d.Set("host_group", *dedicatedHost.Group.ID)

	if err = d.Set("href", dedicatedHost.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("instance_placement_enabled", dedicatedHost.InstancePlacementEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_placement_enabled: %s", err))
	}
	instances := []map[string]interface{}{}
	for _, instancesItem := range dedicatedHost.Instances {
		instancesItemMap := resourceIbmIsDedicatedHostInstanceReferenceToMap(instancesItem)
		instances = append(instances, instancesItemMap)
	}
	if err = d.Set("instances", instances); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instances: %s", err))
	}
	if err = d.Set("lifecycle_state", dedicatedHost.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("memory", intValue(dedicatedHost.Memory)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting memory: %s", err))
	}
	if err = d.Set("name", dedicatedHost.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("profile", *dedicatedHost.Profile.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting profile: %s", err))
	}
	if err = d.Set("provisionable", dedicatedHost.Provisionable); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting provisionable: %s", err))
	}
	if err = d.Set("resource_group", *dedicatedHost.ResourceGroup.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group: %s", err))
	}
	if err = d.Set("resource_type", dedicatedHost.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if err = d.Set("socket_count", intValue(dedicatedHost.SocketCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting socket_count: %s", err))
	}
	if err = d.Set("state", dedicatedHost.State); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}
	supportedInstanceProfiles := []map[string]interface{}{}
	for _, supportedInstanceProfilesItem := range dedicatedHost.SupportedInstanceProfiles {
		supportedInstanceProfilesItemMap := resourceIbmIsDedicatedHostInstanceProfileReferenceToMap(supportedInstanceProfilesItem)
		supportedInstanceProfiles = append(supportedInstanceProfiles, supportedInstanceProfilesItemMap)
	}
	if err = d.Set("supported_instance_profiles", supportedInstanceProfiles); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting supported_instance_profiles: %s", err))
	}
	vcpuMap := resourceIbmIsDedicatedHostVCPUToMap(*dedicatedHost.Vcpu)
	if err = d.Set("vcpu", []map[string]interface{}{vcpuMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vcpu: %s", err))
	}

	if err = d.Set("zone", *dedicatedHost.Zone.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting zone: %s", err))
	}

	return nil
}

func resourceIbmIsDedicatedHostVCPUToMap(vCPU vpcv1.Vcpu) map[string]interface{} {
	vCPUMap := map[string]interface{}{}

	vCPUMap["architecture"] = vCPU.Architecture
	vCPUMap["count"] = intValue(vCPU.Count)

	return vCPUMap
}

func resourceIbmIsDedicatedHostInstanceReferenceToMap(instanceReference vpcv1.InstanceReference) map[string]interface{} {
	instanceReferenceMap := map[string]interface{}{}

	instanceReferenceMap["crn"] = instanceReference.CRN
	if instanceReference.Deleted != nil {
		DeletedMap := resourceIbmIsDedicatedHostInstanceReferenceDeletedToMap(*instanceReference.Deleted)
		instanceReferenceMap["deleted"] = []map[string]interface{}{DeletedMap}
	}
	instanceReferenceMap["href"] = instanceReference.Href
	instanceReferenceMap["id"] = instanceReference.ID
	instanceReferenceMap["name"] = instanceReference.Name

	return instanceReferenceMap
}

func resourceIbmIsDedicatedHostInstanceReferenceDeletedToMap(instanceReferenceDeleted vpcv1.InstanceReferenceDeleted) map[string]interface{} {
	instanceReferenceDeletedMap := map[string]interface{}{}

	instanceReferenceDeletedMap["more_info"] = instanceReferenceDeleted.MoreInfo

	return instanceReferenceDeletedMap
}

func resourceIbmIsDedicatedHostInstanceProfileReferenceToMap(instanceProfileReference vpcv1.InstanceProfileReference) map[string]interface{} {
	instanceProfileReferenceMap := map[string]interface{}{}

	instanceProfileReferenceMap["href"] = instanceProfileReference.Href
	instanceProfileReferenceMap["name"] = instanceProfileReference.Name

	return instanceProfileReferenceMap
}

func resourceIbmIsDedicatedHostUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateDedicatedHostOptions := &vpcv1.UpdateDedicatedHostOptions{}

	updateDedicatedHostOptions.SetID(d.Id())

	hasChange := false

	dedicatedHostPrototypemap := map[string]interface{}{}

	if d.HasChange("name") {

		dedicatedHostPrototypemap["name"] = d.Get("name").(interface{})
		hasChange = true
	}
	if d.HasChange("instance_placement_enabled") {

		dedicatedHostPrototypemap["instance_placement_enabled"] = d.Get("instance_placement_enabled").(interface{})
		hasChange = true
	}
	if d.HasChange("profile") {
		dedicatedHostPrototypemap["profile"] = d.Get("profile").(interface{})
		hasChange = true
	}
	if d.HasChange("resource_group") {
		dedicatedHostPrototypemap["resource_group"] = d.Get("resource_group").(interface{})
		hasChange = true
	}
	if d.HasChange("host_group") {
		dedicatedHostPrototypemap["group"] = d.Get("host_group").(interface{})
		hasChange = true
	}

	if hasChange {
		updateDedicatedHostOptions.SetDedicatedHostPatch(dedicatedHostPrototypemap)
		_, response, err := vpcClient.UpdateDedicatedHostWithContext(context, updateDedicatedHostOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateDedicatedHostWithContext fails %s\n%s", err, response)
			return diag.FromErr(err)
		}
	}

	return resourceIbmIsDedicatedHostRead(context, d, meta)
}

func resourceIbmIsDedicatedHostDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getDedicatedHostOptions := &vpcv1.GetDedicatedHostOptions{}

	getDedicatedHostOptions.SetID(d.Id())

	_, response, err := vpcClient.GetDedicatedHostWithContext(context, getDedicatedHostOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetDedicatedHostWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	updateDedicatedHostOptions := &vpcv1.UpdateDedicatedHostOptions{}
	dedicatedHostPrototypeMap := map[string]interface{}{}
	dedicatedHostPrototypeMap["instance_placement_enabled"] = core.BoolPtr(false)
	updateDedicatedHostOptions.SetID(d.Id())
	updateDedicatedHostOptions.SetDedicatedHostPatch(dedicatedHostPrototypeMap)
	_, updateresponse, err := vpcClient.UpdateDedicatedHostWithContext(context, updateDedicatedHostOptions)
	if err != nil {
		log.Printf("[DEBUG] UpdateDedicatedHostWithContext failed %s\n%s", err, updateresponse)
		return diag.FromErr(err)
	}

	deleteDedicatedHostOptions := &vpcv1.DeleteDedicatedHostOptions{}

	deleteDedicatedHostOptions.SetID(d.Id())

	response, err = vpcClient.DeleteDedicatedHostWithContext(context, deleteDedicatedHostOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteDedicatedHostWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	_, err = isWaitForDedicatedHostDelete(vpcClient, d, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func isWaitForDedicatedHostDelete(instanceC *vpcv1.VpcV1, d *schema.ResourceData, id string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{isDedicatedHostDeleting, isDedicatedHostStable},
		Target:  []string{isDedicatedHostDeleteDone, ""},
		Refresh: func() (interface{}, string, error) {
			getdhoptions := &vpcv1.GetDedicatedHostOptions{
				ID: &id,
			}
			dedicatedhost, response, err := instanceC.GetDedicatedHost(getdhoptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return dedicatedhost, isDedicatedHostDeleteDone, nil
				}
				return nil, "", fmt.Errorf("Error getting dedicated Host: %s\n%s", err, response)
			}
			if *dedicatedhost.State == isDedicatedHostFailed {
				return dedicatedhost, *dedicatedhost.State, fmt.Errorf("The  Dedicated host %s failed to delete: %v", d.Id(), err)
			}
			return dedicatedhost, isDedicatedHostDeleting, nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isWaitForDedicatedHostAvailable(instanceC *vpcv1.VpcV1, id string, timeout time.Duration, d *schema.ResourceData) (interface{}, error) {
	log.Printf("Waiting for dedicated host (%s) to be available.", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{isDedicatedHostStatusPending, isDedicatedHostUpdating, isDedicatedHostWaiting},
		Target:     []string{isDedicatedHostFailed, isDedicatedHostStable, isDedicatedHostSuspended},
		Refresh:    isDedicatedHostRefreshFunc(instanceC, id, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isDedicatedHostRefreshFunc(instanceC *vpcv1.VpcV1, id string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getinsOptions := &vpcv1.GetDedicatedHostOptions{
			ID: &id,
		}
		dhost, response, err := instanceC.GetDedicatedHost(getinsOptions)
		if dhost == nil || err != nil {
			return nil, "", fmt.Errorf("Error getting dedicated host : %s\n%s", err, response)
		}
		d.Set("state", *dhost.State)
		d.Set("lifecycle_state", *dhost.LifecycleState)

		if *dhost.LifecycleState == isDedicatedHostSuspended || *dhost.LifecycleState == isDedicatedHostFailed {

			return dhost, *dhost.LifecycleState, fmt.Errorf("status of dedicated host is %s : %s\n%s", *dhost.LifecycleState, err, response)

		}
		return dhost, *dhost.LifecycleState, nil
	}
}

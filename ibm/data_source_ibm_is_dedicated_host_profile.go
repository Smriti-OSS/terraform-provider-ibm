// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func dataSourceIbmIsDedicatedHostProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsDedicatedHostProfileRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The globally unique name for this virtual server instance profile.",
			},
			"class": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product class this dedicated host profile belongs to.",
			},
			"family": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product family this dedicated host profile belongs toThe enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this dedicated host.",
			},
			"memory": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
						"default": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default value for this profile field.",
						},
						"max": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field.",
						},
						"min": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field.",
						},
						"step": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"values": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted values for this profile field.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"socket_count": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
						"default": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default value for this profile field.",
						},
						"max": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field.",
						},
						"min": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field.",
						},
						"step": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"values": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted values for this profile field.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"supported_instance_profiles": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Array of instance profiles that can be used by instances placed on dedicated hosts with this profile.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this virtual server instance profile.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this virtual server instance profile.",
						},
					},
				},
			},
			"vcpu_architecture": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The VCPU architecture for a dedicated host with this profile.",
						},
					},
				},
			},
			"vcpu_count": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type for this profile field.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The value for this profile field.",
						},
						"default": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The default value for this profile field.",
						},
						"max": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value for this profile field.",
						},
						"min": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value for this profile field.",
						},
						"step": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The increment step value for this profile field.",
						},
						"values": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The permitted values for this profile field.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmIsDedicatedHostProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	getDedicatedHostProfileOptions := &vpcv1.GetDedicatedHostProfileOptions{
		Name: &name,
	}
	dedicatedHostProfile, response, err := vpcClient.GetDedicatedHostProfileWithContext(context, getDedicatedHostProfileOptions)
	if err != nil {
		log.Printf("[DEBUG] ListDedicatedHostProfilesWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	if dedicatedHostProfile == nil {
		return diag.FromErr(fmt.Errorf("No Dedicated Host Profile found with name %s", name))
	}
	d.SetId(dataSourceIbmIsDedicatedHostProfileID(d))

	if err = d.Set("class", dedicatedHostProfile.Class); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting class: %s", err))
	}
	if err = d.Set("family", dedicatedHostProfile.Family); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting family: %s", err))
	}
	if err = d.Set("href", dedicatedHostProfile.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if dedicatedHostProfile.Memory != nil {
		err = d.Set("memory", dataSourceDedicatedHostProfileFlattenMemory(*dedicatedHostProfile.Memory.(*vpcv1.DedicatedHostProfileMemory)))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting memory %s", err))
		}
	}

	if dedicatedHostProfile.SocketCount != nil {
		err = d.Set("socket_count", dataSourceDedicatedHostProfileFlattenSocketCount(*dedicatedHostProfile.SocketCount.(*vpcv1.DedicatedHostProfileSocket)))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting socket_count %s", err))
		}
	}

	if dedicatedHostProfile.SupportedInstanceProfiles != nil {
		err = d.Set("supported_instance_profiles", dataSourceDedicatedHostProfileFlattenSupportedInstanceProfiles(dedicatedHostProfile.SupportedInstanceProfiles))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting supported_instance_profiles %s", err))
		}
	}

	if dedicatedHostProfile.VcpuArchitecture != nil {
		err = d.Set("vcpu_architecture", dataSourceDedicatedHostProfileFlattenVcpuArchitecture(*dedicatedHostProfile.VcpuArchitecture))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting vcpu_architecture %s", err))
		}
	}

	if dedicatedHostProfile.VcpuCount != nil {
		err = d.Set("vcpu_count", dataSourceDedicatedHostProfileFlattenVcpuCount(*dedicatedHostProfile.VcpuCount.(*vpcv1.DedicatedHostProfileVcpu)))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting vcpu_count %s", err))
		}
	}

	return nil

}

// dataSourceIbmIsDedicatedHostProfileID returns a reasonable ID for the list.
func dataSourceIbmIsDedicatedHostProfileID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceDedicatedHostProfileFlattenMemory(result vpcv1.DedicatedHostProfileMemory) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostProfileMemoryToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostProfileMemoryToMap(memoryItem vpcv1.DedicatedHostProfileMemory) (memoryMap map[string]interface{}) {
	memoryMap = map[string]interface{}{}

	if memoryItem.Type != nil {
		memoryMap["type"] = memoryItem.Type
	}
	if memoryItem.Value != nil {
		memoryMap["value"] = memoryItem.Value
	}
	if memoryItem.Default != nil {
		memoryMap["default"] = memoryItem.Default
	}
	if memoryItem.Max != nil {
		memoryMap["max"] = memoryItem.Max
	}
	if memoryItem.Min != nil {
		memoryMap["min"] = memoryItem.Min
	}
	if memoryItem.Step != nil {
		memoryMap["step"] = memoryItem.Step
	}
	if memoryItem.Values != nil {
		memoryMap["values"] = memoryItem.Values
	}

	return memoryMap
}

func dataSourceDedicatedHostProfileFlattenSocketCount(result vpcv1.DedicatedHostProfileSocket) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostProfileSocketCountToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostProfileSocketCountToMap(socketCountItem vpcv1.DedicatedHostProfileSocket) (socketCountMap map[string]interface{}) {
	socketCountMap = map[string]interface{}{}

	if socketCountItem.Type != nil {
		socketCountMap["type"] = socketCountItem.Type
	}
	if socketCountItem.Value != nil {
		socketCountMap["value"] = socketCountItem.Value
	}
	if socketCountItem.Default != nil {
		socketCountMap["default"] = socketCountItem.Default
	}
	if socketCountItem.Max != nil {
		socketCountMap["max"] = socketCountItem.Max
	}
	if socketCountItem.Min != nil {
		socketCountMap["min"] = socketCountItem.Min
	}
	if socketCountItem.Step != nil {
		socketCountMap["step"] = socketCountItem.Step
	}
	if socketCountItem.Values != nil {
		socketCountMap["values"] = socketCountItem.Values
	}

	return socketCountMap
}

func dataSourceDedicatedHostProfileFlattenSupportedInstanceProfiles(result []vpcv1.InstanceProfileReference) (supportedInstanceProfiles []map[string]interface{}) {
	for _, supportedInstanceProfilesItem := range result {
		supportedInstanceProfiles = append(supportedInstanceProfiles, dataSourceDedicatedHostProfileSupportedInstanceProfilesToMap(supportedInstanceProfilesItem))
	}

	return supportedInstanceProfiles
}

func dataSourceDedicatedHostProfileSupportedInstanceProfilesToMap(supportedInstanceProfilesItem vpcv1.InstanceProfileReference) (supportedInstanceProfilesMap map[string]interface{}) {
	supportedInstanceProfilesMap = map[string]interface{}{}

	if supportedInstanceProfilesItem.Href != nil {
		supportedInstanceProfilesMap["href"] = supportedInstanceProfilesItem.Href
	}
	if supportedInstanceProfilesItem.Name != nil {
		supportedInstanceProfilesMap["name"] = supportedInstanceProfilesItem.Name
	}

	return supportedInstanceProfilesMap
}

func dataSourceDedicatedHostProfileFlattenVcpuArchitecture(result vpcv1.DedicatedHostProfileVcpuArchitecture) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostProfileVcpuArchitectureToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostProfileVcpuArchitectureToMap(vcpuArchitectureItem vpcv1.DedicatedHostProfileVcpuArchitecture) (vcpuArchitectureMap map[string]interface{}) {
	vcpuArchitectureMap = map[string]interface{}{}

	if vcpuArchitectureItem.Type != nil {
		vcpuArchitectureMap["type"] = vcpuArchitectureItem.Type
	}
	if vcpuArchitectureItem.Value != nil {
		vcpuArchitectureMap["value"] = vcpuArchitectureItem.Value
	}

	return vcpuArchitectureMap
}

func dataSourceDedicatedHostProfileFlattenVcpuCount(result vpcv1.DedicatedHostProfileVcpu) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceDedicatedHostProfileVcpuCountToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceDedicatedHostProfileVcpuCountToMap(vcpuCountItem vpcv1.DedicatedHostProfileVcpu) (vcpuCountMap map[string]interface{}) {
	vcpuCountMap = map[string]interface{}{}

	if vcpuCountItem.Type != nil {
		vcpuCountMap["type"] = vcpuCountItem.Type
	}
	if vcpuCountItem.Value != nil {
		vcpuCountMap["value"] = vcpuCountItem.Value
	}
	if vcpuCountItem.Default != nil {
		vcpuCountMap["default"] = vcpuCountItem.Default
	}
	if vcpuCountItem.Max != nil {
		vcpuCountMap["max"] = vcpuCountItem.Max
	}
	if vcpuCountItem.Min != nil {
		vcpuCountMap["min"] = vcpuCountItem.Min
	}
	if vcpuCountItem.Step != nil {
		vcpuCountMap["step"] = vcpuCountItem.Step
	}
	if vcpuCountItem.Values != nil {
		vcpuCountMap["values"] = vcpuCountItem.Values
	}

	return vcpuCountMap
}

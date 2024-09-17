package branch

import (
	"maps"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/policy"
)

type policyApplicability struct {
	Default     string
	Conditional string
}

var applicability = policyApplicability{
	Default:     "default",
	Conditional: "conditional",
}

func ResourceBranchPolicyStatusCheck() *schema.Resource {
	resource := genBasePolicyResource(&policyCrudArgs{
		FlattenFunc: statusCheckFlattenFunc,
		ExpandFunc:  statusCheckExpandFunc,
		PolicyType:  StatusCheck,
	})

	settingsSchema := resource.Schema["settings"].Elem.(*schema.Resource).Schema
	maps.Copy(settingsSchema, map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},

		"genre": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"author_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

		"invalidate_on_update": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"applicability": {
			Type:     schema.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				applicability.Default,
				applicability.Conditional,
			}, false),
			Default: applicability.Default,
		},

		"filename_patterns": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"display_name": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "",
		},
	})
	return resource
}

func statusCheckFlattenFunc(d *schema.ResourceData, policyConfig *policy.PolicyConfiguration, projectID *string) error {
	err := baseFlattenFunc(d, policyConfig, projectID)
	if err != nil {
		return err
	}

	policySettings := policyConfig.Settings.(map[string]interface{})

	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})

	settings["name"] = policySettings["statusName"]
	settings["genre"] = policySettings["statusGenre"]
	settings["author_id"] = policySettings["authorId"]
	settings["invalidate_on_update"] = policySettings["invalidateOnSourceUpdate"]
	settings["display_name"] = policySettings["defaultDisplayName"]

	if patterns, ok := policySettings["filenamePatterns"]; ok {
		if patterns != nil {
			settings["filename_patterns"] = policySettings["filenamePatterns"].([]interface{})
		}
	}

	settings["applicability"] = applicability.Default
	if policyApplicability, ok := policySettings["policyApplicability"]; ok {
		if policyApplicability != nil && policyApplicability.(float64) == 1 {
			settings["applicability"] = applicability.Conditional
		}
	}
	_ = d.Set("settings", settingsList)
	return nil
}

func statusCheckExpandFunc(d *schema.ResourceData, typeID uuid.UUID) (*policy.PolicyConfiguration, *string, error) {
	policyConfig, projectID, err := baseExpandFunc(d, typeID)
	if err != nil {
		return nil, nil, err
	}

	settingsList := d.Get("settings").([]interface{})
	settings := settingsList[0].(map[string]interface{})

	policySettings := policyConfig.Settings.(map[string]interface{})
	policySettings["statusName"] = settings["name"].(string)
	policySettings["statusGenre"] = settings["genre"].(string)
	policySettings["authorId"] = settings["author_id"].(string)
	policySettings["invalidateOnSourceUpdate"] = settings["invalidate_on_update"].(bool)
	policySettings["defaultDisplayName"] = settings["display_name"].(string)

	patterns := settings["filename_patterns"].([]interface{})
	patternsArray := make([]string, len(patterns))
	for i, variableGroup := range patterns {
		patternsArray[i] = variableGroup.(string)
	}

	policySettings["filenamePatterns"] = patternsArray

	if v, ok := settings["applicability"].(string); ok {
		if v == applicability.Conditional {
			policySettings["policyApplicability"] = 1
		}
	}

	return policyConfig, projectID, nil
}

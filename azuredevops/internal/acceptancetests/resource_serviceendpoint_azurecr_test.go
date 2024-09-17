//go:build (all || resource_serviceendpoint_azurecr) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_azurecr
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointAzureCR_Spn_Basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_azurecr"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, &[]string{"TEST_ARM_SUBSCRIPTION_ID", "TEST_ARM_SUBSCRIPTION_NAME", "TEST_ARM_TENANT_ID",
				"TEST_ARM_RESOURCE_GROUP", "TEST_ARM_ACR_NAME"})
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclAzureCRSpn(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			}, {
				ResourceName:      tfSvcEpNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceEndpointAzureCR_Spn_Update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_azurecr"
	tfSvcEpNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, &[]string{"TEST_ARM_SUBSCRIPTION_ID", "TEST_ARM_SUBSCRIPTION_NAME", "TEST_ARM_TENANT_ID",
				"TEST_ARM_RESOURCE_GROUP", "TEST_ARM_ACR_NAME"})
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclAzureCRSpn(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			}, {
				ResourceName:      tfSvcEpNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: hclAzureCRSpn(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			}, {
				ResourceName:      tfSvcEpNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfSvcEpNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclAzureCRSpn(projectName, serviceConnectionName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_azurecr" "test" {
  project_id                             = azuredevops_project.test.id
  service_endpoint_authentication_scheme = "ServicePrincipal"
  service_endpoint_name                  = "%s"
  azurecr_spn_tenantid                   = "%s"
  azurecr_subscription_id                = "%s"
  azurecr_subscription_name              = "%s"
  resource_group                         = "%s"
  azurecr_name                           = "%s"
}
`, projectName, serviceConnectionName, os.Getenv("TEST_ARM_TENANT_ID"), os.Getenv("TEST_ARM_SUBSCRIPTION_ID"),
		os.Getenv("TEST_ARM_SUBSCRIPTION_NAME"), os.Getenv("TEST_ARM_RESOURCE_GROUP"), os.Getenv("TEST_ARM_ACR_NAME"))
}

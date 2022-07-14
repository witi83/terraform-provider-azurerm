package lighthouse_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-sdk/resource-manager/managedservices/2019-06-01/registrationassignments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LighthouseAssignmentResource struct{}

func TestAccLighthouseAssignment_basic(t *testing.T) {
	// Multiple tenants are needed to test this acceptance.
	// Second tenant ID needs to be set as a environment variable ARM_TENANT_ID_ALT.
	// ObjectId for user, usergroup or service principal from second Tenant needs to be set as a environment variable ARM_PRINCIPAL_ID_ALT_TENANT.
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	if secondTenantID == "" || principalID == "" {
		t.Skip("Skipping as ARM_TENANT_ID_ALT and/or ARM_PRINCIPAL_ID_ALT_TENANT are not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_lighthouse_assignment", "test")
	r := LighthouseAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(uuid.New().String(), secondTenantID, principalID, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccLighthouseAssignment_requiresImport(t *testing.T) {
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	if secondTenantID == "" || principalID == "" {
		t.Skip("Skipping as ARM_TENANT_ID_ALT and/or ARM_PRINCIPAL_ID_ALT_TENANT are not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_lighthouse_assignment", "test")
	r := LighthouseAssignmentResource{}
	id := uuid.New().String()

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(id, secondTenantID, principalID, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
		{
			Config:      r.requiresImport(id, secondTenantID, principalID, data),
			ExpectError: acceptance.RequiresImportError("azurerm_lighthouse_assignment"),
		},
	})
}

func TestAccLighthouseAssignment_emptyID(t *testing.T) {
	secondTenantID := os.Getenv("ARM_TENANT_ID_ALT")
	principalID := os.Getenv("ARM_PRINCIPAL_ID_ALT_TENANT")
	if secondTenantID == "" || principalID == "" {
		t.Skip("Skipping as ARM_TENANT_ID_ALT and/or ARM_PRINCIPAL_ID_ALT_TENANT are not specified")
	}

	data := acceptance.BuildTestData(t, "azurerm_lighthouse_assignment", "test")
	r := LighthouseAssignmentResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.emptyId(secondTenantID, principalID, data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func (LighthouseAssignmentResource) Exists(ctx context.Context, clients *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := registrationassignments.ParseScopedRegistrationAssignmentID(state.ID)
	if err != nil {
		return nil, err
	}

	options := registrationassignments.GetOperationOptions{
		ExpandRegistrationDefinition: utils.Bool(false),
	}
	resp, err := clients.Lighthouse.AssignmentsClient.Get(ctx, *id, options)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(resp.Model != nil), nil
}

func (LighthouseAssignmentResource) basic(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_role_definition" "contributor" {
  role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
}

resource "azurerm_lighthouse_definition" "test" {
  name               = "acctest-LD-%d"
  description        = "Acceptance Test Lighthouse Definition"
  managing_tenant_id = "%s"
  scope              = data.azurerm_subscription.primary.id

  authorization {
    principal_id       = "%s"
    role_definition_id = data.azurerm_role_definition.contributor.role_definition_id
  }
}

resource "azurerm_lighthouse_assignment" "test" {
  name                     = "%s"
  scope                    = data.azurerm_subscription.primary.id
  lighthouse_definition_id = azurerm_lighthouse_definition.test.id
}
`, data.RandomInteger, secondTenantID, principalID, id)
}

func (r LighthouseAssignmentResource) requiresImport(id string, secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_lighthouse_assignment" "import" {
  name                     = azurerm_lighthouse_assignment.test.name
  lighthouse_definition_id = azurerm_lighthouse_assignment.test.lighthouse_definition_id
  scope                    = azurerm_lighthouse_assignment.test.scope
}
`, r.basic(id, secondTenantID, principalID, data))
}

func (LighthouseAssignmentResource) emptyId(secondTenantID string, principalID string, data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_subscription" "primary" {
}

data "azurerm_role_definition" "contributor" {
  role_definition_id = "b24988ac-6180-42a0-ab88-20f7382dd24c"
}

resource "azurerm_lighthouse_definition" "test" {
  name               = "acctest-LD-%d"
  description        = "Acceptance Test Lighthouse Definition"
  managing_tenant_id = "%s"
  scope              = data.azurerm_subscription.primary.id

  authorization {
    principal_id       = "%s"
    role_definition_id = data.azurerm_role_definition.contributor.role_definition_id
  }
}

resource "azurerm_lighthouse_assignment" "test" {
  scope                    = data.azurerm_subscription.primary.id
  lighthouse_definition_id = azurerm_lighthouse_definition.test.id
}
`, data.RandomInteger, secondTenantID, principalID)
}

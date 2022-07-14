---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_automation_rule"
description: |-
  Manages a Sentinel Automation Rule.
---

# azurerm_sentinel_automation_rule

Manages a Sentinel Automation Rule.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
}

resource "azurerm_log_analytics_solution" "sentinel" {
  solution_name         = "SecurityInsights"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  workspace_resource_id = azurerm_log_analytics_workspace.example.id
  workspace_name        = azurerm_log_analytics_workspace.example.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/SecurityInsights"
  }
}

resource "azurerm_sentinel_automation_rule" "example" {
  name                       = "56094f72-ac3f-40e7-a0c0-47bd95f70336"
  log_analytics_workspace_id = azurerm_log_analytics_solution.sentinel.workspace_resource_id
  display_name               = "automation_rule1"
  order                      = 1
  action_incident {
    order  = 1
    status = "Active"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The UUID which should be used for this Sentinel Automation Rule. Changing this forces a new Sentinel Automation Rule to be created.

* `log_analytics_workspace_id` - (Required) The ID of the Log Analytics Workspace where this Sentinel applies to. Changing this forces a new Sentinel Automation Rule to be created.
  
* `display_name` - (Required) The display name which should be used for this Sentinel Automation Rule.

* `order` - (Required) The order of this Sentinel Automation Rule. Possible values varies between `1` and `1000`.

---

* `action_incident` - (Optional) One or more `action_incident` blocks as defined below.

* `action_playbook` - (Optional) One or more `action_playbook` blocks as defined below.

~> **Note:** Either one `action_incident` block or `action_playbook` block has to be specified.

* `condition` - (Optional) One or more `condition` blocks as defined below.

* `enabled` - (Optional) Whether this Sentinel Automation Rule is enabled? Defaults to `true`.

* `expiration` - (Optional) The time in RFC3339 format of kind `UTC` that determines when this Automation Rule should expire and be disabled.

---

A `action_incident` block supports the following:

* `order` - (Required) The execution order of this action.

* `status` - (Optional) The status to set to the incident. Possible values are: `Active`, `Closed`, `New`.
  
* `classification` - (Optional) The classification of the incident, when closing it. Possible values are: `BenignPositive_SuspiciousButExpected`, `FalsePositive_InaccurateData`, `FalsePositive_IncorrectAlertLogic`, `TruePositive_SuspiciousActivity` and `Undetermined`.
  
  ~> **Note:** The `classification` is required when `status` is `Closed`.

* `classification_comment` - (Optional) The comment why the incident is to be closed.

  ~> **Note:** The `classification_comment` is allowed to set only when `status` is `Closed`.

* `labels` - (Optional) Specifies a list of labels to add to the incident.

* `owner_id` - (Optional) The object ID of the entity this incident is assigned to.

* `severity` - (Optional) The severity to add to the incident.

~> **Note:**: At least one of `status`, `labels`, `owner_id` and `severity` has to be set.

---

A `action_playbook` block supports the following:

* `logic_app_id` - (Required) The ID of the Logic App that defines the playbook's logic.

* `order` - (Required) The execution order of this action.

* `tenant_id` - (Optional) The ID of the Tenant that owns the playbook.

---

A `condition` block supports the following:

* `operator` - (Required) The operator to use for evaluate the condition. Possible values include: `Equals`, `NotEquals`, `Contains`, `NotContains`, `StartsWith`, `NotStartsWith`, `EndsWith`, `NotEndsWith`.

* `property` - (Required) The property to use for evaluate the condition. Possible values include: `AccountAadTenantId`, `AccountAadUserId`, `AccountNTDomain`, `AccountName`, `AccountObjectGuid`, `AccountPUID`, `AccountSid`, `AccountUPNSuffix`, `AzureResourceResourceId`, `AzureResourceSubscriptionId`, `CloudApplicationAppId`, `CloudApplicationAppName`, `DNSDomainName`, `FileDirectory`, `FileHashValue`, `FileName`, `HostAzureID`, `HostNTDomain`, `HostName`, `HostNetBiosName`, `HostOSVersion`, `IPAddress`, `IncidentDescription`, `IncidentProviderName`, `IncidentRelatedAnalyticRuleIds`, `IncidentSeverity`, `IncidentStatus`, `IncidentTactics`, `IncidentTitle`, `IoTDeviceId`, `IoTDeviceModel`, `IoTDeviceName`, `IoTDeviceOperatingSystem`, `IoTDeviceType`, `IoTDeviceVendor`, `MailMessageDeliveryAction`, `MailMessageDeliveryLocation`, `MailMessageP1Sender`, `MailMessageP2Sender`, `MailMessageRecipient`, `MailMessageSenderIP`, `MailMessageSubject`, `MailboxDisplayName`, `MailboxPrimaryAddress`, `MailboxUPN`, `MalwareCategory`, `MalwareName`, `ProcessCommandLine`, `ProcessId`, `RegistryKey`, `RegistryValueData`, `Url`.

* `values` - (Required) Specifies a list of values to use for evaluate the condition.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Sentinel Automation Rule.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Sentinel Automation Rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel Automation Rule.
* `update` - (Defaults to 5 minutes) Used when updating the Sentinel Automation Rule.
* `delete` - (Defaults to 5 minutes) Used when deleting the Sentinel Automation Rule.

## Import

Sentinel Automation Rules can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_automation_rule.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/AutomationRules/rule1
```

package azurerm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func TestAccAzureRMAutomationIntVariable_basic(t *testing.T) {
	resourceName := "azurerm_automation_int_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationIntVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationIntVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationIntVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "1234"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationIntVariable_complete(t *testing.T) {
	resourceName := "azurerm_automation_int_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationIntVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationIntVariable_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationIntVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(resourceName, "value", "12345"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAzureRMAutomationIntVariable_basicCompleteUpdate(t *testing.T) {
	resourceName := "azurerm_automation_int_variable.test"
	ri := tf.AccRandTimeInt()
	location := testLocation()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAzureRMAutomationIntVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAzureRMAutomationIntVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationIntVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "1234"),
				),
			},
			{
				Config: testAccAzureRMAutomationIntVariable_complete(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationIntVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "This variable is created by Terraform acceptance test."),
					resource.TestCheckResourceAttr(resourceName, "value", "12345"),
				),
			},
			{
				Config: testAccAzureRMAutomationIntVariable_basic(ri, location),
				Check: resource.ComposeTestCheckFunc(
					testCheckAzureRMAutomationIntVariableExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "value", "1234"),
				),
			},
		},
	})
}

func testCheckAzureRMAutomationIntVariableExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Automation Int Variable not found: %s", resourceName)
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accountName := rs.Primary.Attributes["automation_account_name"]

		client := testAccProvider.Meta().(*ArmClient).automationVariableClient
		ctx := testAccProvider.Meta().(*ArmClient).StopContext

		if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Automation Int Variable %q (Automation Account Name %q / Resource Group %q) does not exist", name, accountName, resourceGroup)
			}
			return fmt.Errorf("Bad: Get on automationVariableClient: %+v", err)
		}

		return nil
	}
}

func testCheckAzureRMAutomationIntVariableDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ArmClient).automationVariableClient
	ctx := testAccProvider.Meta().(*ArmClient).StopContext

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azurerm_automation_int_variable" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		resourceGroup := rs.Primary.Attributes["resource_group_name"]
		accountName := rs.Primary.Attributes["automation_account_name"]

		if resp, err := client.Get(ctx, resourceGroup, accountName, name); err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Bad: Get on automationVariableClient: %+v", err)
			}
		}

		return nil
	}

	return nil
}

func testAccAzureRMAutomationIntVariable_basic(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku = {
    name = "Basic"
  }
}

resource "azurerm_automation_int_variable" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  value                   = 1234
}
`, rInt, location, rInt, rInt)
}

func testAccAzureRMAutomationIntVariable_complete(rInt int, location string) string {
	return fmt.Sprintf(`
resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%d"
  location = "%s"
}

resource "azurerm_automation_account" "test" {
  name                = "acctestAutoAcct-%d"
  location            = "${azurerm_resource_group.test.location}"
  resource_group_name = "${azurerm_resource_group.test.name}"

  sku = {
    name = "Basic"
  }
}

resource "azurerm_automation_int_variable" "test" {
  name                    = "acctestAutoVar-%d"
  resource_group_name     = "${azurerm_resource_group.test.name}"
  automation_account_name = "${azurerm_automation_account.test.name}"
  description             = "This variable is created by Terraform acceptance test."
  value                   = 12345
}
`, rInt, location, rInt, rInt)
}
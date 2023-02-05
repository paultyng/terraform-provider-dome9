package dome9

import (
	"fmt"
	"github.com/dome9/dome9-sdk-go/services/cloudaccounts/oci"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-dome9/dome9/common/resourcetype"
	"github.com/terraform-providers/terraform-provider-dome9/dome9/common/testing/environmentvariable"
	"github.com/terraform-providers/terraform-provider-dome9/dome9/common/testing/method"
	"github.com/terraform-providers/terraform-provider-dome9/dome9/common/testing/variable"
)

func TestAccResourceCloudAccountOciTempDataBasic(t *testing.T) {
	var cloudAccountOci oci.CloudAccountResponse
	resourceTypeAndName, _, generatedName := method.GenerateRandomSourcesTypeAndName(resourcetype.CloudAccountOCITempData)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccCloudAccountOciEnvVarsPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudAccountOciDestroy,
		Steps: []resource.TestStep{
			{
				// creation test
				Config: testAccCheckCloudAccountOciConfigure(resourceTypeAndName, generatedName, variable.CloudAccountOciCreationResourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudAccountOciExists(resourceTypeAndName, &cloudAccountOci),
					resource.TestCheckResourceAttr(resourceTypeAndName, "name", variable.CloudAccountOciCreationResourceName),
					resource.TestCheckResourceAttr(resourceTypeAndName, "vendor", variable.CloudAccountOciVendor),
					resource.TestCheckResourceAttr(resourceTypeAndName, "organizational_unit_name", os.Getenv(environmentvariable.OrganizationalUnitName)),
				),
			},
		},
	})
}

func testAccCheckCloudAccountOciTempDataDestroy(s *terraform.State) error {

	return nil
}

func testAccCloudAccountOciTempDataEnvVarsPreCheck(t *testing.T) {
	if v := os.Getenv(environmentvariable.OrganizationalUnitName); v == "" {
		t.Fatalf("%s must be set for acceptance tests", environmentvariable.OrganizationalUnitName)
	}
}

func testAccCheckCloudAccountOciExists(resource string, resp *Oci.CloudAccountResponse) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("didn't find resource: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no record ID is set")
		}

		apiClient := testAccProvider.Meta().(*Client)
		receivedCloudAccount, _, err := apiClient.cloudaccountOci.Get(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed fetching resource %s. Recevied error: %s", resource, err)
		}

		*resp = *receivedCloudAccount

		return nil
	}
}

func testAccCheckCloudAccountOciConfigure(resourceTypeAndName, generatedName, resourceName string) string {
	return fmt.Sprintf(`
// oci cloud account creation
%s

data "%s" "%s" {
  id = "${%s.id}"
}
`,
		// oci cloud account
		getCloudAccountOciResourceHCL(generatedName, resourceName),

		// data source variables
		resourcetype.CloudAccountOci,
		generatedName,
		resourceTypeAndName,
	)
}

func getCloudAccountOciResourceHCL(cloudAccountName, generatedAName string) string {
	return fmt.Sprintf(`
resource "%s" "%s" {
	credentials = {
		access_key    = "%s"
		access_secret = "%s"
}
	name          = "%s"
}
`,
		// oci cloud account variables
		resourcetype.CloudAccountOci,
		cloudAccountName,
		os.Getenv(environmentvariable.CloudAccountOciEnvVarAccessKey),
		os.Getenv(environmentvariable.CloudAccountOciEnvVarAccessSecret),
		generatedAName,
	)
}

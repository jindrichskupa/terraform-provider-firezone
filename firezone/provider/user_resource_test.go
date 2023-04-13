package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccUserResourceConfig("one@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firezone_user.test", "email", "one@example.com"),
					resource.TestCheckResourceAttr("firezone_user.test", "role", "unprivileged"),
					resource.TestCheckResourceAttr("firezone_user.test", "id", "example-id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "firezone_user.test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"email", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccUserResourceConfig("two@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firezone_user.test", "email", "two@example.com"),
					resource.TestCheckResourceAttr("firezone_user.test", "role", "unprivileged"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccUserResourceConfig(email string) string {
	return fmt.Sprintf(`
resource "firezone_user" "test" {
  email = %[1]q
	role = %[2]q
}
`, email, "unprivileged")
}

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRuleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccRuleResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firezone_rule.test", "email", "one@example.com"),
					resource.TestCheckResourceAttr("firezone_rule.test", "role", "unprivileged"),
					resource.TestCheckResourceAttr("firezone_rule.test", "id", "example-id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "firezone_rule.test",
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
				Config: testAccRuleResourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("firezone_rule.test", "email", "two@example.com"),
					resource.TestCheckResourceAttr("firezone_rule.test", "role", "unprivileged"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRuleResourceConfig() string {
	return fmt.Sprintf(`
resource "firezone_rule" "test" {
	user_id = "example-id"
	action = "allow"
	destination = "0.0.0.0/0"
	port_range = "0-65535"
	port_type = "tcp"
}
`)
}

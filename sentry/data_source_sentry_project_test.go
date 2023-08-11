package sentry

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSentryProjectDataSource_basic(t *testing.T) {
	rn := "sentry_project.test"
	dn := "data.sentry_project.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSentryProjectDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dn, "organization", rn, "organization"),
					resource.TestCheckResourceAttrPair(dn, "slug", rn, "slug"),
					resource.TestCheckResourceAttrPair(dn, "internal_id", rn, "internal_id"),
					resource.TestCheckResourceAttrPair(dn, "name", rn, "name"),
					resource.TestCheckResourceAttrPair(dn, "has_access", rn, "has_access"),
					resource.TestCheckResourceAttrPair(dn, "is_pending", rn, "is_pending"),
					resource.TestCheckResourceAttrPair(dn, "is_member", rn, "is_member"),
				),
			},
		},
	})
}

func testAccSentryProjectDataSourceConfig() string {
	return testAccSentryOrganizationDataSourceConfig + `
data "sentry_project" "test" {
	organization = sentry_team.test.organization
	slug         = sentry_team.test.id
}
	`
}

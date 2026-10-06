// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"extlib": providerserver.NewProtocol6WithError(New("test")()),
}

// Helper to run test steps with our provider
func providerTest(t *testing.T, steps []resource.TestStep) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// Helper to test the result of an expression in an output.
//
// expr: A terraform expression to check the result of. Generally uses one of our provider functions
// expected: A check to run on the result of that expression.
func exprTest(t *testing.T, expr string, expected knownvalue.Check) {
	providerTest(t, []resource.TestStep{
		{
			Config: fmt.Sprintf(`output "test" {
				value = %s
			}
			`, expr),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownOutputValue("test", expected),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

// Copyright Thayne McCombs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestBoolsetFunction_True(t *testing.T) {
	exprTest(t, "provider::extlib::boolset(true)", knownvalue.SetExact([]knownvalue.Check{knownvalue.StringExact("")}))
}

func TestBoolsetFunction_False(t *testing.T) {
	exprTest(t, "provider::extlib::boolset(false)", knownvalue.SetExact([]knownvalue.Check{}))
}

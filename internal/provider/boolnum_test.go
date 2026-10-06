// Copyright Thayne McCombs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestBoolnumFunction_True(t *testing.T) {
	exprTest(t, "provider::extlib::boolnum(true)", knownvalue.Int32Exact(1))
}

func TestBoolnumFunction_False(t *testing.T) {
	exprTest(t, "provider::extlib::boolnum(false)", knownvalue.Int32Exact(0))
}

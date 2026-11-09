// Copyright Thayne McCombs 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

func TestRepeatFunction_Basic(t *testing.T) {
	exprTest(t, `provider::extlib::repeat("a", 4)`, knownvalue.ListExact([]knownvalue.Check{
		knownvalue.StringExact("a"),
		knownvalue.StringExact("a"),
		knownvalue.StringExact("a"),
		knownvalue.StringExact("a"),
	}))
}

// Doesn't work with opentofu, see github.com/opentofu/opentofu/issues/4819
// func TestRepeatFunction_Null(t *testing.T) {
// 	exprTest(t, `provider::extlib::repeat(tostring(null), 2)`, knownvalue.ListExact([]knownvalue.Check{
// 		knownvalue.Null(),
// 		knownvalue.Null(),
// 	}))
// }

func TestRepeatFunction_Empty(t *testing.T) {
	exprTest(t, `provider::extlib::repeat(1, 0)`, knownvalue.ListSizeExact(0))
}

func TestRepeatFunction_NegSize(t *testing.T) {
	providerTest(t, []resource.TestStep{
		{
			Config:      `output "test" { value = provider::extlib::repeat(1, -1) }`,
			PlanOnly:    true,
			ExpectError: regexp.MustCompile(`Invalid value for "n" parameter: must be non-negative`),
		},
	})
}

func TestRepeatFunction_Unknown(t *testing.T) {
	repeat := NewRepeatFunction()

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{types.DynamicUnknown(), types.Int32Value(3)}),
	}
	resp := function.RunResponse{
		Result: function.NewResultData(types.DynamicUnknown()),
	}

	repeat.Run(t.Context(), req, &resp)
	if resp.Error != nil {
		t.Errorf("Unexpected error: %v", resp.Error)
	}

	v := resp.Result.Value().(types.Dynamic)

	if tuple, ok := v.UnderlyingValue().(types.Tuple); ok {
		elements := tuple.Elements()
		if len(elements) != 3 {
			t.Errorf("length was %d instead of 3", len(elements))
		}
		for i, el := range elements {
			if !el.IsUnknown() {
				t.Errorf("element %d is not unknown", i)
			}
		}
	} else {
		t.Errorf("result (%v)is not a tuple", v)
	}
}

func TestRepeatFunction_DynamicNull(t *testing.T) {
	repeat := NewRepeatFunction()

	req := function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{types.DynamicNull(), types.Int32Value(2)}),
	}
	resp := function.RunResponse{
		Result: function.NewResultData(types.DynamicUnknown()),
	}

	repeat.Run(t.Context(), req, &resp)
	if resp.Error != nil {
		t.Errorf("Unexpected error: %v", resp.Error)
	}

	v := resp.Result.Value().(types.Dynamic)

	if tuple, ok := v.UnderlyingValue().(types.Tuple); ok {
		elements := tuple.Elements()
		if len(elements) != 2 {
			t.Errorf("length was %d instead of 2", len(elements))
		}
		for i, el := range elements {
			if !el.IsNull() {
				t.Errorf("element %d is not null", i)
			}
		}
		for i, eltType := range tuple.ElementTypes(t.Context()) {
			if !eltType.Equal(types.DynamicType) {
				t.Errorf("element type for %d is not dynamic", i)
			}
		}
	} else {
		t.Errorf("result (%v)is not a tuple", v)
	}
}

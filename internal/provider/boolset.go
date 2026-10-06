// Copyright Thayne McCombs 2026
// SPDX-License Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ function.Function = BoolsetFunction{}
)

type BoolsetFunction struct{}

func (f BoolsetFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "boolset"
}

func (f BoolsetFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Convert bool to set",
		// I hate that go doesn't have a way to have multi-line strings that contain "`".
		MarkdownDescription: `Create a set depending on the value of a boolean

If the input is true, create a set containing a single empty string.
If the input is false, create an empty set.

This is useful for conditionally creating 0 or one instances of a dynamic block.`,
		Parameters: []function.Parameter{
			function.BoolParameter{
				Name:                "value",
				MarkdownDescription: "bool to convert",
			},
		},
		Return: function.SetReturn{
			ElementType: types.StringType,
		},
	}
}

func (f BoolsetFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var v bool

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &v))

	if resp.Error != nil {
		return
	}

	var result []string
	if v {
		result = []string{""}
	} else {
		// the difference between nil and an empty slice is significant
		result = []string{}
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result))
}

func NewBoolsetFunction() function.Function {
	return &BoolsetFunction{}
}

// Copyright Thayne McCombs 2026
// SPDX-License Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var (
	_ function.Function = BoolnumFunction{}
)

type BoolnumFunction struct{}

func (f BoolnumFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "boolnum"
}

func (f BoolnumFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Convert bool to number",
		// I hate that go doesn't have a way to have multi-line strings that contain "`".
		MarkdownDescription: "Convert boolean to 0 or 1\n\n" +
			"Convert `false` to `0` and `true` to `1`.",
		Parameters: []function.Parameter{
			function.BoolParameter{
				Name:                "value",
				MarkdownDescription: "bool to convert",
			},
		},
		Return: function.Int32Return{},
	}
}

func (f BoolnumFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var v bool

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &v))

	if resp.Error != nil {
		return
	}

	var result int32 = 0
	if v {
		result = 1
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result))
}

func NewBoolnumFunction() function.Function {
	return &BoolnumFunction{}
}

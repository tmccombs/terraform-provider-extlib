// Copyright Thayne McCombs 2026
// SPDX-License Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RepeatFunction struct{}

type nonNegativeIntValidator struct{}

func (f RepeatFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "repeat"
}

func (f RepeatFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Create a list that repeats a value n times",
		MarkdownDescription: "Create a value that repeats the given value `n` times.",
		Parameters: []function.Parameter{
			function.DynamicParameter{
				Name:                "value",
				MarkdownDescription: "value to repeat",
				AllowNullValue:      true,
				AllowUnknownValues:  true,
			},
			function.Int32Parameter{
				Name:                "n",
				MarkdownDescription: "number of times to repeat the value",
				Validators:          []function.Int32ParameterValidator{nonNegativeIntValidator{}},
			},
		},
		Return: function.DynamicReturn{},
	}
}

func (f RepeatFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var original types.Dynamic
	var n int32

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &original, &n))

	if resp.Error != nil {
		return
	}

	v := original.UnderlyingValue()
	arr := make([]attr.Value, n)

	if v == nil {
		// there is no underlying type for us to use for the list
		v = original
	}

	for i := range arr {
		arr[i] = v
	}

	elementType := v.Type(ctx)

	var list attr.Value
	var diags diag.Diagnostics
	if elementType.Equal(types.DynamicType) {
		// We can't make a list of DynamicType, so make a tuple instead
		tupleTypes := make([]attr.Type, len(arr))
		for i := range tupleTypes {
			tupleTypes[i] = elementType
		}

		list, diags = types.TupleValue(tupleTypes, arr)
	} else {
		list, diags = types.ListValue(elementType, arr)
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, function.FuncErrorFromDiags(ctx, diags))

	if diags.HasError() {
		return
	}

	result := types.DynamicValue(list)

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}

func NewRepeatFunction() function.Function {
	return &RepeatFunction{}
}

func (_ nonNegativeIntValidator) ValidateParameterInt32(ctx context.Context, req function.Int32ParameterValidatorRequest, resp *function.Int32ParameterValidatorResponse) {
	v := req.Value.ValueInt32()
	if v < 0 {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewArgumentFuncError(req.ArgumentPosition, "must be non-negative"))
	}
}

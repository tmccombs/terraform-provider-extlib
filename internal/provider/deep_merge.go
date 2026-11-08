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

var (
	_ function.Function = DeepMergeFunction{}
)

type DeepMergeFunction struct{}

func (f DeepMergeFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "deep_merge"
}

func (f DeepMergeFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Merge two objects recursively",
		// I hate that go doesn't have a way to have multi-line strings that contain "`".
		MarkdownDescription: `Merge the second object into the first one recrusively.
That is, this is the same as merge, but if both objects define an object or map value at the same key, rather than
using the second value, the two objects will be merged`,
		Parameters: []function.Parameter{
			function.DynamicParameter{
				Name:                "base",
				MarkdownDescription: "value to use as the base of the merge",
				// Allow unknown values, because we can potentially partially resolve the object if the value is partially known
				AllowUnknownValues: true,
				// should we validate that it is an object or map?
			},
		},
		VariadicParameter: function.DynamicParameter{
			Name:                "other",
			MarkdownDescription: "addtional objects to merge into the base, in order",
			// Allow unknown values, because we can potentially partially resolve the object if the value is partially known
			AllowUnknownValues: true,
			// should we validate that it is an object or map?
		},
		Return: function.DynamicReturn{},
	}
}

func (f DeepMergeFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var base types.Dynamic
	var rest []types.Dynamic

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &base, &rest))

	if resp.Error != nil {
		return
	}

	// TODO: figure out how to handle if any of the Dynamic objects are themselves unknown
	var result types.Dynamic

	if base.IsUnknown() {
		result = base
		goto end
	}

	{
		value := base.UnderlyingValue()

		for _, obj := range rest {
			if obj.IsUnknown() {
				result = types.DynamicUnknown()
				goto end
			}
			value = deepMerge(result, obj.UnderlyingValue())
		}

		result = types.DynamicValue(value)
	}

end:
	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result))
}

func NewDeepMergeFunction() function.Function {
	return &DeepMergeFunction{}
}

// Warning, this can modify `base`
func deepMerge(base attr.Value, over attr.Value) (attr.Value, diag.Diagnostics) {
	switch b := over.(type) {
	case types.Map:
		return mergeMap(base, b)
	case types.Object:
		return mergeObject(base, b)
	default:
		// For any non-map, non-object value, use the second value
		return over
	}
}

func mergeMap(base attr.Value, over types.Map) (attr.Value, diag.Diagnostics) {
	if over.IsUnknown() {
		// If we don't know the value, we can't merge it
		return over
	}
	switch a := base.(type) {
	case types.Map:
		if a.ElementType() != over.ElementType() {
			// we need to treat the maps as objects
			return mergeMapsAsObjs(a, over)
		}
		// TODO: check that the element types match?
		result := a.Elements()
		for k, v := range over.Elements() {
			if old, ok := result[k]; ok {
				result[k] = deepMerge(old, v)
			} else {
				result[k] = v
			}
		}
		return types.MapValue(over.ElementType(), result)
	// Should we merge the map into an object?
	default:
		// if the original isn't a map or object, just replace it
		return over
	}
}

func mergeMapsAsObjs(ctx context.Context, a types.Map, b types.Map) (types.Map, diag.Diagnostics) {
	aType := a.ElementType()
	var attrTypes map[string]attr.Type
	value := a.Elements()
	for k := range a {
		attrTypes[k] = aType
	}

	bType := b.ElementType()

	for k, v := range b.Elements() {
		if old, ok := value[k]; ok {
			value[k] = deepMerge(old, v)
			// TODO: get context here
			attrTypes[k] = value[k].Type(ctx)
		} else {
			result[k] = v
			attrTypes[k] = bType
		}
	}
}

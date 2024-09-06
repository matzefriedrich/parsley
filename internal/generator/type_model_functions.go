package generator

import (
	"fmt"
	"strings"
)

func RegisterTypeModelFunctions(generator GenericCodeGenerator) error {
	return generator.AddTemplateFunc(
		NamedFunc("HasResults", HasResults),
		NamedFunc("FormattedParameters", FormattedParameters),
		NamedFunc("FormattedCallParameters", FormattedCallParameters),
		NamedFunc("FormattedResultParameters", FormattedResultParameters),
		NamedFunc("FormattedResultTypes", FormattedResultTypes))
}

func HasResults(m Method) bool {
	return len(m.Results) > 0
}

func FormattedParameters(m Method) string {
	formattedParameters := make([]string, len(m.Parameters))
	for i, parameter := range m.Parameters {
		formattedParameters[i] = fmt.Sprintf("%s %s", parameter.Name, parameter.TypeName)
	}
	return strings.Join(formattedParameters, ", ")
}

func FormattedCallParameters(m Method) string {
	formattedParameters := make([]string, len(m.Parameters))
	for i, parameter := range m.Parameters {
		formattedParameters[i] = fmt.Sprintf("%s", parameter.Name)
	}
	return strings.Join(formattedParameters, ", ")
}

func FormattedResultParameters(m Method) string {
	formattedResults := make([]string, len(m.Results))
	for i, result := range m.Results {
		formattedResults[i] = fmt.Sprintf("%s", result.Name)
	}
	return strings.Join(formattedResults, ", ")
}

func FormattedResultTypes(m Method) string {
	formattedResults := make([]string, len(m.Results))
	for i, result := range m.Results {
		formattedResults[i] = fmt.Sprintf("%s", result.TypeName)
	}
	if len(formattedResults) == 0 {
		return ""
	}
	return "(" + strings.Join(formattedResults, ", ") + ")"
}

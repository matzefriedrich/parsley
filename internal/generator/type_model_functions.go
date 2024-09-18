package generator

import (
	"fmt"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"strings"
)

func RegisterTypeModelFunctions(generator GenericCodeGenerator) error {
	return generator.AddTemplateFunc(
		NamedFunc("FormatType", FormatType),
		NamedFunc("FormattedCallParameters", FormattedCallParameters),
		NamedFunc("FormattedParameters", FormattedParameters),
		NamedFunc("FormattedResultParameters", FormattedResultParameters),
		NamedFunc("FormattedResultTypes", FormattedResultTypes),
		NamedFunc("HasResults", HasResults),
		NamedFunc("Signature", Signature),
	)
}

func HasResults(m reflection.Method) bool {
	return len(m.Results) > 0
}

func FormattedParameters(m reflection.Method) string {
	formattedParameters := make([]string, len(m.Parameters))
	for i, parameter := range m.Parameters {
		typeName := FormatType(parameter)
		formattedParameters[i] = fmt.Sprintf("%s %s", parameter.Name, typeName)
	}
	return strings.Join(formattedParameters, ", ")
}

func FormattedCallParameters(m reflection.Method) string {
	formattedParameters := make([]string, len(m.Parameters))
	for i, parameter := range m.Parameters {
		formattedParameters[i] = fmt.Sprintf("%s", parameter.Name)
	}
	return strings.Join(formattedParameters, ", ")
}

func FormattedResultParameters(m reflection.Method) string {
	formattedResults := make([]string, len(m.Results))
	for i, result := range m.Results {
		formattedResults[i] = fmt.Sprintf("%s", result.Name)
	}
	return strings.Join(formattedResults, ", ")
}

func FormattedResultTypes(m reflection.Method) string {
	formattedResults := make([]string, len(m.Results))
	for i, result := range m.Results {
		typeName := FormatType(result)
		formattedResults[i] = fmt.Sprintf("%s", typeName)
	}
	if len(formattedResults) == 0 {
		return ""
	}
	return "(" + strings.Join(formattedResults, ", ") + ")"
}

func Signature(m reflection.Method) string {
	buffer := strings.Builder{}
	buffer.WriteString(fmt.Sprintf("%s", m.Name))
	buffer.WriteString(fmt.Sprintf("(%s)", FormattedParameters(m)))
	if len(m.Results) > 0 {
		buffer.WriteString(fmt.Sprintf(" %s", FormattedResultTypes(m)))
	}
	return buffer.String()
}

func FormatType(parameter reflection.Parameter) string {
	typeName := parameter.TypeName
	if parameter.IsPointer {
		typeName = "*" + typeName
	}
	if parameter.IsArray {
		typeName = "[]" + typeName
	}
	return typeName
}

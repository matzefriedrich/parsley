package generator

import (
	"fmt"
	"github.com/matzefriedrich/parsley/internal"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"strings"
)

// RegisterTypeModelFunctions registers a series of functions for type model processing in the given code generator.
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

// HasResults checks if the given reflection.Method has any result parameters.
func HasResults(m reflection.Method) bool {
	return len(m.Results) > 0
}

// FormattedParameters formats the parameters of the given reflection.Method into a comma-separated string with type info.
func FormattedParameters(m reflection.Method) string {
	formattedParameters := make([]string, len(m.Parameters))
	for i, parameter := range m.Parameters {
		typeName := FormatType(parameter)
		formattedParameters[i] = fmt.Sprintf("%s %s", parameter.Name, typeName)
	}
	return strings.Join(formattedParameters, ", ")
}

// FormattedCallParameters formats the call parameters of the given reflection.Method into a comma-separated string.
func FormattedCallParameters(m reflection.Method) string {
	formattedParameters := make([]string, len(m.Parameters))
	for i, parameter := range m.Parameters {
		name := parameter.Name
		if parameter.IsEllipsis() {
			name = fmt.Sprintf("%s...", name)
		}
		formattedParameters[i] = fmt.Sprintf("%s", name)
	}
	return strings.Join(formattedParameters, ", ")
}

// FormattedResultParameters formats the result parameters of the given reflection.Method into a comma-separated string.
func FormattedResultParameters(m reflection.Method) string {
	if m.Results == nil {
		return ""
	}
	formattedResults := make([]string, len(m.Results))
	for i, result := range m.Results {
		formattedResults[i] = fmt.Sprintf("%s", result.Name)
	}
	return strings.Join(formattedResults, ", ")
}

// FormattedResultTypes formats the result types of the given reflection.Method into a string representation.
func FormattedResultTypes(m reflection.Method) string {
	if m.Results == nil {
		return ""
	}
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

// Signature generates the signature string of a given method.
func Signature(m reflection.Method) string {
	buffer := strings.Builder{}
	buffer.WriteString(fmt.Sprintf("%s", m.Name))
	buffer.WriteString(fmt.Sprintf("(%s)", FormattedParameters(m)))
	if len(m.Results) > 0 {
		buffer.WriteString(fmt.Sprintf(" %s", FormattedResultTypes(m)))
	}
	return buffer.String()
}

const (
	ellipsis = "..."
	star     = "*"
	array    = "[]"
)

// FormatType formats the given reflection.Parameter's type information into a string representation.
func FormatType(parameter reflection.Parameter) string {

	if parameter.Type == nil {
		return "any"
	}

	segments := make([]string, 0)

	s := internal.MakeStack[*reflection.ParameterType]()
	s.Push(parameter.Type)

	for s.IsEmpty() == false {

		t := s.Pop()

		typeName := t.Name
		if len(t.SelectorName) > 0 {
			typeName = fmt.Sprintf("%s.%s", t.SelectorName, typeName)
		}

		if t.IsInterface {
			typeName = fmt.Sprintf("interface{}")
		}

		if t.IsEllipsis {
			typeName = ellipsis + typeName
		}

		if t.IsPointer {
			typeName = star + typeName
		}

		if t.IsArray {
			typeName = array + typeName
		}

		segments = append(segments, typeName)

		if t.Next != nil {
			s.Push(t.Next)
		}
	}

	return strings.Join(segments, "")
}

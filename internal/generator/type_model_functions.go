package generator

import (
	"fmt"
	"github.com/matzefriedrich/parsley/internal"
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

const (
	ellipsis = "..."
	star     = "*"
	array    = "[]"
)

func FormatType(parameter reflection.Parameter) string {

	segments := make([]string, 0)

	s := internal.MakeStack[*reflection.ParameterType]()
	s.Push(parameter.Type)

	for s.IsEmpty() == false {

		t := s.Pop()

		typeName := t.Name
		if len(t.SelectorName) > 0 {
			typeName = fmt.Sprintf("%s.%s", t.SelectorName, typeName)
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

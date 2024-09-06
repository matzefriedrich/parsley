package generator

type TemplateFunction struct {
	Name     string
	Function any
}

func NamedFunc(name string, f any) TemplateFunction {
	return TemplateFunction{
		Name:     name,
		Function: f,
	}
}

type TemplateLoader func(name string) (string, error)

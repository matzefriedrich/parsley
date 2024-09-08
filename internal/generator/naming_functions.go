package generator

import "unicode"

func RegisterNamingFunctions(generator GenericCodeGenerator) error {
	return generator.AddTemplateFunc(
		NamedFunc("asPrivate", MakePrivate),
		NamedFunc("asPublic", MakePublic))
}

func MakePrivate(s string) string {
	if s == "" {
		return ""
	}
	firstRune := []rune(s)[0]
	if unicode.IsUpper(firstRune) {
		s = string(unicode.ToLower(firstRune)) + s[1:]
	}
	return s
}

func MakePublic(s string) string {
	if s == "" {
		return ""
	}
	firstRune := []rune(s)[0]
	if unicode.IsLower(firstRune) {
		s = string(unicode.ToUpper(firstRune)) + s[1:]
	}
	return s
}

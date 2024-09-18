package generator

import (
	"errors"
	"fmt"
	"github.com/matzefriedrich/parsley/internal/reflection"
	"os"
	"path"
)

// GoFileAccessor Creates a new reflection.AstFileAccessor object that reads a source file from the GOFILE variable.
func GoFileAccessor() reflection.AstFileAccessor {

	goFilePath, err := GetGoFilePath()
	if err != nil {
		return func() (*reflection.AstFileSource, error) {
			return nil, newGeneratorError(ErrorFailedToObtainGeneratorSourceFile)
		}
	}

	return reflection.AstFromFile(goFilePath)
}

func GetGoFilePath() (string, error) {

	const goFileVariableName = "GOFILE"
	goFileName := os.Getenv(goFileVariableName)

	if goFileName == "" {
		return "", errors.New(fmt.Sprintf("%s environment variable not set", goFileVariableName))
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	goFilePath := path.Join(cwd, goFileName)
	if _, err = os.Stat(goFilePath); errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	return goFilePath, nil
}

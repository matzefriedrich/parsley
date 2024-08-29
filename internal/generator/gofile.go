package generator

import (
	"errors"
	"fmt"
	"os"
	"path"
)

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

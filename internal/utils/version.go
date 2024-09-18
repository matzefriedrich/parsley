package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/hashicorp/go-version"
)

const (
	VersionString string = "1.0.0"
)

type VersionInfo struct {
	Major int
	Minor int
	Patch int
}

func (v VersionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func (v VersionInfo) LessThan(other VersionInfo) bool {
	a, _ := version.NewVersion(v.String())
	b, _ := version.NewVersion(other.String())
	return a.LessThan(b)
}

func (v VersionInfo) Equal(other VersionInfo) bool {
	a, _ := version.NewVersion(v.String())
	b, _ := version.NewVersion(other.String())
	return a.Equal(b)
}

func ApplicationVersion() (*VersionInfo, error) {
	v, err := tryParseVersionInfo(VersionString)
	if err != nil {
		return nil, errors.New("application version not set")
	}
	return v, nil
}

func tryParseVersionInfo(version string) (*VersionInfo, error) {

	re := regexp.MustCompile("(?:[vV])?(?P<major>\\d+)\\.(?P<minor>\\d+)\\.(?P<patch>\\d+)")
	match := re.FindStringSubmatch(version)
	if match == nil {
		return nil, errors.New("invalid version")
	}

	extracted := map[string]string{}
	names := re.SubexpNames()
	for _, name := range names {
		index := re.SubexpIndex(name)
		if index != -1 && len(name) > 0 {
			extracted[name] = match[index]
		}
	}

	readInt := func(name string) int {
		value, found := extracted[name]
		if found {
			n, err := strconv.Atoi(value)
			if err == nil {
				return n
			}
		}
		return 0
	}

	major := readInt("major")
	minor := readInt("minor")
	patch := readInt("patch")

	return &VersionInfo{Major: major, Minor: minor, Patch: patch}, nil
}

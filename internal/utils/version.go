package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

const (
	VersionString string = "1.0.3"
)

// VersionInfo represents the version details using semantic versioning.
type VersionInfo struct {
	Major int
	Minor int
	Patch int
}

// ComparisonResult represents the result of comparing two values.
type ComparisonResult int

const (

	// LessThan indicates that the first value is less than the second in a comparison operation.
	LessThan ComparisonResult = iota

	// LessThanOrEqual indicates that the first value is either less than or equal to the second in a comparison operation.
	LessThanOrEqual

	// Equal indicates that the first value is equal to the second in a comparison operation.
	Equal

	// GreaterOrEqual indicates that the first value is either greater than or equal to the second in a comparison operation.
	GreaterOrEqual

	// GreaterThan indicates that the first value is greater than the second in a comparison operation.
	GreaterThan
)

// String returns the VersionInfo as a formatted string in the semantic versioning format (Major.Minor.Patch).
func (v VersionInfo) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// LessThan determines if the current VersionInfo is less than the specified VersionInfo based on semantic versioning.
func (v VersionInfo) LessThan(other VersionInfo) bool {
	return v.Compare(other) == LessThan
}

// LessThanOrEqual determines if the current VersionInfo is less than or equal to the specified VersionInfo based on semantic versioning.
func (v VersionInfo) LessThanOrEqual(other VersionInfo) bool {
	return v.Compare(other) == LessThanOrEqual
}

// Equal determines if two VersionInfo instances represent the same version based on semantic versioning.
func (v VersionInfo) Equal(other VersionInfo) bool {
	return v.Compare(other) == Equal
}

// Compare compares the current VersionInfo instance with another VersionInfo instance and returns a ComparisonResult.
func (v VersionInfo) Compare(other VersionInfo) ComparisonResult {

	if isLessThan(v, other) {
		return LessThan
	} else if isEqual(v, other) {
		return Equal
	} else if isGreaterThan(v, other) {
		return GreaterThan
	}

	// If it’s not less or equal, it must be one of the remaining.
	if isLessThanOrEqual(v, other) {
		return LessThanOrEqual
	} else {
		return GreaterOrEqual
	}
}

// ApplicationVersion parses and returns the application's version information. If the version is not set, an error is returned.
func ApplicationVersion() (*VersionInfo, error) {
	v, err := tryParseVersionInfo(VersionString)
	if err != nil {
		return nil, errors.New("application version not set")
	}
	return v, nil
}

func isLessThan(v VersionInfo, other VersionInfo) bool {
	return v.Major < other.Major ||
		(v.Major == other.Major && v.Minor < other.Minor) ||
		(v.Major == other.Major && v.Minor == other.Minor && v.Patch < other.Patch)
}

func isEqual(v VersionInfo, other VersionInfo) bool {
	return v.Major == other.Major && v.Minor == other.Minor && v.Patch == other.Patch
}

func isGreaterThan(v VersionInfo, other VersionInfo) bool {
	return v.Major > other.Major ||
		(v.Major == other.Major && v.Minor > other.Minor) ||
		(v.Major == other.Major && v.Minor == other.Minor && v.Patch > other.Patch)
}

func isLessThanOrEqual(v VersionInfo, other VersionInfo) bool {
	return v.Major < other.Major ||
		(v.Major == other.Major && v.Minor < other.Minor) ||
		(v.Major == other.Major && v.Minor == other.Minor && v.Patch <= other.Patch)
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

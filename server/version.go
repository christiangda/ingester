package server

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

var (
	// ErrBadStringFormat is returned when the string's version ([\"major\"].[minor].[patch]) has a bad format
	ErrBadStringFormat = errors.New("Version: Bad string format")

	// ErrBadStringFormatMajor is returned when the string's version has a bad format and major is a problem
	ErrBadStringFormatMajor = errors.New("Version: The major part of string's version ([\"major\"].[minor].[patch]) cannot be converted to int")

	// ErrBadStringFormatMinor is returned when the string's version has a bad format and minor is a problem
	ErrBadStringFormatMinor = errors.New("Version: The minor part of string's version ([major].[\"minor\"].[patch]) cannot be converted to int")

	// ErrBadStringFormatPatch is returned when the string's version has a bad format and patch is a problem
	ErrBadStringFormatPatch = errors.New("Version: The minor part of string's version ([major].[minor].[\"patch\"]) cannot be converted to int")
)

// Version according Semantic Versioning 2.0.0
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
	Patch int `json:"release"`
}

// NewVersion constructor from string
func NewVersion(version string) (*Version, error) {
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return nil, ErrBadStringFormat
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, ErrBadStringFormatMajor
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, ErrBadStringFormatMinor
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, ErrBadStringFormatPatch
	}

	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

// GetMajor return version
func (v *Version) GetMajor() int {
	return v.Major
}

// GetMinor return the minor version
func (v *Version) GetMinor() int {
	return v.Minor
}

// GetPatch return the patch version
func (v *Version) GetPatch() int {
	return v.Patch
}

// String this is Stringer implementation
// v := Version{1,0,0}
// fmt.println("The version is: ", v)
func (v *Version) String() string {
	version := ""

	if v != nil {
		major := strconv.Itoa(v.Major)
		minor := strconv.Itoa(v.Minor)
		patch := strconv.Itoa(v.Patch)

		version = major + "." + minor + "." + patch
	}

	return version
}

// ToJSON return version object as JSON notation
// v := Version{1,0,0}
// fmt.println("The version is: ", v.ToJSON())
func (v *Version) ToJSON() string {
	json, _ := json.Marshal(v)
	return string(json)
}

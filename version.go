package main


import (
"fmt"
"time"
)

const (
	internalIdentifier           = "f"
	internalBuildTimestamp int64 = 1599092654
	internalBuildNumber    int64 = 01
	internalVersionString        = "0.1.0"
)

func BuildDate() time.Time {
	return time.Unix(internalBuildTimestamp, 0)
}
func BuildNumber() int64 {
	return internalBuildNumber
}
func Version() string {
	return internalVersionString
}

func VersionInfo() string {
	return fmt.Sprintf("%s (v%s, build %d, build date:%v)", internalIdentifier, Version(), BuildNumber(), BuildDate())
}

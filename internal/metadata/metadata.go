package metadata

import (
	"encoding/json"
	"runtime"
)

var m *Metadata

type (
	Metadata struct {
		Rver    string  `json:"version"`
		Hash    string  `json:"commit_hash"`
		Gver    string  `json:"go_version"`
		Date    string  `json:"date"`
		Tag     string  `json:"tag"`
		Runtime Runtime `json:"runtime"`
	}

	Runtime struct {
		Arch string `json:"arch"`
		Goos string `json:"goos"`
	}
)

func init() {
	m = &Metadata{
		Gver: runtime.Version(),
		Runtime: Runtime{
			Arch: runtime.GOARCH,
			Goos: runtime.GOOS,
		},
	}
}

func Set(Version, CommitHash, Date string) {
	m.Date = Date
	m.Hash = CommitHash
	m.Rver = Version
}

func String() string {
	return m.ss()
}

func (m Metadata) ss() string {
	d, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(d)
}

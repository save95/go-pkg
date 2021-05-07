package types

import "github.com/hashicorp/go-version"

// ApiVersion 版本号
type ApiVersion string

func (av ApiVersion) Verify() bool {
	_, err := version.NewVersion(string(av))

	return err == nil
}

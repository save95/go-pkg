package http

import (
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
)

type UserAgent struct {
}

func (ua UserAgent) RandOfPC() string {
	uas := browser.Chrome()

	// 排除手机、Linux系统
	if strings.Contains(uas, "Mobile") || strings.Contains(uas, "Linux") {
		uas = ua.RandOfPC()
	}

	return uas
}

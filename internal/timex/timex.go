package timex

import (
	"time"

	"github.com/mrhpn/go-rest-api/internal/constants"
)

func ToAPIDateTimeFormat(datetime time.Time) string {
	if datetime.IsZero() {
		return ""
	}
	return datetime.UTC().Format(constants.APIDateTimeLayout)
}

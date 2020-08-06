package vo

import (
	"testing"
	"time"
)

func TestDateConvert(t *testing.T) {
	stamp := int64(1569570406)
	tm := time.Unix(stamp, 0)
	date := tm.Format("2006-01-02 15:04:05")
	t.Log(date)

}

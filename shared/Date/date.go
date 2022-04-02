package Date

import "time"

var Now = func() time.Time {
	return time.Now()
}

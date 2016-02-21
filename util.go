package collector

import (
	"time"

	google_protobuf "google/protobuf"
)

func GoogleTimestampToTime(gpt *google_protobuf.Timestamp) time.Time {
	return time.Unix(gpt.Seconds, int64(gpt.Nanos))
}

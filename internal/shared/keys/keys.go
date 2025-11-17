package keys

import "os"

func init() {
	KafkaGroupSuffix = os.Getenv("KAFKA_GROUP_SUFFIX")
	if KafkaGroupSuffix == "" {
		KafkaGroupSuffix = "dev"
	}
}

const (
	LinkCacheKey  = "shorturl:cache:shortlinks:short_code"
	HotLinksKey   = "shorturl:stat:hotlinks"
	GlobalRateKey = "shorturl:stat:global"
)

// kafka
var KafkaGroupSuffix string

const (
	ClickEventTopic  = "link.click"
	CreateEventTopic = "link.create"
)
const (
	ClickConsumerGroup = "analytics-click-processor"
)

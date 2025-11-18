package commonconfig

import "os"

var (
	GlobalPath  string
	GlobalFile  string
	ServicePath string
	ServiceFile string
)

func init() {
	GlobalPath = os.Getenv("GLOBAL_CONFIG_PATH")
	if GlobalPath == "" {
		GlobalPath = "../shared/commonconfig/"
	}
	GlobalFile = os.Getenv("GLOBAL_CONFIG_FILE")
	if GlobalFile == "" {
		GlobalFile = "global_dev"
	}
	ServicePath = os.Getenv("SERVICE_CONFIG_PATH")
	if ServicePath == "" {
		ServicePath = "./config/"
	}
	ServiceFile = os.Getenv("SERVICE_CONFIG_FILE")
	if ServiceFile == "" {
		ServiceFile = "config_dev"
	}
}

type RabbitMqConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	PoolSize     int    `mapstructure:"pool-size"`
	MinIdleConns int    `mapstructure:"min-idle-conns"`
}

type TenantMysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
	Options  string `mapstructure:"options"`
}

type LinkMysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
	Options  string `mapstructure:"options"`
}

type ClickHouseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
	Debug    bool   `mapstructure:"debug"`
}

type GinConfig struct {
}
type TraceConfig struct{}
type ApiGatewayConfig struct {
	Port              int    `mapstructure:"port"`
	ServerCrtFilePath string `mapstructure:"server-crt-file-path"`
	ServerKeyFilePath string `mapstructure:"server-key-file-path"`
}

type TenantServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type LinkServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type AnalyticsServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ServiceInfo struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Desc    string `mapstructure:"desc"`
}

type JaegerConfig struct {
	Host     string `mapstructure:"host"`
	GrpcPort int    `mapstructure:"grpc-port"`
}

type PrometheusConfig struct {
	Port int `mapstructure:"port"`
}

type RateLimiterConfig struct {
	Capacity int64   `mapstructure:"capacity"`
	Rate     float64 `mapstructure:"rate"`
}

type KafkaConfig struct {
	SeedBrokers    []string `mapstructure:"seed-brokers"`
	ConsumerGourp  string   `mapstructure:"consumer-gourp"`
	ConsumerTopics []string `mapstructure:"consumer-topics"`
}

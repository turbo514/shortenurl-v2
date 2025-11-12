package commonconfig

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

type AnalyticsClickHouseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"dbname"`
}

type ApiGatewayConfig struct {
	Port int `mapstructure:"port"`
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

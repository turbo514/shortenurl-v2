package commonconfig

type CommonConfig struct {
	Mq struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"rabbitmq"`
	Redis struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"redis"`
	Service struct {
		TenantName string `mapstructure:"tenant_name"`
		LinkName   string `mapstructure:"link_name"`
	} `mapstructure:"service"`
	Jaeger JaegerConfig `mapstructure:"jaeger"`
}

type ServiceInfo struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Desc    string `mapstructure:"desc"`
}

type JaegerConfig struct {
	Host     string `mapstructure:"host"`
	GrpcPort int    `mapstructure:"grpc_port"`
}

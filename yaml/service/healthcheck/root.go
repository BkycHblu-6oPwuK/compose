package healthcheck

type HealthCheck struct {
	Test        []string `yaml:"test,omitempty"`
	Interval    string   `yaml:"interval,omitempty"`
	Timeout     string   `yaml:"timeout,omitempty"`
	Retries     int      `yaml:"retries,omitempty"`
	StartPeriod string   `yaml:"start_period,omitempty"`
}

type HealthCheckBuilder struct {
	healthCheck HealthCheck
}

func NewHealthCheckBuilder() *HealthCheckBuilder {
	return &HealthCheckBuilder{
		healthCheck: HealthCheck{},
	}
}

func (b *HealthCheckBuilder) SetTest(test ...string) *HealthCheckBuilder {
	b.healthCheck.Test = test
	return b
}

func (b *HealthCheckBuilder) SetInterval(interval string) *HealthCheckBuilder {
	b.healthCheck.Interval = interval
	return b
}

func (b *HealthCheckBuilder) SetTimeout(timeout string) *HealthCheckBuilder {
	b.healthCheck.Timeout = timeout
	return b
}

func (b *HealthCheckBuilder) SetRetries(retries int) *HealthCheckBuilder {
	b.healthCheck.Retries = retries
	return b
}

func (b *HealthCheckBuilder) SetStartPeriod(startPeriod string) *HealthCheckBuilder {
	b.healthCheck.StartPeriod = startPeriod
	return b
}

func (b *HealthCheckBuilder) Build() HealthCheck {
	return b.healthCheck
}

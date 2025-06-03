package logging

type Logging struct {
	Driver  string            `yaml:"driver,omitempty"`
	Options map[string]string `yaml:"options,omitempty"`
}

type LoggingBuilder struct {
	logging Logging
}

func NewLoggingBuilder() *LoggingBuilder {
	return &LoggingBuilder{
		logging: Logging{
			Options: make(map[string]string),
		},
	}
}

func (b *LoggingBuilder) SetDriver(driver string) *LoggingBuilder {
	b.logging.Driver = driver
	return b
}

func (b *LoggingBuilder) SetOption(key, value string) *LoggingBuilder {
	b.logging.Options[key] = value
	return b
}

func (b *LoggingBuilder) SetOptions(options map[string]string) *LoggingBuilder {
	b.logging.Options = options
	return b
}

func (b *LoggingBuilder) Build() Logging {
	return b.logging
}

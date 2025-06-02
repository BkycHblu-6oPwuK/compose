package service

import (
	"docky/yaml/network"
	"docky/yaml/service/build"
	"docky/yaml/service/dependencies"
	"docky/yaml/service/healthcheck"
	"docky/yaml/service/logging"
)

type Service struct {
	Image           string                    `yaml:"image,omitempty"`
	Build           build.Build               `yaml:"build,omitempty"`
	Restart         string                    `yaml:"restart,omitempty"`
	Volumes         []string                  `yaml:"volumes,omitempty"`
	Ports           []string                  `yaml:"ports,omitempty"`
	Environment     map[string]string         `yaml:"environment,omitempty"`
	Dependencies    dependencies.Dependencies `yaml:"depends_on,omitempty"`
	Networks        []string                  `yaml:"networks,omitempty"`
	Command         interface{}               `yaml:"command,omitempty"`
	ExtraHosts      []string                  `yaml:"extra_hosts,omitempty"`
	Secrets         []string                  `yaml:"secrets,omitempty"`
	WorkingDir      string                    `yaml:"working_dir,omitempty"`
	User            string                    `yaml:"user,omitempty"`
	Entrypoint      interface{}               `yaml:"entrypoint,omitempty"`
	Labels          map[string]string         `yaml:"labels,omitempty"`
	Healthcheck     healthcheck.HealthCheck   `yaml:"healthcheck,omitempty"`
	StopGracePeriod string                    `yaml:"stop_grace_period,omitempty"`
	StopSignal      string                    `yaml:"stop_signal,omitempty"`
	Tmpfs           []string                  `yaml:"tmpfs,omitempty"`
	CapAdd          []string                  `yaml:"cap_add,omitempty"`
	CapDrop         []string                  `yaml:"cap_drop,omitempty"`
	Sysctls         map[string]string         `yaml:"sysctls,omitempty"`
	Ulimits         map[string]interface{}    `yaml:"ulimits,omitempty"`
	Privileged      bool                      `yaml:"privileged,omitempty"`
	ReadOnly        bool                      `yaml:"read_only,omitempty"`
	Logging         logging.Logging           `yaml:"logging,omitempty"`
	Ipc             string                    `yaml:"ipc,omitempty"`
	Pid             string                    `yaml:"pid,omitempty"`
	Hostname        string                    `yaml:"hostname,omitempty"`
	MacAddress      string                    `yaml:"mac_address,omitempty"`
	Expose          []string                  `yaml:"expose,omitempty"`
	Devices         []string                  `yaml:"devices,omitempty"`
	Init            bool                      `yaml:"init,omitempty"`
	Platform        string                    `yaml:"platform,omitempty"`
	Profiles        []string                  `yaml:"profiles,omitempty"`
	Runtime         string                    `yaml:"runtime,omitempty"`
	ContainerName   string                    `yaml:"container_name,omitempty"`
}

type ServiceBuilder struct {
	service             Service
	volumeSet           map[string]struct{}
	dependenciesBuilder *dependencies.DependenciesBuilder
	buildBuilder        *build.BuildBuilder
	healthBuilder       *healthcheck.HealthCheckBuilder
}

func NewServiceBuilder() *ServiceBuilder {
	return &ServiceBuilder{
		service: Service{
			Environment:  make(map[string]string),
			Volumes:      []string{},
			Ports:        []string{},
			Dependencies: dependencies.Dependencies{},
			Networks:     []string{},
			ExtraHosts:   []string{},
			Secrets:      []string{},
		},
	}
}

func NewServiceBuilderFrom(service Service) *ServiceBuilder {
	return &ServiceBuilder{
		service: service,
	}
}

func (b *ServiceBuilder) GetDependenciesBuilder() *dependencies.DependenciesBuilder {
	if b.dependenciesBuilder == nil {
		b.dependenciesBuilder = dependencies.NewDependenciesBuilderFrom(b.service.Dependencies)
	}
	return b.dependenciesBuilder
}

func (b *ServiceBuilder) SetImage(image string) *ServiceBuilder {
	b.service.Image = image
	return b
}

func (b *ServiceBuilder) SetRestart(restart string) *ServiceBuilder {
	b.service.Restart = restart
	return b
}

func (b *ServiceBuilder) SetRestartAlways() *ServiceBuilder {
	b.service.Restart = "always"
	return b
}

func (b *ServiceBuilder) AddVolume(volume string) *ServiceBuilder {
	b.service.Volumes = append(b.service.Volumes, volume)
	return b
}

func (b *ServiceBuilder) SetVolume(volume string) *ServiceBuilder {
	if b.volumeSet == nil {
		b.volumeSet = make(map[string]struct{})
		for _, vol := range b.service.Volumes {
			b.volumeSet[vol] = struct{}{}
		}
	}

	if _, exists := b.volumeSet[volume]; exists {
		return b
	}

	b.volumeSet[volume] = struct{}{}
	b.service.Volumes = append(b.service.Volumes, volume)

	return b
}

func (b *ServiceBuilder) AddPort(port string) *ServiceBuilder {
	b.service.Ports = append(b.service.Ports, port)
	return b
}

func (b *ServiceBuilder) AddEnvironment(key, value string) *ServiceBuilder {
	b.service.Environment[key] = value
	return b
}

func (b *ServiceBuilder) AddDefaultNetwork() *ServiceBuilder {
	return b.AddNetwork(network.DefaultName)
}

func (b *ServiceBuilder) AddNetwork(network string) *ServiceBuilder {
	b.service.Networks = append(b.service.Networks, network)
	return b
}

func (b *ServiceBuilder) SetCommand(command interface{}) *ServiceBuilder {
	b.service.Command = command
	return b
}

func (b *ServiceBuilder) SetCommandTailNull() *ServiceBuilder {
	b.service.Command = "tail -f /dev/null"
	return b
}

func (b *ServiceBuilder) AddExtraHost(host string) *ServiceBuilder {
	b.service.ExtraHosts = append(b.service.ExtraHosts, host)
	return b
}

func (b *ServiceBuilder) AddSecret(secret string) *ServiceBuilder {
	b.service.Secrets = append(b.service.Secrets, secret)
	return b
}

func (b *ServiceBuilder) SetContainerName(name string) *ServiceBuilder {
	b.service.ContainerName = name
	return b
}

func (b *ServiceBuilder) SetWorkingDir(value string) *ServiceBuilder {
	b.service.WorkingDir = value
	return b
}

func (b *ServiceBuilder) WithDependenciesBuilder(builder *dependencies.DependenciesBuilder) *ServiceBuilder {
	b.dependenciesBuilder = builder
	return b
}

func (b *ServiceBuilder) WithBuildBuilder(builder *build.BuildBuilder) *ServiceBuilder {
	b.buildBuilder = builder
	return b
}

func (b *ServiceBuilder) WithHealthcheckBuilder(builder *healthcheck.HealthCheckBuilder) *ServiceBuilder {
	b.healthBuilder = builder
	return b
}

func (b *ServiceBuilder) Build() Service {
	if b.dependenciesBuilder != nil {
		b.service.Dependencies = b.dependenciesBuilder.Build()
	}
	if b.buildBuilder != nil {
		b.service.Build = b.buildBuilder.Build()
	}
	if b.healthBuilder != nil {
		b.service.Healthcheck = b.healthBuilder.Build()
	}
	b.volumeSet = nil
	return b.service
}

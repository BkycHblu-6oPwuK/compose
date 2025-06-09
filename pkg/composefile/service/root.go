package service

import (
	"slices"

	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/network"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/service/build"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/service/dependencies"
	//"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/service/healthcheck"
	//"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/service/logging"
)

type Service struct {
	Image         string                    `yaml:"image,omitempty"`
	Build         build.Build               `yaml:"build,omitempty"`
	Restart       string                    `yaml:"restart,omitempty"`
	Volumes       []string                  `yaml:"volumes,omitempty"`
	Ports         []string                  `yaml:"ports,omitempty"`
	Environment   map[string]string         `yaml:"environment,omitempty"`
	Dependencies  dependencies.Dependencies `yaml:"depends_on,omitempty"`
	Networks      []string                  `yaml:"networks,omitempty"`
	Command       any                       `yaml:"command,omitempty"`
	ExtraHosts    []string                  `yaml:"extra_hosts,omitempty"`
	Secrets       []string                  `yaml:"secrets,omitempty"`
	WorkingDir    string                    `yaml:"working_dir,omitempty"`
	ContainerName string                    `yaml:"container_name,omitempty"`
	//User            string                    `yaml:"user,omitempty"`
	//Entrypoint      any                       `yaml:"entrypoint,omitempty"`
	//Labels          map[string]string         `yaml:"labels,omitempty"`
	//Healthcheck     healthcheck.HealthCheck   `yaml:"healthcheck,omitempty"`
	//StopGracePeriod string                    `yaml:"stop_grace_period,omitempty"`
	//StopSignal      string                    `yaml:"stop_signal,omitempty"`
	//Tmpfs           []string                  `yaml:"tmpfs,omitempty"`
	//CapAdd          []string                  `yaml:"cap_add,omitempty"`
	//CapDrop         []string                  `yaml:"cap_drop,omitempty"`
	//Sysctls         map[string]string         `yaml:"sysctls,omitempty"`
	//Ulimits         map[string]any            `yaml:"ulimits,omitempty"`
	//Privileged      bool                      `yaml:"privileged,omitempty"`
	//ReadOnly        bool                      `yaml:"read_only,omitempty"`
	//Logging         logging.Logging           `yaml:"logging,omitempty"`
	//Ipc             string                    `yaml:"ipc,omitempty"`
	//Pid             string                    `yaml:"pid,omitempty"`
	//Hostname        string                    `yaml:"hostname,omitempty"`
	//MacAddress      string                    `yaml:"mac_address,omitempty"`
	//Expose          []string                  `yaml:"expose,omitempty"`
	//Devices         []string                  `yaml:"devices,omitempty"`
	//Init            bool                      `yaml:"init,omitempty"`
	//Platform        string                    `yaml:"platform,omitempty"`
	//Profiles        []string                  `yaml:"profiles,omitempty"`
	//Runtime         string                    `yaml:"runtime,omitempty"`
	Extras map[string]any `yaml:",inline"`
}

type ServiceBuilder struct {
	service             Service
	volumeSet           map[string]struct{}
	dependenciesBuilder *dependencies.DependenciesBuilder
	buildBuilder        *build.BuildBuilder
	//healthBuilder       *healthcheck.HealthCheckBuilder
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

func (b *ServiceBuilder) GetService() *Service {
	return &b.service
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
	b.initVolumeSet()

	if _, exists := b.volumeSet[volume]; exists {
		return b
	}

	b.volumeSet[volume] = struct{}{}
	b.service.Volumes = append(b.service.Volumes, volume)

	return b
}

func (b *ServiceBuilder) RemoveVolume(volume string) *ServiceBuilder {
	b.initVolumeSet()

	if _, exists := b.volumeSet[volume]; !exists {
		return b
	}

	delete(b.volumeSet, volume)

	if i := slices.Index(b.service.Volumes, volume); i != -1 {
		b.service.Volumes = slices.Delete(b.service.Volumes, i, i+1)
	}

	return b
}

func (b *ServiceBuilder) FilterVolumes(filter func(string) bool) *ServiceBuilder {
	b.initVolumeSet()

	filtered := b.service.Volumes[:0]
	newVolumeSet := make(map[string]struct{})

	for _, volume := range b.service.Volumes {
		if filter(volume) {
			filtered = append(filtered, volume)
			newVolumeSet[volume] = struct{}{}
		}
	}

	b.service.Volumes = filtered
	b.volumeSet = newVolumeSet

	return b
}

func (b *ServiceBuilder) initVolumeSet() {
	if b.volumeSet == nil {
		b.volumeSet = make(map[string]struct{})
		for _, vol := range b.service.Volumes {
			b.volumeSet[vol] = struct{}{}
		}
	}
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

func (b *ServiceBuilder) SetCommand(command any) *ServiceBuilder {
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

// func (b *ServiceBuilder) WithHealthcheckBuilder(builder *healthcheck.HealthCheckBuilder) *ServiceBuilder {
// 	b.healthBuilder = builder
// 	return b
// }

func (b *ServiceBuilder) Build() Service {
	if b.dependenciesBuilder != nil {
		b.service.Dependencies = b.dependenciesBuilder.Build()
	}
	if b.buildBuilder != nil {
		b.service.Build = b.buildBuilder.Build()
	}
	// if b.healthBuilder != nil {
	// 	b.service.Healthcheck = b.healthBuilder.Build()
	// }
	b.volumeSet = nil
	return b.service
}

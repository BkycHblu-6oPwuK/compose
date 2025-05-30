package service

import (
	"docky/yaml/build"
	"docky/yaml/network"
	"sync"
)

type Service struct {
	Image         string            `yaml:"image,omitempty"`
	Build         build.Build       `yaml:"build,omitempty"`
	Restart       string            `yaml:"restart,omitempty"`
	Volumes       []string          `yaml:"volumes,omitempty"`
	Ports         []string          `yaml:"ports,omitempty"`
	Environment   map[string]string `yaml:"environment,omitempty"`
	Dependencies  []string          `yaml:"depends_on,omitempty"`
	Networks      []string          `yaml:"networks,omitempty"`
	Command       interface{}       `yaml:"command,omitempty"`
	ExtraHosts    []string          `yaml:"extra_hosts,omitempty"`
	Secrets       []string          `yaml:"secrets,omitempty"`
	ContainerName string            `yaml:"container_name,omitempty"`
}

type ServiceBuilder struct {
	service   Service
	volumeSet map[string]struct{}
	once      sync.Once
}

func NewServiceBuilder() *ServiceBuilder {
	return &ServiceBuilder{
		service: Service{
			Environment:  make(map[string]string),
			Volumes:      []string{},
			Ports:        []string{},
			Dependencies: []string{},
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

func (b *ServiceBuilder) SetImage(image string) *ServiceBuilder {
	b.service.Image = image
	return b
}

func (b *ServiceBuilder) SetBuild(build build.Build) *ServiceBuilder {
	b.service.Build = build
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
	b.once.Do(func() {
		b.volumeSet = make(map[string]struct{})
		for _, vol := range b.service.Volumes {
			b.volumeSet[vol] = struct{}{}
		}
	})

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

func (b *ServiceBuilder) AddDependency(serviceName string) *ServiceBuilder {
	b.service.Dependencies = append(b.service.Dependencies, serviceName)
	return b
}

func (b *ServiceBuilder) SetDependency(key int, value string) *ServiceBuilder {
	b.service.Dependencies[key] = value
	return b
}

func (b *ServiceBuilder) RewriteServiceDependency(search, newValue string) *ServiceBuilder {
	for i, dep := range b.service.Dependencies {
		if dep == search {
			return b.SetDependency(i, newValue)
		}
	}
	return b.AddDependency(newValue)
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

func (b *ServiceBuilder) Build() Service {
	return b.service
}

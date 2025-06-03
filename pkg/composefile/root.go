package composefile

import (
	"docky/pkg/orderedmap"
	"docky/pkg/composefile/network"
	"docky/pkg/composefile/secret"
	"docky/pkg/composefile/service"
	"docky/pkg/composefile/volume"
	"fmt"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type ComposeFile struct {
	Services *orderedmap.OrderedMap[string, service.Service] `yaml:"services"`
	Volumes  map[string]volume.Volume                        `yaml:"volumes,omitempty"`
	Networks map[string]network.Network                      `yaml:"networks,omitempty"`
	Secrets  map[string]secret.Secret                        `yaml:"secrets,omitempty"`
}

type ComposeFileBuilder struct {
	file ComposeFile
}

var (
	once        sync.Once
	currentYaml *ComposeFile
	loadErr     error
)

func NewComposeFileBuilder() *ComposeFileBuilder {
	return &ComposeFileBuilder{
		file: ComposeFile{
			Services: orderedmap.NewOrderedMap[string, service.Service](),
			Volumes:  make(map[string]volume.Volume),
			Networks: make(map[string]network.Network),
			Secrets:  make(map[string]secret.Secret),
		},
	}
}

func NewComposeFileBuilderFrom(compose ComposeFile) *ComposeFileBuilder {
	return &ComposeFileBuilder{
		file: compose,
	}
}

func (b *ComposeFileBuilder) AddService(name string, service service.Service) *ComposeFileBuilder {
	b.file.Services.Set(name, service)
	return b
}

func (b *ComposeFileBuilder) AddNetwork(name string, network network.Network) *ComposeFileBuilder {
	b.file.Networks[name] = network
	return b
}

func (b *ComposeFileBuilder) AddDefaultNetwork() *ComposeFileBuilder {
	defaultNetwork := network.NewNetworkBuilder().
		SetBridgeDriver().
		Build()
	b.file.Networks[network.DefaultName] = defaultNetwork
	return b
}

func (b *ComposeFileBuilder) AddVolume(name string, volume volume.Volume) *ComposeFileBuilder {
	b.file.Volumes[name] = volume
	return b
}

func (b *ComposeFileBuilder) AddSecret(name string, secret secret.Secret) *ComposeFileBuilder {
	b.file.Secrets[name] = secret
	return b
}

func (b *ComposeFileBuilder) HasService(name string) bool {
	return b.file.Services.Has(name)
}

func (b *ComposeFileBuilder) GetService(name string) (service service.Service, exists bool) {
	return b.file.Services.Get(name)
}

func (b *ComposeFileBuilder) RemoveService(name string) *ComposeFileBuilder {
	b.file.Services.Delete(name)
	return b
}

func (b *ComposeFileBuilder) RemoveVolume(name string) *ComposeFileBuilder {
	delete(b.file.Volumes, name)
	return b
}

func (b *ComposeFileBuilder) Build() ComposeFile {
	return b.file
}

func (c *ComposeFile) Save(path string) error {
	out, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0644)
}

func Load(path string) (*ComposeFile, error) {
	once.Do(func() {
		data, err := os.ReadFile(path)
		if err != nil {
			loadErr = fmt.Errorf("ошибка чтения файла: %w", err)
			return
		}

		compose := &ComposeFile{
			Services: orderedmap.NewOrderedMap[string, service.Service](),
			Volumes:  make(map[string]volume.Volume),
			Networks: make(map[string]network.Network),
			Secrets:  make(map[string]secret.Secret),
		}
		err = yaml.Unmarshal(data, compose)
		if err != nil {
			loadErr = fmt.Errorf("ошибка парсинга yaml: %w", err)
			return
		}

		currentYaml = compose
	})

	return currentYaml, loadErr
}

func SetCurrentYaml(c *ComposeFile) {
	currentYaml = c
}

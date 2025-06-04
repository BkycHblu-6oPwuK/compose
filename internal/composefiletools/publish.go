package composefiletools

import (
	"fmt"

	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/service"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/volume"
)

func publishWithBuilder(modifier func(builder *composefile.ComposeFileBuilder) error) error {
	path := config.GetDockerComposeFilePath()
	compose, err := composefile.Load(path)
	if err != nil {
		return err
	}

	builder := composefile.NewComposeFileBuilderFrom(*compose)
	if err := modifier(builder); err != nil {
		return err
	}

	final := builder.Build()
	composefile.SetCurrentYaml(&final)
	return final.Save(path)
}

func PublishMysqlService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(Mysql) {
			appService, exists := b.GetService(App)
			if exists {
				serviceBuilder := service.NewServiceBuilderFrom(appService)
				serviceBuilder.GetDependenciesBuilder().RewriteDependency(Postgres, Mysql)
				b.AddService(App, serviceBuilder.Build())
			}
			b.AddService(Mysql, buildMysqlService()).
				AddVolume(Mysql_data, volume.Volume{})
		}
		if b.HasService(Postgres) {
			b.RemoveService(Postgres).
				RemoveVolume(Postgres_data)
		}
		return nil
	})
}

func PublishPostgresService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(Postgres) {
			appService, exists := b.GetService(App)
			if exists {
				serviceBuilder := service.NewServiceBuilderFrom(appService)
				serviceBuilder.GetDependenciesBuilder().RewriteDependency(Mysql, Postgres)
				b.AddService(App, serviceBuilder.Build())
			}
			b.AddService(Postgres, buildPostgresService()).
				AddVolume(Postgres_data, volume.Volume{})
		}
		if b.HasService(Mysql) {
			b.RemoveService(Mysql).
				RemoveVolume(Mysql_data)
		}
		return nil
	})
}

func PublishNodeService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(Node) {
			b.AddService(Node, buildNodeService())
		}
		return nil
	})
}

func PublishSphinxService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(Sphinx) {
			b.AddService(Sphinx, buildSphinxService()).
				AddVolume(Sphinx_data, volume.Volume{})
		}
		return nil
	})
}

func PublishRedisService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(Redis) {
			b.AddService(Redis, buildRedisService()).
				AddVolume(Redis_data, volume.Volume{})
		}
		return nil
	})
}

func PublishMemcachedService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(Memcached) {
			b.AddService(Memcached, buildMemcachedService())
		}
		return nil
	})
}

func PublishMailhogService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(Mailhog) {
			b.AddService(Mailhog, buildMailHogService())
		}
		return nil
	})
}

func PublishPhpMyAdminService() error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(Mysql) {
			return fmt.Errorf("phpmyadmin работает только с mysql. В docker-compose не найден сервис %s", Mysql)
		}
		if !b.HasService(PhpMyAdmin) {
			b.AddService(PhpMyAdmin, buildPhpMyAdminService())
		}
		return nil
	})
}

func PublishVolumes(serviceNames []string, volumes map[string][]string, modifier func(s *service.Service) (isContinue bool, err error)) error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		for _, serviceName := range serviceNames {
			if curService, exists := b.GetService(serviceName); exists {
				if modifier != nil {
					isContinue, err := modifier(&curService)
					if err != nil {
						return fmt.Errorf("ошибка при модификации сервиса %s: %w", serviceName, err)
					}
					if !isContinue {
						continue
					}
				}
				serviceBuilder := service.NewServiceBuilderFrom(curService)
				if vols, ok := volumes[serviceName]; ok {
					for _, vol := range vols {
						serviceBuilder.SetVolume(vol)
					}
				}
				b.AddService(serviceName, serviceBuilder.Build())
			}
		}
		return nil
	})
}

func PublishDockerfile(serviceName, dockerfile string) error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if curService, exists := b.GetService(serviceName); exists {
			if curService.Build.Dockerfile != "" {
				curService.Build.Dockerfile = dockerfile
				b.AddService(serviceName, curService)
			}
		} else {
			return fmt.Errorf("сервис %s не найден", serviceName)
		}
		return nil
	})
}

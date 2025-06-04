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

func publishDatabaseService(target string, builderFunc func() service.Service) error {
	alternatives := []string{Mysql, Mariadb, Postgres}
	removeVolumes := map[string]string{
		Mysql:    Mysql_data,
		Mariadb:  Mariadb_data,
		Postgres: Postgres_data,
	}

	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		if !b.HasService(target) {
			if appService, exists := b.GetService(App); exists {
				serviceBuilder := service.NewServiceBuilderFrom(appService)
				deps := serviceBuilder.GetDependenciesBuilder()
				for _, alt := range alternatives {
					if alt != target && b.HasService(alt) {
						deps.RewriteDependency(alt, target)
					}
				}
				b.AddService(App, serviceBuilder.Build())
			}
			b.AddService(target, builderFunc()).
				AddVolume(removeVolumes[target], volume.Volume{})
		}

		for _, alt := range alternatives {
			if alt != target && b.HasService(alt) {
				b.RemoveService(alt).
					RemoveVolume(removeVolumes[alt])
			}
		}
		return nil
	})
}

func PublishMysqlService() error {
	return publishDatabaseService(Mysql, buildMysqlService)
}

func PublishMariadbService() error {
	return publishDatabaseService(Mariadb, buildMariadbService)
}

func PublishPostgresService() error {
	return publishDatabaseService(Postgres, buildPostgresService)
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
		publish := func (host string) {
			if !b.HasService(PhpMyAdmin) {
				b.AddService(PhpMyAdmin, buildPhpMyAdminService(host))
			}
		}
		if b.HasService(Mysql) {
			publish(Mysql)
			return nil
		} else if b.HasService(Mariadb) {
			publish(Mariadb)
			return nil
		}
		return fmt.Errorf("phpmyadmin работает только с mysql. В docker-compose не найден сервис %s или %s", Mysql, Mariadb)
	})
}

// volumes map serviceName>>[]string volumes
func PublishVolumes(volumes map[string][]string, modifier func(b *service.ServiceBuilder) (isContinue bool, err error)) error {
	return publishWithBuilder(func(b *composefile.ComposeFileBuilder) error {
		for serviceName, volumes := range volumes {
			if curService, exists := b.GetService(serviceName); exists {
				serviceBuilder := service.NewServiceBuilderFrom(curService)
				if modifier != nil {
					isContinue, err := modifier(serviceBuilder)
					if err != nil {
						return fmt.Errorf("ошибка при модификации сервиса %s: %w", serviceName, err)
					}
					if !isContinue {
						continue
					}
				}
				for _, vol := range volumes {
					serviceBuilder.SetVolume(vol)
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

package helper

import (
	"docky/yaml"
	"docky/yaml/service"
	"docky/yaml/volume"
	"fmt"
)

func publishWithBuilder(modifier func(builder *yaml.ComposeFileBuilder) error) error {
	compose, err := yaml.Load()
	if err != nil {
		return err
	}

	builder := yaml.NewComposeFileBuilderFrom(*compose)
	if err := modifier(builder); err != nil {
		return err
	}

	final := builder.Build()
	yaml.SetCurrentYaml(&final)
	return final.Save()
}

func PublishMysqlService() error {
	return publishWithBuilder(func(b *yaml.ComposeFileBuilder) error {
		if !b.HasService(Mysql) {
			appService, exists := b.GetService(App)
			if exists {
				serviceBuilder := service.NewServiceBuilderFrom(appService).RewriteServiceDependency(Postgres, Mysql)
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
	return publishWithBuilder(func(b *yaml.ComposeFileBuilder) error {
		if !b.HasService(Postgres) {
			appService, exists := b.GetService(App)
			if exists {
				serviceBuilder := service.NewServiceBuilderFrom(appService).RewriteServiceDependency(Mysql, Postgres)
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
	return publishWithBuilder(func(b *yaml.ComposeFileBuilder) error {
		if !b.HasService(Node) {
			b.AddService(Node, buildNodeService())
		}
		return nil
	})
}

func PublishSphinxService() error {
	return publishWithBuilder(func(b *yaml.ComposeFileBuilder) error {
		if !b.HasService(Sphinx) {
			b.AddService(Sphinx, buildSphinxService()).
				AddVolume(Sphinx_data, volume.Volume{})
		}
		return nil
	})
}

func PublishRedisService() error {
	return publishWithBuilder(func(b *yaml.ComposeFileBuilder) error {
		if !b.HasService(Redis) {
			b.AddService(Redis, buildRedisService()).
				AddVolume(Redis_data, volume.Volume{})
		}
		return nil
	})
}

func PublishMemcachedService() error {
	return publishWithBuilder(func(b *yaml.ComposeFileBuilder) error {
		if !b.HasService(Memcached) {
			b.AddService(Memcached, buildMemcachedService())
		}
		return nil
	})
}

func PublishMailhogService() error {
	return publishWithBuilder(func(b *yaml.ComposeFileBuilder) error {
		if !b.HasService(Mailhog) {
			b.AddService(Mailhog, buildMailHogService())
		}
		return nil
	})
}

func PublishPhpMyAdminService() error {
	return publishWithBuilder(func(b *yaml.ComposeFileBuilder) error {
		if !b.HasService(Mysql) {
			return fmt.Errorf("phpmyadmin работает только с mysql. В docker-compose не найден сервис %s", Mysql)
		}
		if !b.HasService(PhpMyAdmin) {
			b.AddService(PhpMyAdmin, buildPhpMyAdminService())
		}
		return nil
	})
}

func PublishVolumes(serviceNames []string, volumes map[string][]string) error {
	return publishWithBuilder(func(b *yaml.ComposeFileBuilder) error {
		for _, serviceName := range serviceNames {
			if curService, exists := b.GetService(serviceName); exists {
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
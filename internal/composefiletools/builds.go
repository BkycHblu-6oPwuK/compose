package composefiletools

import (
	"github.com/BkycHblu-6oPwuK/docky/v2/internal/config"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/service"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/service/build"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/service/dependencies"
	"github.com/BkycHblu-6oPwuK/docky/v2/pkg/composefile/volume"
)

func getBaseBuildBuilder(dockerfile string, args map[string]string) *build.BuildBuilder {
	builder := build.NewBuildBuilder().
		SetContext(GetVarNameString(config.DockerPathVarName)).
		SetDockerfile(dockerfile).
		AddArg(config.UserGroupVarName, GetVarNameString(config.UserGroupVarName))
	for key, value := range args {
		builder.AddArg(key, value)
	}
	return builder
}

func getSimpleDependecyBulder(serviceName string) *dependencies.DependenciesBuilder {
	return dependencies.NewDependenciesBuilder().AddSimple(serviceName)
}

func buildBaseVolume() volume.Volume {
	return volume.NewVolumeBuilder().Build()
}

func buildNginxService() service.Service {
	nginxService := service.NewServiceBuilder().
		WithBuildBuilder(getBaseBuildBuilder(GetNginxComposePath(Dockerfile, false), nil)).
		AddVolume(GetSiteVolumePath()).
		AddPort("80:80").
		AddPort("443:443").
		WithDependenciesBuilder(getSimpleDependecyBulder(App)).
		AddDefaultNetwork().
		SetContainerName(Nginx).
		Build()
	return nginxService
}

func buildAppService(yamlConfig *config.YamlConfig) service.Service {
	appServiceBuilder := service.NewServiceBuilder().
		WithBuildBuilder(getBaseBuildBuilder(GetPhpConfComposePath(Dockerfile, false), nil)).
		AddVolume(GetSiteVolumePath()).
		AddPort("9000:9000").
		AddPort("6001:6001").
		AddExtraHost("host.docker.internal:host-gateway").
		AddDefaultNetwork().
		AddEnvironment("XDEBUG_TRIGGER", "testTrig").
		AddEnvironment("PHP_IDE_CONFIG", "serverName=xdebugServer").
		SetContainerName(App)

	if yamlConfig.DbType != Sqlite {
		appServiceBuilder.WithDependenciesBuilder(getSimpleDependecyBulder(yamlConfig.DbType))
	}
	return appServiceBuilder.Build()
}

func buildMysqlService() service.Service {
	mysqlService := service.NewServiceBuilder().
		SetImage(Mysql+":"+GetVarNameString(config.MysqlVersionVarName)).
		SetRestartAlways().
		AddPort("8102:3306").
		AddVolume(Mysql_data+":/var/lib/mysql").
		AddVolume(GetMysqlCnfPath(false)).
		AddDefaultNetwork().
		AddEnvironment("MYSQL_DATABASE", "site").
		AddEnvironment("MYSQL_ROOT_PASSWORD", "root").
		SetContainerName(Mysql).
		Build()
	return mysqlService
}

func buildMariadbService() service.Service {
	mysqlService := service.NewServiceBuilder().
		SetImage(Mariadb+":"+GetVarNameString(config.MariadbVersionVarName)).
		SetRestartAlways().
		AddPort("8102:3306").
		AddVolume(Mariadb_data+":/var/lib/mysql").
		AddVolume(GetMysqlCnfPath(false)).
		AddDefaultNetwork().
		AddEnvironment("MARIADB_DATABASE", "site").
		AddEnvironment("MARIADB_ROOT_PASSWORD", "root").
		SetContainerName(Mariadb).
		Build()
	return mysqlService
}

func buildPostgresService() service.Service {
	postgresService := service.NewServiceBuilder().
		SetImage(Postgres+":"+GetVarNameString(config.PostgresVersionVarName)).
		SetRestartAlways().
		AddPort("5432:5432").
		AddVolume(Postgres_data+":/var/lib/postgresql/data").
		AddVolume(GetPostgresConfPath(false)).
		AddDefaultNetwork().
		AddEnvironment("POSTGRES_DB", "site").
		AddEnvironment("POSTGRES_PASSWORD", "root").
		AddEnvironment("POSTGRES_USER", "root").
		SetContainerName(Postgres).
		SetCommand([]string{"-c", "config_file=/etc/postgresql/postgresql.conf"}).
		Build()
	return postgresService
}

func buildNodeService() service.Service {
	nodeService := service.NewServiceBuilder().
		WithBuildBuilder(getBaseBuildBuilder(GetVarNameString(config.DockerPathVarName)+"/"+Node+"/"+Dockerfile, map[string]string{
			"NODE_VERSION": GetVarNameString(config.NodeVersionVarName),
		})).
		AddPort("5173:5173").
		AddPort("5174:5174").
		AddVolume(GetSiteVolumePath()).
		WithDependenciesBuilder(getSimpleDependecyBulder(App)).
		AddDefaultNetwork().
		SetCommandTailNull().
		SetWorkingDir(GetVarNameString(config.NodePathVarName)).
		SetContainerName(Node).
		Build()
	return nodeService
}

func buildSphinxService() service.Service {
	sphinxService := service.NewServiceBuilder().
		WithBuildBuilder(getBaseBuildBuilder(GetVarNameString(config.DockerPathVarName)+"/"+Sphinx+"/"+Dockerfile, nil)).
		SetRestartAlways().
		AddPort("9312:9312").
		AddPort("9306:9306").
		AddVolume(Sphinx_data + ":/var/lib/sphinx/data").
		AddDefaultNetwork().
		SetContainerName(Sphinx).
		Build()
	return sphinxService
}

func buildMemcachedService() service.Service {
	memcachedService := service.NewServiceBuilder().
		SetImage(Memcached).
		AddPort("11211:11211").
		AddDefaultNetwork().
		SetContainerName(Memcached).
		SetCommand([]string{"--conn-limit=1024", "--memory-limit=64", "--threads=4"}).
		Build()
	return memcachedService
}

func buildRedisService() service.Service {
	redisService := service.NewServiceBuilder().
		SetImage(Redis).
		AddPort("6379:6379").
		AddVolume(Redis_data + ":/data").
		AddDefaultNetwork().
		SetContainerName(Redis).
		SetCommand([]string{"redis-server", "--appendonly", "yes"}).
		Build()
	return redisService
}

func buildMailHogService() service.Service {
	mailHogService := service.NewServiceBuilder().
		SetImage("mailhog/mailhog").
		AddPort("1025:1025").
		AddPort("8025:8025").
		AddDefaultNetwork().
		SetContainerName(Mailhog).
		Build()
	return mailHogService
}

func buildPhpMyAdminService(host string) service.Service {
	phpMyAdminService := service.NewServiceBuilder().
		SetImage("phpmyadmin/phpmyadmin").
		SetRestartAlways().
		AddPort("8080:80").
		AddEnvironment("PMA_HOST", host).
		AddEnvironment("PMA_PORT", "3306").
		AddEnvironment("PMA_USER", "root").
		AddEnvironment("PMA_PASSWORD", "root").
		WithDependenciesBuilder(getSimpleDependecyBulder(host)).
		AddDefaultNetwork().
		SetContainerName(PhpMyAdmin).
		Build()
	return phpMyAdminService
}

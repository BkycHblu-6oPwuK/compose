package helper

import (
	"docky/config"
	"docky/yaml/build"
	"docky/yaml/service"
	"docky/yaml/volume"
)

func buildBaseBuild(dockerfile string, args map[string]string) build.Build {
	build := build.NewBuildBuilder().
		SetContextDefault().
		SetDockerfile(dockerfile).
		SetBaseArgs()
	for key, value := range args {
		build.AddArg(key, value)
	}
	return build.Build()
}

func buildBaseVolume() volume.Volume {
	return volume.NewVolumeBuilder().Build()
}

func buildNginxService() service.Service {
	nginxService := service.NewServiceBuilder().
		SetBuild(buildBaseBuild("${"+config.DockerPathVarName+"}/"+Nginx+"/"+Dockerfile+"", nil)).
		AddVolume("${" + config.SitePathVarName + "}:" + config.SitePathInContainer).
		AddPort("80:80").
		AddPort("443:443").
		AddDependency(App).
		AddDefaultNetwork().
		SetContainerName(Nginx).
		Build()
	return nginxService
}

func buildAppService(yamlConfig *config.YamlConfig) service.Service {
	phpVersionVarName := "${" + config.PhpVersionVarName + "}"
	appService := service.NewServiceBuilder().
		SetBuild(buildBaseBuild("${"+config.DockerPathVarName+"}/"+App+"/php-"+phpVersionVarName+"/"+Dockerfile, nil)).
		AddVolume("${"+config.SitePathVarName+"}:"+config.SitePathInContainer).
		AddPort("9000:9000").
		AddPort("6001:6001").
		AddExtraHost("host.docker.internal:host-gateway").
		AddDefaultNetwork().
		AddEnvironment("XDEBUG_TRIGGER", "testTrig").
		AddEnvironment("PHP_IDE_CONFIG", "serverName=xdebugServer").
		SetContainerName(App).
		Build()
	if yamlConfig.DbType != Sqlite {
		appService.Dependencies = []string{yamlConfig.DbType}
	}
	return appService
}

func buildMysqlService() service.Service {
	mysqlService := service.NewServiceBuilder().
		SetImage(Mysql+":${"+config.MysqlVersionVarName+"}").
		SetRestartAlways().
		AddPort("8102:3306").
		AddVolume(Mysql_data+":/var/lib/mysql").
		AddVolume("${"+config.DockerPathVarName+"}/"+Mysql+"/my.cnf:/etc/mysql/conf.d/my.cnf").
		AddDefaultNetwork().
		AddEnvironment("MYSQL_DATABASE", "site").
		AddEnvironment("MYSQL_ROOT_PASSWORD", "root").
		SetContainerName(Mysql).
		Build()
	return mysqlService
}

func buildPostgresService() service.Service {
	postgresService := service.NewServiceBuilder().
		SetImage(Postgres+":${"+config.PostgresVersionVarName+"}").
		SetRestartAlways().
		AddPort("5432:5432").
		AddVolume(Postgres_data+":/var/lib/postgresql/data").
		AddVolume("${"+config.DockerPathVarName+"}/"+Postgres+"/postgresql.conf:/etc/postgresql/postgresql.conf").
		AddDefaultNetwork().
		AddEnvironment("POSTGRES_DB", "site").
		AddEnvironment("POSTGRES_PASSWORD", "root").
		AddEnvironment("POSTGRES_USER", Postgres).
		SetContainerName(Postgres).
		SetCommand([]string{"-c", "config_file=/etc/postgresql/postgresql.conf"}).
		Build()
	return postgresService
}

func buildNodeService() service.Service {
	nodeService := service.NewServiceBuilder().
		SetBuild(buildBaseBuild("${"+config.DockerPathVarName+"}/"+Node+"/"+Dockerfile, map[string]string{
			"NODE_VERSION": "${" + config.NodeVersionVarName + "}",
			"NODE_PATH":    "${" + config.NodePathVarName + "}",
		})).
		AddPort("5173:5173").
		AddPort("5174:5174").
		AddVolume("${" + config.SitePathVarName + "}:" + config.SitePathInContainer).
		AddDependency(App).
		AddDefaultNetwork().
		SetCommandTailNull().
		SetContainerName(Node).
		Build()
	return nodeService
}

func buildSphinxService() service.Service {
	sphinxService := service.NewServiceBuilder().
		SetBuild(buildBaseBuild("${"+config.DockerPathVarName+"}/"+Sphinx+"/"+Dockerfile, nil)).
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

func buildPhpMyAdminService() service.Service {
	phpMyAdminService := service.NewServiceBuilder().
		SetImage("phpmyadmin/phpmyadmin").
		SetRestartAlways().
		AddPort("8080:80").
		AddEnvironment("PMA_HOST", Mysql).
		AddEnvironment("PMA_PORT", "3306").
		AddEnvironment("PMA_USER", "root").
		AddEnvironment("PMA_PASSWORD", "root").
		AddDependency(Mysql).
		AddDefaultNetwork().
		SetContainerName(PhpMyAdmin).
		Build()
	return phpMyAdminService
}

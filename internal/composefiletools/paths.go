package composefiletools

import "github.com/BkycHblu-6oPwuK/docky/v2/internal/config"

func GetNginxConfPathInContainer() string {
	return "/etc/nginx/" + ConfDir
}

func GetNginxComposePath(file string, isConfPath bool) string {
	return getBaseConfComposePath(isConfPath) + "/" + Nginx + "/" + file
}

func GetNginxConfVolumePath(toEndJoin string) string {
	hostPath := GetNginxComposePath(ConfDir, true)
	containerPath := GetNginxConfPathInContainer()

	if toEndJoin != "" {
		hostPath += "/" + toEndJoin
		containerPath += "/" + toEndJoin
	}

	return hostPath + ":" + containerPath
}

func GetAppNginxConfVolumePath(toEndJoin string) string {
	hostPath := getBaseConfComposePath(true) + "/" + App + "/" + Nginx
	containerPath := GetNginxConfPathInContainer()

	if toEndJoin != "" {
		hostPath += "/" + toEndJoin
		containerPath += "/" + toEndJoin
	}

	return hostPath + ":" + containerPath
}

func GetNginxSnippetsConfVolumePath(toEndJoin string) string {
	if toEndJoin != "" {
		toEndJoin = "/" + toEndJoin
	}
	return GetNginxConfVolumePath("snippets" + toEndJoin)
}

func GetCertificateConfVolumePath(toEndJoin string) string {
	hostPath := getBaseConfComposePath(true) + "/" + Nginx + "/certs"
	containerPath := "/usr/local/share/ca-certificates"

	if toEndJoin != "" {
		hostPath += "/" + toEndJoin
		containerPath += "/" + toEndJoin
	}

	return hostPath + ":" + containerPath
}

func GetsymlinksConfVolumePath() string {
	return getBaseConfComposePath(true) + "/" + App + "/symlinks:/usr/symlinks_extra"
}

func GetPhpConfComposePath(file string, isConfPath bool) string {
	return getBaseConfComposePath(isConfPath) + "/" + App + "/php-${" + config.PhpVersionVarName + "}/" + file
}

func GetPhpConfVolumePath(file string, isConfPath bool) string {
	return GetPhpConfComposePath(file, isConfPath) + ":/usr/local/etc/php/conf.d/" + file
}

func GetCronConfVolumePath() string {
	return getBaseConfComposePath(true) + "/" + App + "/cron:/var/spool/cron/crontabs"
}

func GetSiteVolumePath() string {
	return GetVarNameString(config.SitePathVarName) + ":" + SitePathInContainer
}

func GetMysqlCnfPath(isConfPath bool) string {
	return getBaseConfComposePath(isConfPath)+"/"+Mysql+"/my.cnf:/etc/mysql/conf.d/my.cnf"
}

func GetPostgresConfPath(isConfPath bool) string {
	return getBaseConfComposePath(isConfPath)+"/"+Postgres+"/postgresql.conf:/etc/postgresql/postgresql.conf"
}

func GetSupervisordConfPath() string {
	return getBaseConfComposePath(true)+"/"+App+"/supervisord.conf:/etc/supervisor/conf.d/supervisord.conf"
}

func getBaseConfComposePath(isConfPath bool) string {
	if isConfPath {
		return GetVarNameString(config.ConfPathVarName)
	} else {
		return GetVarNameString(config.DockerPathVarName)
	}
}

func GetVarNameString(varName string) string {
	return "${" + varName + "}"
}

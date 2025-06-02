package build

import "docky/config"

type Build struct {
	Context    string            `yaml:"context,omitempty"`
	Dockerfile string            `yaml:"dockerfile,omitempty"`
	Args       map[string]string `yaml:"args,omitempty"`
}

type BuildBuilder struct {
	build Build
}

func NewBuildBuilder() *BuildBuilder {
	return &BuildBuilder{
		build: Build{},
	}
}

func (v *BuildBuilder) SetContext(context string) *BuildBuilder {
	v.build.Context = context
	return v
}
func (v *BuildBuilder) SetContextDefault() *BuildBuilder {
	v.build.Context = "${"+config.DockerPathVarName+"}"
	return v
}
func (v *BuildBuilder) SetDockerfile(dockerfile string) *BuildBuilder {
	v.build.Dockerfile = dockerfile
	return v
}
func (v *BuildBuilder) SetBaseArgs() *BuildBuilder {
	v.build.Args = map[string]string{
		config.UserGroupVarName: "${" + config.UserGroupVarName + "}",
	}
	return v
}
func (v *BuildBuilder) AddArg(name, value string) *BuildBuilder {
	v.build.Args[name] = value
	return v
}

func (v *BuildBuilder) Build() Build {
	return v.build
}

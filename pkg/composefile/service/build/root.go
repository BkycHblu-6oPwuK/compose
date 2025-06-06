package build

type Build struct {
	Context    string            `yaml:"context,omitempty"`
	Dockerfile string            `yaml:"dockerfile,omitempty"`
	Args       map[string]string `yaml:"args,omitempty"`
	Extras     map[string]any    `yaml:",inline"`
}

type BuildBuilder struct {
	build Build
}

func NewBuildBuilder() *BuildBuilder {
	return &BuildBuilder{
		build: Build{
			Context:    "",
			Dockerfile: "",
			Args:       make(map[string]string),
		},
	}
}

func (v *BuildBuilder) SetContext(context string) *BuildBuilder {
	v.build.Context = context
	return v
}
func (v *BuildBuilder) SetDockerfile(dockerfile string) *BuildBuilder {
	v.build.Dockerfile = dockerfile
	return v
}
func (v *BuildBuilder) AddArg(name, value string) *BuildBuilder {
	v.build.Args[name] = value
	return v
}

func (v *BuildBuilder) Build() Build {
	return v.build
}

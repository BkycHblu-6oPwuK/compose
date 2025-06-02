package dependencies

import "fmt"

type Dependencies struct {
	Simple   []string
	Advanced map[string]DependenciesCondition
}

type DependenciesCondition struct {
	Condition string `yaml:"condition"`
}

func (d *Dependencies) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var simple []string
	if err := unmarshal(&simple); err == nil {
		d.Simple = simple
		return nil
	}

	var advanced map[string]DependenciesCondition
	if err := unmarshal(&advanced); err == nil {
		d.Advanced = advanced
		return nil
	}

	return fmt.Errorf("depends_on: unsupported format")
}

func (d Dependencies) MarshalYAML() (interface{}, error) {
	if len(d.Advanced) > 0 {
		return d.Advanced, nil
	}
	if len(d.Simple) > 0 {
		return d.Simple, nil
	}
	return nil, nil
}

type DependenciesBuilder struct {
	dependencies Dependencies
}

func NewDependenciesBuilder() *DependenciesBuilder {
	return &DependenciesBuilder{
		dependencies: Dependencies{
			Simple:   []string{},
			Advanced: map[string]DependenciesCondition{},
		},
	}
}

func NewDependenciesBuilderFrom(dependencies Dependencies) *DependenciesBuilder {
	return &DependenciesBuilder{
		dependencies: dependencies,
	}
}

func (b *DependenciesBuilder) AddSimple(service string) *DependenciesBuilder {
	b.dependencies.Simple = append(b.dependencies.Simple, service)
	return b
}

func (b *DependenciesBuilder) SetSimple(services ...string) *DependenciesBuilder {
	b.dependencies.Simple = services
	return b
}

func (b *DependenciesBuilder) AddAdvanced(service string, condition string) *DependenciesBuilder {
	if b.dependencies.Advanced == nil {
		b.dependencies.Advanced = make(map[string]DependenciesCondition)
	}
	b.dependencies.Advanced[service] = DependenciesCondition{
		Condition: condition,
	}
	return b
}

func (b *DependenciesBuilder) RewriteDependency(search, newValue string) *DependenciesBuilder {
	if len(b.dependencies.Simple) > 0 {
		for i, dep := range b.dependencies.Simple {
			if dep == search {
				b.dependencies.Simple[i] = newValue
				return b
			}
		}
		return b.AddSimple(newValue)
	}

	if b.dependencies.Advanced == nil {
		b.dependencies.Advanced = make(map[string]DependenciesCondition)
	}

	if cond, exists := b.dependencies.Advanced[search]; exists {
		delete(b.dependencies.Advanced, search)
		if cond.Condition == "" {
			cond.Condition = "sevice_healthy"
		}
		b.dependencies.Advanced[newValue] = cond
	} else {
		b.dependencies.Advanced[newValue] = DependenciesCondition{
			Condition: "service_healthy",
		}
	}
	return b
}


func (b *DependenciesBuilder) Build() Dependencies {
	return b.dependencies
}

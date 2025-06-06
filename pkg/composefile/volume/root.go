package volume

type Volume struct {
	Driver string         `yaml:"driver,omitempty"`
	Extras map[string]any `yaml:",inline"`
}

type VolumeBuilder struct {
	volume Volume
}

func NewVolumeBuilder() *VolumeBuilder {
	return &VolumeBuilder{
		volume: Volume{},
	}
}

func (v *VolumeBuilder) SetDriver(driver string) *VolumeBuilder {
	v.volume.Driver = driver
	return v
}

func (v *VolumeBuilder) Build() Volume {
	return v.volume
}

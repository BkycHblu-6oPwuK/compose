package network

type Network struct {
	Driver string `yaml:"driver,omitempty"`
}

const (
	DefaultName string = "docky"
)

type NetworkBuilder struct {
	network Network
}

func NewNetworkBuilder() *NetworkBuilder  {
	return &NetworkBuilder{
		network: Network{},
	}
}

func (v *NetworkBuilder) setDriver(driver string) *NetworkBuilder {
	v.network.Driver = driver
	return v
}

func (v *NetworkBuilder) SetBridgeDriver() *NetworkBuilder {
	return v.setDriver("bridge")
}

func (v *NetworkBuilder) Build() Network {
	return v.network
}
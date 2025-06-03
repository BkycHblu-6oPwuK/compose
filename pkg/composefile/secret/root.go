package secret

type Secret struct {
	File string `yaml:"file,omitempty"`
}

type SecretBuilder struct {
	secret Secret
}

func NewSecretBuilder() *SecretBuilder  {
	return &SecretBuilder{
		secret: Secret{},
	}
}

func (v *SecretBuilder) SetFile(file string) *SecretBuilder {
	v.secret.File = file
	return v
}

func (v *SecretBuilder) Build() Secret {
	return v.secret
}
package envcrypto

type Option interface {
	applyToBox(*Box) error
}

type OptionFunc func(*Box) error

func (f OptionFunc) applyToBox(b *Box) error {
	return f(b)
}

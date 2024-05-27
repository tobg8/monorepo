package secret

type text interface {
	~[]byte | ~string
}

func conceal[T text](s T) string {
	if len(s) == 0 {
		return ""
	}
	return "*****"
}

// String allows to avoid displaying secret string values in logs for instance.
type String string

// String implements Stringer.
// Always returns "*****" to avoid leaking secrets.
func (s String) String() string {
	return conceal(s)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (s String) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

// Bytes allows to avoid displaying secret bytes values in logs for instance.
type Bytes []byte

// String implements Stringer.
// Always returns "*****" to avoid leaking secrets.
func (b Bytes) String() string {
	return conceal(b)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (b Bytes) MarshalText() (text []byte, err error) {
	return []byte(b.String()), nil
}

package protoreflecthash

import (
	"crypto/sha256"
	"fmt"
	"math"
)

const (
	// Sorted alphabetically by value.
	boolIdentifier     = `b`
	mapIdentifier      = `d`
	floatIdentifier    = `f`
	intIdentifier      = `i`
	listIdentifier     = `l`
	nilIdentifier      = `n`
	byteIdentifier     = `r`
	unicodeIndentifier = `u`
)

func hashBool(b bool) ([]byte, error) {
	bb := []byte(`0`)
	if b {
		bb = []byte(`1`)
	}
	return hash(boolIdentifier, bb)
}

func hashBytes(bs []byte) ([]byte, error) {
	return hash(byteIdentifier, bs)
}

func hashFloat(f float64) ([]byte, error) {
	var normalizedFloat string

	switch {
	case math.IsInf(f, 1):
		normalizedFloat = "Infinity"
	case math.IsInf(f, -1):
		normalizedFloat = "-Infinity"
	case math.IsNaN(f):
		normalizedFloat = "NaN"
	default:
		var err error
		normalizedFloat, err = floatNormalize(f)
		if err != nil {
			return nil, err
		}
	}

	return hash(floatIdentifier, []byte(normalizedFloat))
}

func hashInt64(i int64) ([]byte, error) {
	return hash(intIdentifier, []byte(fmt.Sprintf("%d", i)))
}

func hashNil() ([]byte, error) {
	return hash(nilIdentifier, []byte(``))
}

func hashUnicode(s string) ([]byte, error) {
	return hash(unicodeIndentifier, []byte(s))
}

func hash(t string, b []byte) ([]byte, error) {
	h := sha256.New()

	if _, err := h.Write([]byte(t)); err != nil {
		return nil, err
	}

	if _, err := h.Write(b); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func floatNormalize(originalFloat float64) (string, error) {
	// Special case 0
	// Note that if we allowed f to end up > .5 or == 0, we'd get the same thing.
	if originalFloat == 0 {
		return "+0:", nil
	}

	// Sign
	f := originalFloat
	s := `+`
	if f < 0 {
		s = `-`
		f = -f
	}
	// Exponent
	e := 0
	for f > 1 {
		f /= 2
		e++
	}
	for f <= .5 {
		f *= 2
		e--
	}
	s += fmt.Sprintf("%d:", e)
	// Mantissa
	if f > 1 || f <= .5 {
		return "", fmt.Errorf("could not normalize float: %f", originalFloat)
	}
	for f != 0 {
		if f >= 1 {
			s += `1`
			f--
		} else {
			s += `0`
		}
		if f >= 1 {
			return "", fmt.Errorf("could not normalize float: %f", originalFloat)
		}
		if len(s) >= 1000 {
			return "", fmt.Errorf("could not normalize float: %f", originalFloat)
		}
		f *= 2
	}
	return s, nil
}

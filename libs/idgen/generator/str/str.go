package str

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const numeric = "1234567890"
const alpha = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const alphaNumeric = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const lowerAlpha = "abcdefghijklmnopqrstuvwxyz"
const symbolic = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz()?!.#%@^=+-&*`~;:<>{}"

func NextSymbolic(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(symbolic))))
		if err != nil {
			return "", err
		}
		ret[i] = symbolic[num.Int64()]
	}

	return string(ret), nil
}

func NextAlphanumeric(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphaNumeric))))
		if err != nil {
			return "", err
		}
		ret[i] = alphaNumeric[num.Int64()]
	}

	return string(ret), nil
}

func NextLowerAlphabetic(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(lowerAlpha))))
		if err != nil {
			return "", err
		}
		ret[i] = lowerAlpha[num.Int64()]
	}

	return string(ret), nil
}

func NextNumeric(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(numeric))))
		if err != nil {
			return "", err
		}
		ret[i] = numeric[num.Int64()]
	}

	return string(ret), nil
}

func NextAlphabetic(n int) (string, error) {
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alpha))))
		if err != nil {
			return "", err
		}
		ret[i] = alpha[num.Int64()]
	}

	return string(ret), nil
}

func NextRandomInt32(min, max int32) (int32, error) {
	if min > max {
		return 0, fmt.Errorf("min should be less than or equal to max")
	}

	rangeSize := int64(max - min + 1)
	if rangeSize <= 0 {
		return 0, fmt.Errorf("range of possible values is zero or negative")
	}

	num, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
	if err != nil {
		return 0, err
	}

	return min + int32(num.Int64()), nil
}

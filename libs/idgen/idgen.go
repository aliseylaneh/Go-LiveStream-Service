package idgen

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"strings"

	"safir/libs/idgen/generator/sortable"
	"safir/libs/idgen/generator/str"
	"safir/libs/idgen/generator/uuid"
)

func NextBytes(len int) ([]byte, error) {
	b := make([]byte, len)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func NextUint32() (uint32, error) {
	b, err := NextBytes(4)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(b), nil
}

func NextUint64() (uint64, error) {
	b, err := NextBytes(8)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(b), nil
}

func NextHexString(len int) (string, error) {
	b, err := NextBytes(len)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func NextAlphabeticString(len int) (string, error) {
	return str.NextAlphabetic(len)
}

func NextLowerAlphabeticString(len int) (string, error) {
	return str.NextLowerAlphabetic(len)
}

func NextAlphanumericString(len int) (string, error) {
	return str.NextAlphanumeric(len)
}

func NextNumericString(len int) (string, error) {
	return str.NextNumeric(len)
}
func NextNumericint32(min, max int32) (int32, error) {
	return str.NextRandomInt32(min, max)
}

func NextSymbolicString(len int) (string, error) {
	return str.NextSymbolic(len)
}

func NextUUID() (string, error) {
	return uuid.NextUUID()
}

// A 40 character hex string globally unique token. Based on secure UUID4
func NextUniqueToken() (string, error) {
	uuid, err := NextUUID()
	if err != nil {
		return "", err
	}
	hex, err := NextHexString(8)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(uuid, "-", "") + hex, nil
}

// A unique token sortable alphabetically
func NextSortableUniqueToken() (string, error) {
	return sortable.NextXid()
}

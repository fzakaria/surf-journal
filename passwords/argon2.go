package passwords

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"strings"
)

type kdfType string

// https://www.rfc-editor.org/rfc/rfc9106.html#name-parameter-choice
// 2), low-memory options

const argon2idKDF kdfType = "argon2id"
const argon2idSaltLen = 16
const argon2idKeyLen = 32
const argon2idTime = 3
const argon2idMemory = 65536 // 2 GiB
const argon2idThreads = 4

// https://github.com/P-H-C/phc-string-format/blob/master/phc-sf-spec.md
func makePHCPrefix() string {
	version := fmt.Sprintf("v=%d", argon2.Version)
	parameters := fmt.Sprintf("m=%d,t=%d,p=%d", argon2idMemory, argon2idTime, argon2idThreads)
	return fmt.Sprintf("$%s$%s$%s$", argon2idKDF, version, parameters)
}

var phcPrefix = makePHCPrefix()

func serialize(salt, hash []byte) string {
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	return fmt.Sprintf("%s%s$%s", phcPrefix, encodedSalt, encodedHash)
}

func deserialize(serialized string) (salt []byte, hash []byte, err error) {
	saltAndHash := strings.TrimPrefix(serialized, phcPrefix)
	if saltAndHash == serialized {
		err = fmt.Errorf("serialized password does not match PHC spec; missing expected prefix %s", phcPrefix)
		return
	}

	parts := strings.Split(saltAndHash, "$")
	if len(parts) != 2 {
		err = fmt.Errorf("serialized password does not match PHC spec; not salt and hash")
		return
	}

	salt, err = base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		err = fmt.Errorf("serialized password has invalid salt: %w", err)
		return
	}
	if len(salt) != 16 {
		err = fmt.Errorf("serialized password's salt length invalid: %v", salt)
		return
	}

	hash, err = base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		err = fmt.Errorf("serialized password has invalid password hash: %w", err)
		return
	}
	if len(hash) != 32 {
		err = fmt.Errorf("serialized password's hash length invalid: %s", hash)
		return
	}

	return
}

func argon2id(salt, password []byte) []byte {
	return argon2.IDKey(password, salt, argon2idTime, argon2idMemory, argon2idThreads, argon2idKeyLen)
}

// CheckPassword checks the serialized password and salt against the
// user's password.
func CheckPassword(serialized, password string) error {
	salt, _, err := deserialize(serialized)
	if err != nil {
		return fmt.Errorf("saved password invalid: %w", err)
	}

	hash := argon2id(salt, []byte(password))
	provided := serialize(salt, hash)

	if subtle.ConstantTimeCompare([]byte(serialized), []byte(provided)) == 0 {
		return fmt.Errorf("passwords do not match")
	}

	return nil
}

// NewSerializedPassword returns a serialized password
func NewSerializedPassword(password string) (string, error) {
	salt := make([]byte, argon2idSaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("password serialization error: could not generate salt %w", err)
	}

	hash := argon2id(salt, []byte(password))

	serialized := serialize(salt, hash)

	return serialized, nil
}

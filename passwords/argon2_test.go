package passwords

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func MustHexDecode(b string) []byte {
	r, err := hex.DecodeString(b)
	if err != nil {
		panic(err)
	}
	return r
}

var password = []byte("hunter2")
var salt = MustHexDecode("4820ad36422cec25b06e66aae89ec313")

func TestSerialize(t *testing.T) {
	hash := argon2id(salt, password)

	serialized := serialize(salt, hash)
	// from argon2_cffi with the same parameters, salt, and password
	expected := "$argon2id$v=19$m=65536,t=3,p=4$SCCtNkIs7CWwbmaq6J7DEw$8Tb2ZejKQnsPeMl9ME4B4ZCiHCHHB+hvZHzsDs2XxF4"

	if expected != serialized {
		t.Fatalf("%s != %s", expected, serialized)
	}
}

func TestRoundTrip(t *testing.T) {
	hash := argon2id(salt, password)

	serialized := serialize(salt, hash)
	deserializedSalt, deserializedHash, err := deserialize(serialized)
	if err != nil {
		t.Fatalf("unexpected error: %+v", err)
	}

	if !bytes.Equal(deserializedSalt, salt) {
		t.Fatalf("deserialized salt %+v != salt %+v", deserializedSalt, salt)
	}

	if !bytes.Equal(deserializedHash, hash) {
		t.Fatalf("deserialized hash %+v != hash %+v", deserializedHash, hash)
	}
}

func FuzzInvalidPrefix(f *testing.F) {
	cases := []string{
		"",
		"jsdklfs",
		"$argon2id$v=19$m=1,t=1,p=1$SCCtNkIs7CWwbmaq6J7DEw$8Tb2ZejKQnsPeMl9ME4B4ZCiHCHHB$8Tb2ZejKQnsPeMl9ME4B4ZCiHCHHB+hvZHzsDs2XxF4",
	}

	for _, c := range cases {
		f.Add(c)
	}

	f.Fuzz(func(t *testing.T, withInvalidPrefix string) {
		_, _, err := deserialize(withInvalidPrefix)
		if err == nil {
			t.Errorf("expected error with invalid prefix: %+v", withInvalidPrefix)
		}
	})
}

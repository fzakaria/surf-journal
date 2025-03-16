package passwords

import "testing"

func FuzzInvalidSaltAndHash(f *testing.F) {
	cases := []string{
		phcPrefix,
		phcPrefix + "$not base 64",
		phcPrefix + "$not base 64$not base 64 either",
	}

	for _, c := range cases {
		f.Add(c)
	}

	f.Fuzz(func(t *testing.T, invalidHashAndSalt string) {
		_, _, err := deserialize(invalidHashAndSalt)
		if err == nil {
			t.Errorf("expected error and with invalid salt and hash: %+v", invalidHashAndSalt)
		}
	})
}

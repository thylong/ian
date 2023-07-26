package share

import (
	"testing"
)

func TestEncryptDecryptFile(t *testing.T) {
	cases := []struct {
		Text         []byte
		Secret       []byte
		ExpectingErr bool
	}{
		{[]byte("a top secret config"), []byte("a very very very very secret key"), false},
		{[]byte(""), []byte("a very very very very secret key"), false},
		{[]byte("a top secret config"), []byte(""), true},
		{[]byte("a flop secret config"), []byte("a not so secret key"), true},
	}

	for _, tc := range cases {
		encrypted, _ := EncryptFile(tc.Text, tc.Secret)
		expectedDecrypted, err := DecryptFile(encrypted, tc.Secret)

		if (err != nil && !tc.ExpectingErr) || (err == nil && tc.ExpectingErr) {
			t.Errorf("Expected nil error, got %#v", err)
		}
		if err == nil && string(expectedDecrypted) != string(tc.Text) {
			t.Errorf("Decrypted text is not the same as original Text, got %s wants %s", expectedDecrypted, tc.Text)
		}
	}
}

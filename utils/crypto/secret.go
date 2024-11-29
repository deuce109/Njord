package crypto

import (
	"encoding/base64"
	"io"
)

type Random struct {
	Reader io.Reader
}

func (r *Random) GetRandomSecret(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := r.Reader.Read(bytes)
	if err != nil {
		return "", err
	} else {
		return base64.RawStdEncoding.EncodeToString(bytes), nil
	}
}

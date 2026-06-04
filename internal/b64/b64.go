package b64

import "encoding/base64"

func Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func Decode(input string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

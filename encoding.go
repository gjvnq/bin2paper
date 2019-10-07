package bin2paper

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"io"
)

func ToB64(src []byte) (dst []byte) {
	base64.StdEncoding.Encode(dst, src)
	return
}

func ToB64str(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func ToHex(src []byte) string {
	return hex.EncodeToString(src)
}

func ReadLine(input io.Reader) (string, error) {
	buf := make([]byte, 1)
	ans := new(bytes.Buffer)
	for {
		n, err := input.Read(buf)
		if n == 0 || buf[0] == '\n' {
			break
		} else if err != nil {
			return "", err
		}
		err = ans.WriteByte(buf[0])
		if err != nil {
			return "", err
		}
	}
	return ans.String(), nil
}

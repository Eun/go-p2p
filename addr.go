package p2p

import (
	"encoding/base32"
	"strings"
	"github.com/vmihailenco/msgpack"
	"crypto/rand"
)

type Addr struct {
	Version     byte
	Identifier  string
	Port        uint16
	Secret      []byte
	FingerPrint []byte
}

func (i *Addr) ParseString(s string) error {
	data, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.Replace(s, "-", "", -1))
	if err != nil {
		return err
	}
	return msgpack.Unmarshal(data, i)
}

func (i *Addr) String() string {
	b, err := msgpack.Marshal(i)
	if err != nil {
		panic(err)
	}

	s := []rune(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b))
	var sb strings.Builder
	last := len(s) - 1
	for i := 0; i < len(s); i++ {
		if i%6 == 0 && i != last && i > 0 {
			sb.WriteRune('-')
		}
		sb.WriteRune(s[i])
	}
	return sb.String()
}

func (i *Addr) Network() string {
return "go-p2p"
}

func generateSecret(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	return b, err
}

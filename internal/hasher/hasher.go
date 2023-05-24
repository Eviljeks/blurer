package hasher

import (
	"crypto/sha1"
	"fmt"
)

type Hasher interface {
	Hash(bytes []byte) string
}

type Sha1Hasher struct {
}

func NewSha1Hasher() *Sha1Hasher {
	return &Sha1Hasher{}
}

func (sh *Sha1Hasher) Hash(bytes []byte) string {
	return fmt.Sprintf("%x", sha1.Sum(bytes))
}

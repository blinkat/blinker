package concat

import (
	"crypto"
	"encoding/binary"
	"hash"
	"io"
)

type Concat struct {
	z, info []byte
	i       uint32
	cache   []byte
	hasher  hash.Hash
}

func NewConcat(hash crypto.Hash, z, alg, pt_u, pt_v, sup_pub, sup_pri []byte) io.Reader {
	buf := make([]byte, len(alg)+len(pt_u)+len(pt_v)+len(sup_pub)+len(sup_pri))
	n := 0
	n += copy(buf, alg)
	n += copy(buf[n:], pt_u)
	n += copy(buf[n:], pt_v)
	n += copy(buf[n:], sup_pub)
	copy(buf[n:], sup_pri)

	hasher := hash.New()

	return &Concat{
		z:      z,
		info:   buf,
		hasher: hasher,
		i:      1,
	}
}

func (c *Concat) Read(out []byte) (int, error) {
	copid := copy(out, c.cache)
	c.cache = c.cache[copid:]

	for copid < len(out) {
		c.hasher.Reset()

		binary.Write(c.hasher, binary.BigEndian, c.i)
		c.hasher.Write(c.z)
		c.hasher.Write(c.info)

		h := c.hasher.Sum(nil)
		c_copied := copy(out[copid:], h)
		copid += c_copied
		c.cache = h[c_copied:]
		c.i += 1
	}
	return copid, nil
}

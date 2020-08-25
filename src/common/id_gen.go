package common

type IdGenerator struct {
	seed    uint64
	maskBit uint32
	index   uint64
}

func (g *IdGenerator) GenId() uint64 {
	id := g.index<<g.maskBit + g.seed
	g.index++
	return id
}

func (g *IdGenerator) Reset() {
	g.index = 0
}

func NewIdGenerator(seed uint64, maskBit uint32) *IdGenerator {
	return &IdGenerator{
		seed:    seed,
		maskBit: maskBit,
		index:   0,
	}
}

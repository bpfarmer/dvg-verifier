package merkle

// Proof comment
type Proof struct {
	r1 []byte
	r2 []byte
}

// IncProof comment
type IncProof struct {
	n []byte
	p []*HashDir
}

// HashDir comment
type HashDir struct {
	h []byte
	l bool
}

package valueobject

type SellerReference string

func NewSellerReference(reference string) SellerReference {
	// some validation for seller reference
	return SellerReference(reference)
}

func (p SellerReference) String() string {
	return string(p)
}

func (i SellerReference) Equals(item SellerReference) bool {
	return i == item
}

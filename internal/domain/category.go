package domain

const (
	ProposalCategory = "proposal"
)

type Subcategory struct {
	ID         uint64
	CategoryID uint64
	Name       string
}

type Category struct {
	ID            uint64
	Name          string
	Subcategories []Subcategory
}

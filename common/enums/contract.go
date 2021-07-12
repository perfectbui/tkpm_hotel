package enums

type ContractIndex int

const (
	_ ContractIndex = iota
	Completed
	Processing
	Cancel
)

func (r ContractIndex) String() string {
	return [...]string{"Completed", "Processing", "Cancel"}[r-1]
}

package enums

type StatusIndex int

const (
	_ StatusIndex = iota
	Available
	Booked
)

func (r StatusIndex) String() string {
	return [...]string{"Available", "Booked"}[r-1]
}

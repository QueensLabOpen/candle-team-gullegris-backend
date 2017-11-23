package utils

type Store struct {
	Games [][]int
}

func NewStore () *Store {
	return &Store{
		Games: [][]int{},
	}
}
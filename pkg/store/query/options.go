package query

type FindOneOptions struct {
	Sort SortOption
}

type FindOptions struct {
	Sort SortOption
}

type SortOption map[string]int

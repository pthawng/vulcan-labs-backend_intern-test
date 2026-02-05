package repository

// CodeRepository abstracts access to promotion code data sources.
type CodeRepository interface {
	Exists(code string) (bool, error)
	LoadAll() (map[string]struct{}, error)
}

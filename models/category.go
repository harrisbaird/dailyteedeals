package models

type Category struct {
	ID          int
	ProductID   int
	Name        string
	Slug        string
	Tags        []string `pg:",array" sql:",notnull"`
	IgnoredTags []string `pg:",array" sql:",notnull"`

	Product *Product
}

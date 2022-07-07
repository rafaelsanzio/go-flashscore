package mongo

type Document interface {
	GetID() string
	SetID(string)
}

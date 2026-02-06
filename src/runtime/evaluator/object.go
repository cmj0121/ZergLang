package evaluator

// ObjectType represents the type of a runtime object.
type ObjectType string

const (
	INTEGER_OBJ ObjectType = "INTEGER"
	STRING_OBJ  ObjectType = "STRING"
	BOOLEAN_OBJ ObjectType = "BOOLEAN"
	NULL_OBJ    ObjectType = "NULL"
)

// Object is the interface for all runtime values.
type Object interface {
	Type() ObjectType
	Inspect() string
}

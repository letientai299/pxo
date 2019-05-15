package internal

import (
	"strings"

	"github.com/knq/snaker"
	"github.com/xo/xo/models"
)

// TemplateType represents a template type.
type TemplateType uint

// the order here will be the alter the output order per file.
const (
	EnumTemplate TemplateType = iota
	ProcTemplate
	TypeTemplate
	ForeignKeyTemplate
	IndexTemplate
	QueryTypeTemplate
	QueryTemplate

	// always last
	XOTemplate
)

// String returns the name for the associated template type.
func (tt TemplateType) String() string {
	var s string
	switch tt {
	case XOTemplate:
		s = "xo_db"
	case EnumTemplate:
		s = "enum"
	case ProcTemplate:
		s = "proc"
	case TypeTemplate:
		s = "type"
	case ForeignKeyTemplate:
		s = "foreignkey"
	case IndexTemplate:
		s = "index"
	case QueryTypeTemplate:
		s = "querytype"
	case QueryTemplate:
		s = "query"
	default:
		panic("unknown TemplateType")
	}
	return s
}

// RelType represents the different types of relational storage (table/view).
type RelType uint

const (
	// Table reltype
	Table RelType = iota

	// View reltype
	View
)

// EscType represents the different escape types.
type EscType uint

const (
	SchemaEsc = iota
	TableEsc
	ColumnEsc
)

// String provides the string representation of RelType.
func (rt RelType) String() string {
	var s string
	switch rt {
	case Table:
		s = "TABLE"
	case View:
		s = "VIEW"
	default:
		panic("unknown RelType")
	}
	return s
}

// EnumValue holds data for a single enum value.
type Schema struct {
	Name    string
	Package string
	Comment string
}

func (s Schema) PackageName() string {
	p := strings.TrimSpace(s.Package)
	if len(p) != 0 {
		return p
	}

	// Any sane db developer should name their db schema in 2 way:
	// CamelCaseDB or snake_case_db, which is well supported by the snaker lib.
	// If the schema is named like "student records" or "users' payment", then
	// they are insane! And sorry, we can't help them.
	return snaker.CamelToSnake(s.Name)
}

func (s *Schema) UnmarshalText(b []byte) error {
	text := string(b)
	index := strings.LastIndex(text, ":")
	if index > 0 {
		s.Name = text[:index]
		s.Package = text[index+1:]
		return nil
	}

	s.Name = text
	return nil
}

// EnumValue holds data for a single enum value.
type EnumValue struct {
	Name    string
	Val     *models.EnumValue
	Comment string
}

// Enum is a template item for a enum.
type Enum struct {
	Name              string
	Schema            string
	Values            []*EnumValue
	Enum              *models.Enum
	Comment           string
	ReverseConstNames bool
}

// Proc is a template item for a stored procedure.
type Proc struct {
	Name       string
	Schema     string
	ProcParams string
	Params     []*Field
	Return     *Field
	Proc       *models.Proc
	Comment    string
}

// Field contains field information.
type Field struct {
	Name    string
	Type    string
	NilType string
	Len     int
	Col     *models.Column
	Comment string
}

// Type is a template item for a type (ie, table/view/custom query).
type Type struct {
	Name             string
	Schema           string
	RelType          RelType
	PrimaryKey       *Field
	PrimaryKeyFields []*Field
	Fields           []*Field
	Table            *models.Table
	Comment          string
}

// ForeignKey is a template item for a foreign relationship on a table.
type ForeignKey struct {
	Name       string
	Schema     string
	Type       *Type
	Field      *Field
	RefType    *Type
	RefField   *Field
	ForeignKey *models.ForeignKey
	Comment    string
}

// Index is a template item for a index into a table.
type Index struct {
	FuncName string
	Schema   string
	Type     *Type
	Fields   []*Field
	Index    *models.Index
	Comment  string
}

// QueryParam is a query parameter for a custom query.
type QueryParam struct {
	Name        string
	Type        string
	Interpolate bool
}

// Query is a template item for a custom query.
type Query struct {
	Schema        string
	Name          string
	Query         []string
	QueryComments []string
	QueryParams   []*QueryParam
	OnlyOne       bool
	Interpolate   bool
	Type          *Type
	Comment       string
}

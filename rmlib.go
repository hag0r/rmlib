package main

import (
	"fmt"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

// RelationDefinition is the base entry point for parsing realtion definitions
// given in the relation model syntax
type RelationDefinition struct {
	// Name is the relation's name
	Name string `@Ident`
	// Attributes is the list of attributes defined for the relation
	Attributes []Attribute `"("@@ ("," @@)* ")"`
}

// GetPKAttributes returns all attributes included in the primary key
func (r *RelationDefinition) GetPKAttributes() (res []Attribute) {
	for _, attr := range r.Attributes {
		if attr.IsPK() {
			res = append(res, attr)
		}
	}
	return
}

// GetFKs returns all attributes that are foreign keys
func (r *RelationDefinition) GetFKs() (res []Attribute) {
	for _, attr := range r.Attributes {
		if attr.IsFK() {
			res = append(res, attr)
		}
	}
	return
}

// Attribute defines an attribute in a relation
type Attribute struct {
	// PK holds the field name if it is a primary key, otherwise nil
	PK *string `("_" @Ident "_"`
	// AttrName holds the field name if it is not a primary key, otherwise nil
	AttrName *string `| @Ident )`
	// FK holds the referenced table and column, otherwise nil
	FK *ForeignKeyDefinition `( "->" @@ )?`
}

// Name returns the attribute's name. Use this instead of direct access to AttrName or PK, which might be nil
func (attr Attribute) Name() *string {
	if attr.IsPK() {
		return attr.PK
	} else {
		return attr.AttrName
	}
}

// IsPK tests if the attribute is part of the primary key in its relation
func (attr Attribute) IsPK() bool {
	return attr.PK != nil
}

// IsFK tests if the attribute is a foreign key
func (attr Attribute) IsFK() bool {
	return attr.FK != nil
}

// ForeignKeyDefinition is the definition of a foreign key
type ForeignKeyDefinition struct {
	// RelationName holds the name of the referenced table
	RelationName string `@Ident`
	// Attribute holds the name of the referenced attribute
	Attribute string `"(" @Ident ")"`
}

// RelationParser encapsules the parser generator and provides
// convenient methods to parse a string into a relation definition
type RelationParser struct {
	parser *participle.Parser
}

// CreateParser creates a new parser instance to parse relation definitions
func CreateParser() (*RelationParser, error) {

	rmLexer := lexer.Must(lexer.Regexp(`(\s+)` +
		// `|(?P<PK>_[a-zA-Z][a-zA-Z0-9]*_)` +
		`|(?P<Ident>[a-zA-Z][a-zA-Z0-9]*)` +
		`|(?P<Operators>->|[,()_])`))

	parser, err := participle.Build(&RelationDefinition{}, participle.Lexer(rmLexer))

	if err != nil {
		return nil, err
	}

	relParser := &RelationParser{parser}
	return relParser, nil
}

// Parse takes a string and tries to parse it into a RelationDefinition
// The given string should be a single line defining one relation and its attributes, like
// R(a, b, _c_,_d_, e -> R2(k))
func (parser *RelationParser) Parse(s string) (*RelationDefinition, error) {
	relation := &RelationDefinition{}

	err := parser.parser.ParseString(s, relation)

	if err != nil {
		return nil, err
	}

	return relation, nil
}

func main() {

	s := "MyRelation(a -> OtherRelation(b),_b_,c,_d_ -> D(c))" //

	parser, err := CreateParser()

	if err != nil {
		fmt.Println("Failed to create parser", err)
		return
	}

	relation, err := parser.Parse(s)

	if err != nil {
		fmt.Println("got an error during parse", err)
		return
	}

	fmt.Printf("Successfully parsed relation: %v\n", relation)

	for _, attr := range relation.Attributes {

		if !attr.IsPK() && !attr.IsFK() {
			fmt.Printf("%s is a plain attribute\n", *attr.Name())
		} else {
			if attr.IsPK() {
				fmt.Printf("%s is a PK \n", *attr.Name())
			}

			if attr.IsFK() {
				fmt.Printf("%s is a FK to attribute %s in relation %s \n", *attr.Name(), attr.FK.Attribute, attr.FK.RelationName)
			}
		}
	}
}

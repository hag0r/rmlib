# RM Lib 
RMLib - short for relational model library - models relations, attributes, as well as primary and foreign keys.

Currently, it parses relation definitions from strings into Go structs so that they can be evaluated by some other tool.

## Syntax
The main syntax for defining relations with attributes is:
```
RelationName ( Attribute1, Attribute2, Attribute3, ... )
```

Each attribute can be a primary key, a foreign, none of them, or both. 

- *Primary Key*: Primary keys are defined by underscores before and after the attribute name: `_mypk_` . There might be multiple attributes marked as primary key in a single relation
- *Foreign key*: Foreign keys are defined by an arrow `->` followed by the relation name and attributed being referenced: `myfk -> R(a)`. `myfk` is an attribute that referenecs attribute `a` in relation `R`.

## Dependencies

We use the great [`https://github.com/alecthomas/participle`](https://github.com/alecthomas/participle) parser library

## Limitations

- Currently, a foreign key can only consist of a single attribute and reference one attribute


## Example 

```golang
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
```

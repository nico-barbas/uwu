package toml

import "fmt"

// NOTE: Poor man's union type..

type Value interface {
	toString() string
}

type (
	Number  float64
	Boolean bool
	String  string
	Array   []Value
)

func (nv Number) toString() string { return "Number" }

func (bv Boolean) toString() string { return "Boolean" }

func (sv String) toString() string { return "String" }

func (a Array) toString() string { return "Array" }

type Table map[string]Value

func (t Table) toString() string { return "Table" }

func (t Table) insertValue(key string, value Value) {
	if _, exist := t[key]; !exist {
		t[key] = value
	} else {
		fmt.Println(`Already key with name "`, key, `"`)
	}
}

func (t Table) insertInlineTable(keyDecl []token) Table {
	currentTable := t
	for _, token := range keyDecl {
		if value, exist := currentTable[token.lexeme]; exist {
			if table, ok := value.(Table); ok {
				currentTable = table
			}
		} else {
			table := make(Table)
			currentTable[token.lexeme] = table
			currentTable = table
		}
	}
	return currentTable
}

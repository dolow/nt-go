# nt-go

nt-go is a paraser for [NestedText](https://github.com/KenKundert/nestedtext) format.

It covers [official teste cases](https://github.com/KenKundert/nestedtext_tests/tree/master/test_cases).


# Usage

## Parsing

Example for parsing and accessing nested text data below;

```
str: hello
text:
  > multi
  > line
list:
  - first str
  -
    - second is
    - list
dict:
  key1: it is str
  key2:
    key2_1: nested
    key2_2: dict

```

```
package main

import (
        "fmt"
        "github.com/dolow/nt-go"
)

func main() {
	var content []byte
	content = someHowGetContent()

        directive := &ntgo.Directive{}
        directive.Marshal(content)

	// accessing root level type
	// type is described as iota
	// you can identify directive type with this field even if you are handling content with unknown structure
	// 0: Unknown, 1: String, 2: Text(multiline text), 3: List, 4: Dictionary
        fmt.Printf("%v\n", directive.Type)

	// accessing dictionary directive element with "str" key
	// every child elements are also instance of Directive
        fmt.Printf("%v\n", directive.Dictionary["str"].String)

	// accessing second line of Text directive 
        fmt.Printf("%v\n", directive.Dictionary["text"].Text[1])

	// accessing first element of list directive
        fmt.Printf("%v\n", directive.Dictionary["list"].List[0].String)
	// of course, nested data is capable
        fmt.Printf("%v\n", directive.Dictionary["list"].List[1].List[1].String)

	// accessing nested dictionary directive elements
        fmt.Printf("%v\n", directive.Dictionary["dict"].Dictionary["key1"].String)
        fmt.Printf("%v\n", directive.Dictionary["dict"].Dictionary["key2"].Dictionary["key2_1"].String)
}

```


## Stringify

Just send `Unmarshal` to directive instance that has already marshalized.

```
directive.Unmarshal()
```

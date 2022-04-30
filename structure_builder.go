package gorouter

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func NewStructureBuilder() *StructureBuilder {
	return &StructureBuilder{
		StructureTree: map[string]*RouteDefine{},
	}
}

//StructureBuilder help to build route structure.
type StructureBuilder struct {
	StructureTree map[string]*RouteDefine
}

//AddOneLine define a route in oneline string.
//example define /parent/sub : /parent/sub/:id,id_2
//example define /parent/sub/sub_2: /parent/sub/sub_2 or  /parent/sub/:id,id_2/sub_2
//Note: the indexes of parent in sub defination must be consistent with the defination of parent.

func (builder *StructureBuilder) AddOneLine(defineString string) error {

	parts := strings.Split(defineString, "/")
	fmt.Println(parts)
	var last *RouteDefine = nil
	errInvalid := errors.New("invalid define")
	for _, part := range parts {

		formattedPart := strings.TrimSpace(part)

		if len(formattedPart) == 0 {

			continue
		}
		if formattedPart[0] == ':' {
			if last == nil {
				return errInvalid
			}
			indexes := strings.Split(formattedPart[1:], ",")
			validIndexes := []string{}

			for _, index := range indexes {
				formattedIndex := strings.TrimSpace(index)
				if len(formattedIndex) > 0 {
					validIndexes = append(validIndexes, formattedIndex)
				}
			}
			numIndex := len(validIndexes)
			if numIndex == 0 || (len(last.Subs) > 0 && len(last.Indexes) != numIndex) {

				return errInvalid

			} else if len(last.Subs) == 0 {
				fmt.Println(last.Indexes, validIndexes)
				for i, existedIndex := range last.Indexes {

					if existedIndex != validIndexes[i] {
						return errInvalid
					}
				}
			}
			last.Indexes = validIndexes

			continue
		}
		if last == nil {
			if def, ok := builder.StructureTree[formattedPart]; ok {

				last = def

			} else {
				last = NewEmptyRouteDefine()
				builder.StructureTree[formattedPart] = last
			}
			continue

		} else if last.SubRoute(formattedPart) == nil {

			newSubRoute := NewEmptyRouteDefine()

			last.Subs[formattedPart] = newSubRoute
			last = newSubRoute

			continue
		}
		last = last.SubRoute(formattedPart)

	}
	return nil
}

func (builder *StructureBuilder) Export() string {
	data, _ := json.Marshal(builder.StructureTree)
	return string(data)
}

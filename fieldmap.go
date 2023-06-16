package fieldmap

import (
	"reflect"
)

// Field ...
type Field interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// FieldMap ...
type FieldMap[T any, F Field] struct {
	mapping  *T
	fields   []F
	children []int64
	parent   []F
}

// New ...
func New[T any, F Field]() (*FieldMap[T, F], error) {
	var mapping T

	var field F
	fieldType := reflect.TypeOf(field)

	getField := func(num int64) F {
		var f F
		val := reflect.ValueOf(&f).Elem()
		val.SetInt(num)
		return f
	}

	ordinal := int64(0)
	var fields []F
	var children []int64
	var parentList []F

	val := reflect.ValueOf(&mapping)
	val = val.Elem()

	var traverseFunc func(val reflect.Value, isFirst bool) error
	traverseFunc = func(val reflect.Value, isFirst bool) error {
		var parent F
		if !isFirst {
			fieldName := val.Type().Field(0).Name
			if fieldName != "Root" {
				// TODO Check missing root
			}
			parent = getField(ordinal + 1)
		}

		for i := 0; i < val.NumField(); i++ {
			f := val.Field(i)

			if f.Kind() == reflect.Struct {
				err := traverseFunc(f, false)
				if err != nil {
					return err
				}
				continue
			}

			if f.Type() != fieldType {
				panic("TODO") // TODO
			}

			if !f.CanSet() {
				panic("TODO") // TODO
			}

			ordinal++

			fields = append(fields, getField(ordinal))
			parentList = append(parentList, parent)

			if !isFirst && i == 0 {
				children = append(children, int64(val.NumField()-1))
			} else {
				children = append(children, 0)
			}
			f.SetInt(ordinal)
		}

		return nil
	}

	err := traverseFunc(val, true)
	if err != nil {
		return nil, err
	}

	return &FieldMap[T, F]{
		mapping:  &mapping,
		fields:   fields,
		children: children,
		parent:   parentList,
	}, nil
}

// GetMapping ...
func (f *FieldMap[T, F]) GetMapping() *T {
	return f.mapping
}

func (*FieldMap[T, F]) indexOf(field F) int64 {
	return reflect.ValueOf(field).Int() - 1
}

// IsStruct ...
func (f *FieldMap[T, F]) IsStruct(field F) bool {
	index := f.indexOf(field)
	return f.children[index] > 0
}

// ChildrenOf ...
func (f *FieldMap[T, F]) ChildrenOf(field F) []F {
	index := f.indexOf(field)
	return f.fields[index+1 : index+1+f.children[index]]
}

// ParentOf ...
func (f *FieldMap[T, F]) ParentOf(field F) F {
	return f.parent[f.indexOf(field)]
}

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
	mapping *T
	fields  []F

	children   []int64
	parentList []F
	fieldNames []string
}

// New ...
func New[T any, F Field]() (*FieldMap[T, F], error) {
	var mapping T
	val := reflect.ValueOf(&mapping)
	val = val.Elem()

	f := &FieldMap[T, F]{
		mapping: &mapping,
	}

	var emptyField F
	ordinal := int64(0)

	err := f.traverse(val, true, &ordinal, emptyField, "")
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (*FieldMap[T, F]) getField(num int64) F {
	var field F
	val := reflect.ValueOf(&field).Elem()
	val.SetInt(num)
	return field
}

func (*FieldMap[T, F]) getFieldType() reflect.Type {
	var field F
	return reflect.TypeOf(field)
}

//revive:disable-next-line:flag-parameter
func (f *FieldMap[T, F]) traverse(
	val reflect.Value, isOuterMost bool, ordinal *int64,
	prevParent F, parentFieldName string,
) error {
	var parent F
	if !isOuterMost {
		fieldName := val.Type().Field(0).Name
		if fieldName != "Root" {
			// TODO Check missing root
			panic("TODO")
		}
		parent = f.getField(*ordinal + 1)
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := val.Type().Field(i).Name

		if field.Kind() == reflect.Struct {
			err := f.traverse(field, false, ordinal, parent, fieldName)
			if err != nil {
				return err
			}
			continue
		}

		if field.Type() != f.getFieldType() {
			panic("TODO") // TODO
		}

		if !field.CanSet() {
			panic("TODO") // TODO
		}

		*ordinal++

		f.fields = append(f.fields, f.getField(*ordinal))

		if !isOuterMost && i == 0 {
			f.children = append(f.children, int64(val.NumField()-1))
			f.parentList = append(f.parentList, prevParent)
			f.fieldNames = append(f.fieldNames, parentFieldName)
		} else {
			f.children = append(f.children, 0)
			f.parentList = append(f.parentList, parent)
			f.fieldNames = append(f.fieldNames, fieldName)
		}
		field.SetInt(*ordinal)
	}

	return nil
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
	return f.parentList[f.indexOf(field)]
}

// AncestorOf ...
func (f *FieldMap[T, F]) AncestorOf(field F) []F {
	var empty F

	result := []F{field}
	for {
		field = f.ParentOf(field)
		if field == empty {
			return result
		}
		result = append(result, field)
	}
}

// GetFieldName ...
func (f *FieldMap[T, F]) GetFieldName(field F) string {
	return f.fieldNames[f.indexOf(field)]
}

// GetFullFieldName ...
func (f *FieldMap[T, F]) GetFullFieldName(field F) string {
	fullName := ""
	for {
		name := f.GetFieldName(field)
		if len(fullName) > 0 {
			fullName = name + "." + fullName
		} else {
			fullName = name
		}

		field = f.ParentOf(field)
		var empty F
		if field == empty {
			return fullName
		}
	}
}

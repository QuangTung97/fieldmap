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
	weights    []float64
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

	_, err := f.traverse(val, true, &ordinal, emptyField)
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

func (f *FieldMap[T, F]) traverse(
	val reflect.Value, isOuter bool, ordinal *int64, prevParent F,
) (float64, error) {
	var parent F
	if !isOuter {
		fieldName := val.Type().Field(0).Name
		if fieldName != "Root" {
			// TODO Check missing root
			panic("TODO")
		}
		parent = f.getField(*ordinal + 1)
	}

	rootWeightIndex := -1
	totalWeight := 0.0
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if field.Kind() == reflect.Struct {
			childWeight, err := f.traverse(field, false, ordinal, parent)
			if err != nil {
				return 0, err
			}
			totalWeight += childWeight
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

		if !isOuter && i == 0 {
			f.children = append(f.children, int64(val.NumField()-1))
			f.parentList = append(f.parentList, prevParent)

			rootWeightIndex = len(f.weights)
			f.weights = append(f.weights, 0.0)
		} else {
			f.children = append(f.children, 0)
			f.parentList = append(f.parentList, parent)

			f.weights = append(f.weights, 1.0)
			totalWeight += 1.0
		}
		field.SetInt(*ordinal)
	}

	if rootWeightIndex >= 0 {
		f.weights[rootWeightIndex] = totalWeight
	}

	return totalWeight, nil
}

// GetMapping ...
func (f *FieldMap[T, F]) GetMapping() *T {
	return f.mapping
}

func (f *FieldMap[T, F]) GetWeight(field F) float64 {
	return f.weights[f.indexOf(field)]
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

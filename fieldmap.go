package fieldmap

import (
	"fmt"
	"reflect"
)

// Field ...
type Field interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// RootField ...
const RootField = "Root"

// FieldMap ...
type FieldMap[T any, F Field] struct {
	options fieldMapOptions

	mapping *T
	fields  []F

	children   []int64
	parentList []F
	fieldNames []string
	structTags map[string][]string
}

type fieldMapOptions struct {
	structTags []string
}

// Option ...
type Option func(opts *fieldMapOptions)

// WithStructTags ...
func WithStructTags(tags ...string) Option {
	return func(opts *fieldMapOptions) {
		opts.structTags = tags
	}
}

func computeOptions(options []Option) fieldMapOptions {
	opts := fieldMapOptions{
		structTags: nil,
	}
	for _, fn := range options {
		fn(&opts)
	}
	return opts
}

// New ...
func New[T any, F Field](options ...Option) *FieldMap[T, F] {
	var mapping T
	val := reflect.ValueOf(&mapping)
	val = val.Elem()

	opts := computeOptions(options)

	f := &FieldMap[T, F]{
		options:    opts,
		mapping:    &mapping,
		structTags: map[string][]string{},
	}

	ordinal := int64(0)
	var info parentInfoData[F]

	f.traverse(val, &ordinal, info)
	return f
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

type parentInfoData[F Field] struct {
	valid    bool
	prevRoot F

	fieldName     string
	fullFieldName string

	structTags map[string]string
}

func (p parentInfoData[F]) isParentField(index int) bool {
	return p.valid && index == 0
}

func (p parentInfoData[F]) computeFullName(currentName string) string {
	if p.valid {
		return p.fullFieldName + "." + currentName
	}
	return currentName
}

func (f *FieldMap[T, F]) findStructTags(
	fieldType reflect.StructField,
	fullFieldName string,
) map[string]string {
	structTags := map[string]string{}

	for _, tag := range f.options.structTags {
		tagVal := fieldType.Tag.Get(tag)
		if len(tagVal) == 0 {
			panic(
				fmt.Sprintf(
					"missing struct tag %q for field %q",
					tag, fullFieldName,
				),
			)
		}
		structTags[tag] = tagVal
	}
	return structTags
}

func (f *FieldMap[T, F]) getRootField(
	val reflect.Value, parentInfo parentInfoData[F], ordinal *int64,
) F {
	if parentInfo.valid {
		panicStr := fmt.Sprintf("missing field %q for field %q", RootField, parentInfo.fullFieldName)

		if val.NumField() == 0 {
			panic(panicStr)
		}

		fieldName := val.Type().Field(0).Name
		if fieldName != RootField {
			panic(panicStr)
		}
		return f.getField(*ordinal + 1)
	}
	var empty F
	return empty
}

func (f *FieldMap[T, F]) handleSingleField(
	val reflect.Value, i int, parentInfo parentInfoData[F],
	rootField F, ordinal *int64,
) {
	field := val.Field(i)
	fieldType := val.Type().Field(i)
	fieldName := fieldType.Name
	fullFieldName := parentInfo.computeFullName(fieldName)

	var currentStructTags map[string]string
	if !parentInfo.isParentField(i) {
		currentStructTags = f.findStructTags(fieldType, fullFieldName)

		if field.Kind() == reflect.Struct {
			newInfo := parentInfoData[F]{
				valid:    true,
				prevRoot: rootField,

				fieldName:     fieldName,
				fullFieldName: fullFieldName,

				structTags: currentStructTags,
			}
			f.traverse(field, ordinal, newInfo)
			return
		}
	}

	if field.Type() != f.getFieldType() {
		panic(fmt.Sprintf("invalid type for field %q", fullFieldName))
	}

	*ordinal++

	f.fields = append(f.fields, f.getField(*ordinal))

	if parentInfo.isParentField(i) {
		f.children = append(f.children, int64(val.NumField()-1))
		f.parentList = append(f.parentList, parentInfo.prevRoot)
		f.fieldNames = append(f.fieldNames, parentInfo.fieldName)

		for _, tag := range f.options.structTags {
			f.structTags[tag] = append(f.structTags[tag], parentInfo.structTags[tag])
		}
	} else {
		f.children = append(f.children, 0)
		f.parentList = append(f.parentList, rootField)
		f.fieldNames = append(f.fieldNames, fieldName)

		for _, tag := range f.options.structTags {
			f.structTags[tag] = append(f.structTags[tag], currentStructTags[tag])
		}
	}
	field.SetInt(*ordinal)
}

func (f *FieldMap[T, F]) traverse(
	val reflect.Value, ordinal *int64, parentInfo parentInfoData[F],
) {
	rootField := f.getRootField(val, parentInfo, ordinal)

	for i := 0; i < val.NumField(); i++ {
		f.handleSingleField(val, i, parentInfo, rootField, ordinal)
	}
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

// AncestorOf includes itself, parent, and all parents of parents
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

// GetStructTag ...
func (f *FieldMap[T, F]) GetStructTag(tag string, field F) string {
	return f.structTags[tag][f.indexOf(field)]
}

// GetFullStructTag ...
func (f *FieldMap[T, F]) GetFullStructTag(tag string, field F) string {
	fullTag := ""
	for {
		tagName := f.GetStructTag(tag, field)
		if len(fullTag) > 0 {
			fullTag = tagName + "." + fullTag
		} else {
			fullTag = tagName
		}

		field = f.ParentOf(field)
		var empty F
		if field == empty {
			return fullTag
		}
	}
}

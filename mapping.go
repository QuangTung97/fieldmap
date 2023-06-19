package fieldmap

import "fmt"

// Mapper ...
type Mapper[F1 Field, F2 Field] struct {
	parentOf func(source F1) F1
	fieldMap map[F1][][]F2
}

// Mapping ...
type Mapping[F1, F2 Field] struct {
	From   F1
	ToList []F2
}

// NewMapping ...
func NewMapping[F1 Field, F2 Field](from F1, toList ...F2) Mapping[F1, F2] {
	if len(toList) == 0 {
		panic("missing destination fields")
	}
	return Mapping[F1, F2]{From: from, ToList: toList}
}

type emptyStruct struct{}

// NewMapper ...
func NewMapper[T1 any, F1 Field, T2 any, F2 Field](
	source *FieldMap[T1, F1],
	dest *FieldMap[T2, F2],
	mappings ...Mapping[F1, F2],
) *Mapper[F1, F2] {
	fieldMap := map[F1][][]F2{}
	dedupSets := map[F1]map[F2]emptyStruct{}

	getDedupSet := func(source F1) map[F2]emptyStruct {
		s, ok := dedupSets[source]
		if !ok {
			s = map[F2]emptyStruct{}
		}
		dedupSets[source] = s
		return s
	}

	var addMappingList func(mappings []Mapping[F1, F2])
	addMappingList = func(mappings []Mapping[F1, F2]) {
		for _, m := range mappings {
			set := getDedupSet(m.From)
			if len(m.ToList) == 1 {
				for _, to := range m.ToList {
					_, existed := set[to]
					if existed {
						panic(fmt.Sprintf(
							"duplicated destination field %q for source field %q",
							dest.GetFullFieldName(to),
							source.GetFullFieldName(m.From),
						))
					}
					set[to] = emptyStruct{}
				}
			}
			fieldMap[m.From] = append(fieldMap[m.From], m.ToList)
		}
	}

	addMappingList(mappings)

	return &Mapper[F1, F2]{
		parentOf: source.ParentOf,
		fieldMap: fieldMap,
	}
}
func (m *Mapper[F1, F2]) findMappedFieldsForSourceField(
	sourceField F1, resultSet map[F2]emptyStruct, result []F2,
) []F2 {
	var empty F1

	for {
		for _, destFields := range m.fieldMap[sourceField] {
			for _, f := range destFields {
				_, existed := resultSet[f]
				if existed {
					continue
				}
				resultSet[f] = emptyStruct{}
				result = append(result, f)
			}
			if len(destFields) > 0 {
				return result
			}
		}

		sourceField = m.parentOf(sourceField)
		if sourceField == empty {
			return result
		}
	}
}

// FindMappedFields ...
func (m *Mapper[F1, F2]) FindMappedFields(sourceFields []F1) []F2 {
	var result []F2
	resultSet := map[F2]emptyStruct{}

	for _, sourceField := range sourceFields {
		result = m.findMappedFieldsForSourceField(sourceField, resultSet, result)
	}

	return result
}

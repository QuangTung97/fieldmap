package fieldmap

// Mapper ...
type Mapper[F1 Field, F2 Field] struct {
	fieldMap map[F1][]F2
}

// Mapping ...
type Mapping[F1, F2 Field] struct {
	From   F1
	ToList []F2
}

// NewMapping ...
func NewMapping[F1 Field, F2 Field](from F1, toList ...F2) Mapping[F1, F2] {
	if len(toList) == 0 {
		panic("TODO") // TODO
	}
	return Mapping[F1, F2]{From: from, ToList: toList}
}

type emptyStruct struct{}

// NewMapper ...
func NewMapper[T1 any, F1 Field, T2 any, F2 Field](
	source *FieldMap[T1, F1],
	_ *FieldMap[T2, F2],
	mappings ...Mapping[F1, F2],
) *Mapper[F1, F2] {
	fieldMap := map[F1][]F2{}
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
			if source.IsStruct(m.From) {
				var childMappings []Mapping[F1, F2]
				for _, child := range source.ChildrenOf(m.From) {
					childMappings = append(childMappings, NewMapping(child, m.ToList...))
				}
				addMappingList(childMappings)
			}

			set := getDedupSet(m.From)
			for _, to := range m.ToList {
				_, existed := set[to]
				if existed {
					return
				}
				set[to] = emptyStruct{}
			}
			fieldMap[m.From] = append(fieldMap[m.From], m.ToList...)
		}
	}

	addMappingList(mappings)

	return &Mapper[F1, F2]{
		fieldMap: fieldMap,
	}
}

// FindMappedFields ...
func (m *Mapper[F1, F2]) FindMappedFields(sourceFields []F1) []F2 {
	var result []F2
	set := map[F2]emptyStruct{}

	for _, sourceField := range sourceFields {
		for _, f := range m.fieldMap[sourceField] {
			_, existed := set[f]
			if existed {
				continue
			}
			set[f] = emptyStruct{}
			result = append(result, f)
		}
	}

	return result
}

func (m *Mapper[F1, F2]) GetWeight(destField F2) float64 {
	return 0
}

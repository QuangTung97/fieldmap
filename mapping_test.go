package fieldmap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type sourceField int
type destField int

type sourceDataSimple struct {
	Sku  sourceField
	Name sourceField
	Body sourceField
}

type destDataSimple struct {
	Info   destField
	Detail destField
}

func TestMapping_Simple_Structs(t *testing.T) {
	sourceFm, err := New[sourceDataSimple, sourceField]()
	assert.Equal(t, nil, err)

	destFm, err := New[destDataSimple, destField]()
	assert.Equal(t, nil, err)

	source := sourceFm.GetMapping()
	dest := destFm.GetMapping()

	m := NewMapper(
		sourceFm, destFm,
		NewMapping(source.Sku, dest.Info),
		NewMapping(source.Name, dest.Info),
		NewMapping(source.Body, dest.Detail),
	)

	assert.Equal(t, 0, len(m.FindMappedFields(nil)))

	assert.Equal(t, []destField{dest.Info}, m.FindMappedFields([]sourceField{source.Sku}))
	assert.Equal(t, []destField{dest.Info}, m.FindMappedFields([]sourceField{source.Sku, source.Name}))

	assert.Equal(t, []destField{dest.Info, dest.Detail}, m.FindMappedFields([]sourceField{source.Sku, source.Body}))
}

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

type sourceSellerInfo struct {
	Root sourceField

	Logo sourceField
}

type sourceSeller struct {
	Root sourceField

	ID   sourceField
	Name sourceField
	Info sourceSellerInfo
}

type sourceDataComplex struct {
	Sku      sourceField
	Name     sourceField
	Body     sourceField
	Seller   sourceSeller
	ImageURL sourceField
}

type destInfo struct {
	Root destField

	Sku  destField
	Name destField
}

type destDetail struct {
	Root destField

	Body destField
}

type destDataComplex struct {
	Info       destInfo
	Detail     destDetail
	SearchText destField
}

func TestMapping_Complex(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		sourceFm, err := New[sourceDataComplex, sourceField]()
		assert.Equal(t, nil, err)

		destFm, err := New[destDataComplex, destField]()
		assert.Equal(t, nil, err)

		source := sourceFm.GetMapping()
		dest := destFm.GetMapping()

		m := NewMapper(
			sourceFm, destFm,
			NewMapping(source.Sku, dest.Info.Root),
			NewMapping(source.Name, dest.Info.Root),
			NewMapping(source.Seller.Root, dest.Detail.Root),
			NewMapping(source.Body, dest.Detail.Body),
			NewMapping(source.ImageURL, dest.Detail.Root),
		)

		assert.Equal(t, 0, len(m.FindMappedFields(nil)))

		assert.Equal(t, []destField{dest.Info.Root}, m.FindMappedFields([]sourceField{source.Sku}))
		assert.Equal(t, []destField{dest.Info.Root}, m.FindMappedFields([]sourceField{source.Sku, source.Name}))

		assert.Equal(t, []destField{dest.Detail.Body}, m.FindMappedFields([]sourceField{source.Body}))

		assert.Equal(t, []destField{dest.Detail.Root}, m.FindMappedFields([]sourceField{source.Seller.Root}))

		assert.Equal(t, []destField{dest.Detail.Root}, m.FindMappedFields([]sourceField{source.Seller.ID}))
		assert.Equal(t, []destField{dest.Detail.Root}, m.FindMappedFields([]sourceField{source.Seller.Name}))

		assert.Equal(t, []destField{dest.Detail.Root}, m.FindMappedFields([]sourceField{source.Seller.Info.Root}))

		assert.Equal(t, []destField{dest.Detail.Root}, m.FindMappedFields([]sourceField{source.Seller.Info.Logo}))

		assert.Equal(t, []destField{dest.Detail.Root},
			m.FindMappedFields([]sourceField{source.Seller.ID, source.Seller.Info.Logo}))

		assert.Equal(t, []destField{dest.Info.Root, dest.Detail.Root},
			m.FindMappedFields([]sourceField{source.Sku, source.Seller.Info.Logo}))
	})

	t.Run("one field to multiple dest fields using logical AND", func(t *testing.T) {
		sourceFm, err := New[sourceDataComplex, sourceField]()
		assert.Equal(t, nil, err)

		destFm, err := New[destDataComplex, destField]()
		assert.Equal(t, nil, err)

		source := sourceFm.GetMapping()
		dest := destFm.GetMapping()

		m := NewMapper(
			sourceFm, destFm,
			NewMapping(source.Sku, dest.Info.Root, dest.SearchText),
		)

		assert.Equal(t, []destField{dest.Info.Root, dest.SearchText}, m.FindMappedFields([]sourceField{source.Sku}))
	})

	t.Run("one field to multiple dest fields using logical OR", func(t *testing.T) {
		sourceFm, err := New[sourceDataComplex, sourceField]()
		assert.Equal(t, nil, err)

		destFm, err := New[destDataComplex, destField]()
		assert.Equal(t, nil, err)

		source := sourceFm.GetMapping()
		dest := destFm.GetMapping()

		m := NewMapper(
			sourceFm, destFm,
			NewMapping(source.Name, dest.Info.Root),
			NewMapping(source.Sku, dest.Info.Root),
			NewMapping(source.Sku, dest.SearchText),
		)

		assert.Equal(t, 2.0, m.GetWeight(dest.Info.Root))
		assert.Equal(t, 1.0, m.GetWeight(dest.Info.Root))

		assert.Equal(t, []destField{dest.SearchText}, m.FindMappedFields([]sourceField{source.Sku}))
	})

	t.Run("multiple source fields to one dest field", func(t *testing.T) {
		sourceFm, err := New[sourceDataComplex, sourceField]()
		assert.Equal(t, nil, err)

		destFm, err := New[destDataComplex, destField]()
		assert.Equal(t, nil, err)

		source := sourceFm.GetMapping()
		dest := destFm.GetMapping()

		m := NewMapper(
			sourceFm, destFm,
			NewMapping(source.Sku, dest.SearchText),
			NewMapping(source.Name, dest.SearchText),
			NewMapping(source.Seller.Name, dest.SearchText),
		)

		assert.Equal(t, []destField{dest.SearchText}, m.FindMappedFields([]sourceField{source.Sku}))
	})
}

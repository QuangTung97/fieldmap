package fieldmap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type field int

type simpleData struct {
	Sku      field
	Name     field
	ImageURL field
}

type sellerData struct {
	Root field

	ID   field
	Name field
	Logo field
}

type productData struct {
	Sku      field
	Name     field
	Seller   sellerData
	ImageURL field
}

func TestFieldMap__GetMapping(t *testing.T) {
	t.Run("simple-struct", func(t *testing.T) {
		fm, err := InitFieldMap[simpleData, field]()
		assert.Equal(t, nil, err)

		p := fm.GetMapping()

		assert.Equal(t, field(1), p.Sku)
		assert.Equal(t, field(2), p.Name)
		assert.Equal(t, field(3), p.ImageURL)
	})

	t.Run("complex-struct", func(t *testing.T) {
		fm, err := InitFieldMap[productData, field]()
		assert.Equal(t, nil, err)

		p := fm.GetMapping()

		assert.Equal(t, field(1), p.Sku)
		assert.Equal(t, field(2), p.Name)
		assert.Equal(t, field(3), p.Seller.Root)
		assert.Equal(t, field(4), p.Seller.ID)
		assert.Equal(t, field(5), p.Seller.Name)
		assert.Equal(t, field(6), p.Seller.Logo)
		assert.Equal(t, field(7), p.ImageURL)

		assert.Equal(t, false, fm.IsStruct(p.Sku))
		assert.Equal(t, true, fm.IsStruct(p.Seller.Root))
		assert.Equal(t, false, fm.IsStruct(p.Seller.ID))

		assert.Equal(t, []field{4, 5, 6}, fm.ChildrenOf(p.Seller.Root))

		assert.Equal(t, p.Seller.Root, fm.ParentOf(p.Seller.ID))
		assert.Equal(t, p.Seller.Root, fm.ParentOf(p.Seller.Name))
		assert.Equal(t, field(0), fm.ParentOf(p.Sku))
		assert.Equal(t, field(0), fm.ParentOf(p.Name))
	})
}

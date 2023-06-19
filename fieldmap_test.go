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

type sellerAttr struct {
	Root field

	Code field
	Name field
}

type sellerData struct {
	Root field

	ID   field
	Name field
	Logo field
	Attr sellerAttr
}

type productData struct {
	Sku      field
	Name     field
	Seller   sellerData
	ImageURL field
}

func TestFieldMap__GetMapping(t *testing.T) {
	t.Run("simple-struct", func(t *testing.T) {
		fm, err := New[simpleData, field]()
		assert.Equal(t, nil, err)

		p := fm.GetMapping()

		assert.Equal(t, field(1), p.Sku)
		assert.Equal(t, field(2), p.Name)
		assert.Equal(t, field(3), p.ImageURL)

		// Get Weights
		assert.Equal(t, 1.0, fm.GetWeight(p.Sku))
		assert.Equal(t, 1.0, fm.GetWeight(p.Name))
		assert.Equal(t, 1.0, fm.GetWeight(p.ImageURL))
	})

	t.Run("complex-struct", func(t *testing.T) {
		fm, err := New[productData, field]()
		assert.Equal(t, nil, err)

		p := fm.GetMapping()

		assert.Equal(t, field(1), p.Sku)
		assert.Equal(t, field(2), p.Name)
		assert.Equal(t, field(3), p.Seller.Root)
		assert.Equal(t, field(4), p.Seller.ID)
		assert.Equal(t, field(5), p.Seller.Name)
		assert.Equal(t, field(6), p.Seller.Logo)
		assert.Equal(t, field(7), p.Seller.Attr.Root)
		assert.Equal(t, field(8), p.Seller.Attr.Code)
		assert.Equal(t, field(9), p.Seller.Attr.Name)
		assert.Equal(t, field(10), p.ImageURL)

		assert.Equal(t, false, fm.IsStruct(p.Sku))
		assert.Equal(t, true, fm.IsStruct(p.Seller.Root))
		assert.Equal(t, false, fm.IsStruct(p.Seller.ID))

		assert.Equal(t, []field{4, 5, 6, 7}, fm.ChildrenOf(p.Seller.Root))

		assert.Equal(t, p.Seller.Root, fm.ParentOf(p.Seller.ID))
		assert.Equal(t, p.Seller.Root, fm.ParentOf(p.Seller.Name))

		assert.Equal(t, field(0), fm.ParentOf(p.Seller.Root))

		assert.Equal(t, field(0), fm.ParentOf(p.Sku))
		assert.Equal(t, field(0), fm.ParentOf(p.Name))

		assert.Equal(t, []field{p.Seller.Name, p.Seller.Root}, fm.AncestorOf(p.Seller.Name))

		// Check Weights
		assert.Equal(t, 1.0, fm.GetWeight(p.Sku))
		assert.Equal(t, 2.0, fm.GetWeight(p.Seller.Attr.Root))
		assert.Equal(t, 5.0, fm.GetWeight(p.Seller.Root))
		assert.Equal(t, 1.0, fm.GetWeight(p.Seller.Name))
	})
}

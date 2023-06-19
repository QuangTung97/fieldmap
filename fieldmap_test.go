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
	})

	t.Run("simple struct get field name", func(t *testing.T) {
		fm, err := New[simpleData, field]()
		assert.Equal(t, nil, err)

		p := fm.GetMapping()

		assert.Equal(t, "Sku", fm.GetFieldName(p.Sku))
		assert.Equal(t, "Name", fm.GetFieldName(p.Name))
		assert.Equal(t, "ImageURL", fm.GetFieldName(p.ImageURL))

		assert.Equal(t, "ImageURL", fm.GetFullFieldName(p.ImageURL))
	})

	t.Run("complex struct get field name", func(t *testing.T) {
		fm, err := New[productData, field]()
		assert.Equal(t, nil, err)

		p := fm.GetMapping()

		assert.Equal(t, "Sku", fm.GetFieldName(p.Sku))
		assert.Equal(t, "Name", fm.GetFieldName(p.Name))
		assert.Equal(t, "Seller", fm.GetFieldName(p.Seller.Root))
		assert.Equal(t, "ID", fm.GetFieldName(p.Seller.ID))
		assert.Equal(t, "Name", fm.GetFieldName(p.Seller.Attr.Name))

		assert.Equal(t, "Seller.ID", fm.GetFullFieldName(p.Seller.ID))
		assert.Equal(t, "Seller.Attr.Code", fm.GetFullFieldName(p.Seller.Attr.Code))
	})
}

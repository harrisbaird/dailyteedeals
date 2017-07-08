// This file is generated by SQLBoiler (https://github.com/vattle/sqlboiler)
// and is meant to be re-generated in place and/or deleted at any time.
// DO NOT EDIT

package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testProducts(t *testing.T) {
	t.Parallel()

	query := Products(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testProductsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = product.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Products(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testProductsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Products(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Products(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testProductsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ProductSlice{product}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Products(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testProductsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := ProductExists(tx, product.ID)
	if err != nil {
		t.Errorf("Unable to check if Product exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ProductExistsG to return true, but got false.")
	}
}
func testProductsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx); err != nil {
		t.Error(err)
	}

	productFound, err := FindProduct(tx, product.ID)
	if err != nil {
		t.Error(err)
	}

	if productFound == nil {
		t.Error("want a record, got nil")
	}
}
func testProductsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Products(tx).Bind(product); err != nil {
		t.Error(err)
	}
}

func testProductsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Products(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testProductsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	productOne := &Product{}
	productTwo := &Product{}
	if err = randomize.Struct(seed, productOne, productDBTypes, false, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}
	if err = randomize.Struct(seed, productTwo, productDBTypes, false, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = productOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = productTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Products(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testProductsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	productOne := &Product{}
	productTwo := &Product{}
	if err = randomize.Struct(seed, productOne, productDBTypes, false, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}
	if err = randomize.Struct(seed, productTwo, productDBTypes, false, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = productOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = productTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Products(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
func productBeforeInsertHook(e boil.Executor, o *Product) error {
	*o = Product{}
	return nil
}

func productAfterInsertHook(e boil.Executor, o *Product) error {
	*o = Product{}
	return nil
}

func productAfterSelectHook(e boil.Executor, o *Product) error {
	*o = Product{}
	return nil
}

func productBeforeUpdateHook(e boil.Executor, o *Product) error {
	*o = Product{}
	return nil
}

func productAfterUpdateHook(e boil.Executor, o *Product) error {
	*o = Product{}
	return nil
}

func productBeforeDeleteHook(e boil.Executor, o *Product) error {
	*o = Product{}
	return nil
}

func productAfterDeleteHook(e boil.Executor, o *Product) error {
	*o = Product{}
	return nil
}

func productBeforeUpsertHook(e boil.Executor, o *Product) error {
	*o = Product{}
	return nil
}

func productAfterUpsertHook(e boil.Executor, o *Product) error {
	*o = Product{}
	return nil
}

func testProductsHooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &Product{}
	o := &Product{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, productDBTypes, false); err != nil {
		t.Errorf("Unable to randomize Product object: %s", err)
	}

	AddProductHook(boil.BeforeInsertHook, productBeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	productBeforeInsertHooks = []ProductHook{}

	AddProductHook(boil.AfterInsertHook, productAfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	productAfterInsertHooks = []ProductHook{}

	AddProductHook(boil.AfterSelectHook, productAfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	productAfterSelectHooks = []ProductHook{}

	AddProductHook(boil.BeforeUpdateHook, productBeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	productBeforeUpdateHooks = []ProductHook{}

	AddProductHook(boil.AfterUpdateHook, productAfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	productAfterUpdateHooks = []ProductHook{}

	AddProductHook(boil.BeforeDeleteHook, productBeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	productBeforeDeleteHooks = []ProductHook{}

	AddProductHook(boil.AfterDeleteHook, productAfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	productAfterDeleteHooks = []ProductHook{}

	AddProductHook(boil.BeforeUpsertHook, productBeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	productBeforeUpsertHooks = []ProductHook{}

	AddProductHook(boil.AfterUpsertHook, productAfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	productAfterUpsertHooks = []ProductHook{}
}
func testProductsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Products(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testProductsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx, productColumnsWithoutDefault...); err != nil {
		t.Error(err)
	}

	count, err := Products(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testProductToOneDesignUsingDesign(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Product
	var foreign Design

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, productDBTypes, false, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, designDBTypes, false, designColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Design struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.DesignID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Design(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ProductSlice{&local}
	if err = local.L.LoadDesign(tx, false, (*[]*Product)(&slice)); err != nil {
		t.Fatal(err)
	}
	if local.R.Design == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Design = nil
	if err = local.L.LoadDesign(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Design == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testProductToOneSiteUsingSite(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Product
	var foreign Site

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, productDBTypes, false, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, siteDBTypes, false, siteColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Site struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.SiteID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Site(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ProductSlice{&local}
	if err = local.L.LoadSite(tx, false, (*[]*Product)(&slice)); err != nil {
		t.Fatal(err)
	}
	if local.R.Site == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Site = nil
	if err = local.L.LoadSite(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Site == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testProductToOneSetOpDesignUsingDesign(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Product
	var b, c Design

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, productDBTypes, false, strmangle.SetComplement(productPrimaryKeyColumns, productColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, designDBTypes, false, strmangle.SetComplement(designPrimaryKeyColumns, designColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, designDBTypes, false, strmangle.SetComplement(designPrimaryKeyColumns, designColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Design{&b, &c} {
		err = a.SetDesign(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Design != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Products[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.DesignID != x.ID {
			t.Error("foreign key was wrong value", a.DesignID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.DesignID))
		reflect.Indirect(reflect.ValueOf(&a.DesignID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.DesignID != x.ID {
			t.Error("foreign key was wrong value", a.DesignID, x.ID)
		}
	}
}
func testProductToOneSetOpSiteUsingSite(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Product
	var b, c Site

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, productDBTypes, false, strmangle.SetComplement(productPrimaryKeyColumns, productColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, siteDBTypes, false, strmangle.SetComplement(sitePrimaryKeyColumns, siteColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, siteDBTypes, false, strmangle.SetComplement(sitePrimaryKeyColumns, siteColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Site{&b, &c} {
		err = a.SetSite(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Site != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Products[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.SiteID != x.ID {
			t.Error("foreign key was wrong value", a.SiteID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.SiteID))
		reflect.Indirect(reflect.ValueOf(&a.SiteID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.SiteID != x.ID {
			t.Error("foreign key was wrong value", a.SiteID, x.ID)
		}
	}
}
func testProductsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = product.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testProductsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := ProductSlice{product}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testProductsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Products(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	productDBTypes = map[string]string{`Active`: `boolean`, `ActiveAt`: `timestamp without time zone`, `Deal`: `boolean`, `DesignID`: `integer`, `ExpiresAt`: `timestamp without time zone`, `ID`: `integer`, `ImageBackground`: `text`, `ImageUpdatedAt`: `date`, `LastChance`: `boolean`, `Prices`: `hstore`, `SiteID`: `integer`, `Slug`: `text`, `Tags`: `ARRAYtext`, `URL`: `text`}
	_              = bytes.MinRead
)

func testProductsUpdate(t *testing.T) {
	t.Parallel()

	if len(productColumns) == len(productPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Products(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	if err = product.Update(tx); err != nil {
		t.Error(err)
	}
}

func testProductsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(productColumns) == len(productPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	product := &Product{}
	if err = randomize.Struct(seed, product, productDBTypes, true, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Products(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, product, productDBTypes, true, productPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(productColumns, productPrimaryKeyColumns) {
		fields = productColumns
	} else {
		fields = strmangle.SetComplement(
			productColumns,
			productPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(product))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := ProductSlice{product}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testProductsUpsert(t *testing.T) {
	t.Parallel()

	if len(productColumns) == len(productPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	product := Product{}
	if err = randomize.Struct(seed, &product, productDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = product.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Product: %s", err)
	}

	count, err := Products(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &product, productDBTypes, false, productPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	if err = product.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Product: %s", err)
	}

	count, err = Products(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

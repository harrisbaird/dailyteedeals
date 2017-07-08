// This file is generated by SQLBoiler (https://github.com/vattle/sqlboiler)
// and is meant to be re-generated in place and/or deleted at any time.
// DO NOT EDIT

package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries"
	"github.com/vattle/sqlboiler/queries/qm"
	"github.com/vattle/sqlboiler/strmangle"
	"gopkg.in/nullbio/null.v6"
)

// Site is an object representing the database table.
type Site struct {
	ID           int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name         string      `boil:"name" json:"name" toml:"name" yaml:"name"`
	Slug         string      `boil:"slug" json:"slug" toml:"slug" yaml:"slug"`
	DomainName   string      `boil:"domain_name" json:"domain_name" toml:"domain_name" yaml:"domain_name"`
	AffiliateURL null.String `boil:"affiliate_url" json:"affiliate_url,omitempty" toml:"affiliate_url" yaml:"affiliate_url,omitempty"`
	DealScraper  bool        `boil:"deal_scraper" json:"deal_scraper" toml:"deal_scraper" yaml:"deal_scraper"`
	FullScraper  bool        `boil:"full_scraper" json:"full_scraper" toml:"full_scraper" yaml:"full_scraper"`
	Active       bool        `boil:"active" json:"active" toml:"active" yaml:"active"`
	DisplayOrder int         `boil:"display_order" json:"display_order" toml:"display_order" yaml:"display_order"`

	R *siteR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L siteL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

// siteR is where relationships are stored.
type siteR struct {
	Products ProductSlice
}

// siteL is where Load methods for each relationship are stored.
type siteL struct{}

var (
	siteColumns               = []string{"id", "name", "slug", "domain_name", "affiliate_url", "deal_scraper", "full_scraper", "active", "display_order"}
	siteColumnsWithoutDefault = []string{"id", "name", "slug", "domain_name", "affiliate_url"}
	siteColumnsWithDefault    = []string{"deal_scraper", "full_scraper", "active", "display_order"}
	sitePrimaryKeyColumns     = []string{"id"}
)

type (
	// SiteSlice is an alias for a slice of pointers to Site.
	// This should generally be used opposed to []Site.
	SiteSlice []*Site
	// SiteHook is the signature for custom Site hook methods
	SiteHook func(boil.Executor, *Site) error

	siteQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	siteType                 = reflect.TypeOf(&Site{})
	siteMapping              = queries.MakeStructMapping(siteType)
	sitePrimaryKeyMapping, _ = queries.BindMapping(siteType, siteMapping, sitePrimaryKeyColumns)
	siteInsertCacheMut       sync.RWMutex
	siteInsertCache          = make(map[string]insertCache)
	siteUpdateCacheMut       sync.RWMutex
	siteUpdateCache          = make(map[string]updateCache)
	siteUpsertCacheMut       sync.RWMutex
	siteUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)
var siteBeforeInsertHooks []SiteHook
var siteBeforeUpdateHooks []SiteHook
var siteBeforeDeleteHooks []SiteHook
var siteBeforeUpsertHooks []SiteHook

var siteAfterInsertHooks []SiteHook
var siteAfterSelectHooks []SiteHook
var siteAfterUpdateHooks []SiteHook
var siteAfterDeleteHooks []SiteHook
var siteAfterUpsertHooks []SiteHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Site) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range siteBeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Site) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range siteBeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Site) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range siteBeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Site) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range siteBeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Site) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range siteAfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Site) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range siteAfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Site) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range siteAfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Site) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range siteAfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Site) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range siteAfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddSiteHook registers your hook function for all future operations.
func AddSiteHook(hookPoint boil.HookPoint, siteHook SiteHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		siteBeforeInsertHooks = append(siteBeforeInsertHooks, siteHook)
	case boil.BeforeUpdateHook:
		siteBeforeUpdateHooks = append(siteBeforeUpdateHooks, siteHook)
	case boil.BeforeDeleteHook:
		siteBeforeDeleteHooks = append(siteBeforeDeleteHooks, siteHook)
	case boil.BeforeUpsertHook:
		siteBeforeUpsertHooks = append(siteBeforeUpsertHooks, siteHook)
	case boil.AfterInsertHook:
		siteAfterInsertHooks = append(siteAfterInsertHooks, siteHook)
	case boil.AfterSelectHook:
		siteAfterSelectHooks = append(siteAfterSelectHooks, siteHook)
	case boil.AfterUpdateHook:
		siteAfterUpdateHooks = append(siteAfterUpdateHooks, siteHook)
	case boil.AfterDeleteHook:
		siteAfterDeleteHooks = append(siteAfterDeleteHooks, siteHook)
	case boil.AfterUpsertHook:
		siteAfterUpsertHooks = append(siteAfterUpsertHooks, siteHook)
	}
}

// OneP returns a single site record from the query, and panics on error.
func (q siteQuery) OneP() *Site {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single site record from the query.
func (q siteQuery) One() (*Site, error) {
	o := &Site{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for sites")
	}

	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}

	return o, nil
}

// AllP returns all Site records from the query, and panics on error.
func (q siteQuery) AllP() SiteSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all Site records from the query.
func (q siteQuery) All() (SiteSlice, error) {
	var o []*Site

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Site slice")
	}

	if len(siteAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// CountP returns the count of all Site records in the query, and panics on error.
func (q siteQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all Site records in the query.
func (q siteQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count sites rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q siteQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q siteQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if sites exists")
	}

	return count > 0, nil
}

// ProductsG retrieves all the product's products.
func (o *Site) ProductsG(mods ...qm.QueryMod) productQuery {
	return o.Products(boil.GetDB(), mods...)
}

// Products retrieves all the product's products with an executor.
func (o *Site) Products(exec boil.Executor, mods ...qm.QueryMod) productQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"products\".\"site_id\"=?", o.ID),
	)

	query := Products(exec, queryMods...)
	queries.SetFrom(query.Query, "\"products\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"products\".*"})
	}

	return query
}

// LoadProducts allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (siteL) LoadProducts(e boil.Executor, singular bool, maybeSite interface{}) error {
	var slice []*Site
	var object *Site

	count := 1
	if singular {
		object = maybeSite.(*Site)
	} else {
		slice = *maybeSite.(*[]*Site)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &siteR{}
		}
		args[0] = object.ID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &siteR{}
			}
			args[i] = obj.ID
		}
	}

	query := fmt.Sprintf(
		"select * from \"products\" where \"site_id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)
	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load products")
	}
	defer results.Close()

	var resultSlice []*Product
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice products")
	}

	if len(productAfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}
	if singular {
		object.R.Products = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.SiteID {
				local.R.Products = append(local.R.Products, foreign)
				break
			}
		}
	}

	return nil
}

// AddProductsG adds the given related objects to the existing relationships
// of the site, optionally inserting them as new records.
// Appends related to o.R.Products.
// Sets related.R.Site appropriately.
// Uses the global database handle.
func (o *Site) AddProductsG(insert bool, related ...*Product) error {
	return o.AddProducts(boil.GetDB(), insert, related...)
}

// AddProductsP adds the given related objects to the existing relationships
// of the site, optionally inserting them as new records.
// Appends related to o.R.Products.
// Sets related.R.Site appropriately.
// Panics on error.
func (o *Site) AddProductsP(exec boil.Executor, insert bool, related ...*Product) {
	if err := o.AddProducts(exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddProductsGP adds the given related objects to the existing relationships
// of the site, optionally inserting them as new records.
// Appends related to o.R.Products.
// Sets related.R.Site appropriately.
// Uses the global database handle and panics on error.
func (o *Site) AddProductsGP(insert bool, related ...*Product) {
	if err := o.AddProducts(boil.GetDB(), insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// AddProducts adds the given related objects to the existing relationships
// of the site, optionally inserting them as new records.
// Appends related to o.R.Products.
// Sets related.R.Site appropriately.
func (o *Site) AddProducts(exec boil.Executor, insert bool, related ...*Product) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.SiteID = o.ID
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"products\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"site_id"}),
				strmangle.WhereClause("\"", "\"", 2, productPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.SiteID = o.ID
		}
	}

	if o.R == nil {
		o.R = &siteR{
			Products: related,
		}
	} else {
		o.R.Products = append(o.R.Products, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &productR{
				Site: o,
			}
		} else {
			rel.R.Site = o
		}
	}
	return nil
}

// SitesG retrieves all records.
func SitesG(mods ...qm.QueryMod) siteQuery {
	return Sites(boil.GetDB(), mods...)
}

// Sites retrieves all the records using an executor.
func Sites(exec boil.Executor, mods ...qm.QueryMod) siteQuery {
	mods = append(mods, qm.From("\"sites\""))
	return siteQuery{NewQuery(exec, mods...)}
}

// FindSiteG retrieves a single record by ID.
func FindSiteG(id int, selectCols ...string) (*Site, error) {
	return FindSite(boil.GetDB(), id, selectCols...)
}

// FindSiteGP retrieves a single record by ID, and panics on error.
func FindSiteGP(id int, selectCols ...string) *Site {
	retobj, err := FindSite(boil.GetDB(), id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindSite retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSite(exec boil.Executor, id int, selectCols ...string) (*Site, error) {
	siteObj := &Site{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"sites\" where \"id\"=$1", sel,
	)

	q := queries.Raw(exec, query, id)

	err := q.Bind(siteObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from sites")
	}

	return siteObj, nil
}

// FindSiteP retrieves a single record by ID with an executor, and panics on error.
func FindSiteP(exec boil.Executor, id int, selectCols ...string) *Site {
	retobj, err := FindSite(exec, id, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *Site) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *Site) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *Site) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *Site) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no sites provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(siteColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	siteInsertCacheMut.RLock()
	cache, cached := siteInsertCache[key]
	siteInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			siteColumns,
			siteColumnsWithDefault,
			siteColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(siteType, siteMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(siteType, siteMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"sites\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"sites\" DEFAULT VALUES"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		if len(wl) != 0 {
			cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into sites")
	}

	if !cached {
		siteInsertCacheMut.Lock()
		siteInsertCache[key] = cache
		siteInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(exec)
}

// UpdateG a single Site record. See Update for
// whitelist behavior description.
func (o *Site) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single Site record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *Site) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the Site, and panics on error.
// See Update for whitelist behavior description.
func (o *Site) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the Site.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *Site) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return err
	}
	key := makeCacheKey(whitelist, nil)
	siteUpdateCacheMut.RLock()
	cache, cached := siteUpdateCache[key]
	siteUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(
			siteColumns,
			sitePrimaryKeyColumns,
			whitelist,
		)

		if len(whitelist) == 0 {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return errors.New("models: unable to update sites, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"sites\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, sitePrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(siteType, siteMapping, append(wl, sitePrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update sites row")
	}

	if !cached {
		siteUpdateCacheMut.Lock()
		siteUpdateCache[key] = cache
		siteUpdateCacheMut.Unlock()
	}

	return o.doAfterUpdateHooks(exec)
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q siteQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q siteQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for sites")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o SiteSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o SiteSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o SiteSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SiteSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sitePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"sites\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, sitePrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in site slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *Site) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *Site) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *Site) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *Site) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no sites provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(siteColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
	buf := strmangle.GetBuffer()

	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	siteUpsertCacheMut.RLock()
	cache, cached := siteUpsertCache[key]
	siteUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := strmangle.InsertColumnSet(
			siteColumns,
			siteColumnsWithDefault,
			siteColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		update := strmangle.UpdateColumnSet(
			siteColumns,
			sitePrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert sites, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(sitePrimaryKeyColumns))
			copy(conflict, sitePrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"sites\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(siteType, siteMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(siteType, siteMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert sites")
	}

	if !cached {
		siteUpsertCacheMut.Lock()
		siteUpsertCache[key] = cache
		siteUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(exec)
}

// DeleteP deletes a single Site record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Site) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single Site record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *Site) DeleteG() error {
	if o == nil {
		return errors.New("models: no Site provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single Site record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *Site) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single Site record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Site) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Site provided for delete")
	}

	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), sitePrimaryKeyMapping)
	sql := "DELETE FROM \"sites\" WHERE \"id\"=$1"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from sites")
	}

	if err := o.doAfterDeleteHooks(exec); err != nil {
		return err
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q siteQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q siteQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no siteQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from sites")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o SiteSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o SiteSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no Site slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o SiteSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SiteSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no Site slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	if len(siteBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sitePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"sites\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sitePrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from site slice")
	}

	if len(siteAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *Site) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *Site) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *Site) ReloadG() error {
	if o == nil {
		return errors.New("models: no Site provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Site) Reload(exec boil.Executor) error {
	ret, err := FindSite(exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *SiteSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *SiteSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SiteSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty SiteSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SiteSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	sites := SiteSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), sitePrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"sites\".* FROM \"sites\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, sitePrimaryKeyColumns, len(*o))

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&sites)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in SiteSlice")
	}

	*o = sites

	return nil
}

// SiteExists checks if the Site row exists.
func SiteExists(exec boil.Executor, id int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"sites\" where \"id\"=$1 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, id)
	}

	row := exec.QueryRow(sql, id)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if sites exists")
	}

	return exists, nil
}

// SiteExistsG checks if the Site row exists.
func SiteExistsG(id int) (bool, error) {
	return SiteExists(boil.GetDB(), id)
}

// SiteExistsGP checks if the Site row exists. Panics on error.
func SiteExistsGP(id int) bool {
	e, err := SiteExists(boil.GetDB(), id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// SiteExistsP checks if the Site row exists. Panics on error.
func SiteExistsP(exec boil.Executor, id int) bool {
	e, err := SiteExists(exec, id)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

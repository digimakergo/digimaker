// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package entity

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Folder is an object representing the database table.
type Folder struct {
	DataID    int         `boil:"id" json:"id" toml:"id" yaml:"id"`
	Title     null.String `boil:"title" json:"title,omitempty" toml:"title" yaml:"title,omitempty"`
	Summary   null.String `boil:"summary" json:"summary,omitempty" toml:"summary" yaml:"summary,omitempty"`
	Published null.Int    `boil:"published" json:"published,omitempty" toml:"published" yaml:"published,omitempty"`
	Modified  null.Int    `boil:"modified" json:"modified,omitempty" toml:"modified" yaml:"modified,omitempty"`
	RemoteID  null.String `boil:"remote_id" json:"remote_id,omitempty" toml:"remote_id" yaml:"remote_id,omitempty"`

	R *folderR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L folderL  `boil:"-" json:"-" toml:"-" yaml:"-"`
	*Location
}

func (c *Folder) Fields() map[string]Folder {
	return nil
}

func (c *Folder) Field(name string) interface{} {
	var result = nil
	switch name {
	case "id", "DataID":
		result = c.DataID
	case "title", "Title":
		result = c.Title
	case "summary", "Summary":
		result = c.Summary
	case "published", "Published":
		result = c.Published
	case "modified", "Modified":
		result = c.Modified
	case "remote_id", "RemoteID":
		result = c.RemoteID
	default:
	}
	return result
}

var FolderColumns = struct {
	DataID    string
	Title     string
	Summary   string
	Published string
	Modified  string
	RemoteID  string
}{
	DataID:    "id",
	Title:     "title",
	Summary:   "summary",
	Published: "published",
	Modified:  "modified",
	RemoteID:  "remote_id",
}

// Generated where

var FolderWhere = struct {
	DataID    whereHelperint
	Title     whereHelpernull_String
	Summary   whereHelpernull_String
	Published whereHelpernull_Int
	Modified  whereHelpernull_Int
	RemoteID  whereHelpernull_String
}{
	DataID:    whereHelperint{field: `id`},
	Title:     whereHelpernull_String{field: `title`},
	Summary:   whereHelpernull_String{field: `summary`},
	Published: whereHelpernull_Int{field: `published`},
	Modified:  whereHelpernull_Int{field: `modified`},
	RemoteID:  whereHelpernull_String{field: `remote_id`},
}

// FolderRels is where relationship names are stored.
var FolderRels = struct {
}{}

// folderR is where relationships are stored.
type folderR struct {
}

// NewStruct creates a new relationship struct
func (*folderR) NewStruct() *folderR {
	return &folderR{}
}

// folderL is where Load methods for each relationship are stored.
type folderL struct{}

var (
	folderColumns               = []string{"id", "title", "summary", "published", "modified", "remote_id"}
	folderColumnsWithoutDefault = []string{"title", "summary", "published", "modified", "remote_id"}
	folderColumnsWithDefault    = []string{"id"}
	folderPrimaryKeyColumns     = []string{"id"}
)

type (
	// FolderSlice is an alias for a slice of pointers to Folder.
	// This should generally be used opposed to []Folder.
	FolderSlice []*Folder
	// FolderHook is the signature for custom Folder hook methods
	FolderHook func(context.Context, boil.ContextExecutor, *Folder) error

	folderQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	folderType                 = reflect.TypeOf(&Folder{})
	folderMapping              = queries.MakeStructMapping(folderType)
	folderPrimaryKeyMapping, _ = queries.BindMapping(folderType, folderMapping, folderPrimaryKeyColumns)
	folderInsertCacheMut       sync.RWMutex
	folderInsertCache          = make(map[string]insertCache)
	folderUpdateCacheMut       sync.RWMutex
	folderUpdateCache          = make(map[string]updateCache)
	folderUpsertCacheMut       sync.RWMutex
	folderUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

var folderBeforeInsertHooks []FolderHook
var folderBeforeUpdateHooks []FolderHook
var folderBeforeDeleteHooks []FolderHook
var folderBeforeUpsertHooks []FolderHook

var folderAfterInsertHooks []FolderHook
var folderAfterSelectHooks []FolderHook
var folderAfterUpdateHooks []FolderHook
var folderAfterDeleteHooks []FolderHook
var folderAfterUpsertHooks []FolderHook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *Folder) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range folderBeforeInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *Folder) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range folderBeforeUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *Folder) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range folderBeforeDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *Folder) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range folderBeforeUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *Folder) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range folderAfterInsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *Folder) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range folderAfterSelectHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *Folder) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range folderAfterUpdateHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *Folder) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range folderAfterDeleteHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *Folder) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range folderAfterUpsertHooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

// AddFolderHook registers your hook function for all future operations.
func AddFolderHook(hookPoint boil.HookPoint, folderHook FolderHook) {
	switch hookPoint {
	case boil.BeforeInsertHook:
		folderBeforeInsertHooks = append(folderBeforeInsertHooks, folderHook)
	case boil.BeforeUpdateHook:
		folderBeforeUpdateHooks = append(folderBeforeUpdateHooks, folderHook)
	case boil.BeforeDeleteHook:
		folderBeforeDeleteHooks = append(folderBeforeDeleteHooks, folderHook)
	case boil.BeforeUpsertHook:
		folderBeforeUpsertHooks = append(folderBeforeUpsertHooks, folderHook)
	case boil.AfterInsertHook:
		folderAfterInsertHooks = append(folderAfterInsertHooks, folderHook)
	case boil.AfterSelectHook:
		folderAfterSelectHooks = append(folderAfterSelectHooks, folderHook)
	case boil.AfterUpdateHook:
		folderAfterUpdateHooks = append(folderAfterUpdateHooks, folderHook)
	case boil.AfterDeleteHook:
		folderAfterDeleteHooks = append(folderAfterDeleteHooks, folderHook)
	case boil.AfterUpsertHook:
		folderAfterUpsertHooks = append(folderAfterUpsertHooks, folderHook)
	}
}

// One returns a single folder record from the query.
func (q folderQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Folder, error) {
	o := &Folder{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "entity: failed to execute a one query for dm_folder")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

// All returns all Folder records from the query.
func (q folderQuery) All(ctx context.Context, exec boil.ContextExecutor) (FolderSlice, error) {
	var o []*Folder

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "entity: failed to assign all query results to Folder slice")
	}

	if len(folderAfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx, exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

// Count returns the count of all Folder records in the query.
func (q folderQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "entity: failed to count dm_folder rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q folderQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "entity: failed to check if dm_folder exists")
	}

	return count > 0, nil
}

// Folders retrieves all the records using an executor.
func Folders(mods ...qm.QueryMod) folderQuery {
	mods = append(mods, qm.From("`dm_folder`"))
	return folderQuery{NewQuery(mods...)}
}

// FindFolder retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindFolder(ctx context.Context, exec boil.ContextExecutor, dataID int, selectCols ...string) (*Folder, error) {
	folderObj := &Folder{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from `dm_folder` where `id`=?", sel,
	)

	q := queries.Raw(query, dataID)

	err := q.Bind(ctx, exec, folderObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "entity: unable to select from dm_folder")
	}

	return folderObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Folder) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("entity: no dm_folder provided for insertion")
	}

	var err error

	if err := o.doBeforeInsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(folderColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	folderInsertCacheMut.RLock()
	cache, cached := folderInsertCache[key]
	folderInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			folderColumns,
			folderColumnsWithDefault,
			folderColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(folderType, folderMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(folderType, folderMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO `dm_folder` (`%s`) %%sVALUES (%s)%%s", strings.Join(wl, "`,`"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO `dm_folder` () VALUES ()%s%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			cache.retQuery = fmt.Sprintf("SELECT `%s` FROM `dm_folder` WHERE %s", strings.Join(returnColumns, "`,`"), strmangle.WhereClause("`", "`", 0, folderPrimaryKeyColumns))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "entity: unable to insert into dm_folder")
	}

	var lastID int64
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.ID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == folderMapping["ID"] {
		goto CacheNoHooks
	}

	identifierCols = []interface{}{
		o.DataID,
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, identifierCols...)
	}

	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "entity: unable to populate default values for dm_folder")
	}

CacheNoHooks:
	if !cached {
		folderInsertCacheMut.Lock()
		folderInsertCache[key] = cache
		folderInsertCacheMut.Unlock()
	}

	return o.doAfterInsertHooks(ctx, exec)
}

// Update uses an executor to update the Folder.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Folder) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
	}
	key := makeCacheKey(columns, nil)
	folderUpdateCacheMut.RLock()
	cache, cached := folderUpdateCache[key]
	folderUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			folderColumns,
			folderPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("entity: unable to update dm_folder, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE `dm_folder` SET %s WHERE %s",
			strmangle.SetParamNames("`", "`", 0, wl),
			strmangle.WhereClause("`", "`", 0, folderPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(folderType, folderMapping, append(wl, folderPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "entity: unable to update dm_folder row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entity: failed to get rows affected by update for dm_folder")
	}

	if !cached {
		folderUpdateCacheMut.Lock()
		folderUpdateCache[key] = cache
		folderUpdateCacheMut.Unlock()
	}

	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
}

// UpdateAll updates all rows with the specified column values.
func (q folderQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "entity: unable to update all for dm_folder")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entity: unable to retrieve rows affected for dm_folder")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o FolderSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("entity: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), folderPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE `dm_folder` SET %s WHERE %s",
		strmangle.SetParamNames("`", "`", 0, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, folderPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "entity: unable to update all in folder slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entity: unable to retrieve rows affected all in update all folder")
	}
	return rowsAff, nil
}

var mySQLFolderUniqueColumns = []string{
	"id",
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Folder) Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("entity: no dm_folder provided for upsert")
	}

	if err := o.doBeforeUpsertHooks(ctx, exec); err != nil {
		return err
	}

	nzDefaults := queries.NonZeroDefaultSet(folderColumnsWithDefault, o)
	nzUniques := queries.NonZeroDefaultSet(mySQLFolderUniqueColumns, o)

	if len(nzUniques) == 0 {
		return errors.New("cannot upsert with a table that cannot conflict on a unique column")
	}

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzUniques {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	folderUpsertCacheMut.RLock()
	cache, cached := folderUpsertCache[key]
	folderUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			folderColumns,
			folderColumnsWithDefault,
			folderColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			folderColumns,
			folderPrimaryKeyColumns,
		)

		if len(update) == 0 {
			return errors.New("entity: unable to upsert dm_folder, could not build update column list")
		}

		ret = strmangle.SetComplement(ret, nzUniques)
		cache.query = buildUpsertQueryMySQL(dialect, "dm_folder", update, insert)
		cache.retQuery = fmt.Sprintf(
			"SELECT %s FROM `dm_folder` WHERE %s",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
			strmangle.WhereClause("`", "`", 0, nzUniques),
		)

		cache.valueMapping, err = queries.BindMapping(folderType, folderMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(folderType, folderMapping, ret)
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

	result, err := exec.ExecContext(ctx, cache.query, vals...)

	if err != nil {
		return errors.Wrap(err, "entity: unable to upsert for dm_folder")
	}

	var lastID int64
	var uniqueMap []uint64
	var nzUniqueCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	o.DataID = int(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == folderMapping["id"] {
		goto CacheNoHooks
	}

	uniqueMap, err = queries.BindMapping(folderType, folderMapping, nzUniques)
	if err != nil {
		return errors.Wrap(err, "entity: unable to retrieve unique values for dm_folder")
	}
	nzUniqueCols = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), uniqueMap)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, nzUniqueCols...)
	}

	err = exec.QueryRowContext(ctx, cache.retQuery, nzUniqueCols...).Scan(returns...)
	if err != nil {
		return errors.Wrap(err, "entity: unable to populate default values for dm_folder")
	}

CacheNoHooks:
	if !cached {
		folderUpsertCacheMut.Lock()
		folderUpsertCache[key] = cache
		folderUpsertCacheMut.Unlock()
	}

	return o.doAfterUpsertHooks(ctx, exec)
}

// Delete deletes a single Folder record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Folder) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("entity: no Folder provided for delete")
	}

	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), folderPrimaryKeyMapping)
	sql := "DELETE FROM `dm_folder` WHERE `id`=?"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "entity: unable to delete from dm_folder")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entity: failed to get rows affected by delete for dm_folder")
	}

	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q folderQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("entity: no folderQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "entity: unable to delete all from dm_folder")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entity: failed to get rows affected by deleteall for dm_folder")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o FolderSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("entity: no Folder slice provided for delete all")
	}

	if len(o) == 0 {
		return 0, nil
	}

	if len(folderBeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), folderPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM `dm_folder` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, folderPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "entity: unable to delete all from folder slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "entity: failed to get rows affected by deleteall for dm_folder")
	}

	if len(folderAfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Folder) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindFolder(ctx, exec, o.DataID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *FolderSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := FolderSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), folderPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT `dm_folder`.* FROM `dm_folder` WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 0, folderPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "entity: unable to reload all in FolderSlice")
	}

	*o = slice

	return nil
}

// FolderExists checks if the Folder row exists.
func FolderExists(ctx context.Context, exec boil.ContextExecutor, dataID int) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from `dm_folder` where `id`=? limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, dataID)
	}

	row := exec.QueryRowContext(ctx, sql, dataID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "entity: unable to check if dm_folder exists")
	}

	return exists, nil
}

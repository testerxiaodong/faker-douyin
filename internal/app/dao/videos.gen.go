// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"faker-douyin/internal/app/model/entity"
)

func newVideo(db *gorm.DB, opts ...gen.DOOption) video {
	_video := video{}

	_video.videoDo.UseDB(db, opts...)
	_video.videoDo.UseModel(&entity.Video{})

	tableName := _video.videoDo.TableName()
	_video.ALL = field.NewAsterisk(tableName)
	_video.ID = field.NewInt64(tableName, "id")
	_video.CreatedAt = field.NewTime(tableName, "created_at")
	_video.UpdatedAt = field.NewTime(tableName, "updated_at")
	_video.DeletedAt = field.NewField(tableName, "deleted_at")
	_video.AuthorID = field.NewInt64(tableName, "author_id")
	_video.Title = field.NewString(tableName, "title")
	_video.PlayURL = field.NewString(tableName, "play_url")
	_video.CoverURL = field.NewString(tableName, "cover_url")

	_video.fillFieldMap()

	return _video
}

type video struct {
	videoDo

	ALL       field.Asterisk
	ID        field.Int64
	CreatedAt field.Time
	UpdatedAt field.Time
	DeletedAt field.Field
	AuthorID  field.Int64
	Title     field.String
	PlayURL   field.String
	CoverURL  field.String

	fieldMap map[string]field.Expr
}

func (v video) Table(newTableName string) *video {
	v.videoDo.UseTable(newTableName)
	return v.updateTableName(newTableName)
}

func (v video) As(alias string) *video {
	v.videoDo.DO = *(v.videoDo.As(alias).(*gen.DO))
	return v.updateTableName(alias)
}

func (v *video) updateTableName(table string) *video {
	v.ALL = field.NewAsterisk(table)
	v.ID = field.NewInt64(table, "id")
	v.CreatedAt = field.NewTime(table, "created_at")
	v.UpdatedAt = field.NewTime(table, "updated_at")
	v.DeletedAt = field.NewField(table, "deleted_at")
	v.AuthorID = field.NewInt64(table, "author_id")
	v.Title = field.NewString(table, "title")
	v.PlayURL = field.NewString(table, "play_url")
	v.CoverURL = field.NewString(table, "cover_url")

	v.fillFieldMap()

	return v
}

func (v *video) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := v.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (v *video) fillFieldMap() {
	v.fieldMap = make(map[string]field.Expr, 8)
	v.fieldMap["id"] = v.ID
	v.fieldMap["created_at"] = v.CreatedAt
	v.fieldMap["updated_at"] = v.UpdatedAt
	v.fieldMap["deleted_at"] = v.DeletedAt
	v.fieldMap["author_id"] = v.AuthorID
	v.fieldMap["title"] = v.Title
	v.fieldMap["play_url"] = v.PlayURL
	v.fieldMap["cover_url"] = v.CoverURL
}

func (v video) clone(db *gorm.DB) video {
	v.videoDo.ReplaceConnPool(db.Statement.ConnPool)
	return v
}

func (v video) replaceDB(db *gorm.DB) video {
	v.videoDo.ReplaceDB(db)
	return v
}

type videoDo struct{ gen.DO }

type IVideoDo interface {
	gen.SubQuery
	Debug() IVideoDo
	WithContext(ctx context.Context) IVideoDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IVideoDo
	WriteDB() IVideoDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IVideoDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IVideoDo
	Not(conds ...gen.Condition) IVideoDo
	Or(conds ...gen.Condition) IVideoDo
	Select(conds ...field.Expr) IVideoDo
	Where(conds ...gen.Condition) IVideoDo
	Order(conds ...field.Expr) IVideoDo
	Distinct(cols ...field.Expr) IVideoDo
	Omit(cols ...field.Expr) IVideoDo
	Join(table schema.Tabler, on ...field.Expr) IVideoDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IVideoDo
	RightJoin(table schema.Tabler, on ...field.Expr) IVideoDo
	Group(cols ...field.Expr) IVideoDo
	Having(conds ...gen.Condition) IVideoDo
	Limit(limit int) IVideoDo
	Offset(offset int) IVideoDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IVideoDo
	Unscoped() IVideoDo
	Create(values ...*entity.Video) error
	CreateInBatches(values []*entity.Video, batchSize int) error
	Save(values ...*entity.Video) error
	First() (*entity.Video, error)
	Take() (*entity.Video, error)
	Last() (*entity.Video, error)
	Find() ([]*entity.Video, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*entity.Video, err error)
	FindInBatches(result *[]*entity.Video, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*entity.Video) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IVideoDo
	Assign(attrs ...field.AssignExpr) IVideoDo
	Joins(fields ...field.RelationField) IVideoDo
	Preload(fields ...field.RelationField) IVideoDo
	FirstOrInit() (*entity.Video, error)
	FirstOrCreate() (*entity.Video, error)
	FindByPage(offset int, limit int) (result []*entity.Video, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IVideoDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (v videoDo) Debug() IVideoDo {
	return v.withDO(v.DO.Debug())
}

func (v videoDo) WithContext(ctx context.Context) IVideoDo {
	return v.withDO(v.DO.WithContext(ctx))
}

func (v videoDo) ReadDB() IVideoDo {
	return v.Clauses(dbresolver.Read)
}

func (v videoDo) WriteDB() IVideoDo {
	return v.Clauses(dbresolver.Write)
}

func (v videoDo) Session(config *gorm.Session) IVideoDo {
	return v.withDO(v.DO.Session(config))
}

func (v videoDo) Clauses(conds ...clause.Expression) IVideoDo {
	return v.withDO(v.DO.Clauses(conds...))
}

func (v videoDo) Returning(value interface{}, columns ...string) IVideoDo {
	return v.withDO(v.DO.Returning(value, columns...))
}

func (v videoDo) Not(conds ...gen.Condition) IVideoDo {
	return v.withDO(v.DO.Not(conds...))
}

func (v videoDo) Or(conds ...gen.Condition) IVideoDo {
	return v.withDO(v.DO.Or(conds...))
}

func (v videoDo) Select(conds ...field.Expr) IVideoDo {
	return v.withDO(v.DO.Select(conds...))
}

func (v videoDo) Where(conds ...gen.Condition) IVideoDo {
	return v.withDO(v.DO.Where(conds...))
}

func (v videoDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) IVideoDo {
	return v.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (v videoDo) Order(conds ...field.Expr) IVideoDo {
	return v.withDO(v.DO.Order(conds...))
}

func (v videoDo) Distinct(cols ...field.Expr) IVideoDo {
	return v.withDO(v.DO.Distinct(cols...))
}

func (v videoDo) Omit(cols ...field.Expr) IVideoDo {
	return v.withDO(v.DO.Omit(cols...))
}

func (v videoDo) Join(table schema.Tabler, on ...field.Expr) IVideoDo {
	return v.withDO(v.DO.Join(table, on...))
}

func (v videoDo) LeftJoin(table schema.Tabler, on ...field.Expr) IVideoDo {
	return v.withDO(v.DO.LeftJoin(table, on...))
}

func (v videoDo) RightJoin(table schema.Tabler, on ...field.Expr) IVideoDo {
	return v.withDO(v.DO.RightJoin(table, on...))
}

func (v videoDo) Group(cols ...field.Expr) IVideoDo {
	return v.withDO(v.DO.Group(cols...))
}

func (v videoDo) Having(conds ...gen.Condition) IVideoDo {
	return v.withDO(v.DO.Having(conds...))
}

func (v videoDo) Limit(limit int) IVideoDo {
	return v.withDO(v.DO.Limit(limit))
}

func (v videoDo) Offset(offset int) IVideoDo {
	return v.withDO(v.DO.Offset(offset))
}

func (v videoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IVideoDo {
	return v.withDO(v.DO.Scopes(funcs...))
}

func (v videoDo) Unscoped() IVideoDo {
	return v.withDO(v.DO.Unscoped())
}

func (v videoDo) Create(values ...*entity.Video) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Create(values)
}

func (v videoDo) CreateInBatches(values []*entity.Video, batchSize int) error {
	return v.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (v videoDo) Save(values ...*entity.Video) error {
	if len(values) == 0 {
		return nil
	}
	return v.DO.Save(values)
}

func (v videoDo) First() (*entity.Video, error) {
	if result, err := v.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Video), nil
	}
}

func (v videoDo) Take() (*entity.Video, error) {
	if result, err := v.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Video), nil
	}
}

func (v videoDo) Last() (*entity.Video, error) {
	if result, err := v.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Video), nil
	}
}

func (v videoDo) Find() ([]*entity.Video, error) {
	result, err := v.DO.Find()
	return result.([]*entity.Video), err
}

func (v videoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*entity.Video, err error) {
	buf := make([]*entity.Video, 0, batchSize)
	err = v.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (v videoDo) FindInBatches(result *[]*entity.Video, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return v.DO.FindInBatches(result, batchSize, fc)
}

func (v videoDo) Attrs(attrs ...field.AssignExpr) IVideoDo {
	return v.withDO(v.DO.Attrs(attrs...))
}

func (v videoDo) Assign(attrs ...field.AssignExpr) IVideoDo {
	return v.withDO(v.DO.Assign(attrs...))
}

func (v videoDo) Joins(fields ...field.RelationField) IVideoDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Joins(_f))
	}
	return &v
}

func (v videoDo) Preload(fields ...field.RelationField) IVideoDo {
	for _, _f := range fields {
		v = *v.withDO(v.DO.Preload(_f))
	}
	return &v
}

func (v videoDo) FirstOrInit() (*entity.Video, error) {
	if result, err := v.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Video), nil
	}
}

func (v videoDo) FirstOrCreate() (*entity.Video, error) {
	if result, err := v.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Video), nil
	}
}

func (v videoDo) FindByPage(offset int, limit int) (result []*entity.Video, count int64, err error) {
	result, err = v.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = v.Offset(-1).Limit(-1).Count()
	return
}

func (v videoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = v.Count()
	if err != nil {
		return
	}

	err = v.Offset(offset).Limit(limit).Scan(result)
	return
}

func (v videoDo) Scan(result interface{}) (err error) {
	return v.DO.Scan(result)
}

func (v videoDo) Delete(models ...*entity.Video) (result gen.ResultInfo, err error) {
	return v.DO.Delete(models)
}

func (v *videoDo) withDO(do gen.Dao) *videoDo {
	v.DO = *do.(*gen.DO)
	return v
}

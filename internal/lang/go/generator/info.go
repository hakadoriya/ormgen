package generator

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/hakadoriya/z.go/otelz/tracez"

	"github.com/hakadoriya/ormgen/internal/contexts"
	"github.com/hakadoriya/ormgen/internal/lang/go/source"
)

type TableInfo struct {
	StructName      string
	TableName       string
	Columns         []*ColumnInfo
	PrimaryKeys     []*ColumnInfo
	NotPrimaryKeys  []*ColumnInfo
	HasOneTagsKeys  []string
	HasOneTags      map[string][]*ColumnInfo
	HasManyTagsKeys []string
	HasManyTags     map[string][]*ColumnInfo
}

type ColumnInfo struct {
	FieldName   string
	FieldType   string
	ColumnName  string
	PK          bool
	HasOneTags  []string
	HasManyTags []string
}

//nolint:cyclop,funlen
func BuildTableInfo(ctx context.Context, structSource *source.StructSource) *TableInfo {
	ctx, span := tracez.Start(ctx)
	defer span.End()

	cfg := contexts.GenerateConfig(ctx)

	// NOTE: fill other fields later
	//
	//nolint:exhaustruct
	tableInfo := &TableInfo{
		StructName:     structSource.TypeSpec.Name.Name,
		TableName:      structSource.ExtractTableName(ctx, cfg.GoColumnTag),
		Columns:        make([]*ColumnInfo, 0, len(structSource.StructType.Fields.List)),
		PrimaryKeys:    make([]*ColumnInfo, 0, len(structSource.StructType.Fields.List)),
		NotPrimaryKeys: make([]*ColumnInfo, 0, len(structSource.StructType.Fields.List)),
		HasOneTags:     make(map[string][]*ColumnInfo),
		HasManyTags:    make(map[string][]*ColumnInfo),
	}
	for _, field := range structSource.StructType.Fields.List {
		if field.Tag != nil {
			tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
			// db tag
			switch columnName := tag.Get(cfg.GoColumnTag); columnName {
			case "", "-":
				contexts.Trace(ctx).Debug(fmt.Sprintf("SKIP: %s: field.Names=%s, columnName=%q", structSource.Position.String(), field.Names, columnName))
				// noop
			default:
				contexts.Trace(ctx).Debug(fmt.Sprintf("%s: field.Names=%s, columnName=%q", structSource.Position.String(), field.Names, columnName))
				// NOTE: fill other fields later
				//
				//nolint:exhaustruct
				columnInfo := &ColumnInfo{
					FieldName:  field.Names[0].Name,
					FieldType:  fieldName(field.Type).String(),
					ColumnName: columnName,
				}
				// pk tag
				switch pk := tag.Get(cfg.GoPKTag); pk {
				case "true":
					contexts.Trace(ctx).Debug(fmt.Sprintf("SKIP: %s: field.Names=%s, pk=%q", structSource.Position.String(), field.Names, pk))
					columnInfo.PK = true
					tableInfo.PrimaryKeys = append(tableInfo.PrimaryKeys, columnInfo)
				default:
					tableInfo.NotPrimaryKeys = append(tableInfo.NotPrimaryKeys, columnInfo)
				}
				// hasOne tag
				for _, hasOneTag := range strings.Split(tag.Get(cfg.GoHasOneTag), ",") {
					if hasOneTag != "" {
						contexts.Trace(ctx).Debug(fmt.Sprintf("%s: field.Names=%s, hasOneTag=%q", structSource.Position.String(), field.Names, hasOneTag))
						columnInfo.HasOneTags = append(columnInfo.HasOneTags, hasOneTag)
						tableInfo.HasOneTagsKeys = append(tableInfo.HasOneTagsKeys, hasOneTag)
						tableInfo.HasOneTags[hasOneTag] = append(tableInfo.HasOneTags[hasOneTag], columnInfo)
					}
				}
				// hasMany tag
				for _, hasManyTag := range strings.Split(tag.Get(cfg.GoHasManyTag), ",") {
					if hasManyTag != "" {
						contexts.Trace(ctx).Debug(fmt.Sprintf("%s: field.Names=%s, hasManyTag=%q", structSource.Position.String(), field.Names, hasManyTag))
						columnInfo.HasManyTags = append(columnInfo.HasManyTags, hasManyTag)
						tableInfo.HasManyTagsKeys = append(tableInfo.HasManyTagsKeys, hasManyTag)
						tableInfo.HasManyTags[hasManyTag] = append(tableInfo.HasManyTags[hasManyTag], columnInfo)
					}
				}

				tableInfo.Columns = append(tableInfo.Columns, columnInfo)
			}
		}
	}

	slices.Sort(tableInfo.HasOneTagsKeys)
	tableInfo.HasOneTagsKeys = slices.Compact(tableInfo.HasOneTagsKeys)

	slices.Sort(tableInfo.HasManyTagsKeys)
	tableInfo.HasManyTagsKeys = slices.Compact(tableInfo.HasManyTagsKeys)

	return tableInfo
}

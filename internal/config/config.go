package config

import (
	"log/slog"
	"strings"

	"github.com/hakadoriya/z.go/cliz"
	"github.com/hakadoriya/z.go/contextz"
	"github.com/hakadoriya/z.go/errorz"

	"github.com/hakadoriya/ormgen/internal/logs"
)

type GenerateConfig struct {
	// Common
	Trace    bool   `cli:"trace,   env=ORMGEN_TRACE,                      description=Enable trace mode"`
	Debug    bool   `cli:"debug,   env=ORMGEN_DEBUG,                      description=Enable debug mode"`
	Dialect  string `cli:"dialect, env=ORMGEN_DIALECT,  default=postgres, description=dialect for DML"`
	Language string `cli:"lang,    env=ORMGEN_LANGUAGE, default=go,       description=programming language to generate ORM"`
	// Go
	GoColumnTag                  string `cli:"go-column-tag,                     env=ORMGEN_GO_COLUMN_TAG,                     default=db,          description=column annotation key for Go struct tag"`
	GoPKTag                      string `cli:"go-pk-tag,                         env=ORMGEN_GO_PK_TAG,                         default=pk,          description=primary key annotation key for Go struct tag"`
	GoHasOneTag                  string `cli:"go-has-one-tag,                    env=ORMGEN_GO_HAS_ONE_TAG,                    default=hasOne,      description=\"hasOne\" annotation key for Go struct tag"`
	GoHasManyTag                 string `cli:"go-has-many-tag,                   env=ORMGEN_GO_HAS_MANY_TAG,                   default=hasMany,     description=\"hasMany\" annotation key for Go struct tag"`
	GoTableNameMethod            string `cli:"go-table-name-method,              env=ORMGEN_GO_TABLE_NAME_METHOD,              default=TableName,   description=method name for table"`
	GoColumnNameMethodPrefix     string `cli:"go-column-name-method-prefix,      env=ORMGEN_GO_COLUMN_NAME_METHOD_PREFIX,      default=ColumnName_, description=method name for columns"`
	GoColumnsNameMethod          string `cli:"go-columns-name-method,            env=ORMGEN_GO_COLUMNS_NAME_METHOD,            default=ColumnsName, description=method prefix for column"`
	GoSliceTypeSuffix            string `cli:"go-slice-type-suffix,              env=ORMGEN_GO_SLICE_TYPE_SUFFIX,              default=Slice,       description=suffix for slice type"`
	GoORMOutputPath              string `cli:"go-orm-output-path,                env=ORMGEN_GO_ORM_OUTPUT_PATH,                default=ormgen,      description=output path of ORM."`
	GoORMOutputPackageImportPath string `cli:"go-orm-output-package-import-path, env=ORMGEN_GO_ORM_OUTPUT_PACKAGE_IMPORT_PATH, default=,            description=package import path of ORM output directory. If empty, try to detect automatically."`
	GoORMPackageName             string `cli:"go-orm-package-name,               env=ORMGEN_GO_ORM_PACKAGE_NAME,               default=,            description=package name for ORM. If empty, use the base name of the output path."`
	GoORMStructPackageImportPath string `cli:"go-orm-struct-package-import-path, env=ORMGEN_GO_ORM_STRUCT_PACKAGE_IMPORT_PATH,                      description=package import path of ORM target struct. If empty, try to detect automatically."`
	GoORMInterfaceName           string `cli:"go-orm-interface-name,             env=ORMGEN_GO_ORM_INTERFACE_NAME,             default=ORM,         description=interface type name for ORM"`
	GoORMStructName              string `cli:"go-orm-struct-name,                env=ORMGEN_GO_ORM_STRUCT_NAME,                default=_ORM,        description=struct name for ORM"`
}

func GeneratePreHookExec(c *cliz.Command, args []string) (err error) {
	cfg := new(GenerateConfig)
	if err := cliz.UnmarshalOptions(c, cfg); err != nil {
		return errorz.Errorf("cliz.UnmarshalOptions: %w", err)
	}

	if cfg.Trace {
		logs.Trace = logs.NewTrace(c.Stdout(), slog.LevelDebug, slog.String("logName", "trace"))
		logs.Stdout = logs.NewStdout(c.Stdout(), slog.LevelDebug, slog.String("logName", "stdout"))
		logs.Stderr = logs.NewStderr(c.Stderr(), slog.LevelDebug, slog.String("logName", "stderr"))
		cfg.Debug = true
	}

	if cfg.Debug {
		logs.Stdout = logs.NewStdout(c.Stdout(), slog.LevelDebug, slog.String("logName", "stdout"))
		logs.Stderr = logs.NewStderr(c.Stderr(), slog.LevelDebug, slog.String("logName", "stderr"))

		cmd := c.GetExecutedCommandNames()
		logs.Stdout.Debug(strings.Join(cmd, " ")+" "+strings.Join(args, " "), slog.Any("cmd", cmd), slog.Any("args", args), slog.Any("config", cfg))
	}

	c.SetContext(contextz.WithValue(c.Context(), cfg))

	return nil
}

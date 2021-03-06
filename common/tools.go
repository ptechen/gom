package common

import (
	"fmt"
	"path"
	"strings"

	"xorm.io/core"
)

var created = []string{"created_at"}
var updated = []string{"updated_at"}
var deleted = []string{"deleted_at"}

func InStringSlice(f string, a []string) bool {
	for _, s := range a {
		if f == s {
			return true
		}
	}
	return false
}

func getTypeAndImports(table *core.Table, column *core.Column) (t string, imports map[string]string) {
	t = sqlType2TypeString(column.SQLType)

	imports = map[string]string{}

	if Configs().Tinyint2Bool && strings.HasPrefix(column.Name, "is_") &&
		column.SQLType.Name == "TINYINT" && column.SQLType.DefaultLength == 1 {
		t = "bool"
		return
	}

	for k, v := range Configs().ReplaceType {
		//如果字段类型等于原类型
		if t == path.Base(k) {
			t, pkg := getTypeAndPkg(v)
			if pkg != "" {
				imports[pkg] = pkg
			}
			return t, imports
		}
	}

	//如果字段名字在集合
	if v, ok := Configs().ColumnType[column.Name]; ok {
		t, pkg := getTypeAndPkg(v)
		if pkg != "" {
			imports[pkg] = pkg
		}
		return t, imports
	} else if v, ok := Configs().ColumnType[fmt.Sprintf("%s.%s", table.Name, column.Name)]; ok {
		t, pkg := getTypeAndPkg(v)
		if pkg != "" {
			imports[pkg] = pkg
		}
		return t, imports
	}

	if t == "time.Time" {
		imports["time"] = "time"
	}

	return
}

func getTypeAndPkg(v string) (t string, pkg string) {
	//获取新类型名
	t = path.Base(v)
	//获取新包名
	if path.Ext(v) != "" {
		pkg = strings.TrimSuffix(v, path.Ext(v))
	}

	return
}

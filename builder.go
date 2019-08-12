package framework

import (
	. "database/sql"
	"errors"
	"html"
	"reflect"
	"strconv"
	"strings"
)

type Database struct {
	sel               string
	from              string
	where             []map[string]string
	whereResult       string
	join              string
	groupBy           string
	having            string
	orderBy           []string
	limit             string
	query             string
	call              bool
	DB                *DB
	Option            string
	BeforeOption      string
	transatction      *Tx
	removeSpecialChar bool
}

func (sql *Database) Select(value string) *Database {
	sql.sel = value
	return sql
}

func (sql *Database) From(value string) *Database {
	sql.call = false
	sql.from = value
	sql.removeSpecialChar = true
	return sql
}
func (sql *Database) RemoveSpecialChar(value bool) *Database {
	sql.removeSpecialChar = value
	return sql
}

func whereProccess(field string, value interface{}) string {
	field = strings.TrimSpace(field)
	fields := strings.Split(field, " ")
	row := fields[0]
	op := "="
	var where string
	if len(fields) > 2 {
		op = fields[1] + " " + fields[2]
	} else if len(fields) > 1 {
		op = fields[1]
	}
	var reflectValue = reflect.ValueOf(value)

	var val string
	var i int
	switch reflectValue.Kind() {
	case reflect.String:
		val = strings.TrimSpace(reflectValue.String())
		switch op {
		case "sql":
			where = row + val
		case "raw":
			where = val
		default:
			where = row + " " + op + " '" + RemoveSpecialChar(val).(string) + "'"

		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i = int(reflectValue.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i = int(reflectValue.Int())
	}
	if val == "" {
		val = strconv.Itoa(i)
		where = row + " " + op + " " + val
	}
	return where

}
func (sql *Database) Where(field string, value interface{}) *Database {
	where := map[string]string{}
	where["op"] = "AND"
	where["value"] = whereProccess(field, value)
	sql.where = append(sql.where, where)
	return sql
}
func (sql *Database) WhereOr(field string, value interface{}) *Database {
	where := map[string]string{}
	where["op"] = "OR"
	where["value"] = whereProccess(field, value)
	sql.where = append(sql.where, where)
	return sql
}
func (sql *Database) StartGroup(op string) *Database {
	op = strings.ToUpper(op)
	if op != "AND" || op != "OR" {
		op = "AND"
	}
	where := map[string]string{}
	where["op"] = op
	where["value"] = "("
	where["groupstart"] = "1"
	sql.where = append(sql.where, where)
	return sql
}

func (sql *Database) EndGroup() *Database {
	where := map[string]string{}
	where["op"] = ""
	where["value"] = ")"
	sql.where = append(sql.where, where)
	return sql
}

func (sql *Database) WhereIn(field string, values []string) *Database {
	where := map[string]string{}
	where["op"] = "AND"
	where["value"] = field + " in('" + strings.Join(values, "','") + "')"
	sql.where = append(sql.where, where)
	return sql
}

func (sql *Database) WhereInOr(field string, values []string) *Database {
	where := map[string]string{}
	where["op"] = "OR"
	where["value"] = field + " in('" + strings.Join(values, "','") + "')"
	sql.where = append(sql.where, where)
	return sql
}

func (sql *Database) WhereNotIn(field string, values []string) *Database {
	where := map[string]string{}
	where["op"] = "AND"
	where["value"] = field + " not in('" + strings.Join(values, "','") + "')"
	sql.where = append(sql.where, where)
	return sql
}

func (sql *Database) WhereNotInOr(field string, values []string) *Database {
	where := map[string]string{}
	where["op"] = "OR"
	where["value"] = field + " not in('" + strings.Join(values, "','") + "')"
	sql.where = append(sql.where, where)
	return sql
}

func (sql *Database) WhereBetween(field, start, end string) *Database {
	where := map[string]string{}
	where["op"] = "AND"
	where["value"] = field + " BETWEEN '" + start + "' AND '" + end + "'"
	sql.where = append(sql.where, where)
	return sql
}

func (sql *Database) WhereBetweenOr(field, start, end string) *Database {
	where := map[string]string{}
	where["op"] = "OR"
	where["value"] = field + " BETWEEN '" + start + "' AND '" + end + "'"
	sql.where = append(sql.where, where)
	return sql
}

func (sql *Database) Join(table string, on string, join string) *Database {
	sql.join += join + " join " + table + " on " + on + "\n"
	return sql

}

func (sql *Database) OrderBy(orderBy string, dir string) *Database {
	sql.orderBy = append(sql.orderBy, orderBy+" "+dir)
	return sql
}

func (sql *Database) GroupBy(groupBy string) *Database {
	sql.groupBy = groupBy
	return sql
}

func (sql *Database) Limit(limit int, start int) *Database {
	sql.limit = "LIMIT " + strconv.Itoa(start) + "," + strconv.Itoa(limit)
	return sql
}

func (sql *Database) Having(str string) *Database {
	sql.having = str
	return sql
}

func whereBuild(sql *Database) {

	sql.whereResult = ""
	opShow := true
	for i, v := range sql.where {
		if i == 0 {
			opShow = false
		}

		if opShow {
			sql.whereResult += " " + v["op"] + " "
		}
		sql.whereResult += v["value"]
		opShow = true
		if v["groupstart"] == "1" {
			opShow = false
		}
	}
}

func get(sql *Database) {
	whereBuild(sql)
	if sql.sel == "" {
		sql.sel = "*"
	}
	sql.query = "select " + sql.sel + "\nfrom " + sql.from
	if sql.join != "" {
		sql.query += " \n" + sql.join
	}
	if sql.whereResult != "" {
		sql.query += "\nWHERE " + sql.whereResult
	}

	if sql.groupBy != "" {
		sql.query += "\nGROUP BY " + sql.groupBy
	}

	if sql.having != "" {
		sql.query += "\nHaving " + sql.having
	}

	if sql.orderBy != nil {
		sql.query += "\nORDER BY " + strings.Join(sql.orderBy, ",")
	}

	if sql.limit != "" {
		sql.query += "\n" + sql.limit
	}
}

func (sql *Database) Call(procedure string, value []string) *Database {
	var values = ""
	if len(value) != 0 {
		if sql.removeSpecialChar {
			for i, v := range value {
				value[i] = RemoveSpecialChar(v).(string)
			}
		}
		values = "('" + strings.Join(value, "','") + "')"
	}
	sql.query = "call " + procedure + " " + values
	sql.call = true
	return sql
}

func (sql *Database) Result() ([]map[string]interface{}, error) {
	if sql.from == "" && sql.call == false {
		return nil, errors.New("nothing table selected")
	}
	if sql.call == false {
		get(sql)
	}
	query := sql.query
	sql.call = false
	var err error
	var rows *Rows
	//var stmt *Stmt
	if sql.DB == nil {
		sql.DB, err = MysqlConnect()
		if err != nil {
			return nil, err
		}
	}
	//stmt, err = sql.DB.Prepare(query)
	//if err != nil {
	//	return nil, err
	//}
	//rows, err = stmt.Query()
	rows, err = sql.DB.Query(query)
	if err != nil {
		return nil, err
	}
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	result := []map[string]interface{}{}
	if count == 0 {
		return nil, nil
	}
	for rows.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}
		data := map[string]interface{}{}
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
				v = html.UnescapeString(v.(string))
			} else {
				v = val
			}
			if v == nil {
				data[col] = ""
			} else {
				data[col] = v
			}
		}
		result = append(result, data)
	}
	//stmt.Close()
	rows.Close()

	return result, nil

}

func (sql *Database) Row() (map[string]interface{}, error) {
	if sql.call == false {
		sql.Limit(1, 0)
	}
	db, err := sql.Result()
	if err != nil {
		return nil, err
	}

	if len(db) >= 1 {
		return db[0], nil
	}
	return nil, nil

}
func insert(querySql string, value []interface{}, sql *Database) (interface{}, error) {
	var err error
	var stmt *Stmt
	if sql.DB == nil {
		sql.DB, err = MysqlConnect()
		if err != nil {
			return 0, err
		}
	}
	//defer sql.DB.Close()
	querySql += " " + sql.Option
	sql.query = querySql
	if sql.transatction != nil {
		stmt, err = sql.transatction.Prepare(querySql)
	} else {
		stmt, err = sql.DB.Prepare(querySql)
	}
	if err != nil {
		return 0, err
	}
	//defer stmt.Close()
	res, err := stmt.Exec(value...)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	stmt.Close()
	return id, nil
}
func (sql *Database) Insert(query map[string]interface{}) (interface{}, error) {
	if sql.from == "" {
		return nil, errors.New("nothing table selected")
	}
	querySql := "INSERT INTO " + sql.from
	tag := ""
	field := ""
	value := []interface{}{}
	for i, v := range query {
		tag += "?,"
		field += i + ","
		if sql.removeSpecialChar {
			v = RemoveSpecialChar(v)
		}
		value = append(value, v)

	}
	tag = tag[0 : len(tag)-1]
	field = field[0 : len(field)-1]
	querySql += "(" + field + ") values " + "(" + tag + ")"

	return insert(querySql, value, sql)

}

func (sql *Database) InsertBatch(query []map[string]interface{}) (interface{}, error) {
	if sql.from == "" {
		return nil, errors.New("nothing table selected")
	}
	querySql := "INSERT INTO " + sql.from
	value := []interface{}{}
	field := JoinMapKey(query[0], ",")
	fieldArray := strings.Split(field, ",")
	tag := strings.Repeat("?,", len(fieldArray))
	tag = tag[0 : len(tag)-1]
	tag = "(" + tag + ")"
	tags := strings.Repeat(tag+",", len(query))
	tags = tags[0 : len(tags)-1]
	for _, v := range query {
		for _, v2 := range fieldArray {
			if sql.removeSpecialChar {
				v[v2] = RemoveSpecialChar(v[v2])
			}
			value = append(value, v[v2])
		}
	}

	querySql += "(" + field + ") values " + tags

	return insert(querySql, value, sql)
}

func (sql *Database) Update(query map[string]interface{}) error {
	if sql.from == "" {
		return errors.New("nothing table selected")
	}

	if query == nil {
		return errors.New("Query invalid")
	}
	whereBuild(sql)
	var join = ""
	if sql.join != "" {
		join = " \n" + sql.join
	}
	querySql := "UPDATE " + sql.from + join + " SET "
	var set []string
	value := []interface{}{}
	for i, v := range query {
		set = append(set, i+"=?")
		if sql.removeSpecialChar {
			v = RemoveSpecialChar(v)
		}
		value = append(value, v)
	}

	querySql += strings.Join(set, ",")
	sql.query = querySql
	if sql.whereResult != "" {
		sql.query += "\nWHERE " + sql.whereResult
	}

	return UpdateProses(sql, value)
}

func (sql *Database) UpdateBatch(query []map[string]interface{}, id string) error {
	id = strings.TrimSpace(id)
	if sql.from == "" {
		return errors.New("nothing table selected")
	}

	if query == nil {
		return errors.New("query invalid")
	}
	var join = ""
	if sql.join != "" {
		join = " \n" + sql.join
	}
	querySql := "UPDATE " + sql.from + join + " SET "
	var set map[string][]string
	set = map[string][]string{}
	var value map[string][]interface{}
	value = map[string][]interface{}{}
	values := []interface{}{}
	whereIn := []string{}
	for i, v := range query {
		if v[id] == nil {
			return errors.New("primary key for update, not found")
		}
		var reflectValue = reflect.ValueOf(v[id])
		var valId string
		var idInt int
		switch reflectValue.Kind() {
		case reflect.String:
			valId = strings.TrimSpace(reflectValue.String())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			idInt = int(reflectValue.Uint())
			valId = strconv.Itoa(idInt)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			idInt = int(reflectValue.Int())
			valId = strconv.Itoa(idInt)
		}
		whereIn = append(whereIn, "'"+valId+"'")

		for i2, v2 := range v {
			if i2 == id {
				continue
			}
			if i == 0 {
				set[i2] = append(set[i2], i2+" = (CASE "+id+"\n")
			}
			set[i2] = append(set[i2], "WHEN '"+valId+"' THEN ?\n")
			if sql.removeSpecialChar {
				v2 = RemoveSpecialChar(v2)
			}
			value[i2] = append(value[i2], v2)
		}
	}
	for i, v := range set {
		querySql += strings.Join(v, "") + " END),\n"
		for _, v2 := range value[i] {
			values = append(values, v2)
		}

	}
	querySql = strings.TrimRight(querySql, ",\n")
	querySql += "where " + id + " in(" + strings.Join(whereIn, ",") + ")"
	sql.query = querySql

	return UpdateProses(sql, values)
}

func UpdateProses(sql *Database, value []interface{}) error {
	var err error
	var stmt *Stmt
	if sql.DB == nil {
		sql.DB, err = MysqlConnect()
		if err != nil {
			return err
		}
	}
	if sql.transatction != nil {
		stmt, err = sql.transatction.Prepare(sql.query)
	} else {
		stmt, err = sql.DB.Prepare(sql.query)
	}
	//defer stmt.Close()
	_, err = stmt.Exec(value...)
	if err != nil {
		return err
	}
	stmt.Close()
	return nil
}

func (sql *Database) Delete() error {
	var err error

	if sql.from == "" {
		return errors.New("nothing table selected")
	}
	if sql.DB == nil {
		sql.DB, err = MysqlConnect()
		if err != nil {
			return err
		}
	}
	whereBuild(sql)
	querySql := "DELETE FROM " + sql.from + " "
	sql.query = querySql
	if sql.whereResult != "" {
		sql.query += "\nWHERE " + sql.whereResult
	}
	if sql.transatction != nil {
		_, err = sql.transatction.Exec(sql.query)
	} else {
		_, err = sql.DB.Exec(sql.query)
	}

	if err != nil {
		return err
	}
	return nil
}
func (sql *Database) Close() {
	if sql.DB != nil {
		sql.DB.Close()
	}
}

func (sql *Database) Clear() {
	var tx = sql.transatction
	var db = sql.DB
	p := reflect.ValueOf(sql).Elem()
	p.Set(reflect.Zero(p.Type()))
	sql.transatction = tx
	sql.DB = db
}

func (sql *Database) QueryView() string {
	return sql.query
}

func (sql *Database) Transaction() error {
	var err error
	if sql.DB == nil {
		sql.DB, err = MysqlConnect()
		if err != nil {
			return err
		}
	}
	tx, err := sql.DB.Begin()
	sql.transatction = tx
	return err
}

func (sql *Database) Rollback() error {
	var err error
	if sql.transatction != nil {
		err = sql.transatction.Rollback()
	}
	sql.transatction = nil
	return err
}

func (sql *Database) Commit() error {
	var err error
	if sql.transatction != nil {
		err = sql.transatction.Commit()
	}
	sql.transatction = nil
	return err
}

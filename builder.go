package framework

import (
	"fmt"
	"strings"
	"reflect"
	"strconv"
	. "database/sql"
	"errors"
)

type Database struct {
	sel string
	from string
	where string
	whereOr []string
	whereAnd []string
	whereInOr []string
	whereInAnd []string
	join string
	groupBy string
	orderBy []string
	limit string
	query string
	call bool
	DB *DB
	row *Rows
	stmt *Stmt
	Option string
	BeforeOption string
}

func(sql *Database) Select(value string)  {
	sql.sel = value
}

func(sql *Database) From(value string)  {
	sql.call = false
	sql.from = value
}

func whereProccess(field string,value interface{}) string{
	field = strings.TrimSpace(field)
	fields := strings.Split(field," ")
	row := fields[0]
	op := "="
	var where string
	if(len(fields)>1){
		op = fields[1]
	}
	var	reflectValue = reflect.ValueOf(value)

	var val string
	var i int
	switch reflectValue.Kind() {
	case reflect.String :
		val = strings.TrimSpace(reflectValue.String())
		where = row+" "+op+" '"+val+"'"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64 :
		i = int(reflectValue.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64 :
		i = int(reflectValue.Int())
	}
	if val == "" {
		val = strconv.Itoa(i)
		where = row+" "+op+" "+val
	}
	return where

}
func(sql *Database) Where(field string,value interface{})  {
	sql.whereAnd = append(sql.whereAnd,whereProccess(field,value))
}
func(sql *Database) WhereOr(field string,value interface{})  {
	sql.whereOr = append(sql.whereOr,whereProccess(field,value))
}

func(sql *Database) WhereIn(field string, values []string)  {
	sql.whereInAnd = append(sql.whereInAnd,field+" in('"+strings.Join(values,"','")+"')")
}

func(sql *Database) WhereInOr(field string, values []string)  {
	sql.whereInOr = append(sql.whereInOr,field+" in('"+strings.Join(values,"','")+"')")
}

func(sql *Database) Join(table string,on string,join string)  {
	sql.join += join+" join "+table+" on "+on+"\n"

}

func(sql *Database) OrderBy(orderBy string, dir string){
	sql.orderBy = append(sql.orderBy, orderBy+" "+dir)
}

func(sql *Database) GroupBy(groupBy string){
	sql.groupBy = groupBy
}

func(sql *Database) Limit(limit int,start int){
	sql.limit = "LIMIT "+strconv.Itoa(start)+","+strconv.Itoa(limit)
}

func whereBuild(sql *Database){

	sql.where = ""
	if len(sql.whereOr)>= 1{
		sql.where += strings.Join(sql.whereOr," \nOR ")
	}
	if len(sql.whereInOr)>= 1{
		if sql.where != ""{
			sql.where += " \nOR "
		}
		sql.where += strings.Join(sql.whereInOr," \nOR ")
	}


	if len(sql.whereAnd)>= 1{
		if sql.where != ""{
			sql.where += " \nAND "
		}
		sql.where += strings.Join(sql.whereAnd," \nAND ")
	}

	if len(sql.whereInAnd)>= 1{
		if sql.where != ""{
			sql.where += " \nAND "
		}
		sql.where += strings.Join(sql.whereInAnd," \nAND ")
	}
}

func get(sql *Database) {
	whereBuild(sql)
	if sql.sel == ""{
		sql.sel="*"
	}
	sql.query = "select "+sql.sel+"\nfrom "+sql.from
	if sql.join != "" {
		sql.query+=" \n"+sql.join
	}
	if sql.where != "" {
		sql.query+="\nWHERE "+sql.where
	}

	if sql.groupBy != "" {
		sql.query+="\nGROUP BY "+sql.groupBy
	}

	if sql.orderBy != nil {
		sql.query+="\nORDER BY "+strings.Join(sql.orderBy,",")
	}

	if sql.limit != "" {
		sql.query+="\n"+sql.limit
	}
}

func(sql *Database) Call(procedure string,value []string){
	values := "('"+strings.Join(value,"','")+"')"
	sql.query = "call "+procedure+" "+values
	sql.call = true
}

func(sql *Database) Result() ([]map[string]string,error) {
	if sql.from == "" && sql.call == false{
		return  nil,errors.New("nothing table selected")
	}
	if sql.call == false {
		get(sql)
	}
	query := sql.query
	sql.call = false
	var err error
	if sql.DB == nil {
		sql.DB, err = MysqlConnect()
		if err != nil {
			return nil,err
		}
	}
	//defer sql.DB.Close()
	rows,err  := sql.DB.Query(query)
	sql.row = rows
	if err != nil {
		return nil,err
	}
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	result := []map[string]string{}
	if count == 0 {
		return nil, nil
	}
	for rows.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		data := map[string]string{}
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if (ok) {
				v = string(b)
			} else {
				v = val
			}
			if v == nil {
				data[col] = ""
			}else{
				data[col] = v.(string)
			}
		}
		result = append(result, data)
	}

	return result,nil

}

func(sql *Database) Row() (map[string]string,error) {
	if sql.call == false {
		sql.Limit(1,0)
	}
	db,err := sql.Result()
	if err!=nil {
		return nil,err
	}

	if(len(db)>=1){
		return db[0],nil
	}
	return nil,nil

}
func insert(querySql string,value []interface{},sql *Database) (interface{},error) {
	var err error
	if sql.DB == nil {
		sql.DB, err = MysqlConnect()
		if err != nil {
			return 0,err
		}
	}
	//defer sql.DB.Close()
	querySql+=" "+sql.Option
	sql.query=querySql
	stmt, err := sql.DB.Prepare(querySql)
	sql.stmt = stmt
	if err != nil {
		return  0,err
	}
	//defer stmt.Close()
	res, err := stmt.Exec(value...)
	if err != nil {
		return nil,err
	}
	id,err := res.LastInsertId()
	if err != nil {
		return nil,err
	}
	return id,nil
}
func(sql *Database) Insert(query map[string]string) (interface{},error){
	if sql.from == ""{
		return  nil,errors.New("nothing table selected")
	}
	querySql := "INSERT INTO "+sql.from
	tag := ""
	field := ""
	value := []interface{}{}
	for i,v := range query{
		tag+="?,"
		field+=i+","
		value = append(value, v)
	}
	tag = tag[0:len(tag)-1]
	field = field[0:len(field)-1]
	querySql+="("+field+") values "+"("+tag+")"

	return insert(querySql,value,sql)


}

func(sql *Database) MultiInsert(query []map[string]string) (interface{},error){
	if sql.from == ""{
		return  nil,errors.New("nothing table selected")
	}
	querySql := "INSERT INTO "+sql.from
	value := []interface{}{}
	field := JoinMapKey(query[0],",")
	fieldArray := strings.Split(field,",")
	tag := strings.Repeat("?,",len(fieldArray))
	tag = tag[0:len(tag)-1]
	tag = "("+tag+")"
	tags := strings.Repeat(tag+",",len(query))
	tags = tags[0:len(tags)-1]
	for _,v := range query{
		for _,v2 := range fieldArray{
			value = append(value,v[v2])
		}
	}

	querySql+="("+field+") values "+tags


	return insert(querySql,value,sql)
}

func(sql *Database) Update(query map[string]string) error {
	if sql.from == ""{
		return  errors.New("nothing table selected")
	}

	if query == nil{
		return  errors.New("Query invalid")
	}
	whereBuild(sql)
	querySql := "UPDATE "+sql.from+" SET "
	var set []string
	value := []interface{}{}
	for i,v := range query{
		set = append(set,i+"=?")
		value = append(value,v)
	}

	querySql+= strings.Join(set,",")
	sql.query = querySql
	if sql.where != "" {
		sql.query+="\nWHERE "+sql.where
	}

	return  UpdateProses(sql,value)
}

func UpdateProses(sql *Database,value []interface{}) error {
	var err error
	if sql.DB == nil {
		sql.DB, err = MysqlConnect()
		if err != nil {
			return err
		}
	}
	stmt, err := sql.DB.Prepare(sql.query)
	sql.stmt = stmt
	if err != nil {
		return  err
	}
	//defer stmt.Close()
	res, err := stmt.Exec(value...)
	if err != nil {
		return err
	}
	fmt.Println(res)

	return nil
}

func (sql *Database) Delete() error{
	var err error

	if sql.from == ""{
		return  errors.New("nothing table selected")
	}
	if sql.DB == nil {
		sql.DB, err = MysqlConnect()
		if err != nil {
			return err
		}
	}
	whereBuild(sql)
	querySql := "DELETE FROM "+sql.from+" "
	sql.query = querySql
	if sql.where != "" {
		sql.query+="\nWHERE "+sql.where
	}

	_, err = sql.DB.Exec(sql.query)
	if err != nil {
		return  err
	}

	return nil
}
func(sql *Database) Close() {
	if sql.DB != nil{
		sql.DB.Close()
	}
	if sql.row != nil{
		sql.row.Close()
	}
	if sql.stmt != nil{
		sql.stmt.Close()
	}
}

func(sql *Database) Clear() {
	p := reflect.ValueOf(sql).Elem()
	p.Set(reflect.Zero(p.Type()))
}

func(sql *Database) QueryView() string {
	return sql.query
}
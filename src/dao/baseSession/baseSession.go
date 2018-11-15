package baseSession

import (
	"errors"
	"fmt"
	"github.com/smallnest/rpcx/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"reflect"
	"setting"
	"strings"
	"time"
)

/**
mgo使用参考网址：
http://my.oschina.net/ffs/blog/300148
http://blog.csdn.net/varding/article/details/17200299
http://chenzhou123520.iteye.com/blog/1637629
官方资料：https://labix.org/mgo
*/

type MgoOperators = string

const MGO_UPDATE_SET MgoOperators = "$set"
const MGO_UPDATE_UNSET MgoOperators = "$unset"
const MGO_UPDATE_INC MgoOperators = "$inc"

const MGO_UPDATE_pushAll MgoOperators = "$pushAll"
const MGO_UPDATE_push MgoOperators = "$push"
const MGO_UPDATE_pull MgoOperators = "$pull"

const MGO_SELECT_AND MgoOperators = "$and"
const MGO_SELECT_OR MgoOperators = "$or"

const MGO_SELECT_REGEX MgoOperators = "$regex"
const MGO_SELECT_OPTIONS MgoOperators = "$options"
const MGO_SELECT_OPTIONS_IGNORECASE MgoOperators = "i"

const MGO_SELECT_EQ MgoOperators = "$eq"     //=($eq)
const MGO_SELECT_NE MgoOperators = "$ne"     //!=($ne)
const MGO_SELECT_GT MgoOperators = "$gt"     //>($gt)
const MGO_SELECT_LT MgoOperators = "$lt"     //<($lt)
const MGO_SELECT_GTE MgoOperators = "$gte"   //>=($gte)
const MGO_SELECT_LTE MgoOperators = "$lte"   //<=($lte)
const MGO_SELECT_IN MgoOperators = "$in"     //in($in)
const MGO_SELECT_NOTIN MgoOperators = "$nin" //not in ($nin)
const MGO_SELECT_EXISTS MgoOperators = "$exists"
const MGO_SELECT_eleMatch MgoOperators = "$elemMatch"
const MGO_SELECT_Slice MgoOperators = "$slice"

const MGO_SELECT_geoWithin MgoOperators = "$geoWithin"
const MGO_SELECT_geoWithin_centerSphere MgoOperators = "$centerSphere"

const MGO_MROutputOptions_REPLACE MgoOperators = "replace"
const MGO_MROutputOptions_MERGE MgoOperators = "merge"
const MGO_MROutputOptions_REDUCE MgoOperators = "reduce"

const MGO_AGG_PREFIX_DOLLAR MgoOperators = "$"
const MGO_AGG_MATCH MgoOperators = "$match"
const MGO_AGG_SORT MgoOperators = "$sort"
const MGO_AGG_SKIP MgoOperators = "$skip"
const MGO_AGG_LIMIT MgoOperators = "$limit"
const MGO_AGG_DAYOFWEEK MgoOperators = "$dayOfWeek"
const MGO_AGG_GROUP MgoOperators = "$group"
const MGO_AGG_SUM MgoOperators = "$sum"
const MGO_AGG_ADD MgoOperators = "$add"
const MGO_AGG_GEONEAR MgoOperators = "$geoNear"
const MGO_AGG_GEONEAR_NEAR MgoOperators = "near"
const MGO_AGG_GEONEAR_distanceField MgoOperators = "distanceField"
const MGO_AGG_GEONEAR_Query MgoOperators = "query"
const PIPELINE_MAX_NUM int = 10
const DEFAULT_LIMIT_NUM int = 10
const TAG_BSON = "bson"
const TAG_JOIN_SPLIT_SYMBAL = "."
const LOG_FREQUENCY int = 50

//https://docs.mongodb.com/manual/reference/operator/query/centerSphere/

type MGO_Int_Array []int

type QueryOption struct {
	Operator MgoOperators
	MdDefine string
}

type BaseSession struct {
	session        *mgo.Session
	sessionGridFs  *mgo.Session
	collection     *mgo.Collection
	gridFS         *mgo.GridFS
	type_          reflect.Type
	collectionName string
}

var DataBaseName string

var globalMgoSession *mgo.Session

var jobAddUpCounter = make(map[string]int, DEFAULT_LIMIT_NUM)

func DBInit() {
	testDir, _ := os.Getwd()
	log.Info("baseSession path=" + testDir)
	InitBaseSession(setting.DBConfig.DatabaseName, fmt.Sprintf("mongodb://%s:%s@%s/%s", setting.DBConfig.UserName, setting.DBConfig.Password, setting.DBConfig.Host, setting.DBConfig.DatabaseName))
	//if strings.Contains(testDir, "xm1804ServiceApp") {
	//	//源码启动
	//	InitBaseSession("xm1804", "mongodb://app1804:pwd1804@192.168.1.23:30017/xm1804")
	//} else if strings.Contains(testDir, "xm1804ServiceDao") {
	//	//make启动
	//	InitBaseSession("xm1804", "mongodb://app1804:pwd1804@192.168.1.23:30017/xm1804")
	//} else if strings.Contains(testDir, "/src/") {
	//	//测试用例启动
	//	InitBaseSession("xm1804test", "mongodb://xm1804test:pwdtest@192.168.1.23:30017/xm1804test")
	//} else if strings.Contains(testDir, "xm1804Shared") {
	//	//运行测试dao时会用到
	//	InitBaseSession("xm1804", "mongodb://app1804:pwd1804@192.168.1.23:30017/xm1804")
	//}
	//RecoverDir = TestDir
	//if index := strings.Index(TestDir, AppDir); index > 0 {
	//	TestDir = TestDir[:index+len(AppDir)]
	//	if err != nil {
	//		panic(err)
	//	} else {
	//		os.Chdir(TestDir)
	//		configPath = shell.GetCurrentPath() + "/conf/"
	//		aliPayConfigPath = configPath + "aliPay/"
	//		LoadAliPayConfig()
	//		os.Chdir(RecoverDir)
	//	}
	//}

	//这个初始化仅用于测试用例调用。
	//InitBaseSession("xm1804", "mongodb://app1804:pwd1804@192.168.1.23:30017/xm1804")
	//InitBaseSession("xm1804", "mongodb://127.0.0.1:27017/xm1804")
}

func InitBaseSession(databaseName, url string) {
	//log.Info("初始化数据库:%s",databaseName)
	var err error
	globalMgoSession, err = mgo.Dial(url)
	if err != nil {
		panic(fmt.Sprintf("InitBaseSession发生错误：%v url:%s", err, url))
	}
	globalMgoSession.SetPoolLimit(setting.MongoDBMaxConnectionNum) //设置线程池上线
	DataBaseName = databaseName
}

func (*BaseSession) GetSession() *mgo.Session {
	session := globalMgoSession.Clone()
	return session
}

func (this *BaseSession) Init(object interface{}, collection string) {
	this.collectionName = collection
	this.type_ = reflect.TypeOf(object)
}

func OrderDesc(field string) string {
	return "-" + field
}

//func (this *BaseSession) ConnectSessionGridFS(collection string) {
//	this.session = this.GetSession()
//	this.gridFS = this.GetGridFS(DataBaseName, collection)
//}

func (this *BaseSession) ConnectSessionGridFS() {
	this.sessionGridFs = this.GetSession()
	this.gridFS = this.GetGridFS(DataBaseName, this.collectionName)
}

func (this *BaseSession) CreateFile(filename string, data []byte) {
	this.ConnectSessionGridFS()
	defer this.closeSession()
	if this.gridFS == nil {
		panic(errors.New("gridFS未初始化"))
	} else if file, err := this.gridFS.Create(filename); err == nil {
		if _, werr := file.Write(data); werr == nil {
			if cerr := file.Close(); cerr == nil {
				return
			} else {
				panic(cerr)
			}
		} else {
			panic(werr)
		}
	} else {
		panic(err)
	}
}

func (this *BaseSession) ReadFile(filename string) *[]byte {
	this.ConnectSessionGridFS()
	defer this.closeSession()
	if this.gridFS == nil {
		panic(errors.New("gridFS未初始化"))
	} else if file, err := this.gridFS.Open(filename); err == nil {
		defer file.Close()
		var data []byte = make([]byte, file.Size())
		if _, rerr := file.Read(data); rerr == nil {

			return &data
		} else {
			//panic(rerr)
			log.Error(rerr, filename)
			return nil
		}
	} else {
		//panic(err)
		log.Error(err, filename)
		return nil
	}
}

func (this *BaseSession) IsExistFile(filename string) bool {
	this.ConnectSessionGridFS()
	defer this.closeSession()
	if this.gridFS == nil {
		panic(errors.New("gridFS未初始化"))
	} else if file, err := this.gridFS.Open(filename); err == nil {
		defer file.Close()
		return true
	} else {
		return false
	}
}

func (this *BaseSession) RemoveFile(filename string) bool {
	this.ConnectSessionGridFS()
	defer this.closeSession()
	if this.gridFS == nil {
		panic(errors.New("gridFS未初始化"))
	} else if err := this.gridFS.Remove(filename); err == nil {
		return true
	} else {
		log.Error(err)
		return false
	}
}

func (this *BaseSession) GetCollection(dbname string, collection string) *mgo.Collection {
	return this.session.DB(dbname).C(collection)
}

func (this *BaseSession) GetGridFS(dbname string, gridfs string) *mgo.GridFS {
	return this.sessionGridFs.DB(dbname).GridFS(gridfs)
}

func (this *BaseSession) closeSession() {
	if this.session != nil {
		this.session.Close()
	} else if this.sessionGridFs != nil {
		this.sessionGridFs.Close()
	} else {
		panic(errors.New("没有可关闭的session"))
	}
}

/*
使用方法：
	user:=dao.SelectByID(user.Id_)
	user.(*mhj_models.User)
*/
func (this *BaseSession) SelectByID(id bson.ObjectId) interface{} {
	this.initConnect()
	defer this.closeSession()
	result := reflect.New(this.type_).Interface()
	err := this.collection.FindId(id).One(result)
	if err != nil {
		log.Error(err)
		//panic(err)
		return nil
	}
	return result
}

func (this *BaseSession) SelectOne(query interface{}) interface{} {
	this.initConnect()
	defer this.closeSession()
	result := reflect.New(this.type_).Interface()
	err := this.collection.Find(query).One(result)
	if err != nil {
		log.Error(err, query)
		//panic(err)
		return nil
	}
	return result
}

func (this *BaseSession) SelectAll(query interface{}, result interface{}) bool {
	this.initConnect()
	defer this.closeSession()
	err := this.collection.Find(query).All(result)
	if err != nil {
		log.Error(err, query)
		//panic(err)
		return false
	}
	return true
}

func (this *BaseSession) InsertModel(m interface{}) bool {
	this.initConnect()
	defer this.closeSession()
	rv := reflect.ValueOf(m)
	if rv.Kind() == reflect.Ptr {
		prv := rv.Elem()
		if prv.Kind() == reflect.Ptr {
			panic("投入的对象不应该是双重指针的对象！！！！")
		}
	}
	err := this.collection.Insert(m)
	if err != nil {
		log.Error(err)
	}
	return err == nil
}

func (this *BaseSession) InsertManyModel(ptrArray interface{}) bool {
	this.initConnect()
	defer this.closeSession()
	rv := reflect.ValueOf(ptrArray)
	if rv.Kind() != reflect.Ptr {
		panic("投入的对象应该是指针对象数组！！！！ 当前类型：" + rv.Kind().String())
	} else if rv.Elem().Kind() != reflect.Array && rv.Elem().Kind() != reflect.Slice {
		panic("投入的对象应该是指针对象数组！！！！ 当前类型：" + rv.Kind().String())
	}
	objs := make([]interface{}, rv.Elem().Len())
	for i := 0; i < rv.Elem().Len(); i++ {
		objs[i] = rv.Elem().Index(i).Interface()
	}
	err := this.collection.Insert(objs...)
	if err != nil {
		log.Error(err)
	}
	return err == nil
}

func (this *BaseSession) UpdateModel(id bson.ObjectId, m interface{}) bool {
	this.initConnect()
	defer this.closeSession()
	rv := reflect.ValueOf(m)
	if rv.Kind() == reflect.Ptr {
		prv := rv.Elem()
		if prv.Kind() == reflect.Ptr {
			panic("投入的对象不应该是双重指针的对象！！！！")
		}
	}
	err := this.collection.Update(bson.M{"_id": id}, m)
	if err != nil {
		log.Error(err)
	}
	return err == nil
}

func (this *BaseSession) FindRef(ref *mgo.DBRef) interface{} {
	this.initConnect()
	defer this.closeSession()
	if ref.Id == nil {
		return nil
	}
	result := reflect.New(this.type_).Interface()
	if ref.Database != DataBaseName {
		ref.Database = DataBaseName
	}
	if err := this.session.FindRef(ref).One(result); err == nil {
		return result
	} else {
		log.Error(err)
		return nil
	}
}

func GetIndexName(field string) string {
	return fmt.Sprintf("_%s_", strings.Replace(strings.Replace(field, ".", "_", -1), ":", "_", -1))

}

func (this *BaseSession) EnsureIndex(index ...mgo.Index) {
	this.initConnect()
	defer this.closeSession()
	for _, i := range index {
		this.collection.EnsureIndex(i)
	}
}

func (this *BaseSession) NewIndex(fieldName string, unique bool) mgo.Index {
	////检查索引匹配
	//moarray := getFieldMoArray(this.type_)
	//if !arrayString.Contains(*moarray, fieldList[0]) {
	//	log.Panic(errors.New(fmt.Sprintf("Index 字段不匹配！ field=%s Type=%s \n", fieldName, this.type_.Name())))
	//	return mgo.Index{}
	//}
	fieldList := strings.Split(fieldName, ".")
	if !this.indexNameCheck(fieldName, fieldList, this.type_) {
		log.Panic(errors.New(fmt.Sprintf("Index 字段不匹配！ field=%s Type=%s \n", fieldName, this.type_.Name())))
	}

	return mgo.Index{
		Key:        []string{fieldName},
		Unique:     unique,
		DropDups:   false,
		Background: false, // See notes.
		Sparse:     false,
		Name:       GetIndexName(fieldName),
	}
}

func (this *BaseSession) indexNameCheck(fieldName string, field []string, type_ reflect.Type) bool {
	if field[0] == "" {
		panic("错误的参数投入 ")
	}
	var structField *reflect.StructField
	if type_.Kind() == reflect.Ptr {
		type_ = type_.Elem()
	}
	for i := 0; i < type_.NumField(); i++ {
		//fmt.Println(type_.Field(i).Name,type_.Field(i).Type,type_.Field(i).Type.Kind())
		if type_.Field(i).Tag.Get("bson") == field[0] {
			s := type_.Field(i)
			structField = &s
			break
		} else if type_.Field(i).Type.Kind() == reflect.Struct && this.indexNameCheck(fieldName, field, type_.Field(i).Type) {
			return true
		} else if type_.Field(i).Tag.Get("bson") == "" {
			if type_.Name() != "Time" {
				log.Warn("处理 %s 类型 %s field %s 未指定属性 bson", fieldName, type_.Name(), type_.Field(i).Name)
			}
		}
	}
	if structField != nil {
		if len(field) > 1 {
			return this.indexNameCheck(fieldName, field[1:], structField.Type)
		} else {
			return true
		}
	}
	return false
}

func (this *BaseSession) NewMultiIndex(unique bool, fieldName ...string) mgo.Index {
	//moarray := getFieldMoArray(this.type_)
	//for _, f := range fieldName {
	//	if !arrayString.Contains(*moarray, f) {
	//		log.Panic(errors.New(fmt.Sprintf("Index 字段不匹配！ field=%s Type=%s \n", fieldName, this.type_.Name())))
	//		return mgo.Index{}
	//	}
	//}
	for _, f := range fieldName {
		fieldList := strings.Split(f, ".")
		if !this.indexNameCheck(f, fieldList, this.type_) {
			log.Panic(errors.New(fmt.Sprintf("Index 字段不匹配！ field=%s Type=%s \n", fieldName, this.type_.Name())))
		}
	}

	return mgo.Index{
		Key:        fieldName,
		Unique:     unique,
		DropDups:   false,
		Background: false, // See notes.
		Sparse:     false,
		Name:       GetIndexName(fieldName[0]),
	}
}

func (this *BaseSession) New2dsphereIndex(fieldName string, unique bool) mgo.Index {

	//moarray := getFieldMoArray(this.type_)
	//if !arrayString.Contains(*moarray, fieldName) {
	//	log.Panic(errors.New(fmt.Sprintf("Index 字段不匹配！ field=%s Type=%s \n", fieldName, this.type_.Name())))
	//	return mgo.Index{}
	//}
	fieldList := strings.Split(fieldName, ".")
	if !this.indexNameCheck(fieldName, fieldList, this.type_) {
		log.Panic(errors.New(fmt.Sprintf("Index 字段不匹配！ field=%s Type=%s \n", fieldName, this.type_.Name())))
	}

	return mgo.Index{
		Key:        []string{fmt.Sprintf("$2dsphere:%s", fieldName)},
		Unique:     unique,
		DropDups:   false,
		Background: false, // See notes.
		Sparse:     false,
		Name:       GetIndexName(fieldName),
	}
}

//func getFieldMoArray(m reflect.Type) *[]string {
//	w := []string{}
//
//	var fieldName string
//	for i := 0; i < m.NumField(); i++ {
//		tags := strings.Split(string(m.Field(i).Tag), "\"")
//		if len(tags) > 1 {
//			fieldName = tags[1]
//
//		} else {
//			fieldName = m.Field(i).Name
//
//		}
//		appendField(m.Field(i), "", fieldName, &w)
//	}
//	return &w
//}

func createChildField(parentFieldName string, childtype reflect.Type, w *[]string) {
	s := reflect.New(childtype).Elem()
	if s.Type().Kind() == reflect.Array || s.Type().Kind() == reflect.Slice {
		return
	}
	typeOfT := childtype
	//fmt.Println("createChildField",s, typeOfT)
	fieldName := ""

	for i := 0; i < s.NumField(); i++ {

		tags := strings.Split(string(typeOfT.Field(i).Tag), "\"")
		if len(tags) > 1 {
			fieldName = tags[1]

		} else {
			fieldName = typeOfT.Field(i).Name
		}
		appendField(typeOfT.Field(i), "", parentFieldName+fieldName, w)
	}
}

func appendField(field reflect.StructField, parentFieldName, fieldName string, w *[]string) {
	if field.Type.String() == "mgo.DBRef" {
		//联级查询
		*w = append(*w, fieldName)
		*w = append(*w, fieldName+".$ref")
		*w = append(*w, fieldName+".$id")
		*w = append(*w, fieldName+".$db")
	} else {
		*w = append(*w, fieldName)

		typename := fmt.Sprint(field.Type)
		if strings.Contains(typename, "mhj_models") {
			if parentFieldName == "" {
				parentFieldName = fieldName + "."
			} else {
				parentFieldName = parentFieldName + fieldName + "."
			}

			if strings.Index(typename, "*") == 0 {
				createChildField(parentFieldName, field.Type.Elem(), w)
			} else {
				createChildField(parentFieldName, field.Type, w)
			}
		}
	}
}

func (this *BaseSession) Run(run func(collection *mgo.Collection) interface{}) interface{} {
	this.initConnect()
	defer this.closeSession()
	return run(this.collection)
}

func (this *BaseSession) FindCount(query interface{}) (int, error) {
	this.initConnect()
	defer this.closeSession()
	return this.collection.Find(query).Count()
}

func (this *BaseSession) FindOne(query interface{}, one interface{}) error {
	this.initConnect()
	defer this.closeSession()
	return this.collection.Find(query).One(one)
}

func (this *BaseSession) FindAll(query interface{}, one interface{}) error {
	this.initConnect()
	defer this.closeSession()
	return this.collection.Find(query).All(one)
}

func (this *BaseSession) Update(selector interface{}, update interface{}) error {
	this.initConnect()
	defer this.closeSession()
	return this.collection.Update(selector, update)
}

func (this *BaseSession) UpdateAll(selector interface{}, update interface{}) error {
	this.initConnect()
	defer this.closeSession()
	info, err := this.collection.UpdateAll(selector, update)
	if err != nil {
		log.Info(info)
		log.Error(err)
		return err
	}
	return nil
}

func (this *BaseSession) Remove(selector interface{}) error {
	this.initConnect()
	defer this.closeSession()
	return this.collection.Remove(selector)
}

func (this *BaseSession) RemoveAll(selector bson.M) error {
	this.initConnect()
	defer this.closeSession()
	info, err := this.collection.RemoveAll(selector)
	if err != nil {
		log.Error(err, info)
	}
	return err
}

func (this *BaseSession) initConnect() {
	this.session = this.GetSession()
	this.collection = this.GetCollection(DataBaseName, this.collectionName)
}

func (this *BaseSession) FindPage(selector bson.M, sortFieldName string, time time.Time, result interface{}) error { //通过时间分页
	this.initConnect()
	defer this.closeSession()

	if selector == nil {
		selector = bson.M{}
	}
	selector[sortFieldName] = bson.M{MGO_SELECT_LT: time}

	//selector[sortFieldName]=bson.M{MGO_SELECT_GT:time}
	log.Info(time)
	return this.collection.Find(selector).Sort(OrderDesc(sortFieldName)).Limit(10).All(result)
}

func (this *BaseSession) FindPageByCount(selector bson.M, count int, result interface{}) error {
	this.initConnect()
	defer this.closeSession()

	return this.collection.Find(selector).Skip(count).Limit(10).All(result)
}

func (this *BaseSession) FindPageSortByField(selector bson.M, count int, sortFieldName string, result interface{}) error {
	this.initConnect()
	defer this.closeSession()

	return this.collection.Find(selector).Sort(sortFieldName).Skip(count).Limit(10).All(result)
}

func (this *BaseSession) FindPageSortByManyField(selector bson.M, count int, result interface{}, sortFieldName ...string) error {
	this.initConnect()
	defer this.closeSession()

	return this.collection.Find(selector).Sort(sortFieldName...).Skip(count).Limit(DEFAULT_LIMIT_NUM).All(result)
}

func (this *BaseSession) FindOneAndSelect(query interface{}, selectquery interface{}, one interface{}) error {
	this.initConnect()
	defer this.closeSession()
	return this.collection.Find(query).Select(selectquery).One(one)
}

func (this *BaseSession) MapReduce(job *mgo.MapReduce, result interface{}, sortFieldNames ...string) error {
	this.initConnect()
	defer this.closeSession()
	mapreduceinfo, err := this.collection.Find(nil).Sort(sortFieldNames...).MapReduce(job, result)
	if _, ok := jobAddUpCounter[job.Map]; !ok {
		if job.Map != "" {
			jobAddUpCounter[job.Map] = 0
		}
	}
	jobAddUpCounter[job.Map]++
	count := jobAddUpCounter[job.Map]
	if (count-1)%LOG_FREQUENCY == 0 {
		log.Infof("Run a total of %d times, the most recent record is \n%#v", count, mapreduceinfo)
	}
	return err
}

//生成一个对应的dbref
func (this *BaseSession) MakeDBRef(m interface{}) *mgo.DBRef {
	return &mgo.DBRef{Database: DataBaseName, Id: m, Collection: this.collectionName}
}

func (this *BaseSession) FindAndModify(selector interface{}, change mgo.Change) (interface{}, error) {
	this.initConnect()
	defer this.closeSession()

	result := reflect.New(this.type_).Interface()
	_, err := this.collection.Find(selector).Apply(change, result)
	//fmt.Println(result)
	if err != nil {
		return nil, err
	} else {
		return result, nil
	}

}

func (this *BaseSession) FindBySortAndModify(selector interface{}, change mgo.Change, sortField string, limit int) (interface{}, error) {
	this.initConnect()
	defer this.closeSession()

	result := reflect.New(this.type_).Interface()
	_, err := this.collection.Find(selector).Sort(sortField).Limit(limit).Apply(change, result)
	//fmt.Println(result)
	if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

//待修改
//func (this *BaseSession) FindAndModifyRef(selector,ref *mgo.DBRef,change *mgo.Change) interface{}{
//	this.initConnect()
//	defer this.closeSession()
//
//	if ref.Id==nil{
//		return nil
//	}
//	result := reflect.New(this.type_).Interface()
//	if ref.Database != DataBaseName {
//		ref.Database = DataBaseName
//	}
//
//
//
//	c :=this.session.DB(ref.Database).C(ref.Collection)
//
//	if _,err := c.Find(bson.D{{"_id", ref.Id},}).Apply(*change,result);err ==nil{
//		return result
//	}else{
//		log.Error(err)
//		return nil
//	}
//
//	if _,err := this.session.FindRef(ref).Apply(*change,result);err ==nil{
//		return result
//	}else{
//		log.Error(err)
//		return nil
//	}
//
//}

func (this *BaseSession) SelectByLimit(query interface{}, limit int, result interface{}) bool {
	this.initConnect()
	defer this.closeSession()
	err := this.collection.Find(query).Limit(limit).All(result)
	if err != nil {
		log.Error(err, query)
		return false
	}
	return true
}

func (this *BaseSession) SelectBySort(query interface{}, sort string, limit int, result interface{}) bool {
	this.initConnect()
	defer this.closeSession()
	err := this.collection.Find(query).Sort(sort).Limit(limit).All(result)
	if err != nil {
		log.Error(err, query)
		//panic(err)
		return false
	}
	return true
}

func (this *BaseSession) Upsert(selector interface{}, data interface{}) error {
	this.initConnect()
	defer this.closeSession()

	_, err := this.collection.Upsert(selector, data)
	return err
}

func (this *BaseSession) Aggregate(pipeline interface{}, results interface{}) error {
	this.initConnect()
	defer this.closeSession()

	return this.collection.Pipe(pipeline).Iter().All(results)
}

func (this *BaseSession) CreateMatchFields(structobj interface{}, options ...QueryOption) bson.M {
	result := make(bson.M, DEFAULT_LIMIT_NUM)
	v := reflect.ValueOf(structobj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for _, option := range options {
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).Tag.Get(TAG_BSON) == option.MdDefine {
				result[option.MdDefine] = bson.M{option.Operator: v.Field(i).Interface()}
				break
			}
		}
	}
	return result
}

func Assemblytag(parentTag string, childTag ...string) string {
	l := make([]string, 0, len(childTag)+1)
	if parentTag == "" {
		l = append(l, childTag...)
		return strings.Join(l, TAG_JOIN_SPLIT_SYMBAL)
	}
	l = append(l, parentTag)
	l = append(l, childTag...)
	return strings.Join(l, TAG_JOIN_SPLIT_SYMBAL)
}

//func (this *BaseSession) CreateMatchFields_V1(structobj interface{}, parentTag string, options ...QueryOption) bson.M {
//	var result bson.M
//	if parentTag == "" {
//		result = make(bson.M, DEFAULT_LIMIT_NUM)
//	} else {
//		//查询子结构传递单个option
//		result = make(bson.M, 1)
//	}
//
//	v := reflect.ValueOf(structobj)
//	v = reflect.Indirect(v)
//	for _, option := range options {
//		for i := 0; i < v.NumField(); i++ {
//			topTag := Assemblytag(parentTag, v.Type().Field(i).Tag.Get(TAG_BSON))
//			if topTag == option.MdDefine {
//				result[option.MdDefine] = bson.M{option.Operator: v.Field(i).Interface()}
//				break
//			}
//
//			if v.Field(i).Kind() == reflect.Struct || v.Field(i).Kind() == reflect.Ptr {
//				// Time结构不进行递归
//				_, ok := v.Field(i).Interface().(time.Time)
//				if ok {
//					continue
//				}
//
//				tags := strings.Split(option.MdDefine, TAG_JOIN_SPLIT_SYMBAL)
//				// 不包含此嵌套结构体的字段前缀,跳过递归匹配
//				if !arrayString.Contains(tags, topTag) {
//					continue
//				}
//
//				currentObj := v.Field(i).Interface()
//				// 递归嵌套结构体
//				match := this.CreateMatchFields_V1(currentObj, topTag, option)
//				hit := len(match) != 0
//				if hit {
//					result[option.MdDefine] = match[option.MdDefine]
//					break
//				} else {
//					continue
//				}
//			}
//		}
//	}
//	return result
//}

func (this *BaseSession) InsertWithoutOrdered(ptrArray interface{}) (error, *mgo.BulkResult) {
	this.initConnect()
	defer this.closeSession()
	rv := reflect.ValueOf(ptrArray)
	if rv.Kind() != reflect.Ptr {
		panic("投入的对象应该是指针对象数组！！！！ 当前类型：" + rv.Kind().String())
	} else if rv.Elem().Kind() != reflect.Array && rv.Elem().Kind() != reflect.Slice {
		panic("投入的对象应该是指针对象数组！！！！ 当前类型：" + rv.Kind().String())
	}
	objs := make([]interface{}, rv.Elem().Len())
	for i := 0; i < rv.Elem().Len(); i++ {
		objs[i] = rv.Elem().Index(i).Interface()
	}
	//bulk := this.GetCollection(DataBaseName, this.collectionName).Bulk()
	bulk := this.collection.Bulk()
	bulk.Unordered()
	bulk.Insert(objs...)
	bulkRes, bulkErr := bulk.Run()
	if bulkErr != nil {
		return bulkErr, bulkRes
	}
	return nil, bulkRes
}

func (this *BaseSession) SelectByPagination(selector interface{}, pageIndex int, pageSize int) interface{} {
	if pageIndex <= 0 {
		pageIndex = 1
	}
	return this.Run(func(collection *mgo.Collection) interface{} {
		slice := reflect.MakeSlice(reflect.SliceOf(this.type_), pageSize, pageSize)
		x := reflect.New(slice.Type()).Interface()
		err := collection.Find(selector).Skip((pageIndex - 1) * pageSize).Limit(pageSize).All(x)
		if err != nil {
			log.Error(err)
			return nil
		}
		return x
	})
}

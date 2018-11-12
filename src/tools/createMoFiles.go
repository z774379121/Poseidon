package main

import (
	"bufio"
	"fmt"
	"log"
	"models"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"xm/common/disk"
	"xm/common/shell"
)

type modules struct {
	songs models.Songs
}

var pathMO = "/mdb/src/models/modelsDefine"
var path string

func main() {
	fmt.Println("模型枚举生成工具!")

	path = shell.GetCurrentPath()
	fmt.Printf("当前路径：%s\n", path)
	path = strings.Replace(path, "\n", "", -1) + pathMO

	fmt.Printf("检测路径：%s\n", path)
	if !disk.IsDirExists(path) {
		fmt.Printf("未找到modelsDefine文件夹，请确认路径是否正确！\n")
		return
	}

	m := modules{}
	s := reflect.ValueOf(&m).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		createMOFile(typeOfT.Field(i).Type.String(), f)
	}
}

var digitsRegexp = regexp.MustCompile(`(\d+)\D+(\d+)`)

func createMOFile(typename string, m reflect.Value) {
	typeOfT := m.Type()
	classname := substr(typename, strings.Index(typename, ".")+1, 0xffff)

	filename := path + "/md" + classname + ".go"
	if disk.IsFileExists(filename) {
		os.Remove(filename)
	}
	//fmt.Printf("写入文件：%s\n", filename)
	if f, err1 := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666); err1 != nil {
		fmt.Printf("创建文件%s失败\n", filename)
	} else {

		w := bufio.NewWriter(f) //创建新的 Writer 对象

		if _, err3 := w.WriteString("package modelsDefine\n"); err3 != nil {
			fmt.Printf("写入文件失败%s失败\n", filename)
		}

		w.WriteString("/*\n		模型枚举文件，请勿修改！\n*/\n")

		//linestring := ""
		fieldName := ""
		for i := 0; i < m.NumField(); i++ {
			//fmt.Print("TypeName=", typeOfT.Field(i).Type.Name())
			tags := strings.Split(string(typeOfT.Field(i).Tag), "\"")
			if len(tags) > 1 {
				fieldName = tags[1]

			} else {
				fieldName = typeOfT.Field(i).Name

			}
			writeField(classname, typeOfT.Field(i), "", fieldName, w, filename)

		}
		w.Flush()
		f.Close()
	}
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func createChildField(classname, parentFieldName string, childtype reflect.Type, w *bufio.Writer, filename string) {
	s := reflect.New(childtype).Elem()
	if s.Type().Kind() <= reflect.String {
		return
	}
	typeOfT := childtype
	//fmt.Println(s, typeOfT)
	fieldName := ""
	for i := 0; i < s.NumField(); i++ {
		tags := typeOfT.Field(i).Tag
		if bsonName := tags.Get("bson"); len(bsonName) > 0 {
			fieldName = bsonName
		} else {
			fieldName = typeOfT.Field(i).Name
		}
		writeField(classname, typeOfT.Field(i), "", parentFieldName+fieldName, w, filename)
	}
}

func writeField(classname string, field reflect.StructField, parentFieldName, fieldName string, w *bufio.Writer, filename string) {
	linestring := ""
	if strings.Contains(field.Type.String(), "mgo.DBRef") {
		//联级查询
		linestring = fmt.Sprintf("const MD%s_%s string = \"%s\"\n", classname, field.Name, fieldName)
		if _, err3 := w.WriteString(linestring); err3 != nil {
			fmt.Printf("写入文件失败%s失败\n", filename)
		}

		linestring = fmt.Sprintf("const MD%s_%s_ref string = \"%s.$ref\"\n", classname, field.Name, fieldName)
		if _, err3 := w.WriteString(linestring); err3 != nil {
			fmt.Printf("写入文件失败%s失败\n", filename)
		}
		//fmt.Printf("写入内容： %s", linestring)
		linestring = fmt.Sprintf("const MD%s_%s_id string = \"%s.$id\"\n", classname, field.Name, fieldName)
		if _, err3 := w.WriteString(linestring); err3 != nil {
			fmt.Printf("写入文件失败%s失败\n", filename)
		}
		//fmt.Printf("写入内容： %s", linestring)
		linestring = fmt.Sprintf("const MD%s_%s_db string = \"%s.$db\"\n", classname, field.Name, fieldName)
		if _, err3 := w.WriteString(linestring); err3 != nil {
			fmt.Printf("写入文件失败%s失败\n", filename)
		}
		//fmt.Printf("写入内容： %s", linestring)
	} else {
		linestring = fmt.Sprintf("const MD%s_%s string = \"%s\"\n", classname, field.Name, fieldName)
		if _, err3 := w.WriteString(linestring); err3 != nil {
			fmt.Printf("写入文件失败%s失败\n", filename)
		}
		//fmt.Printf("写入内容： %s", linestring)

		typename := field.Type.String()
		if strings.Contains(typename, "Tile") {
			//fmt.Println(typename)
		}
		if strings.Contains(typename, "models") {
			stypeName := strings.Split(typename, ".")
			clsName := ""
			if stypeName[1] == fieldName {
				clsName = classname
			} else {
				if parentFieldName == "" {
					parentFieldName = fieldName + "."
				} else {
					parentFieldName = parentFieldName + fieldName + "."
				}
				clsName = classname + "_" + field.Name
			}

			if strings.Index(typename, "*") == 0 {
				createChildField(clsName, parentFieldName, field.Type.Elem(), w, filename)
			} else {
				createChildField(clsName, parentFieldName, field.Type, w, filename)
			}
		}
	}
}

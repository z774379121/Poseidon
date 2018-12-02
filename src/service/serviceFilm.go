package service

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"github.com/z774379121/untitled1/src/xm/common"
	"net/http"
	"github.com/z774379121/untitled1/src/dao"
	"os"
	"io"
	"bytes"
)

func UpLoad(context echo.Context) error {
	avatar, err := context.FormFile("avatar")
	if err != nil {
		return err
	}

	// Source
	fmt.Println(avatar.Filename, common.FileSize(avatar.Size))
	src, err := avatar.Open()
	if err != nil {
		return err
	}

	defer src.Close()

	des := make([]byte, avatar.Size)
	n, err2 := src.Read(des)
	if err2 != nil {
		fmt.Println(err2)
		return context.String(http.StatusBadRequest, "读取失败")
	}
	fmt.Println(n)
	// Destination
	daoFile := dao.NewDaoBTFile()
	ok, name := daoFile.UploadBTFile(bson.NewObjectId(), &des)
	if !ok {
		return context.String(http.StatusBadRequest, "插入到数据库失败")
	}
	daoFileContent := dao.NewDaoBTFileContent()
	if daoFileContent.UpdateRealFileName(name, avatar.Filename) {
		return context.JSON(http.StatusOK, map[string]interface{}{
			"msg":"ok",
			"name":name,
		})
	}
	return context.String(http.StatusBadRequest, "更新名字失败")
}

func Download(context echo.Context) error {
	filename := context.Param("filename")
	if !bson.IsObjectIdHex(filename) {
		return context.String(http.StatusForbidden, "invaild filename")
	}

	daoFile := dao.NewDaoBTFile()
	img := daoFile.DownloadBTFile(filename)
	if img == nil {
		return context.String(http.StatusForbidden, "not found")
	}

	daoFileContent := dao.NewDaoBTFileContent()
	imgInfo := daoFileContent.SelectByName(filename)
	filename = filename + ".torrent"
	if imgInfo.RealFileName != "" {
		filename = imgInfo.RealFileName
	}
	fmt.Println(imgInfo.RealFileName)
	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()
	defer os.Remove(filename)

	// Copy
	if _, err = io.Copy(dst, bytes.NewReader(*img)); err != nil {
		return err
	}
	context.Response().Header().Set("Content-Type", "application/octet-stream")
	context.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	return context.File(filename)
}
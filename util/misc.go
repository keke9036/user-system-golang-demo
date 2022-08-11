// @Description
// @Author weitao.yin@shopee.com
// @Since 2022/6/13

package util

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"strings"
)

func GenerateSid() string {
	return uuid.NewV4().String()
}

func IntToStr(uid uint64) string {
	return strconv.FormatInt(int64(uid), 10)
}

func StrToInt(str string) (uint64, error) {
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint64(val), nil
}

func ReadFileByte(file *multipart.FileHeader) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, src); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Md5(source string) string {
	hash := md5.Sum([]byte(source))
	return hex.EncodeToString(hash[:])
}

func GetFileExt(name string) string {
	return path.Ext(name)
}

func SaveFile(file *[]byte, dst string, dir string) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return errors.New("UploadFile: failed to create save directory")
		}
	}
	err = ioutil.WriteFile(dst, *file, 0644)
	return err
}

func BuildFileName(name string) string {
	ext := GetFileExt(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = Md5(fileName)

	return fileName + ext
}

func DeleteStringElement(list []string, ele string) []string {
	result := make([]string, 0)
	for _, v := range list {
		if v != ele {
			result = append(result, v)
		}
	}
	return result
}

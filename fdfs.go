package main

/*
#cgo LDFLAGS: -lfdfs -lfastcommon -lfdfsclient  -lpthread -ldl -L/usr/local/lib
#cgo CFLAGS: -I/usr/local/include/fastcommon -I/usr/local/include/fastdfs
#include "fdfs.h"
*/
import "C"
import "fmt"
import "path/filepath"
import "errors"

//上传文件到fdfs
//conf:配置文件的文件,example:/etc/fdfs/client.conf
//imagePath:需要上传文件的完整路径,example:"/root/Desktop/logo.jpg"
func FdfsUploadFile(conf string, imagePath string) (result map[string]string, err error) {

	result = make(map[string]string)
	var resData C.responseData = C.upload_file(C.CString(conf), C.CString(imagePath))

	// fmt.Println("upload file msg:", C.GoString(resData.msg)) ////当成功的时候，是返回图片的id,example:group1/M00/00/00/wKgBP1NxvSqH9qNuAAAED6CzHYE179.jpg ,当失败的时候是返回错误消息
	// fmt.Println("upload file result:", resData.result)       //１表示成功，０表示失败

	if resData.result == 0 { //上传失败, 返回例子:{"r":{"r":false}}

		err = errors.New(C.GoString(resData.msg))
		return
	} else { //上传成功, 返回例子:"filename":"scree.jpg","group":"group1","url":"M00\/00\/00\/wKgBP1NxvIPf81_1AABJuZk6wJM879.jpg"

		filename := filepath.Base(imagePath)

		strPath := C.GoString(resData.msg)

		result["filename"] = filename
		result["path"] = strPath

		return result, nil
	}

}

//删除fdfs文件
//conf:配置文件的文件,example:/etc/fdfs/client.conf
//imagePath:图片的id,example:"group1/M00/00/00/wKgBP1Nx1S7bSab8AAAED6CzHYE352.jpg"
func FdfsDeleteFile(conf string, fileId string) (result map[string]interface{}, err error) {

	result = make(map[string]interface{})
	var resData C.responseData = C.delete_file(C.CString(conf), C.CString(fileId))

	fmt.Println("upload file msg:", C.GoString(resData.msg)) ////当成功的时候，是返回图片的id,example:group1/M00/00/00/wKgBP1NxvSqH9qNuAAAED6CzHYE179.jpg ,当失败的时候是返回错误消息
	fmt.Println("upload file result:", resData.result)       //１表示成功，０表示失败

	if resData.result == 0 { //删除文件失败,
		err = errors.New(C.GoString(resData.msg))
		return
	} else { //删除文件成功,
		result["msg"] = C.GoString(resData.msg)
		return result, nil
	}

}

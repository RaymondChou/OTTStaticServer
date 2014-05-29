package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"
)

type RespT struct {
	Ret int    `json:"ret"`
	Url string `json:"url"`
}

//This is where the action happens.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		response(w, RespT{Ret: 105})

	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		r.ParseMultipartForm(50 << 20)

		access_token := r.Header.Get("access_token")
		if access_token == "" || access_token != "9iui23o48gnklnvyeqiu313pob042" {
			response(w, RespT{Ret: 104})
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			response(w, RespT{Ret: 100})
			return
		}
		defer file.Close()
		fileName := handler.Filename
		filePath := "/tmp/" + fmt.Sprintf("%s", time.Now().Unix()) + fileName
		dst, err := os.Create(filePath)
		defer dst.Close()

		if err != nil {
			response(w, RespT{Ret: 101})
			return
		}
		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			response(w, RespT{Ret: 102})
			return
		}

		//FdfsUploadFile
		result, err := FdfsUploadFile("/etc/fdfs/client.conf", filePath)
		if err != nil {
			// fdfs save error
			os.Remove(filePath)
			response(w, RespT{Ret: 103})
		} else {
			os.Remove(filePath)
			url := "http://127.0.0.1:8080/" + result["path"]
			response(w, RespT{Ret: 0, Url: url})
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func response(w http.ResponseWriter, responseData interface{}) {
	b, _ := json.Marshal(responseData)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.HandleFunc("/upload", uploadHandler)

	//Listen on port 9090
	http.ListenAndServe(":9090", nil)
	//上传文件
	//str,err:=FdfsUploadFile("/etc/fdfs/client.conf","README.md")

	//删除文件
	// str, err := FdfsDeleteFile("/etc/fdfs/client.conf", "group1/M00/00/01/wKgBP1N26sD3333335g63rAAAFcd0RcrU7478.md")
	// fmt.Println("file upload result:", str)
	// fmt.Println("file upload err:", err)
}

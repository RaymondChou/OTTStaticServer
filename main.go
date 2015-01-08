package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

const (
	OTTServerVerifyTokenAPI = "http://172.18.14.20:8080/ott/verify_token"
)

type RespT struct {
	Ret int    `json:"ret"`
	Url string `json:"url"`
}

type ResultT struct {
	Ret int `json:"ret"`
}

func safeHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				// logging
				fmt.Println("WARN: panic fired in %v.panic - %v", fn, e)
				fmt.Println(string(debug.Stack()))
			}
		}()
		fn(w, r)
	}
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
		if access_token == "" || !checkToken(access_token) {
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
			url := "http://112.4.28.92:8081/" + result["path"]
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

func checkToken(access_token string) bool {
	v := make(url.Values)
	v.Set("access_token", access_token)

	if res, err := http.PostForm(OTTServerVerifyTokenAPI, v); err == nil {
		result, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Printf("%v", err)
			return false
		}

		var resultT ResultT

		err = json.Unmarshal(result, &resultT)
		if err != nil {
			fmt.Printf("%v", err)
			return false
		}

		if resultT.Ret != 0 {
			return false
		} else {
			return true
		}

	} else {
		fmt.Printf("%v", err)
		return false
	}

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	http.HandleFunc("/upload", safeHandler(uploadHandler))

	//Listen on port 9090
	http.ListenAndServe(":8899", nil)
	//上传文件
	//str,err:=FdfsUploadFile("/etc/fdfs/client.conf","README.md")

	//删除文件
	// str, err := FdfsDeleteFile("/etc/fdfs/client.conf", "group1/M00/00/01/wKgBP1N26sD3333335g63rAAAFcd0RcrU7478.md")
	// fmt.Println("file upload result:", str)
	// fmt.Println("file upload err:", err)
}

// func FdfsUploadFile(a string, b string) (c map[string]string, d error) {
// 	return
// }

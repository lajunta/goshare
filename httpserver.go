package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
)

var (
	uploadPath, downloadPath string
)

func desktop() string {
	var desk string
	currentUser, _ := user.Current()
	home := currentUser.HomeDir
	d1 := path.Join(home, "桌面")
	d2 := path.Join(home, "Desktop")
	if os.MkdirAll(d1, 0777) == nil {
		desk = d2
	} else {
		desk = d1
	}
	return desk
}

func checkDir() {
	downloadPath = path.Join(desktop(), "下载")
	uploadPath = path.Join(desktop(), "上传")
	os.MkdirAll(downloadPath, 0777)
	os.MkdirAll(uploadPath, 0777)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("views/index.html")
	zys := listZuoye()
	t.Execute(w, zys)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tpl, _ := template.ParseFiles("views/upload.html")
		tpl.Execute(w, nil)
	} else {
		ufile, handler, _ := r.FormFile("uploadFile")
		f, err := os.OpenFile(uploadPath+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, ufile)
		ipArray := strings.Split(r.RemoteAddr, ":")
		remoteIp := ipArray[0]
		userName := r.FormValue("userName")
		fileName := handler.Filename
		//fmt.Fprintf(w, "http header : %v \n", handler.Header)

		zy := Zuoye{remoteIp, userName, fileName, "0"}
		err = zy.save()
		if err != nil {
			fmt.Fprint(w, "无法保存")
		} else {
			fmt.Fprintf(w, "%v  已经保存", zy)
		}
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	println(ip)
	err := rmZuoye(ip)
	if err == nil {
		fmt.Fprint(w, " 删除成功 ")
	} else {
		fmt.Fprint(w, " 删除失败 ")
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("ip")
	println(ip)
	zy := getZuoye(ip)
	fmt.Fprintf(w, "%v", zy)
}

func server() {
	checkDir()
	println(downloadPath)
	http.HandleFunc("/", indexHandler)
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir(downloadPath))))
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/delete", deleteHandler)
	http.HandleFunc("/get", getHandler)
	http.ListenAndServe(":8080", nil)
}

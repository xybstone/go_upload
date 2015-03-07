package main

import (
	"github.com/Unknwon/goconfig"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	CFG_PATH = "app.conf"
)

var Cfg *goconfig.ConfigFile
var path string
var port string

var uploadTemplate = template.Must(template.ParseFiles("index.html"))

func indexHandle(w http.ResponseWriter, r *http.Request) {
	if err := uploadTemplate.Execute(w, nil); err != nil {
		log.Fatal("Execute: ", err.Error())
		return
	}
}

func uploadHandle(w http.ResponseWriter, r *http.Request) {
	file, fileHead, err := r.FormFile("file")
	if err != nil {
		log.Fatal("FormFile: ", err.Error())
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal("Close: ", err.Error())
			return
		}
	}()

	bytes, err := ioutil.ReadAll(file)

	if err != nil {
		log.Fatal("ReadAll: ", err.Error())
		return
	}
	newfile, err := os.Create(path + fileHead.Filename)
	defer newfile.Close()
	if err == nil {
		_, err = newfile.Write(bytes)
		if err == nil {
			return
		}
	} else {
		log.Println("Creat Err: ", err.Error())
	}
}

func init() {
	var err error
	Cfg, err = goconfig.LoadConfigFile(CFG_PATH)
	if err != nil {
		panic(err)
	}

	path = Cfg.MustValue("path", "rootpath")
	port = Cfg.MustValue("port", "port")
}

func main() {
	http.HandleFunc("/", indexHandle)
	http.HandleFunc("/upload", uploadHandle)
	http.ListenAndServe(port, nil)
}

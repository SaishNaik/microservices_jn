package handlers

import (
	"fmt"
	"io"
	"strconv"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/building-microservices-youtube/product-images/files"
)

// Files is a handler for reading and writing files
type Files struct {
	log   hclog.Logger
	store files.Storage
}

// NewFiles creates a new File handler
func NewFiles(s files.Storage, l hclog.Logger) *Files {
	return &Files{store: s, log: l}
}

// UploadResT implements the http.Handler interface
func (f *Files) UploadRest(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handle POST", "id", id, "filename", fn)

	// no need to check for invalid id or filename as the mux router will not send requests
	// here unless they have the correct parameters

	f.saveFile(id, fn, rw, r.Body)
}

// UploadMultipart something
func (f *Files) UploadMultipart(rw http.ResponseWriter, r *http.Request){
	err := r.ParseMultipartForm(128 * 1024)
	if err != nil{
		f.log.Error("Bad Request","error",err)
		http.Error(rw,"Expected multipart form data",http.StatusBadRequest)
		return
	}

	id,idErr := strconv.Atoi(r.FormValue("id"))
	if idErr != nil{
		f.log.Error("Bad Request","error",err)
		http.Error(rw,"Expected multipart form data",http.StatusBadRequest)
		return
	}

	f.log.Info("Process form for id",id)
	ff,mh,err := r.FormFile("file")
	if err != nil{
		f.log.Error("Bad Request","error",err)
		http.Error(rw,"Expected multipart form data",http.StatusBadRequest)
		return
	}
	fmt.Println(mh.Filename)
	f.saveFile(r.FormValue("id"),mh.Filename,rw,ff)
}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Error("Invalid path", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}

// saveFile saves the contents of the request to a file
func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r io.ReadCloser) {
	f.log.Info("Save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	fmt.Println(fp)
	err := f.store.Save(fp, r)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}

package api

import (
	"os"
	"net/http"
	"strings"
	"path/filepath"
	"encoding/json"
	"io/ioutil"
	"encoding/base64"

	"github.com/gophish/gophish/models"
)

func ls(dir string) (map[string][]string, error) {
	dirlist := make(map[string][]string)
	files, err := os.ReadDir(dir)
	if err != nil {
		return dirlist, err
	}
	var file_names []string
	var dir_names []string
	for _, file := range files {
		if file.IsDir() {
			dir_names = append(dir_names, file.Name())
		} else {
			file_names = append(file_names, file.Name())
		}
	}
	dirlist["files"] = file_names
	dirlist["dirs"] = dir_names
	return dirlist, err
}

func (as *Server) StaticContent(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	path := v.Get("path")
	// Checking for path traversal
	if strings.Contains(path, "..") {
		JSONResponse(w, models.Response{Success: false, Message: "Path blocked"}, http.StatusBadRequest)
		return
	}

	full_path := filepath.Join("./static/endpoint", path)

	switch {
	// Listing dirs in path
	case r.Method == "GET":
		dirlist, err := ls(full_path)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, dirlist, http.StatusOK)
	// Upload file
	case r.Method == "PUT":
		var filename string
		var fileBytes []byte
		// Get filename and contents from body
		if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
			// Multipart body
			r.ParseMultipartForm(10 << 20)
			file, handler, err := r.FormFile("file")
			if err != nil {
				JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
				return
			}
			defer file.Close()
			filename = handler.Filename

			fileBytes, err = ioutil.ReadAll(file)
			if err != nil {
				JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
				return
			}

		} else {
			// JSON body
			var body map[string]interface{}
			err := json.NewDecoder(r.Body).Decode(&body)
			if err != nil {
				JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
				return
			}
			filename = body["filename"].(string)
			fileBytes, err = base64.StdEncoding.DecodeString(body["file"].(string))
			if err != nil {
				JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
				return
			}
		}
		// Catch path traversal
		if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
			JSONResponse(w, models.Response{Success: false, Message: "Invalid filename"}, http.StatusBadRequest)
			return
		}
		// Write uploaded file
		file_path := filepath.Join(full_path, filename)
		err := os.WriteFile(file_path, fileBytes, 0700)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "File uploaded"}, http.StatusCreated)

	// Create directory
	case r.Method == "POST":
		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		var dirname, ok = body["dirname"].(string)
		if !ok {
			JSONResponse(w, models.Response{Success: false, Message: "dirname not in request"}, http.StatusBadRequest)
			return
		} else if strings.Contains(dirname, "..") {
			JSONResponse(w, models.Response{Success: false, Message: "Path blocked"}, http.StatusBadRequest)
			return
		}
		dirpath := filepath.Join(full_path, dirname)
		err = os.Mkdir(dirpath, 0700)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "Directory created"}, http.StatusCreated)

	// Delete file
	case r.Method == "DELETE":
		stat, err := os.Stat(full_path)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusNotFound)
			return
		} else if stat.IsDir() {
			JSONResponse(w, models.Response{Success: false, Message: "Cannot delete directory"}, http.StatusBadRequest)
			return
		}
		err = os.Remove(full_path)
		if err != nil {
			JSONResponse(w, models.Response{Success: false, Message: err.Error()}, http.StatusInternalServerError)
			return
		}
		JSONResponse(w, models.Response{Success: true, Message: "File deteled"}, http.StatusOK)
		return
	}
}

package prettyFileServe

import (
	"fmt"
	"net/http"
	"path/filepath"
	"os"
	"html/template"
	"log"
	"io/ioutil"
)

type StorageHandlerDefinition struct {
	UrlBase      string
	InternalPath string
}

func New(urlBase string, internalPath string) *StorageHandlerDefinition {
	var storage = new(StorageHandlerDefinition)
	storage.InternalPath = internalPath
	storage.UrlBase = urlBase
	return storage
}

func (handler *StorageHandlerDefinition) CreateHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		realPath := resolveRealPath(handler, request.URL.Path)
		fileInfo, err := os.Stat(realPath)
		if (err != nil) {
			if (os.IsNotExist(err)) {
				//TODO REDIRECT 404
				fmt.Fprintf(response, "File %s not found", realPath)
			}
		}else {
			if (fileInfo.IsDir()) {
				if (hasIndexFile(realPath)) {
					serveFile(response, request, filepath.Join(realPath, "index.html"))
				}else {
					serveDirectoryView(response, request, realPath)
				}
			}else {
				serveFile(response, request, realPath)
			}
		}
	}
}

func serveFile(response http.ResponseWriter, request *http.Request, filePath string) {
	http.ServeFile(response, request, filePath)
}

func serveDirectoryView(response http.ResponseWriter, request *http.Request, directoryPath string) {
	listTemplate := template.Must(template.New("list-template").Parse(getListTemplate()))
	children, _ := ioutil.ReadDir(directoryPath)
	err := listTemplate.Execute(response, children)
	if (err != nil) {
		log.Fatal(err)
	}
}

func hasIndexFile(folderPath string) bool {
	fileInfo , err := os.Stat(filepath.Join(folderPath, "index.html"))
	if (err != nil) {
		//isNOtExistOrAnother problem lets assume that we won't be able to use it
		return false;
	}else {
		return fileInfo.IsDir() == false
	}
}

func resolveRealPath(handler *StorageHandlerDefinition, requestedPath string) string {
	base := filepath.Clean(handler.UrlBase)
	value := requestedPath[len(base):len(requestedPath)]
	return filepath.Join(handler.InternalPath, value)
}

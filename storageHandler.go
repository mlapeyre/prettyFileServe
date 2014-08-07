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
	return func(response http.ResponseWriter, request *http.Request)  {
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
					sendFile(response, request, filepath.Join(realPath, "index.html"))
				}else {
					sendDirectoryView(response,request,realPath)
				}
			}else {
				sendFile(response, request, realPath)
			}
		}
	}
}

func sendFile(response http.ResponseWriter, request *http.Request, filePath string){
	http.ServeFile(response, request, filePath)
}

func sendDirectoryView(response http.ResponseWriter, request *http.Request, directoryPath string) {
	//TODO CHANGE ME ASAP
	listTemplate, err := template.ParseFiles("/home/martial/projects/go-path/go-server/resources/templates/listFolder/template.html")
	if(err!=nil){
		log.Fatal(err)
	}else{
		children, _ := ioutil.ReadDir(directoryPath)
		err2 := listTemplate.Execute(response, children)
		if(err2 != nil){
			log.Fatal(err2)
		}
	}

	//for _, child := range children {
	//	fmt.Fprintf(w, "Name :%s", child.Name(),request)
	//}
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

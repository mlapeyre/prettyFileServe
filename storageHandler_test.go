package prettyFileServe

import (
	"testing"
	"io/ioutil"
	"os"
	"path/filepath"
	"net/http/httptest"
	"net/http"
)

var testData = []struct {
	basePattern  string
	urlRequest   string
	internalPath string
	resolvedPath string
}{
	{"/browse/", "/browse", "/home/martial/", "/home/martial"},
	{"/browse/", "/browse/", "/home/martial/", "/home/martial"},
	{"/browse/", "/browse/mama", "/home/martial/", "/home/martial/mama"},
	{"/browse/", "/browse/mama", "/home/martial", "/home/martial/mama"},
	{"/browse/", "/browse/mama/", "/home/martial", "/home/martial/mama"},
	{"/browse/film", "/browse/film/mama/", "/home/martial", "/home/martial/mama"},
}

//========================================================================================
//  Tests
//========================================================================================
func TestResolveRealPath(t *testing.T) {
	for _, element := range testData {
		storageHandler := New(element.basePattern, element.internalPath)
		value := resolveRealPath(storageHandler, element.urlRequest)
		assertEquals(value, element.resolvedPath, t)
	}
}

func TestShouldHasIndexBeTrueWhenTheFolderContainsAnIndexFile(t *testing.T) {
	tmpPath, _ := ioutil.TempDir(os.TempDir(), "go-test")
	defer os.RemoveAll(tmpPath)
	os.Create(filepath.Join(tmpPath, "index.html"))
	assertTrue(hasIndexFile(tmpPath), t)
}

func TestShouldHasIndexBeFalseWhenTheFolderDoesntContainAnIndexFile(t *testing.T) {
	tmpPath, _ := ioutil.TempDir(os.TempDir(), "go-test")
	defer os.RemoveAll(tmpPath)
	assertFalse(hasIndexFile(tmpPath), t)
}

func TestIndexShouldBeReturnedIfTheParanetFolderIsRequested(t *testing.T) {
	tmpPath, server := createAndStartServer(t)
	defer os.RemoveAll(tmpPath)
	defer server.Close()
	createFileWithContent(filepath.Join(tmpPath,"index.html"),"content",t)
	pageContent := doGetAsString(server.URL,t)
	assertEquals("content",pageContent,t)
}
func TestCorrectFileShouldBeReturnedIfAvailable(t *testing.T) {
	tmpPath, server := createAndStartServer(t)
	defer os.RemoveAll(tmpPath)
	defer server.Close()
	createFileWithContent(filepath.Join(tmpPath,"random.txt"),"my txt content",t)
	pageContent := doGetAsString(server.URL + "/random.txt",t)
	assertEquals("my txt content",pageContent,t)
}

//========================================================================================
//  Utils
//========================================================================================

func createFileWithContent(fullPath string, content string,t *testing.T){
	_, err := os.Create(fullPath)
	failIfNotNil(err,t)
	contentAsByte := []byte(content)
	err2 := ioutil.WriteFile(fullPath, contentAsByte, 0777)
	failIfNotNil(err2,t)
}

func createAndStartServer(t *testing.T) (string,* httptest.Server) {
	tmpPath, err := ioutil.TempDir(os.TempDir(), "go-test")
	failIfNotNil(err,t)
	handlerDefinition := New("/", tmpPath)
	server := httptest.NewServer(http.HandlerFunc(handlerDefinition.CreateHandler()))
	return tmpPath,server
}

func doGetAsString(url string, t *testing.T) (string) {
	res, err := http.Get(url)
	failIfNotNil(err,t)
	content, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	failIfNotNil(err,t)
	return string(content[:])
}

func failIfNotNil(err error,t *testing.T){
	if err != nil {
		t.Fatal(err)
		t.Fail()
	}
}

func assertTrue(predicate bool, t *testing.T) {
	if (predicate != true) {
		t.Fail()
	}
}
func assertFalse(predicate bool, t *testing.T) {
	assertTrue(!predicate, t)
}

func assertEquals(first string, second string, t *testing.T) {
	if (first == second) {
		return;
	}else {
		t.Errorf("%s != %s", first, second)
		t.Fail()
	}
}

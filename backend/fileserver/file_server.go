package fileserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
)

const FRONTEND_DIR = "frontend"
const DEFAULT_OBJECT = "index.html"

func GetRootDir() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		log.Println(err)
		return "", err
	}
	rootDir := path.Dir(execPath)
	return rootDir, nil
}

func ServeFrontend(w http.ResponseWriter, r *http.Request) {
	rootDir, err := GetRootDir()

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to determine root dir: %v", err), http.StatusInternalServerError)
	}
	basePath := path.Join(rootDir, FRONTEND_DIR)

	requestedPath := path.Join(basePath, path.Clean(r.URL.Path))

	_, err = os.Stat(requestedPath)
	if err == nil {
		log.Printf("GET %s\n", requestedPath)
		http.ServeFile(w, r, requestedPath)
		return
	} else if !os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("Something is wrong, file does and doesn't exist: %v", err))
	}
	log.Printf("GET %s, rewriting to %s\n", requestedPath, path.Join(basePath, DEFAULT_OBJECT))
	http.ServeFile(w, r, path.Join(basePath, DEFAULT_OBJECT))
	return
}

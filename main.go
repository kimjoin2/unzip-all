package main

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func main() {
	if len(os.Args) < 2 {
		println("input zip files path")
		os.Exit(1)
	}

	path := os.Args[1]
	pathSep := string(os.PathSeparator)
	if path[len(path)-1:] != pathSep {
		path = path + pathSep
	}
	fmt.Printf("Searching from %s\n", path)

	files, err := ioutil.ReadDir(path)
	if err != nil{
		panic(err)
		os.Exit(1)
	}

	var fileNames []string
	zipExt := ".zip"
	for _, file := range files{
		name := file.Name()
		if file.IsDir(){
			continue
		} else if !strings.HasSuffix(name, zipExt){
			continue
		}
		fileNames = append(fileNames, name)
	}

	dst := "output"
	fileCount := len(fileNames)

	var wg sync.WaitGroup
	wg.Add(fileCount)

	for index, fileName := range fileNames {
		fmt.Printf("%d / %d unzip.. %s\n", index+1, fileCount, fileName)
		fileName := fileName
		go func(){
			unzipFile(path, fileName, dst + pathSep + fileName[:len(fileName)-len(zipExt)])
			wg.Done()
		}()
	}
	wg.Wait()
	println("done.")
}

func unzipFile(zipBasePath string, fileName string, destination string) {
	archive, err := zip.OpenReader(zipBasePath + fileName)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(destination, f.Name)
		//fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(destination)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}
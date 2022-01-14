package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type PathInfo struct {
	Tree Tree `json:"tree"`
	Meta Meta `json:"meta"`
}

type Meta struct {
	Files int64 `json:"files"`
	Dirs  int64 `json:"dirs"`
}

type Tree struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Contents []Tree `json:"contents"`
}

var (
	depth   int
	typeOut string
)

const (
	humanOut = "h"
	jsonOut  = "j"
	dir      = "dir"
	file     = "file"
)

func init() {
	flag.IntVar(&depth, "depth", -1, "path reading depth, depth >= 0")
	flag.StringVar(&typeOut, "type", humanOut, "type of output:h(human readable), j(json)")
}

func main() {
	flag.CommandLine.Parse(os.Args[2:])
	meta := new(Meta)
	dirname := os.Args[1]
	trees, err := scanDir(dirname, meta, true, depth)
	if err != nil {
		fmt.Println("error:", err)
	}
	pathInfo := &PathInfo{*trees, *meta}
	dataOutput(pathInfo)
}

//Function for data output
//Функция для вывода данных
func dataOutput(pathInfo *PathInfo) {
	if typeOut == jsonOut {
		pathOutput, err := json.Marshal(pathInfo)
		if err != nil {
			return
		}
		fmt.Printf("%s\n", pathOutput)
	} else if typeOut == humanOut {
		lenName := 0
		fmt.Printf("%s:\n", pathInfo.Tree.Name)
		humanReadOut(&pathInfo.Tree, lenName)
		fmt.Printf("files: %d\ndirs: %d\n", pathInfo.Meta.Files, pathInfo.Meta.Dirs)
	} else {
		fmt.Println("Unknown type of output.")
	}
}

//Function for human readable text
//Функция для человекочитаемого формата вывода
func humanReadOut(tree *Tree, lenName int) (*Tree, *int) {
	lenName += 5
	for _, treeElem := range tree.Contents {
		fmt.Printf("%s|--%s\n", strings.Repeat(" ", lenName), treeElem.Name)
		humanReadOut(&treeElem, lenName)
	}
	return tree, &lenName
}

//Function for filling of PathInfo struct
//Функция для полного заполнения структуры PathInfo
func scanDir(dirname string, meta *Meta, needTypes bool, depth int) (*Tree, error) {

	tree := new(Tree)
	fi, err := os.Stat(dirname)
	if err != nil {
		return nil, err
	}
	//Accounting dirs and files
	//Подсчёт кол-ва папок и файлов, запоминание их имен и типа
	if !fi.IsDir() {
		tree.Name = fi.Name()
		meta.Files++
		if needTypes {
			tree.Type = dir
		}
		return tree, nil
	} else {
		tree.Name = fi.Name()
		meta.Dirs++
		if needTypes {
			tree.Type = file
		}
	}
	//Count of depth scan
	//Учёт глубины прохождения
	if depth != 0 {
		depth--
	} else {
		return tree, nil
	}

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	//Read content of directory
	//Чтение содержимого папки
	for _, file := range files {
		dirname := dirname + "/" + file.Name()
		childTree, err := scanDir(dirname, meta, needTypes, depth)
		if err != nil {
			return nil, err
		}
		tree.Contents = append(tree.Contents, *childTree)
	}
	return tree, nil
}

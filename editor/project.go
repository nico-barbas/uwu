package editor

import (
	"io/fs"
	"os"
)

var exceptionList = []string{
	"./.git",
	"./.vscode",
}

type (
	project struct {
		rootInfo os.FileInfo
		root     *folder
		current  *folder
		previous *folder
	}

	projectNode interface {
		name() string
	}

	folder struct {
		entry fs.DirEntry
		nodes map[string]projectNode
	}

	file struct {
		entry fs.DirEntry
	}
)

func openProject(path string) project {
	var err error
	proj := project{}
	proj.rootInfo, err = os.Stat(path)
	if err != nil {
		panic(err)
	}
	proj.root = &folder{
		nodes: make(map[string]projectNode),
	}
	proj.current = proj.root

	proj.readDir(path)

	return proj
}

func (p *project) readDir(path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		switch f.IsDir() {
		case true:
			dirName := f.Name()
			dirPath := path
			if path[len(path)-1] != '/' {
				dirPath += "/"
			}
			dirPath += dirName
			if !isPathException(dirPath) {
				p.current.addSubFolder(f)
				p.previous = p.current
				p.current = p.current.nodes[dirName].(*folder)
				p.readDir(dirPath)
				p.current = p.previous
			}
		case false:
			p.current.addFile(f)
		}
	}
}

func (f *folder) addSubFolder(entry fs.DirEntry) {
	f.nodes[entry.Name()] = &folder{
		entry: entry,
		nodes: make(map[string]projectNode),
	}
}

func (f *folder) addFile(entry fs.DirEntry) {
	// No need to check for existing one since the
	// OS garantees that filenames are unique
	f.nodes[entry.Name()] = file{
		entry: entry,
	}
}

func (f folder) name() string {
	return f.entry.Name()
}

func (f file) name() string {
	return f.entry.Name()
}

func isPathException(path string) bool {
	for _, e := range exceptionList {
		if e == path {
			return true
		}
	}
	return false
}

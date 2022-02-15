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
		path() string
		getEntry() fs.DirEntry
	}

	folder struct {
		entry    fs.DirEntry
		nodes    map[string]projectNode
		nodePath string
	}

	file struct {
		entry    fs.DirEntry
		nodePath string
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
		nodes:    make(map[string]projectNode),
		nodePath: path,
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

func (p *project) findNode(name string) projectNode {
	return p.root.findChild(name)
}

func (f *folder) addSubFolder(entry fs.DirEntry) {
	// No need to check for existing one since the
	// OS garantees that filenames are unique
	f.nodes[entry.Name()] = &folder{
		entry:    entry,
		nodes:    make(map[string]projectNode),
		nodePath: f.nodePath + "/" + entry.Name(),
	}
}

func (f *folder) addFile(entry fs.DirEntry) {
	// No need to check for existing one since the
	// OS garantees that filenames are unique
	f.nodes[entry.Name()] = file{
		entry:    entry,
		nodePath: f.nodePath + "/" + entry.Name(),
	}
}

func (f *folder) findChild(name string) (child projectNode) {
search:
	for k, v := range f.nodes {
		if k == name {
			child = v
			break search
		}
		switch f := v.(type) {
		case *folder:
			child = f.findChild(name)
			if child != nil {
				break search
			}
		}
	}
	return
}

func (f folder) name() string {
	return f.entry.Name()
}

func (f folder) path() string {
	return f.nodePath
}

func (f folder) getEntry() fs.DirEntry {
	return f.entry
}

func (f file) name() string {
	return f.entry.Name()
}

func (f file) path() string {
	return f.nodePath
}

func (f file) getEntry() fs.DirEntry {
	return f.entry
}

func isPathException(path string) bool {
	for _, e := range exceptionList {
		if e == path {
			return true
		}
	}
	return false
}

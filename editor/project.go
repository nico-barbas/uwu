package editor

import (
	"io/fs"
	"os"
)

type (
	project struct {
		rootInfo os.FileInfo
		root     folder
		current  *folder
		previous *folder
	}

	projectNode interface{}

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
	proj.root = folder{
		nodes: make(map[string]projectNode),
	}
	proj.current = &proj.root

	files, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		switch f.IsDir() {
		case true:
		case false:
			proj.current.addFile(f)
		}
	}

	return proj
}

func (f folder) readFolder(entry fs.DirEntry) {

}

func (f *folder) addFile(entry fs.DirEntry) {
	// No need to check for existing one since the
	// OS garantees that filenames are unique
	f.nodes[entry.Name()] = file{
		entry: entry,
	}
}

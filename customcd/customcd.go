package customcd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	matchedFolders []byte
	wg             sync.WaitGroup
)

var (
	path = os.Getenv("HOME")
)

// SearchPath searches for the real paths of folders inside the root.
// Fetches the paths which are matched to the prefix.
func SearchPath(prefix string) {
	runtime.GOMAXPROCS(1)
	running := make(chan struct{}, 10)
	f, err := os.Open(path)
	if err != nil {
		exit(running, nil, err)
	}
	fSlice, err := f.Readdir(-1)
	if err != nil {
		exit(running, nil, err)
	}
	folders := make(chan []byte, len(fSlice))
	for _, f := range fSlice {
		if !strings.HasPrefix(f.Name(), ".") {
			wg.Add(1)
			running <- struct{}{}
			go func(f os.FileInfo) {
				defer wg.Done()
				defer func() { <-running }()
				foldersList, err := findFolders([]byte{}, fmt.Sprintf("%s/%s", path, f.Name()), prefix)
				if err != nil {
					exit(running, folders, err)
				}
				folders <- foldersList
				return
			}(f)
		}
	}
	wg.Wait()
	for i := 0; i < len(fSlice); i++ {
		select {
		case ff := <-folders:
			matchedFolders = append(matchedFolders, ff...)

		default:
		}
	}
	close(folders)
	if len(matchedFolders) == 0 {
		return
	}
	fmt.Println(string(matchedFolders[1:]))
}

func findFolders(folders []byte, path, prefix string) ([]byte, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return folders, err
	}
	if !stat.IsDir() || (stat.Mode()&(1<<2) == 0) {
		return folders, nil
	}
	if strings.HasPrefix(filepath.Base(path), prefix) {
		b := []byte("," + path)
		folders = append(folders, b...)
	}
	fSlice, err := ioutil.ReadDir(path)
	if err != nil {
		return folders, err
	}

	for _, ff := range fSlice {
		if ff.IsDir() && !strings.HasPrefix(ff.Name(), ".") {
			folders, err = findFolders(folders, fmt.Sprintf("%s/%s", path, ff.Name()), prefix)
			if err != nil {
				return folders, err
			}
		}
	}
	return folders, nil
}

func exit(running chan struct{}, folders chan []byte, err error) {
	if running != nil {
		close(running)
	}
	if folders != nil {
		close(folders)
	}
	fmt.Printf("Error: %v", err)
	os.Exit(0)
}

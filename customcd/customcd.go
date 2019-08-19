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

// SearchPath searches for the real paths of folders inside the root.
// Fetches the paths which are matched to the prefix.
func SearchPath(prefix string) {
	var (
		wg             sync.WaitGroup
		matchedFolders []byte
	)

	if prefix == "" {
		exit(nil, nil, &notify{
			mesgType: WARN,
			message:  "Invalid prefix",
		})
	}
	runtime.GOMAXPROCS(1)

	running := make(chan struct{}, 10)
	path := os.Getenv("HOME")
	f, err := os.Open(path)
	if err != nil {
		exit(running, nil, &notify{
			mesgType: ERROR,
			message:  err.Error(),
		})
	}
	defer f.Close()

	fSlice, err := f.Readdir(-1)
	if err != nil {
		exit(running, nil, &notify{
			mesgType: ERROR,
			message:  err.Error(),
		})
	}
	folders := make(chan []byte, len(fSlice))
	for _, f := range fSlice {
		if !strings.HasPrefix(f.Name(), ".") && f.IsDir() {
			wg.Add(1)
			running <- struct{}{}
			go func(f os.FileInfo) {
				defer wg.Done()
				defer func() { <-running }()
				foldersList, err := findFolders([]byte{}, fmt.Sprintf("%s/%s", path, f.Name()), prefix)
				if err != nil {
					exit(running, folders, &notify{
						mesgType: ERROR,
						message:  err.Error(),
					})
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
	if len(matchedFolders) == 0 {
		exit(running, folders, &notify{
			mesgType: WARN,
			message:  "No folders found. Please verify the prefix",
		})
	}
	fmt.Println(string(matchedFolders[1:]))
	exit(running, folders, nil)
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

func exit(running chan struct{}, folders chan []byte, notif *notify) {
	if running != nil {
		close(running)
	}
	if folders != nil {
		close(folders)
	}
	if notif == nil {
		return
	}

	if notif.mesgType == WARN {
		fmt.Printf("%sWarning: %s%s", YELLOW, WHITE, notif.message)
		os.Exit(0)
	}

	if notif.mesgType == ERROR {
		fmt.Printf("%sError: %s%s", RED, WHITE, notif.message)
		os.Exit(0)
	}
}

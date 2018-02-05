package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	displayAll = 0
)

// Args cli arguments definition
type Args struct {
	pictDir string
	display int
}

func main() {
	os.Exit(run())
}

func run() int {
	args, err := parseArgs()
	if err != nil {
		log.Print(err)
		fmt.Printf("usage: binggo --pictdir /path/to/downloads\n")
		fmt.Printf("usage: binggo --pictdir /path/to/downloads --display 1\n")
		fmt.Printf("usage: --display option is available only on macOS\n")
		return 1
	}
	urls, err := getPictureUrls()
	if err != nil {
		log.Print(err)
		return 2
	}
	pictDirName := args.pictDir
	for _, url := range urls {
		err := downloadPicture(url, pictDirName)
		if err != nil {
			log.Print(err)
			return 3
		}
	}
	files, err := getWallpaperFile(pictDirName)
	if err != nil {
		log.Print(err)
		return 4
	}
	sort.Sort(files)
	err = changeWallpaper(args.display, pictDirName, files[0])
	if err != nil {
		log.Print(err)
		return 5
	}
	return 0
}

func parseArgs() (*Args, error) {
	var pictDir string
	flag.StringVar(&pictDir, "pictdir", "", "directory path for download")
	var display int
	flag.IntVar(&display, "display", displayAll, "target display number.")
	flag.Parse()

	if pictDir == "" {
		return nil, errors.New("pictdir is empty")
	}
	if pictDir[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			return nil, err
		}
		pictDir = strings.Replace(pictDir, "~/", usr.HomeDir+"/", 1)
	}
	pictDir, err := validatePictDir(pictDir)
	if err != nil {
		return nil, err
	}

	args := &Args{
		pictDir: pictDir,
		display: display,
	}
	return args, err
}

// validate pict directory path
func validatePictDir(path string) (string, error) {
	pictDirPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	f, err := os.Stat(pictDirPath)
	if err != nil {
		return "", fmt.Errorf("[%s] is not exists.\n%v", pictDirPath, err)
	}
	if !f.IsDir() {
		return "", fmt.Errorf("[%s] is not a directory", pictDirPath)
	}
	return pictDirPath, nil
}

// get all urls of picture from bing top page
func getPictureUrls() ([]string, error) {
	resp, err := http.Get("http://www.bing.com")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s := string(body)
	regex1 := regexp.MustCompile(`url:'\\([^']*)'`)
	results := regex1.FindAllStringSubmatch(s, -1)

	var urls []string
	for _, r := range results {
		if len(r) == 2 {
			p := strings.Replace(r[1], "\\", "", -1)
			urls = append(urls, fmt.Sprintf("http://bing.com%s", p))
		}
	}
	return urls, nil
}

// download bing picture
func downloadPicture(url string, dirName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	splited := strings.Split(url, "/")
	filename := splited[len(splited)-1]
	out, err := os.Create(dirName + string(os.PathSeparator) + filename)
	if err != nil {
		return err
	}
	defer out.Close()
	out.Write(body)
	return nil
}

// PictFiles are wallpaper pictures
type PictFiles []os.FileInfo

// Len for sort.Interface
func (files PictFiles) Len() int {
	return len(files)
}

// Swap for sort.Interface
func (files PictFiles) Swap(i, j int) {
	files[i], files[j] = files[j], files[i]
}

// Less for sort.Interface
func (files PictFiles) Less(i, j int) bool {
	return files[i].ModTime().UnixNano() > files[j].ModTime().UnixNano()
}

// getWallpaperFiles collects wallpaper picture file
func getWallpaperFile(pictDir string) (PictFiles, error) {
	var files PictFiles
	files, err := ioutil.ReadDir(pictDir)
	if err != nil {
		return nil, err
	}
	return files, nil
}

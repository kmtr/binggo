package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const PICT_DIR = "PICT_DIR"

func main() {
	params := parseArgs()
	if params == nil {
		return
	}
	pictDirName, err := validatePictDir(params[PICT_DIR])
	if err != nil {
		log.Fatal(err)
	}
	urls, err := getPictureUrls()
	if err != nil {
		log.Fatal(err)
	}
	for _, url := range urls {
		downloadPicture(url, pictDirName)
	}
	files := getWallpaperFile(pictDirName)
	sort.Sort(files)
	changeWallpaper(pictDirName, files[0])
}

func parseArgs() map[string]string {
	var pictDir string
	flag.StringVar(&pictDir, "pictdir", "", "directory path for download")

	flag.Parse()

	if pictDir == "" {
		fmt.Printf("usage: binggo --pictdir=/path/to/downloads\n")
		return nil
	}
	if pictDir[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		pictDir = strings.Replace(pictDir, "~/", usr.HomeDir+"/", 1)
	}
	params := map[string]string{
		PICT_DIR: pictDir,
	}
	return params
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
func downloadPicture(url string, dirName string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	splited := strings.Split(url, "/")
	filename := splited[len(splited)-1]
	out, err := os.Create(dirName + string(os.PathSeparator) + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	out.Write(body)
}

type PictFiles []os.FileInfo

func (files PictFiles) Len() int {
	return len(files)
}

func (files PictFiles) Swap(i, j int) {
	files[i], files[j] = files[j], files[i]
}

func (f PictFiles) Less(i, j int) bool {
	return f[i].ModTime().UnixNano() > f[j].ModTime().UnixNano()
}

// get wallpaper files
func getWallpaperFile(pictDir string) PictFiles {
	var files PictFiles
	files, err := ioutil.ReadDir(pictDir)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

// change wallpaper
func changeWallpaper(pictDir string, f os.FileInfo) {
	script := `
	tell application "Finder"
	set desktop picture to POSIX file "%s/%s"
	end tell
	`
	cmd := exec.Command("osascript", "-e",
		fmt.Sprintf(script, pictDir, f.Name()))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

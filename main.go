package main

import (
	"flag"
	"fmt"
	"github.com/nleeper/goment"
	"io/ioutil"
	"mime"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
)

var (
	mimeGroup       = flag.String("m", "", "required")
	readDirectory   = flag.String("r", "", "required")
	outputDirectory = flag.String("o", "", "required")
	dryRun          = flag.Bool("d", false, "dry-run")
)

func main() {
	flag.Parse()

	if *readDirectory == "" || *mimeGroup == "" || *outputDirectory == "" {
		fmt.Fprintln(os.Stderr, "Required flag is not specified.")
		os.Exit(1)
	}

	if *dryRun {
		fmt.Fprintln(os.Stdout, "------ DRY RUN MODE ------")
	}

	files, err := ioutil.ReadDir(*readDirectory)
	if err != nil {
		panic("Could not read directory.")
	}

	mimePattern, _ := regexp.Compile("^" + *mimeGroup + "/")
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		mimeType := mime.TypeByExtension(filepath.Ext(file.Name()))
		if !mimePattern.MatchString(mimeType) {
			continue
		}

		date, _ := goment.New(file.ModTime())
		basedir := path.Join(path.Clean(*outputDirectory), date.Format("YYYY-MM"))
		if _, err := os.Stat(basedir); os.IsNotExist(err) && !*dryRun {
			os.Mkdir(basedir, 777)
		}

		containerDir := path.Join(basedir, date.Format("MM-DD"))
		if _, err := os.Stat(containerDir); os.IsNotExist(err) && !*dryRun{
			os.Mkdir(containerDir, 777)
		}

		from := path.Join(*readDirectory, file.Name())
		to := path.Join(containerDir, file.Name())

		for i := 1; true; i++ {
			if _, err := os.Stat(to); os.IsNotExist(err) {
				break
			}

			to = path.Join(containerDir, strconv.Itoa(i) + "_" + file.Name())
		}

		if *dryRun {
			fmt.Fprintln(os.Stdout, "Move: " +from+ " -> " +to)
		} else {
			os.Rename(from, to)
		}
	}

	os.Exit(0)
}

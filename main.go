package main

import (
	"flag"
	"fmt"
	"github.com/plally/workshopdl/internal/steam"
	"github.com/ulikunitz/xz/lzma"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	workshopPattern = regexp.MustCompile(`https?://steamcommunity\.com/workshop/filedetails/\?id=([0-9]+)`)
	idPattern       = regexp.MustCompile("([0-9]+)")
)

var (
	flagExtract = flag.Bool("g", true, "extract downloaded addon with gmad")
)

func main() {
	flag.Parse()
	flag.Usage = func() {
		os.Stderr.Write([]byte("workshopdl [addonid]\n"))
		os.Stderr.Write([]byte("Downloads an addon from the workship given an addonid or url\n\n"))
		flag.PrintDefaults()
	}

	if len(flag.Args()) < 1 {
		flag.Usage()
	}
	for _, input := range flag.Args() {
		foundId := findWorkshopId(input)
		err := DownloadAddons(foundId)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
	}

}

func findWorkshopId(input string) string {
	matches := workshopPattern.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}

	matches = idPattern.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}
func DownloadAddons(workshopId ...string) error {
	fmt.Printf("Fetching file details for %v\n", workshopId)
	fileDetails, err := steam.GetPublishedFileDetails(workshopId)
	if err != nil {
		return err
	}

	for _, file := range fileDetails.PublishedFileDetails {
		fmt.Printf("Downloading %v\n", file.Title)

		downloadedFilename := strings.ReplaceAll(file.Filename, "/", "_")
		err := downloadFile(file.FileURL, downloadedFilename)
		if err != nil {
			return err
		}

		if *flagExtract {
			fmt.Printf("File downloaded \"%v\". Extracting with gmad\n", downloadedFilename)
			err := gmadExtract(downloadedFilename)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("File downloaded \"%v\" \n", downloadedFilename)
		}
	}

	return nil
}

func downloadFile(fileUrl, out string) error {
	resp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	decodedBody, err := lzma.NewReader(resp.Body)
	if err != nil {
		return err
	}

	f, err := os.Create(out)
	_, err = io.Copy(f, decodedBody)
	if err != nil {
		return err
	}
	return nil
}

func findGmad() (string, error) {
	gmadPath := os.Getenv("GMAD_PATH")
	if len(gmadPath) > 0 {
		return gmadPath, nil
	}
	return exec.LookPath("gmad")
}

func gmadExtract(filename string) error {
	gmad, err := findGmad()
	if err != nil {
		fmt.Println("You should define the environment variable GMAD_PATH or add gmad to your PATH")
		return fmt.Errorf("gmad executable not found: %w", err)
	}
	cmd := exec.Command(gmad, "extract", "-file", filename)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)


type AnimeEntry struct {
	Name 		string	`json:"name"`
	Size 		string	`json:"size"`
	Downloads 	int		`json:"downloads"`
	TorrentLink string	`json:"torrent_link"`
	Date 		string	`json:"date"`
}


type ScrapperInstance struct {
	C *colly.Collector
	URL string
}

func (i *ScrapperInstance) Init(domain string, URL string) {
	// create a tmp file to temporarily store the data
	// don't judge me, i have no idea for now
	os.Remove("/tmp/nyanime"); // gettting rid of old data
	file , err := os.Create("/tmp/nyanime");
	if err != nil {
		log.Println("Error creating file: ", err);
		os.Exit(1);
	}

	i.C = colly.NewCollector(
		colly.AllowedDomains(domain),
	)
	i.URL = URL;


	i.C.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})


	i.C.OnHTML("table.torrent-list tbody tr", func(e *colly.HTMLElement) {
		var element AnimeEntry;
		var err error;
		element.Name = e.ChildAttr("td:nth-child(2) a:nth-child(2)", "title");
		if (element.Name == "" ) {
			element.Name = e.ChildAttr("td:nth-child(2) a:nth-child(1)", "title");
		}
		element.TorrentLink = URL + e.ChildAttr("td:nth-child(3) a:nth-child(1)", "href");
		element.Size = e.ChildText("td:nth-child(4)");
		element.Date = e.ChildText("td:nth-child(5)");
		element.Downloads, err = strconv.Atoi(e.ChildText("td:nth-child(8)"));
		if err != nil {
			log.Println("Error parsing downloads: ", err)
			element.Downloads = 0;
		}

		// write to /tmp/nyanime instead of stdout
		res := fmt.Sprintf("%s | %s | %s | %s", element.TorrentLink, element.Name, element.Size, element.Date);
		file.WriteString(res + "\n");
	})
}


func (i *ScrapperInstance) Search(anime_name string) {
	i.C.Visit(i.URL + "?f=0&c=0_0&s=seeders&o=desc&q=" + anime_name);
	cmd := exec.Command("sh", "-c", "cat /tmp/nyanime | fzf --reverse --with-nth=2.. --cycle")

    output, err := cmd.Output()
    if err != nil {
        fmt.Println("aborting")
		os.Exit(0)
    }

    torrent_link := strings.TrimSpace(strings.Split(string(output), " |")[0])
	outfile, err := os.CreateTemp("/tmp", "nyanime.*.torrent");
	if err != nil {
		log.Println("Error creating .torrent file, aborting");
		os.Exit(1);
	}
	defer outfile.Close();
	cmd = exec.Command("wget", torrent_link, "-O", outfile.Name());
	err = cmd.Run();
	if err != nil {
		log.Println("Error downloading .torrent file: ", err);
		os.Exit(1);
	}

	cmd = exec.Command("xdg-open", outfile.Name());
	err = cmd.Run();
	if err != nil {
		log.Println("Error opening .torrent file: ", err);
		os.Exit(1);
	}
}
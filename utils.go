package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)



func usage() {
	fmt.Println("Usage: nyanime [...options]");
	fmt.Println("Options:");
	fmt.Println("  -n | --name <anime name> Search for anime by its name");
	fmt.Println("  -h | --help              Show this help message");
	fmt.Println("  -c | --clean             Show this help message");
	os.Exit(0);
}



func clean_tmp() {
	files, err := os.ReadDir("/tmp");
	if err != nil {
		log.Println("Error reading /tmp: ", err);
		os.Exit(1);
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "nyanime") {
			err := os.Remove("/tmp/" + file.Name());
			if err != nil {
				log.Println("Error removing file: ", err);
			}
		}
	}
	os.Exit(0);
}


func parse_args() string {
	anime_name := "";
	argc := len(os.Args);
	if (argc < 2) {
		usage();
	}
	for i, arg := range os.Args {
		switch arg {
		case "--name", "-n":
			if (i+1 < argc) {
				anime_name = os.Args[i+1];
				if (anime_name == "") {
					usage();
				}
			} else {
				usage();
			}

		case "--help", "-h":
			usage();

		case "--clean", "-c":
			clean_tmp();
		}
	}

	// url encoding the query
	anime_name = url.QueryEscape(anime_name);
	return anime_name;
}
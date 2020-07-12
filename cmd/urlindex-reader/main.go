package main

import (
	"bufio"
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xavivars/urlindex/domain"
	"os"
	"strings"
)

var file string
func init() {
	log.SetFormatter(&log.JSONFormatter{})

	flag.StringVar(&file, "file", "", "File with the compiled FST")
	flag.StringVar(&file, "f", "", "File with the compiled FST")
}

func main() {

	flag.Parse()

	i, err := domain.NewUrlIndex("vistaprint", "en-ie", file, 1, false)

	if err != nil {
		log.Fatal("Could not initialize index")
	}

	reader := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter text: ")
	for reader.Scan() {

		text := reader.Text()
		fmt.Println(text)

		if strings.Compare(text, "bye") == 0 {
			break
		}

		fmt.Println(i.Exists(text))

		fmt.Print("Enter text: ")
	}
}

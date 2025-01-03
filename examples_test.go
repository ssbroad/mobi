package mobi_test

import (
	"math/rand"
	"os"
	"time"

	"github.com/ssbroad/mobi"
	"golang.org/x/text/language"
)

func ExampleBook() {
	// Create simple book with chapter
	ch := mobi.Chapter{
		Title:  "Chapter 1",
		Chunks: mobi.Chunks(`Lorem ipsum dolor sit amet, consetetur sadipscing elitr.`),
	}
	mb := mobi.Book{
		Title:       "De vita Caesarum librus",
		Authors:     []string{"Sueton"},
		CreatedDate: time.Now(),
		Language:    language.Italian,
		Chapters:    []mobi.Chapter{ch},
		UniqueID:    rand.Uint32(),
	}

	// Convert book to PalmDB database
	db := mb.Realize()

	// Write database to file
	f, _ := os.Create("test.azw3")
	defer f.Close()
	err := db.Write(f)
	if err != nil {
		panic(err)
	}
}

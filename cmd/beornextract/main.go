package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/jessevdk/go-flags"
	"github.com/squk/lotr/cmd/beornextract/types"
)

type Options struct {
	RawConversion bool `short:"r" long:"raw" description:"Enable to keep the original text from HallOfBeorn dump. Enable to prep for ALEP pipeline."`
}

var opts = Options{
	RawConversion: false,
}

func main() {
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		panic(err)
	}

	fmt.Println("LOTR CARD PARSE")
	f, err := bazel.Runfile(".")
	if err != nil {
		panic(err)
	}
	err = filepath.Walk(f,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path, info.Size())
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	// Open our jsonFile
	// jsonFile, err := os.Open("cmd/beornextract/data/Bot.Cards.json")
	jsonFile, err := os.Open("cmd/beornextract/data/Export.Cards.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
		return
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	cards := []types.HallOfBeornCard{}
	json.Unmarshal(byteValue, &cards)

	// Open a file for writing
	csvFile, err := os.Create("/Users/christian/Downloads/lotr-lcg-set-generator.csv")
	defer csvFile.Close()
	if err != nil {
		// Handle error
	}

	// Create a writer
	w := csv.NewWriter(csvFile)

	// Write some rows
	for _, card := range cards {
		if card.EncounterSet != "" {
			continue // skip non=player cards
		}

		w.Write(
			[]string{
				card.Octgnid,
				"", // hidden
				"", // hidden
				card.EncounterSet,
				strconv.Itoa(card.Position),
				strconv.Itoa(card.Quantity),
				card.Name,
				fmt.Sprintf(
					"%t",
					card.IsUnique,
				),
				card.TypeName,
				card.SphereName,
				card.Traits,
				findKeywords(card.Text),
				card.Cost,
				card.EngagementCost,
				strconv.Itoa(card.Threat),
				strconv.Itoa(card.Willpower),
				strconv.Itoa(card.Attack),
				strconv.Itoa(card.Defense),
				strconv.Itoa(card.Health),
				card.QuestPoints,
				strconv.Itoa(card.VictoryPoints),
				"", // Special Icon
				transformText(card.Name, card.Text),
				card.Flavor,
			},
		)
	}

	// Close the writer
	w.Flush()
	// spew.Dump(c)
}

func transformText(name, text string) string {
	if opts.RawConversion {
		return text
	}

	out := strings.ReplaceAll(text, name, "[name]") // insert name tag
	out = strip.StripTags(out)
	out = keywordPattern.ReplaceAllLiteralString(out, "")
	return strings.TrimSpace(out)
}

var keywordPattern = regexp.MustCompile(`((?:(?:[A-Z][a-z]+(\.|\s[0-9]+\.)\s*)+))`)

func findKeywords(text string) string {
	return strings.TrimSpace(keywordPattern.FindString(text))
}

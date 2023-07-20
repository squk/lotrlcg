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
	"github.com/squk/lotrlcg/cmd/beornextract/types"
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
	jsonFile, err := os.Open("cmd/beornextract/data/Bot.Cards.json")
	// jsonFile, err := os.Open("cmd/beornextract/data/Export.Cards.json")
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

	// header
	w.Write([]string{"Card GUID", "Updated", "Diff", "Card Number", "Quantity", "Encounter Set", "Name", "Unique", "Type", "Sphere", "Traits", "Keywords", "Cost", "Engagement Cost", "Threat", "Willpower", "Attack", "Defense", "Health", "Quest Points", "Victory Points", "Special Icon", "Text", "Shadow", "Flavour", "Printed Card Number", "Encounter Set Number", "Encounter Set Icon", "Flags", "Artist", "PanX", "PanY", "Scale", "Portrait Shadow", "Side B", "Unique", "Type", "Sphere", "Traits", "Keywords", "Cost", "Engagement Cost", "Threat", "Willpower", "Attack", "Defense", "Health", "Quest Points", "Victory Points", "Special Icon", "Text", "Shadow", "Flavour", "Printed Card Number", "Encounter Set Number", "Encounter Set Icon", "Flags", "Artist", "PanX", "PanY", "Scale", "Portrait Shadow", "Removed for Easy Mode", "Additional Encounter Sets", "Adventure", "Collection Icon", "Copyright", "Card Back", "Version"})

	// Write some rows
	for _, card := range cards {
		playerCard := true
		if card.EncounterSet != "" {
			playerCard = false
			continue // skip non=player cards for now
		}

		if card.Octgnid == "" {
			card.Octgnid = card.Name + card.SphereCode + card.TypeCode + strconv.Itoa(card.Position)
		}

		threat := strconv.Itoa(card.Threat)
		victoryPoints := strconv.Itoa(card.VictoryPoints)
		if playerCard {
			threat = "" // triggers redundant warning in AleP
			// AleP wants hero threat in cost
			card.Cost = strconv.Itoa(card.Threat)

			// triggers redundant warning in AleP
			if card.VictoryPoints == 0 {
				victoryPoints = ""
			}

		}

		w.Write(
			[]string{
				card.Octgnid,
				"", // hidden
				"", // hidden
				strconv.Itoa(card.Position),
				strconv.Itoa(card.Quantity),
				card.EncounterSet,
				card.Name,
				strconv.Itoa(b2i(card.IsUnique)),
				card.TypeName,
				card.SphereName,
				card.Traits,
				findKeywords(card.Text),
				card.Cost,
				card.EngagementCost,
				threat,
				strconv.Itoa(card.Willpower),
				strconv.Itoa(card.Attack),
				strconv.Itoa(card.Defense),
				strconv.Itoa(card.Health),
				card.QuestPoints,
				victoryPoints,
				"", // Special Icon
				transformText(card.Name, card.Text),
				card.Shadow,
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

var keywordPattern = regexp.MustCompile(`^((?:(?:[A-Z][a-z]+(\.|\s[0-9]+\.)\s*)+))`)

func findKeywords(text string) string {
	return strings.TrimSpace(keywordPattern.FindString(text))
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

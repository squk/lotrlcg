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

		willpower := getXVal(card.Willpower)
		attack := getXVal(card.Attack)
		defense := getXVal(card.Defense)
		health := getXVal(card.Health)
		is_unique := strconv.Itoa(b2i(card.IsUnique))

		threat := getXVal(card.Threat)
		victoryPoints := getXVal(card.VictoryPoints)
		cost := card.Cost
		questPoints := getXVal(card.QuestPoints)
		if playerCard {
			threat = "" // triggers redundant warning in AleP
			if card.TypeCode == "hero" {
				// AleP wants hero threat in cost
				cost = strconv.Itoa(card.Threat)
			}

			if card.TypeCode != "player-side-quest" {
				questPoints = ""
			}

			if card.TypeCode == "event" || card.TypeCode == "attachment" {
				willpower, attack, defense, health = "", "", "", ""
				is_unique = ""
			}

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
				is_unique,
				card.TypeName,
				card.SphereName,
				card.Traits,
				findKeywords(card.Text),
				cost,
				card.EngagementCost,
				threat,
				willpower, attack, defense, health,
				questPoints,
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

var keywordPattern = regexp.MustCompile(`^((?:(?:[A-Z][a-z]+(\.|\s[0-9]+\.)\s*)+))`)
var paragraphPattern = regexp.MustCompile(`(\r\n|\r|\n)+`)

func transformText(name, text string) string {

	if opts.RawConversion {
		return text
	}

	boldList := []string{"Travel"}
	out := strings.ReplaceAll(text, name, "[name]") // insert name tag
	out = strip.StripTags(out)
	out = keywordPattern.ReplaceAllString(out, "")
	for _, str := range boldList {
		out = strings.ReplaceAll(out, str, "[b]"+str+"[/b]")
	}

	out = regexp.MustCompile(`may trigger this (?:action|response)`).ReplaceAllString(out, "may trigger this effect")
	out = regexp.MustCompile(`(\bheal[^.]+?\b)on(\b)`).ReplaceAllString(out, "${1}from${2}")
	out = regexp.MustCompile(`(Traits?)`).ReplaceAllString(out, "{${1}}")

	//  make all newline groups exactly two newlines
	out = paragraphPattern.ReplaceAllLiteralString(out, "\r\n\r\n")
	return strings.TrimSpace(out)
}

func findKeywords(text string) string {
	return strings.TrimSpace(keywordPattern.FindString(text))
}

func getXVal(val int) string {
	if val == 254 {
		return "X"
	}

	return strconv.Itoa(val)
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

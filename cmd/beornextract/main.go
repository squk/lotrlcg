package main

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"github.com/satori/go.uuid"

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
			fmt.Println(path)
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	// Open our jsonFile
	// jsonFile, err := os.Open("cmd/beornextract/data/ringsdb.json")
	jsonFile, err := os.Open("cmd/beornextract/data/Export.Cards.json")
	defer jsonFile.Close()
	// jsonFile, err := os.Open("cmd/beornextract/data/Bot.Cards.json") // has incorrect data
	if err != nil {
		fmt.Println(err)
		return
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	cards := []types.HallOfBeornCard{}
	json.Unmarshal(byteValue, &cards)

	cyclesJson, err := os.Open("cmd/beornextract/data/cycles.json")
	defer cyclesJson.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	byteValue, _ = ioutil.ReadAll(cyclesJson)
	cycles := types.CycleMappings{}
	json.Unmarshal(byteValue, &cycles)

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
			// continue // skip non=player cards for now
		}

		if card.Octgnid == "" {
			card.Octgnid = deterministicUUID(card.Name + card.SphereCode + card.TypeCode + strconv.Itoa(card.Position))
		}
		processedText := transformText(card.Name, card.Text)
		sideAText := extractSideAText(processedText)
		sideBText := extractSideBText(processedText)
		sideBName := ""
		sideBType := ""
		if sideBText != "" {
			sideBName = card.Name
			sideBType = card.TypeName
		}
		willpower := getXVal(card.Willpower)
		attack := getXVal(card.Attack)
		defense := getXVal(card.Defense)
		health := getXVal(card.Health)
		is_unique := b2s(card.IsUnique)
		sphere := card.SphereName

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

			if card.TypeCode == "contract" {
				sphere = ""
			}

			if card.TypeCode != "player-side-quest" {
				questPoints = ""
			}

			if card.TypeCode == "event" || card.TypeCode == "contract" || card.TypeCode == "attachment" {
				willpower, attack, defense, health = "", "", "", ""
			}
			if card.TypeCode == "event" || card.TypeCode == "contract" {
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
				sphere,
				card.Traits,
				findKeywords(card.Text),
				cost,
				card.EngagementCost,
				threat,
				willpower, attack, defense, health,
				questPoints,
				victoryPoints,
				"", // Special Icon
				sideAText,
				card.Shadow,
				fixFlavor(card.Flavor),
				strconv.Itoa(card.Position),                   // printed card number
				"",                                // encounter set number
				card.EncounterSet,                 // encounter set icon
				"",                                // flags
				"",                                // artist
				"",                                // pan X
				"",                                // pan Y
				"",                                // scale
				"",                                // portait shadow
				sideBName,                         // Side B
				"",                                // is_unique,
				sideBType,                         // card.TypeName,
				"",                                // card.SphereName,
				"",                                // card.Traits,
				"",                                // findKeywords(card.Text),
				"",                                // cost,
				"",                                // EngagementCost,
				strconv.Itoa(card.ThreatStrength), // threat,
				"",                                // willpower
				"",                                // attack
				"",                                // defense
				"",                                // health
				strconv.Itoa(card.QuestPoints),    // questPoints
				strconv.Itoa(card.VictoryPoints),  // victoryPoints,
				"",                                // Special Icon
				sideBText,
				"",                   // flavor
				"",                   // shadow
				"",                   // printed card number
				"",                   // encounter set number
				"",                   // encounter set icon
				"",                   // flags
				"",                   // artist
				"",                   // pan X
				"",                   // pan Y
				"",                   // scale
				"",                   // portrait shadow
				"",                   // remove for easy
				"",                   // additional encounter set
				"",                   // adventure
				"",                   // collection icon
				"©FFG ©Middle-earth", // copyright
				"",                   //Deck Rules
				"",                   //Selected
				"",                   //Changed
				"",                   //Discord Bot
				"",
				"", //Current Snapshot
				cycles.GetCycleFromPack(card.PackName),
				card.PackCode,
				card.PackName,
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
	out := strings.ReplaceAll(text, name, "[name]") // insert name tag
	out = strip.StripTags(out)
	out = keywordPattern.ReplaceAllString(out, "")

	boldList := []string{"Travel"}
	for _, str := range boldList {
		out = strings.ReplaceAll(out, str, "[b]"+str+"[/b]")
	}

	out = regexp.MustCompile(`may trigger this (?:action|response)`).ReplaceAllString(out, "may trigger this effect")
	// heal on -> heal from
	out = regexp.MustCompile(`(\bheal[^.]+?\b)on(\b)`).ReplaceAllString(out, "${1}from${2}")

	// surround traits
	reggy := fmt.Sprintf(`([^(?:The One|Thrór's)]\s)(%s)\s(trait|or|and|cards?|ally|allies|attachments?|events?|heroes|hero|contracts?|characters?|enemy|enemies|location|or)(\.?)`, strings.Join(types.TraitsList, "|"))
	out = regexp.MustCompile(reggy).ReplaceAllString(out, "$1{${2}} $3")

	out = regexp.MustCompile(`\s([Tt]raits?)\s`).ReplaceAllString(out, " {${1}} ")

	out = regexp.MustCompile(`(\s)[^\[](tactics|leadership|spirit|lore)[^\]]([\s,\.])`).ReplaceAllString(out, "$1[$2]$3")
	//  make all newline groups exactly two newlines
	out = paragraphPattern.ReplaceAllLiteralString(out, "\r\n\r\n")
	return strings.TrimSpace(out)
}

func extractSideAText(text string) string {
	sideA := "Side A"
	sideB := "Side B"
	idxA := strings.Index(text, sideA)
	idxB := strings.Index(text, sideB)
	if idxA != -1 && idxB != -1 {
		return strings.TrimSpace(text[idxA+len(sideA) : idxB])
	}
	return text
}

func extractSideBText(text string) string {
	sideB := "Side B"
	if idx := strings.Index(text, sideB); idx != -1 {
		return strings.TrimSpace(text[idx+len(sideB):])
	}
	return ""
}

func deterministicUUID(uniqueVal string) string {
	md5hash := md5.New()
	md5hash.Write([]byte(uniqueVal))
	md5string := hex.EncodeToString(md5hash.Sum(nil))
	uuid, err := uuid.FromBytes([]byte(md5string[0:16]))
	if err != nil {
		log.Fatal(err)
	}
	return uuid.String()
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

func fixFlavor(val string) string {
	return regexp.MustCompile(`\s-+`).ReplaceAllString(val, " --")
}

func b2s(b bool) string {
	if b {
		return "1"
	}
	return ""
}

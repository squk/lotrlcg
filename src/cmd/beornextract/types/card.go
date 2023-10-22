package types

type HallOfBeornCard struct {
	PackCode      string `json:"pack_code,omitempty"`
	PackName      string `json:"pack_name,omitempty"`
	IsOfficial    bool   `json:"is_official,omitempty"`
	TypeCode      string `json:"type_code,omitempty"`
	TypeName      string `json:"type_name,omitempty"`
	SphereCode    string `json:"sphere_code,omitempty"`
	SphereName    string `json:"sphere_name,omitempty"`
	Position      int    `json:"position,omitempty"`
	Code          string `json:"code,omitempty"`
	Cost          string `json:"cost,omitempty"`
	Name          string `json:"name,omitempty"`
	Traits        string `json:"traits,omitempty"`
	Text          string `json:"text,omitempty"`
	Flavor        string `json:"flavor,omitempty"`
	Shadow        string `json:"shadow,omitempty"`
	IsUnique      bool   `json:"is_unique,omitempty"`
	Threat        int    `json:"threat,omitempty"`
	Willpower     int    `json:"willpower,omitempty"`
	Attack        int    `json:"attack,omitempty"`
	Defense       int    `json:"defense,omitempty"`
	Health        int    `json:"health,omitempty"`
	Quantity      int    `json:"quantity,omitempty"`
	DeckLimit     int    `json:"deck_limit,omitempty"`
	Illustrator   string `json:"illustrator,omitempty"`
	Octgnid       string `json:"octgnid,omitempty"`
	HasErrata     bool   `json:"has_errata,omitempty"`
	URL           string `json:"url,omitempty"`
	QuestPoints   int    `json:"quest,omitempty"`
	VictoryPoints int    `json:"victory,omitempty"`
	Imagesrc      string `json:"imagesrc,omitempty"`

	EncounterSet   string `json:"encounter_set,omitempty"`
	EngagementCost string `json:"engagement_cost,omitempty"`
	ThreatStrength int    `json:"threat_strength,omitempty"`
}

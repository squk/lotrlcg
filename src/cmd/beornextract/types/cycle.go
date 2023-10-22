package types

type CycleMappings []struct {
	Cycle string   `json:"cycle"`
	Packs []string `json:"packs"`
}

func (c *CycleMappings) GetCycleFromPack(pack string) string {
	for _, mapping := range *c {
		for _, cyclePack := range mapping.Packs {
			if cyclePack == pack {
				return mapping.Cycle
			}
		}
	}

	return "IDK lol";
}

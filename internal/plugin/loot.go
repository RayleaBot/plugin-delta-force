package plugin

import (
	"fmt"
	"strings"
	"time"
)

type lootEngine struct {
	catalog *lootCatalog
	random  randomSource
}

type lootResult struct {
	Container container
	Profile   profile
	Items     []item
	Highest   string
}

func newLootEngine(catalog *lootCatalog, random randomSource) *lootEngine {
	return &lootEngine{catalog: catalog, random: random}
}

func (engine *lootEngine) roll(query string) (lootResult, error) {
	owner, exists := engine.catalog.findContainer(query)
	if !exists {
		return lootResult{}, fmt.Errorf("unsupported container %q", strings.TrimSpace(query))
	}
	count := owner.MinItems
	if width := owner.MaxItems - owner.MinItems + 1; width > 1 {
		offset, err := engine.random.Intn(width)
		if err != nil {
			return lootResult{}, err
		}
		count += offset
	}
	lootProfile := engine.catalog.profiles[owner.Profile]
	selected := make(map[string]struct{}, count)
	items := make([]item, 0, count)
	highest := rarityOrder[0]
	for len(items) < count {
		rarity, err := engine.pickRarity(lootProfile)
		if err != nil {
			return lootResult{}, err
		}
		picked, err := engine.pickItem(owner, rarity, selected)
		if err != nil {
			return lootResult{}, err
		}
		selected[picked.ID] = struct{}{}
		items = append(items, picked)
		if rarityRank(picked.Rarity) > rarityRank(highest) {
			highest = picked.Rarity
		}
	}
	return lootResult{Container: owner, Profile: lootProfile, Items: items, Highest: highest}, nil
}

func (engine *lootEngine) pickRarity(lootProfile profile) (string, error) {
	value, err := engine.random.Intn(weightScale)
	if err != nil {
		return "", err
	}
	accumulator := 0
	for _, rarity := range rarityOrder {
		accumulator += lootProfile.Weights[rarity]
		if value < accumulator {
			return rarity, nil
		}
	}
	return "", fmt.Errorf("loot profile %q did not select a rarity", lootProfile.ID)
}

func (engine *lootEngine) pickItem(owner container, rarity string, selected map[string]struct{}) (item, error) {
	for _, candidateRarity := range fallbackRarities(rarity) {
		candidates := engine.catalog.candidates(owner, candidateRarity, selected)
		if len(candidates) == 0 {
			continue
		}
		byCategory := make(map[string][]item)
		categories := make([]string, 0)
		for _, candidate := range candidates {
			if len(byCategory[candidate.Category]) == 0 {
				categories = append(categories, candidate.Category)
			}
			byCategory[candidate.Category] = append(byCategory[candidate.Category], candidate)
		}
		categoryIndex, err := engine.random.Intn(len(categories))
		if err != nil {
			return item{}, err
		}
		categoryItems := byCategory[categories[categoryIndex]]
		itemIndex, err := engine.random.Intn(len(categoryItems))
		if err != nil {
			return item{}, err
		}
		return categoryItems[itemIndex], nil
	}
	return item{}, fmt.Errorf("container %q has no remaining loot candidates", owner.ID)
}

func fallbackRarities(selected string) []string {
	index := rarityRank(selected)
	result := []string{selected}
	for offset := 1; offset < len(rarityOrder); offset++ {
		if lower := index - offset; lower >= 0 {
			result = append(result, rarityOrder[lower])
		}
		if higher := index + offset; higher < len(rarityOrder) {
			result = append(result, rarityOrder[higher])
		}
	}
	return result
}

func rarityRank(rarity string) int {
	for index, candidate := range rarityOrder {
		if candidate == rarity {
			return index
		}
	}
	return -1
}

func lootRenderData(result lootResult, actorID, actorName, eventID string, now time.Time, catalog *lootCatalog) map[string]any {
	items := make([]map[string]any, 0, len(result.Items))
	for index, candidate := range result.Items {
		items = append(items, map[string]any{
			"index":          fmt.Sprintf("%02d", index+1),
			"name":           candidate.Name,
			"rarity":         rarityLabels[candidate.Rarity],
			"rarity_key":     candidate.Rarity,
			"category":       categoryLabels[candidate.Category],
			"category_short": categoryShort(candidate.Category),
		})
	}
	if strings.TrimSpace(actorName) == "" {
		actorName = actorID
	}
	serial := strings.TrimSpace(eventID)
	if len(serial) > 12 {
		serial = serial[len(serial)-12:]
	}
	if serial == "" {
		serial = now.Format("150405.000")
	}
	return map[string]any{
		"title":             result.Container.Name,
		"profile":           result.Profile.Label,
		"grade":             rarityLabels[result.Highest] + "收获",
		"grade_key":         result.Highest,
		"items":             items,
		"item_count":        len(items),
		"actor_name":        actorName,
		"actor_id":          actorID,
		"time":              now.Format("2006-01-02 15:04:05"),
		"serial":            serial,
		"container_source":  catalog.containerSource,
		"item_source":       catalog.itemSource,
		"simulation_notice": "娱乐模拟概率，不代表游戏实际掉落",
	}
}

func categoryShort(category string) string {
	switch category {
	case "electronics":
		return "电子"
	case "household":
		return "家居"
	case "equipment":
		return "装备"
	case "keycard":
		return "门卡"
	default:
		label := categoryLabels[category]
		runes := []rune(label)
		if len(runes) > 2 {
			return string(runes[:2])
		}
		return label
	}
}

func formatLoot(result lootResult) string {
	lines := []string{"摸到了「" + result.Container.Name + "」", "容器档位：" + result.Profile.Label}
	for _, candidate := range result.Items {
		lines = append(lines, fmt.Sprintf("- [%s] %s · %s", rarityLabels[candidate.Rarity], candidate.Name, categoryLabels[candidate.Category]))
	}
	lines = append(lines, "娱乐模拟概率，不代表游戏实际掉落。")
	return strings.Join(lines, "\n")
}

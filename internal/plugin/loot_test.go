package plugin

import "testing"

type sequenceRandom struct {
	values []int
	index  int
}

func (source *sequenceRandom) Intn(max int) (int, error) {
	if source.index >= len(source.values) {
		return 0, nil
	}
	value := source.values[source.index]
	source.index++
	if value < 0 {
		value = -value
	}
	return value % max, nil
}

func TestLootRollUsesContainerProfileAndReturnsUniqueItems(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	random := &sequenceRandom{values: []int{
		2,
		9999, 0, 0,
		9000, 1, 0,
		6000, 2, 0,
		700, 3, 0,
	}}
	result, err := newLootEngine(catalog, random).roll("航空箱")
	if err != nil {
		t.Fatal(err)
	}
	if result.Container.Name != "航空储物箱" || result.Profile.ID != "s" {
		t.Fatalf("unexpected roll owner/profile: %#v", result)
	}
	if len(result.Items) != 4 {
		t.Fatalf("item count = %d, want 4", len(result.Items))
	}
	if result.Highest != "red" {
		t.Fatalf("highest rarity = %q, want red", result.Highest)
	}
	seen := map[string]struct{}{}
	for _, candidate := range result.Items {
		if _, exists := seen[candidate.ID]; exists {
			t.Fatalf("duplicate item selected: %q", candidate.ID)
		}
		seen[candidate.ID] = struct{}{}
	}
}

func TestRarityBoundaries(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	profile := catalog.profiles["a"]
	tests := []struct {
		value int
		want  string
	}{
		{0, "white"},
		{499, "white"},
		{500, "green"},
		{1999, "green"},
		{2000, "blue"},
		{5999, "blue"},
		{6000, "purple"},
		{9199, "purple"},
		{9200, "gold"},
		{9949, "gold"},
		{9950, "red"},
		{9999, "red"},
	}
	for _, test := range tests {
		t.Run(test.want, func(t *testing.T) {
			engine := newLootEngine(catalog, &sequenceRandom{values: []int{test.value}})
			got, err := engine.pickRarity(profile)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Fatalf("pickRarity(%d) = %q, want %q", test.value, got, test.want)
			}
		})
	}
}

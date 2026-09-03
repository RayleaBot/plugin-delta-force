package plugin

import "testing"

func TestEmbeddedCatalogIsComplete(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(catalog.containers), 28; got != want {
		t.Fatalf("container count = %d, want %d", got, want)
	}
	if len(catalog.items) < 70 {
		t.Fatalf("item count = %d, want at least 70", len(catalog.items))
	}
	for _, alias := range []string{"航空箱", "大保险柜", "白衣", "医疗堆", "机箱"} {
		if _, ok := catalog.findContainer(alias); !ok {
			t.Errorf("alias %q did not resolve", alias)
		}
	}
}

func TestContainerNormalizationIgnoresSpacingAndSeparators(t *testing.T) {
	if got, want := normalizeContainer(" 航空-储物 箱 "), normalizeContainer("航空储物箱"); got != want {
		t.Fatalf("normalizeContainer() = %q, want %q", got, want)
	}
}

func TestEveryWeightedRarityHasCandidates(t *testing.T) {
	catalog, err := loadCatalog()
	if err != nil {
		t.Fatal(err)
	}
	for _, owner := range catalog.containers {
		profile := catalog.profiles[owner.Profile]
		for rarity, weight := range profile.Weights {
			if weight == 0 {
				continue
			}
			if got := len(catalog.candidates(owner, rarity, nil)); got == 0 {
				t.Errorf("container %q has no %s candidates", owner.Name, rarity)
			}
		}
	}
}

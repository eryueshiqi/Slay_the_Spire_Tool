package app

import (
	"os"
	"path/filepath"
	"testing"

	"slay_the_spire_tool/internal/model"
)

func TestAppBasicFlow(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	writeTestData(t, tmpDir)

	a, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}

	roles, err := a.GetRoleList()
	if err != nil {
		t.Fatalf("GetRoleList() error = %v, want nil", err)
	}
	if len(roles) == 0 {
		t.Fatal("GetRoleList() returned empty roles")
	}

	decks, err := a.GetDeckList(model.AnyRole)
	if err != nil {
		t.Fatalf("GetDeckList() error = %v, want nil", err)
	}
	if len(decks) != 1 {
		t.Fatalf("GetDeckList() len = %d, want 1", len(decks))
	}
	if decks[0].Name != "Ironclad Starter Deck" {
		t.Fatalf("GetDeckList()[0].Name = %s, want Ironclad Starter Deck", decks[0].Name)
	}

	if err := a.ChooseDeck(decks[0].Id); err != nil {
		t.Fatalf("ChooseDeck() error = %v, want nil", err)
	}

	completion, err := a.GetDeckCompletionStatus()
	if err != nil {
		t.Fatalf("GetDeckCompletionStatus() error = %v, want nil", err)
	}
	if completion.Completed || completion.MissingCount != 1 {
		t.Fatalf("GetDeckCompletionStatus() = %+v, want completed=false missing=1", completion)
	}

	stats, err := a.GetSelectedDeckTypeStats()
	if err != nil {
		t.Fatalf("GetSelectedDeckTypeStats() error = %v, want nil", err)
	}
	if len(stats) == 0 {
		t.Fatal("GetSelectedDeckTypeStats() returned empty stats, want fixed type stats")
	}

	deckCards, err := a.GetSelectedDeckCardStats(model.AnyType, "", "all")
	if err != nil {
		t.Fatalf("GetSelectedDeckCardStats() error = %v, want nil", err)
	}
	if len(deckCards) != 1 {
		t.Fatalf("GetSelectedDeckCardStats() len = %d, want 1", len(deckCards))
	}
	if deckCards[0].ID != "anger" || deckCards[0].Count != 1 {
		t.Fatalf("GetSelectedDeckCardStats()[0] = %+v, want anger x1", deckCards[0])
	}

	detail, err := a.GetCardDetail("anger")
	if err != nil {
		t.Fatalf("GetCardDetail() error = %v, want nil", err)
	}
	if detail.ID != "anger" {
		t.Fatalf("GetCardDetail().ID = %s, want anger", detail.ID)
	}

	cards, err := a.GetCardList("", model.AnyType, model.AnyRole, model.AnyCost, model.AnyRarity, true)
	if err != nil {
		t.Fatalf("GetCardList() error = %v, want nil", err)
	}
	if len(cards) != 1 {
		t.Fatalf("GetCardList(isBase=true) len = %d, want 1", len(cards))
	}

	allCards, err := a.GetCardList("", model.AnyType, model.AnyRole, model.AnyCost, model.AnyRarity, false)
	if err != nil {
		t.Fatalf("GetCardList(isBase=false) error = %v, want nil", err)
	}
	if len(allCards) != 1 {
		t.Fatalf("GetCardList(isBase=false) len = %d, want 1", len(allCards))
	}
	if allCards[0].ID != "aggression" {
		t.Fatalf("GetCardList(isBase=false) first id = %s, want aggression", allCards[0].ID)
	}

	if err := a.AddCardToUserDeck("anger"); err != nil {
		t.Fatalf("AddCardToUserDeck() error = %v, want nil", err)
	}

	completion, err = a.GetDeckCompletionStatus()
	if err != nil {
		t.Fatalf("GetDeckCompletionStatus() after add error = %v, want nil", err)
	}
	if !completion.Completed || completion.MissingCount != 0 {
		t.Fatalf("GetDeckCompletionStatus() after add = %+v, want completed=true missing=0", completion)
	}

	userDeck, err := a.GetUserDeck()
	if err != nil {
		t.Fatalf("GetUserDeck() error = %v, want nil", err)
	}
	if len(userDeck.Cards) != 1 {
		t.Fatalf("GetUserDeck().Cards len = %d, want 1", len(userDeck.Cards))
	}

	filtered, err := a.GetCardList("aggression", model.Power, model.Ironclad, model.One, model.Rare, false)
	if err != nil {
		t.Fatalf("filtered GetCardList() error = %v, want nil", err)
	}
	if len(filtered) != 1 || filtered[0].ID != "aggression" {
		t.Fatalf("filtered GetCardList() = %+v, want aggression", filtered)
	}
	if filtered[0].Need {
		t.Fatalf("filtered card Need = true, want false: %+v", filtered[0])
	}
}

func TestRemoveDuplicateOneByOne(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	writeTestData(t, tmpDir)

	a, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}

	if err := a.ChooseDeck(1); err != nil {
		t.Fatalf("ChooseDeck() error = %v, want nil", err)
	}

	if err := a.AddCardToUserDeck("anger"); err != nil {
		t.Fatalf("AddCardToUserDeck first error = %v, want nil", err)
	}
	if err := a.AddCardToUserDeck("anger"); err != nil {
		t.Fatalf("AddCardToUserDeck second error = %v, want nil", err)
	}

	deck, err := a.GetUserDeck()
	if err != nil {
		t.Fatalf("GetUserDeck() error = %v, want nil", err)
	}
	if got := countByID(deck, "anger"); got != 2 {
		t.Fatalf("countByID(anger) = %d, want 2", got)
	}

	if err := a.RemoveCardFromUserDeck("anger"); err != nil {
		t.Fatalf("RemoveCardFromUserDeck first error = %v, want nil", err)
	}
	deck, err = a.GetUserDeck()
	if err != nil {
		t.Fatalf("GetUserDeck() after first remove error = %v, want nil", err)
	}
	if got := countByID(deck, "anger"); got != 1 {
		t.Fatalf("countByID(anger) after first remove = %d, want 1", got)
	}

	if err := a.RemoveCardFromUserDeck("anger"); err != nil {
		t.Fatalf("RemoveCardFromUserDeck second error = %v, want nil", err)
	}
	deck, err = a.GetUserDeck()
	if err != nil {
		t.Fatalf("GetUserDeck() after second remove error = %v, want nil", err)
	}
	if got := countByID(deck, "anger"); got != 0 {
		t.Fatalf("countByID(anger) after second remove = %d, want 0", got)
	}

	if err := a.RemoveCardFromUserDeck("anger"); err == nil {
		t.Fatal("RemoveCardFromUserDeck third error = nil, want non-nil")
	}
}

func TestListReturnsEmptySliceNotNil(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	writeTestData(t, tmpDir)

	a, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}

	decks, err := a.GetDeckList(model.Defect)
	if err != nil {
		t.Fatalf("GetDeckList() error = %v, want nil", err)
	}
	if decks == nil {
		t.Fatal("GetDeckList() returned nil slice, want empty slice")
	}

	cards, err := a.GetCardList("not-exist", model.AnyType, model.AnyRole, model.AnyCost, model.AnyRarity, true)
	if err != nil {
		t.Fatalf("GetCardList() error = %v, want nil", err)
	}
	if cards == nil {
		t.Fatal("GetCardList() returned nil slice, want empty slice")
	}
}

func TestChooseDeckToggleUnselect(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	writeTestData(t, tmpDir)

	a, err := New(tmpDir)
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}

	if err := a.ChooseDeck(1); err != nil {
		t.Fatalf("ChooseDeck first error = %v, want nil", err)
	}
	if err := a.AddCardToUserDeck("anger"); err != nil {
		t.Fatalf("AddCardToUserDeck error = %v, want nil", err)
	}

	deck, err := a.GetUserDeck()
	if err != nil {
		t.Fatalf("GetUserDeck() error = %v, want nil", err)
	}
	if got := len(deck.Cards); got != 1 {
		t.Fatalf("GetUserDeck().Cards len = %d, want 1", got)
	}

	if err := a.ChooseDeck(1); err != nil {
		t.Fatalf("ChooseDeck second(toggle off) error = %v, want nil", err)
	}

	completion, err := a.GetDeckCompletionStatus()
	if err != nil {
		t.Fatalf("GetDeckCompletionStatus() after toggle off error = %v, want nil", err)
	}
	if completion.Completed {
		t.Fatalf("GetDeckCompletionStatus().Completed after toggle off = true, want false")
	}

	stats, err := a.GetSelectedDeckTypeStats()
	if err != nil {
		t.Fatalf("GetSelectedDeckTypeStats() after toggle off error = %v, want nil", err)
	}
	if len(stats) != 0 {
		t.Fatalf("GetSelectedDeckTypeStats() len after toggle off = %d, want 0", len(stats))
	}

	deck, err = a.GetUserDeck()
	if err != nil {
		t.Fatalf("GetUserDeck() after toggle off error = %v, want nil", err)
	}
	if got := len(deck.Cards); got != 0 {
		t.Fatalf("GetUserDeck().Cards len after toggle off = %d, want 0", got)
	}

	cards, err := a.GetCardList("anger", model.AnyType, model.AnyRole, model.AnyCost, model.AnyRarity, true)
	if err != nil {
		t.Fatalf("GetCardList() after toggle off error = %v, want nil", err)
	}
	if len(cards) != 1 {
		t.Fatalf("GetCardList() len = %d, want 1", len(cards))
	}
	if cards[0].Need || cards[0].HighLight {
		t.Fatalf("card flags after toggle off should be false, got Need=%v HighLight=%v", cards[0].Need, cards[0].HighLight)
	}
}

func TestNewDataDirNotFound(t *testing.T) {
	t.Parallel()

	_, err := New(filepath.Join(t.TempDir(), "missing-data"))
	if err == nil {
		t.Fatal("New() error = nil, want non-nil")
	}
}

func writeTestData(t *testing.T, dir string) {
	t.Helper()

	cardsJSON := `{
  "cards": [
    {
      "id": "anger",
      "name": "Anger",
      "alias": "Anger",
      "image": "assets/images/anger.png",
      "type": 0,
      "role": 0,
      "cost": 0,
      "rarity": 1,
      "is_base": true
    },
    {
      "id": "aggression",
      "name": "Aggression",
      "alias": "Aggression",
      "image": "assets/images/aggression.png",
      "type": 2,
      "role": 0,
      "cost": 1,
      "rarity": 3,
      "is_base": false
    }
  ]
}`
	decksJSON := `{
  "decks": [
    {
      "id": 1,
      "name": "Ironclad Starter Deck",
      "role": 0,
      "cards": [
        {
          "id": "anger",
          "image": "assets/images/anger.png",
          "need": true,
          "highlight": false
        }
      ]
    }
  ]
}`

	if err := os.WriteFile(filepath.Join(dir, "cards.json"), []byte(cardsJSON), 0o644); err != nil {
		t.Fatalf("write cards.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "decks.json"), []byte(decksJSON), 0o644); err != nil {
		t.Fatalf("write decks.json: %v", err)
	}
}

func countByID(deck model.Deck, id string) int {
	count := 0
	for _, card := range deck.Cards {
		if card.ID == id {
			count++
		}
	}
	return count
}

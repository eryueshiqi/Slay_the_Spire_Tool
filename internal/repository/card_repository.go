package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"slay_the_spire_tool/internal/model"
	"slices"
	"strings"
)

type CardRepository struct {
	Cards     []model.Card
	RoleCards []model.RoleCards
	Decks     []model.Deck

	Deck     model.Deck
	UserDeck model.Deck

	hasSelectedDeck bool
	selectedDeckID  uint
}

type cardsPayload struct {
	Cards []model.Card `json:"cards"`
}

type decksPayload struct {
	Decks []model.Deck `json:"decks"`
}

func (c *CardRepository) GetRoleList(ctx context.Context) ([]model.Role, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	resp := make([]model.Role, 0, 6)
	resp = append(resp,
		model.Ironclad,
		model.Silent,
		model.Regent,
		model.Necrobinder,
		model.Defect,
		model.Colorless,
	)
	return resp, nil
}

func (c *CardRepository) GetCardList(ctx context.Context, str string, theType model.Type, role model.Role, cost model.Cost, rarity model.Rarity, isBase bool) ([]model.CardInfo, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	resp := make([]model.CardInfo, 0, len(c.Cards))

	query := strings.TrimSpace(strings.ToLower(str))

	for _, card := range c.Cards {
		if query != "" &&
			!strings.Contains(strings.ToLower(card.ID), query) &&
			!strings.Contains(strings.ToLower(card.Name), query) &&
			!strings.Contains(strings.ToLower(card.Alias), query) {
			continue
		}
		if role != model.AnyRole && card.Role != role {
			continue
		}
		if theType != model.AnyType && card.Type != theType {
			continue
		}
		if cost != model.AnyCost && card.Cost != cost {
			continue
		}
		if rarity != model.AnyRarity && card.Rarity != rarity {
			continue
		}
		if card.IsBase != isBase {
			continue
		}

		cardInfo := model.CardInfo{
			ID:         card.ID,
			Image:      card.Image,
			Need:       c.isTargetCard(card.ID),
			HighLight:  c.userDeckCountByID(card.ID) > 0,
			NeedDelete: false,
		}
		resp = append(resp, cardInfo)
	}

	return resp, nil
}

func (c *CardRepository) GetDeckList(ctx context.Context, role model.Role) ([]model.Deck, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if role == model.AnyRole {
		resp := make([]model.Deck, 0, len(c.Decks))
		resp = append(resp, c.Decks...)
		return resp, nil
	}

	resp := make([]model.Deck, 0)
	for _, deck := range c.Decks {
		if deck.Role == role {
			resp = append(resp, deck)
		}
	}
	return resp, nil
}

func (c *CardRepository) GetDeck(ctx context.Context, id uint) (model.Deck, error) {
	if err := ctx.Err(); err != nil {
		return model.Deck{}, err
	}

	for _, deck := range c.Decks {
		if deck.Id == id {
			return deck, nil
		}
	}
	return model.Deck{}, errors.New("牌组不存在")
}

func (c *CardRepository) GetCard(ctx context.Context, id string, role model.Role) (model.CardInfo, error) {
	if err := ctx.Err(); err != nil {
		return model.CardInfo{}, err
	}

	info := model.CardInfo{
		ID:        id,
		Need:      c.isTargetCard(id),
		HighLight: c.userDeckCountByID(id) > 0,
	}

	for _, roleCard := range c.RoleCards {
		if roleCard.Role == role || roleCard.Role == model.Colorless {
			for _, card := range roleCard.Cards {
				if card.ID == id {
					info.Image = card.Image
					info.Need = false
				}
			}
		}
	}

	if info.Image == "" {
		return model.CardInfo{}, errors.New("卡牌不存在")
	}

	for _, deckCard := range c.Deck.Cards {
		if deckCard.ID == id {
			info.Need = true
		}
	}

	return info, nil
}

func (c *CardRepository) GetCardDetail(ctx context.Context, id string) (model.Card, error) {
	if err := ctx.Err(); err != nil {
		return model.Card{}, err
	}

	card, ok := c.findCardByID(id)
	if !ok {
		return model.Card{}, errors.New("卡牌不存在")
	}
	return card, nil
}

func (c *CardRepository) GetSelectedDeckTypeStats(ctx context.Context) ([]model.CardTypeStat, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if !c.hasSelectedDeck {
		return make([]model.CardTypeStat, 0), nil
	}

	countByType := map[model.Type]int{
		model.Attack: 0,
		model.Skill:  0,
		model.Power:  0,
		model.Quest:  0,
		model.Status: 0,
		model.Curse:  0,
	}

	for _, deckCard := range c.Deck.Cards {
		card, ok := c.findCardByID(deckCard.ID)
		if !ok {
			continue
		}
		countByType[card.Type]++
	}

	orderedTypes := []model.Type{
		model.Attack,
		model.Skill,
		model.Power,
		model.Quest,
		model.Status,
		model.Curse,
	}
	resp := make([]model.CardTypeStat, 0, len(orderedTypes))
	for _, t := range orderedTypes {
		resp = append(resp, model.CardTypeStat{
			Type:  t,
			Count: countByType[t],
		})
	}

	return resp, nil
}

func (c *CardRepository) GetSelectedDeckCardStats(ctx context.Context, theType model.Type, keyword string, upgradeFilter string) ([]model.DeckCardStat, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if !c.hasSelectedDeck {
		return make([]model.DeckCardStat, 0), nil
	}

	countByID := make(map[string]int)
	for _, deckCard := range c.Deck.Cards {
		countByID[deckCard.ID]++
	}

	query := strings.TrimSpace(strings.ToLower(keyword))
	resp := make([]model.DeckCardStat, 0, len(countByID))

	for cardID, count := range countByID {
		card, ok := c.findCardByID(cardID)
		if !ok {
			continue
		}
		if theType != model.AnyType && card.Type != theType {
			continue
		}
		if query != "" &&
			!strings.Contains(strings.ToLower(card.ID), query) &&
			!strings.Contains(strings.ToLower(card.Name), query) &&
			!strings.Contains(strings.ToLower(card.Alias), query) {
			continue
		}
		if upgradeFilter == "base" && !card.IsBase {
			continue
		}
		if upgradeFilter == "upgraded" && card.IsBase {
			continue
		}

		resp = append(resp, model.DeckCardStat{
			ID:     card.ID,
			Name:   card.Name,
			Image:  card.Image,
			Type:   card.Type,
			IsBase: card.IsBase,
			Count:  count,
		})
	}

	slices.SortFunc(resp, func(a, b model.DeckCardStat) int {
		if a.Type != b.Type {
			if a.Type < b.Type {
				return -1
			}
			return 1
		}
		if a.Name != b.Name {
			return strings.Compare(a.Name, b.Name)
		}
		return strings.Compare(a.ID, b.ID)
	})

	return resp, nil
}

func (c *CardRepository) GetDeckCompletionStatus(ctx context.Context) (model.DeckCompletion, error) {
	if err := ctx.Err(); err != nil {
		return model.DeckCompletion{}, err
	}

	if !c.hasSelectedDeck {
		return model.DeckCompletion{
			Completed:    false,
			MissingCount: 0,
		}, nil
	}

	targetCounts := make(map[string]int)
	for _, deckCard := range c.Deck.Cards {
		targetCounts[deckCard.ID]++
	}

	userCounts := make(map[string]int)
	for _, userCard := range c.UserDeck.Cards {
		if userCard.NeedDelete {
			continue
		}
		userCounts[userCard.ID]++
	}

	missingCount := 0
	for cardID, needCount := range targetCounts {
		current := userCounts[cardID]
		if current < needCount {
			missingCount += needCount - current
		}
	}

	return model.DeckCompletion{
		Completed:    missingCount == 0,
		MissingCount: missingCount,
	}, nil
}

func (c *CardRepository) ChooseDeck(ctx context.Context, id uint) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if c.hasSelectedDeck && c.selectedDeckID == id {
		c.Deck = model.Deck{
			Cards: make([]model.CardInfo, 0),
		}
		c.UserDeck = model.Deck{
			Cards: make([]model.CardInfo, 0),
		}
		c.hasSelectedDeck = false
		c.selectedDeckID = 0
		return nil
	}

	for _, deck := range c.Decks {
		if deck.Id == id {
			c.Deck = deck
			c.UserDeck = model.Deck{
				Role:  deck.Role,
				Cards: make([]model.CardInfo, 0),
			}
			c.hasSelectedDeck = true
			c.selectedDeckID = deck.Id
			return nil
		}
	}

	return errors.New("牌组不存在")
}

func (c *CardRepository) ClearUserDeck(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.UserDeck = model.Deck{
		Role:  c.UserDeck.Role,
		Cards: make([]model.CardInfo, 0),
	}
	return nil
}

func (c *CardRepository) GetUserDeck(ctx context.Context) (model.Deck, error) {
	if err := ctx.Err(); err != nil {
		return model.Deck{}, err
	}

	liveCards := make([]model.CardInfo, 0, len(c.UserDeck.Cards))
	for _, card := range c.UserDeck.Cards {
		if !card.NeedDelete {
			liveCards = append(liveCards, card)
		}
	}

	slices.SortFunc(liveCards, func(a, b model.CardInfo) int {
		if a.Need != b.Need {
			if a.Need {
				return -1
			}
			if b.Need {
				return 1
			}
		}

		return strings.Compare(a.ID, b.ID)
	})

	c.UserDeck.Cards = liveCards
	return c.UserDeck, nil
}

func (c *CardRepository) AddCardToUserDeck(ctx context.Context, cardID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	info := model.CardInfo{
		ID: cardID,
	}

	foundCard := false
	for _, roleCard := range c.RoleCards {
		if roleCard.Role == c.UserDeck.Role || roleCard.Role == model.Colorless {
			for _, card := range roleCard.Cards {
				if card.ID == cardID {
					info.Image = card.Image
					foundCard = true
					break
				}
			}
		}
	}

	if !foundCard {
		return errors.New("卡牌不存在")
	}

	for _, deckCard := range c.Deck.Cards {
		if deckCard.ID == cardID {
			info.Need = true
		}
	}

	c.UserDeck.Cards = append(c.UserDeck.Cards, info)
	return nil
}

func (c *CardRepository) RemoveCardFromUserDeck(ctx context.Context, cardID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	found := false
	for i := range c.UserDeck.Cards {
		if c.UserDeck.Cards[i].ID == cardID && !c.UserDeck.Cards[i].NeedDelete {
			c.UserDeck.Cards[i].Image = ""
			c.UserDeck.Cards[i].NeedDelete = true
			found = true
			break
		}
	}

	if !found {
		return errors.New("用户卡组中不存在该卡牌")
	}

	return nil
}

func (c *CardRepository) isTargetCard(cardID string) bool {
	for _, deckCard := range c.Deck.Cards {
		if deckCard.ID == cardID {
			return true
		}
	}
	return false
}

func (c *CardRepository) userDeckCountByID(cardID string) int {
	count := 0
	for _, userCard := range c.UserDeck.Cards {
		if userCard.ID == cardID && !userCard.NeedDelete {
			count++
		}
	}
	return count
}

func (c *CardRepository) findCardByID(cardID string) (model.Card, bool) {
	for _, card := range c.Cards {
		if card.ID == cardID {
			return card, true
		}
	}
	return model.Card{}, false
}

func NewCardRepositoryFromFS(dataFS fs.FS) (*CardRepository, error) {
	return newCardRepository(dataFS)
}

func NewCardRepository(filePath string) (*CardRepository, error) {
	return newCardRepository(os.DirFS(filePath))
}

func newCardRepository(dataFS fs.FS) (*CardRepository, error) {
	bytes, err := fs.ReadFile(dataFS, path.Join("cards.json"))
	if err != nil {
		return nil, fmt.Errorf("读取 cards.json 失败: %w", err)
	}

	var cardsResp cardsPayload
	err = json.Unmarshal(bytes, &cardsResp)
	if err != nil {
		return nil, fmt.Errorf("解析 cards.json 失败: %w", err)
	}

	cards := cardsResp.Cards
	var ironcladCard []model.Card
	var silentCard []model.Card
	var regentCard []model.Card
	var necrobinderCard []model.Card
	var defectCard []model.Card
	var colorlessCard []model.Card
	for _, card := range cards {
		if card.Role == model.Ironclad {
			ironcladCard = append(ironcladCard, card)
		}
		if card.Role == model.Silent {
			silentCard = append(silentCard, card)
		}
		if card.Role == model.Regent {
			regentCard = append(regentCard, card)
		}
		if card.Role == model.Necrobinder {
			necrobinderCard = append(necrobinderCard, card)
		}
		if card.Role == model.Defect {
			defectCard = append(defectCard, card)
		}
		if card.Role == model.Colorless {
			colorlessCard = append(colorlessCard, card)
		}
	}
	roleCards := []model.RoleCards{
		{
			Role:  model.Ironclad,
			Cards: ironcladCard,
		},
		{
			Role:  model.Silent,
			Cards: silentCard,
		},
		{
			Role:  model.Regent,
			Cards: regentCard,
		},
		{
			Role:  model.Necrobinder,
			Cards: necrobinderCard,
		},
		{
			Role:  model.Defect,
			Cards: defectCard,
		},
		{
			Role:  model.Colorless,
			Cards: colorlessCard,
		},
	}

	bytes, err = fs.ReadFile(dataFS, path.Join("decks.json"))
	if err != nil {
		return nil, fmt.Errorf("读取 decks.json 失败: %w", err)
	}

	var decksResp decksPayload
	err = json.Unmarshal(bytes, &decksResp)
	if err != nil {
		return nil, fmt.Errorf("解析 decks.json 失败: %w", err)
	}

	decks := decksResp.Decks

	return &CardRepository{
		Cards:           cards,
		RoleCards:       roleCards,
		Decks:           decks,
		Deck:            model.Deck{Cards: make([]model.CardInfo, 0)},
		UserDeck:        model.Deck{Cards: make([]model.CardInfo, 0)},
		hasSelectedDeck: false,
		selectedDeckID:  0,
	}, nil
}

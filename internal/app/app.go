package app

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"slay_the_spire_tool/internal/model"
	"slay_the_spire_tool/internal/repository"
	"slay_the_spire_tool/internal/service"
)

// App is the backend entry object that can be bound to frontend by Wails.
type App struct {
	cardService *service.CardService
}

func New(dataDir string) (*App, error) {
	return NewWithDataFS(os.DirFS(dataDir))
}

func NewWithDataFS(dataFS fs.FS) (*App, error) {
	cardRepo, err := repository.NewCardRepositoryFromFS(dataFS)
	if err != nil {
		return nil, fmt.Errorf("初始化卡牌仓储失败: %w", err)
	}
	cardSvc := service.NewCardService(cardRepo)

	return &App{
		cardService: cardSvc,
	}, nil
}

func (a *App) GetRoleList() ([]model.Role, error) {
	return a.cardService.GetRoleList(context.Background())
}

func (a *App) GetCardList(str string, theType model.Type, role model.Role, cost model.Cost, rarity model.Rarity, isBase bool) ([]model.CardInfo, error) {
	return a.cardService.GetCardList(context.Background(), str, theType, role, cost, rarity, isBase)
}

func (a *App) GetDisplayLabels() model.DisplayLabels {
	return model.BuildDisplayLabelsZH()
}

func (a *App) GetDeckList(role model.Role) ([]model.Deck, error) {
	return a.cardService.GetDeckList(context.Background(), role)
}

func (a *App) GetDeck(id uint) (model.Deck, error) {
	return a.cardService.GetDeck(context.Background(), id)
}

func (a *App) GetCard(id string, role model.Role) (model.CardInfo, error) {
	return a.cardService.GetCard(context.Background(), id, role)
}

func (a *App) GetCardDetail(id string) (model.Card, error) {
	return a.cardService.GetCardDetail(context.Background(), id)
}

func (a *App) ChooseDeck(id uint) error {
	return a.cardService.ChooseDeck(context.Background(), id)
}

func (a *App) GetSelectedDeckTypeStats() ([]model.CardTypeStat, error) {
	return a.cardService.GetSelectedDeckTypeStats(context.Background())
}

func (a *App) GetSelectedDeckCardStats(theType model.Type, keyword string, upgradeFilter string) ([]model.DeckCardStat, error) {
	return a.cardService.GetSelectedDeckCardStats(context.Background(), theType, keyword, upgradeFilter)
}

func (a *App) GetDeckCompletionStatus() (model.DeckCompletion, error) {
	return a.cardService.GetDeckCompletionStatus(context.Background())
}

func (a *App) GetUserDeck() (model.Deck, error) {
	return a.cardService.GetUserDeck(context.Background())
}

func (a *App) ClearUserDeck() error {
	return a.cardService.ClearUserDeck(context.Background())
}

func (a *App) AddCardToUserDeck(cardID string) error {
	return a.cardService.AddCardToUserDeck(context.Background(), cardID)
}

func (a *App) RemoveCardFromUserDeck(cardID string) error {
	return a.cardService.RemoveCardFromUserDeck(context.Background(), cardID)
}

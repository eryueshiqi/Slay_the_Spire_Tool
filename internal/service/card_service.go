package service

import (
	"context"
	"slay_the_spire_tool/internal/model"
	"slay_the_spire_tool/internal/repository"
)

type CardService struct {
	cardRepo repository.CardInterface
}

func NewCardService(cardRepo repository.CardInterface) *CardService {
	return &CardService{cardRepo: cardRepo}
}

func (s *CardService) GetRoleList(ctx context.Context) ([]model.Role, error) {
	return s.cardRepo.GetRoleList(ctx)
}

func (s *CardService) GetCardList(ctx context.Context, str string, theType model.Type, role model.Role, cost model.Cost, rarity model.Rarity, isBase bool) ([]model.CardInfo, error) {
	return s.cardRepo.GetCardList(ctx, str, theType, role, cost, rarity, isBase)
}

func (s *CardService) GetDeckList(ctx context.Context, role model.Role) ([]model.Deck, error) {
	return s.cardRepo.GetDeckList(ctx, role)
}

func (s *CardService) GetDeck(ctx context.Context, id uint) (model.Deck, error) {
	return s.cardRepo.GetDeck(ctx, id)
}

func (s *CardService) GetCard(ctx context.Context, id string, role model.Role) (model.CardInfo, error) {
	return s.cardRepo.GetCard(ctx, id, role)
}

func (s *CardService) GetCardDetail(ctx context.Context, id string) (model.Card, error) {
	return s.cardRepo.GetCardDetail(ctx, id)
}

func (s *CardService) ChooseDeck(ctx context.Context, id uint) error {
	return s.cardRepo.ChooseDeck(ctx, id)
}

func (s *CardService) GetSelectedDeckTypeStats(ctx context.Context) ([]model.CardTypeStat, error) {
	return s.cardRepo.GetSelectedDeckTypeStats(ctx)
}

func (s *CardService) GetSelectedDeckCardStats(ctx context.Context, theType model.Type, keyword string, upgradeFilter string) ([]model.DeckCardStat, error) {
	return s.cardRepo.GetSelectedDeckCardStats(ctx, theType, keyword, upgradeFilter)
}

func (s *CardService) GetDeckCompletionStatus(ctx context.Context) (model.DeckCompletion, error) {
	return s.cardRepo.GetDeckCompletionStatus(ctx)
}

func (s *CardService) GetUserDeck(ctx context.Context) (model.Deck, error) {
	return s.cardRepo.GetUserDeck(ctx)
}

func (s *CardService) ClearUserDeck(ctx context.Context) error {
	return s.cardRepo.ClearUserDeck(ctx)
}

func (s *CardService) AddCardToUserDeck(ctx context.Context, cardID string) error {
	return s.cardRepo.AddCardToUserDeck(ctx, cardID)
}

func (s *CardService) RemoveCardFromUserDeck(ctx context.Context, cardID string) error {
	return s.cardRepo.RemoveCardFromUserDeck(ctx, cardID)
}

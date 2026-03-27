package repository

import (
	"context"
	"slay_the_spire_tool/internal/model"
)

type CardInterface interface {
	GetCardList(ctx context.Context, str string, theType model.Type, role model.Role, cost model.Cost, rarity model.Rarity, isBase bool) ([]model.CardInfo, error) // 获取卡牌列表
	GetDeckList(ctx context.Context, role model.Role) ([]model.Deck, error)                                                                                        // 获取卡组列表

	GetRoleList(ctx context.Context) ([]model.Role, error)                           // 获取角色列表
	GetCard(ctx context.Context, id string, role model.Role) (model.CardInfo, error) // 获取卡牌信息
	GetCardDetail(ctx context.Context, id string) (model.Card, error)                // 获取卡牌详细信息
	GetDeck(ctx context.Context, id uint) (model.Deck, error)                        // 获取卡组信息 函数实现中包括对比修改CardInfo, Need和HighLight影响前端展示
	GetSelectedDeckTypeStats(ctx context.Context) ([]model.CardTypeStat, error)      // 获取当前已选择牌组的类型统计
	GetSelectedDeckCardStats(ctx context.Context, theType model.Type, keyword string, upgradeFilter string) ([]model.DeckCardStat, error)
	GetDeckCompletionStatus(ctx context.Context) (model.DeckCompletion, error)
	ChooseDeck(ctx context.Context, id uint) error                   // 选择卡组
	GetUserDeck(ctx context.Context) (model.Deck, error)             // 获取用户卡组信息 内容为用户选择的卡组 函数实现中包含排序和清除无用的卡牌
	ClearUserDeck(ctx context.Context) error                         // 清空用户卡组
	AddCardToUserDeck(ctx context.Context, cardID string) error      // 添加卡牌到用户卡组
	RemoveCardFromUserDeck(ctx context.Context, cardID string) error // 从用户卡组中移除卡牌
}

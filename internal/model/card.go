package model

type Cost uint8

const (
	Zero Cost = iota
	One
	Two
	Three
	X
	AnyCost Cost = 255
)

type Type uint8

const (
	Attack Type = iota
	Skill
	Power
	Quest
	Status
	Curse
	AnyType Type = 255
)

type Rarity uint8

const (
	Basic Rarity = iota
	Common
	Uncommon
	Rare
	Ancient
	AnyRarity Rarity = 255
)

type Role uint8

const (
	Ironclad Role = iota
	Silent
	Regent
	Necrobinder
	Defect
	Colorless
	AnyRole Role = 255
)

type Card struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Alias  string `json:"alias"`
	Image  string `json:"image"`
	Type   Type   `json:"type"`
	Role   Role   `json:"role"`
	Cost   Cost   `json:"cost"`
	Rarity Rarity `json:"rarity"`
	IsBase bool   `json:"is_base"`
}

type RoleCards struct {
	Role  Role   `json:"role"`
	Cards []Card `json:"cards"`
}

type CardInfo struct {
	ID         string `json:"id"`
	Image      string `json:"image"`
	Need       bool   `json:"need"`
	HighLight  bool   `json:"highlight"`
	NeedDelete bool   `json:"need_delete,omitempty"`
}

type Deck struct {
	Id    uint       `json:"id"`
	Name  string     `json:"name"`
	Role  Role       `json:"role"`
	Cards []CardInfo `json:"cards"`
}

type CardTypeStat struct {
	Type  Type `json:"type"`
	Count int  `json:"count"`
}

type DeckCardStat struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	Type   Type   `json:"type"`
	IsBase bool   `json:"is_base"`
	Count  int    `json:"count"`
}

type EnumLabel struct {
	Value int    `json:"value"`
	Key   string `json:"key"`
	Label string `json:"label"`
}

type DisplayLabels struct {
	Roles    []EnumLabel `json:"roles"`
	Types    []EnumLabel `json:"types"`
	Costs    []EnumLabel `json:"costs"`
	Rarities []EnumLabel `json:"rarities"`
}

type DeckCompletion struct {
	Completed    bool `json:"completed"`
	MissingCount int  `json:"missing_count"`
}

func BuildDisplayLabelsZH() DisplayLabels {
	return DisplayLabels{
		Roles: []EnumLabel{
			{Value: int(Ironclad), Key: "Ironclad", Label: "铁甲战士"},
			{Value: int(Silent), Key: "Silent", Label: "静默猎手"},
			{Value: int(Regent), Key: "Regent", Label: "摄政"},
			{Value: int(Necrobinder), Key: "Necrobinder", Label: "缚灵师"},
			{Value: int(Defect), Key: "Defect", Label: "故障机器人"},
			{Value: int(Colorless), Key: "Colorless", Label: "无色"},
			{Value: int(AnyRole), Key: "AnyRole", Label: "全部角色"},
		},
		Types: []EnumLabel{
			{Value: int(Attack), Key: "Attack", Label: "攻击"},
			{Value: int(Skill), Key: "Skill", Label: "技能"},
			{Value: int(Power), Key: "Power", Label: "能力"},
			{Value: int(Quest), Key: "Quest", Label: "任务"},
			{Value: int(Status), Key: "Status", Label: "状态"},
			{Value: int(Curse), Key: "Curse", Label: "诅咒"},
			{Value: int(AnyType), Key: "AnyType", Label: "全部类型"},
		},
		Costs: []EnumLabel{
			{Value: int(Zero), Key: "Zero", Label: "0费"},
			{Value: int(One), Key: "One", Label: "1费"},
			{Value: int(Two), Key: "Two", Label: "2费"},
			{Value: int(Three), Key: "Three", Label: "3费"},
			{Value: int(X), Key: "X", Label: "X费"},
			{Value: int(AnyCost), Key: "AnyCost", Label: "全部费用"},
		},
		Rarities: []EnumLabel{
			{Value: int(Basic), Key: "Basic", Label: "基础"},
			{Value: int(Common), Key: "Common", Label: "普通"},
			{Value: int(Uncommon), Key: "Uncommon", Label: "罕见"},
			{Value: int(Rare), Key: "Rare", Label: "稀有"},
			{Value: int(Ancient), Key: "Ancient", Label: "远古"},
			{Value: int(AnyRarity), Key: "AnyRarity", Label: "全部稀有度"},
		},
	}
}

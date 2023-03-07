package main

type WagerRecord struct {
	StartTimestamp        uint64  `json:"startTimestamp"`        //
	EndTimestamp          uint64  `json:"endTimestamp"`          //
	PlayerID              string  `json:"playerId"`              //
	Nickname              string  `json:"nickname"`              //
	EcName                string  `json:"ecName"`                //
	RawEcName             string  `json:"rawEcName"`             //
	RecordType            string  `json:"recordType"`            // 類型 Deposit: 上分, Draw: 下分, Game: 遊戲, Club: 俱樂部返水領取紀錄, CompetitionOrder:比賽場入場費, CompetitionReward 比賽場獎金
	GameName              string  `json:"gameName"`              // 簡中
	GameRecordID          string  `json:"gameRecordId"`          // GameRecords._key
	SessionRecordID       string  `json:"sessionRecordId"`       //
	Seq                   int64   `json:"seq"`                   // deprecated: Session Record Seq
	KKBeforBalance        float64 `json:"kkBeforBalance"`        // 初始金額
	BeforeBetBalance      float64 `json:"beforeBetBalance"`      // 攜入金額
	KKDepositAmount       float64 `json:"kkDepositAmount"`       //
	KKDrawAmount          float64 `json:"kkDrawAmount"`          //
	DepositAmount         float64 `json:"depositAmount"`         //
	DrawAmount            float64 `json:"drawAmount"`            //
	BetAmount             float64 `json:"betAmount"`             //
	BetBalance            float64 `json:"betBalance"`            //
	ValidBetAmount        float64 `json:"validBetAmount"`        //
	MemberJuiceAmounts    float64 `json:"memberJuiceAmount"`     // 返水金額
	MemberHostJuiceAmount float64 `json:"memberHostJuiceAmount"` // 玩家成為房主返水金額
	Profit                float64 `json:"profit"`                // 玩家營利
	PayoutAmount          float64 `json:"payoutAmount"`          //
	Balance               float64 `json:"balance"`               //
	EcUserID              string  `json:"ecUserId"`              //
	EcSiteID              string  `json:"ecSiteId"`              //
	GameID                string  `json:"gameId"`                //
	GameType              uint8   `json:"gameType"`              //
	ScoreType             uint8   `json:"scoreType"`             //
}

type GameRecord struct {
	Key                    string         `json:"_key,omitempty"`         // GR-{UnixTime}-{RedisIncr}-{InstanceIncr}
	EcSiteId               string         `json:"ecSiteId"`               // 平台識別碼
	GameType               int            `json:"gameType"`               // 遊戲編號
	MergedGameRecordId     string         `json:"mergedGameRecordId"`     // 合併注單關聯Id
	GameId                 string         `json:"gameId"`                 // 遊戲Id
	ThemeId                string         `json:"themeId"`                // 遊戲大廳Id
	RoomId                 string         `json:"roomId"`                 // 遊戲房間Id
	MemberCount            int            `json:"memberCount"`            // 真實玩家數量
	MemberData             MemberDataList `json:"memberData"`             // 真實玩家資料
	MemberIncomes          int64          `json:"memberIncomes"`          // 真實玩家總獲利金額
	MemberOutcomes         int64          `json:"memberOutcomes"`         // 真實玩家總虧損金額
	MemberJuiceAmounts     int64          `json:"memberJuiceAmounts"`     // 真實玩家被系統抽水總金額
	MemberHostJuiceAmounts int64          `json:"memberHostJuiceAmounts"` // 真實玩家成為房主時的返水金額
	BotCount               int            `json:"botCount"`               // 機器人數量
	BotIncome              int64          `json:"botIncome"`              // 機器人總獲利金額
	BotOutcome             int64          `json:"botOutcome"`             // 機器人總虧損金額
	BotJuiceAmount         int64          `json:"botJuiceAmount"`         // 機器人被系統抽水總金額
	JuiceRatio             float64        `json:"juiceRatio"`             // 系統抽水百分比
	History                interface{}    `json:"history"`                // 遊戲歷程
	StartTimestamp         uint64         `json:"startTimestamp"`         // 遊戲開始時間
	EndTimestamp           uint64         `json:"endTimestamp"`           // 遊戲結束時間
	ScoreType              uint8          `json:"scoreType"`              // 比賽房分數類型 0: 積分類型, 1: 現金類型
	Signature              []byte         `json:"signature"`              // 數位簽章
}

type MemberDataList []*MemberData

func (list MemberDataList) Len() int {
	return len(list)
}

func (list MemberDataList) Less(i, j int) bool {
	return list[j].SessionId < list[i].SessionId
}

func (list MemberDataList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

type MemberData struct {
	EcUserId              string `json:"ecUserId"`              // 真實玩家 EcSite UserId
	PlayerId              string `json:"playerId"`              // 真實玩家 UUId
	SessionId             string `json:"sessionId"`             // 真實玩家本次連線的 SessionId
	WagerId               string `json:"wagerId"`               // 真實玩家注單關聯Id ex: gameRecordId + playerId
	Nickname              string `json:"nickname"`              // 真實玩家暱稱
	Icon                  int    `json:"icon"`                  // 真實玩家大頭貼
	IconFrame             int    `json:"iconFrame"`             // 頭像匡編號列表
	KKBeforBalance        int64  `json:"kkBeforBalance"`        // 玩家當前KK錢包總餘額(不含輸贏結果)
	KKAfterBalance        int64  `json:"kkAfterBalance"`        // 玩家當前KK錢包總餘額(含輸贏結果)
	BeforeBetBalance      int64  `json:"beforeBetBalance"`      // 真實玩家下注前的身上餘額
	BetBalance            int64  `json:"betBalance"`            // 真實玩家下注後的身上餘額
	BetAmount             int64  `json:"betAmount"`             // 真實玩家投注金額
	ValidBetAmount        int64  `json:"validBetAmount"`        // 真實玩家有效投注金額
	MemberIncome          int64  `json:"memberIncome"`          // 真實玩家獲利金額(已經扣除系統抽水金額)
	MemberOutcome         int64  `json:"memberOutcome"`         // 真實玩家虧損金額
	MemberJuiceAmount     int64  `json:"memberJuiceAmount"`     // 真實玩家被系統抽水金額
	MemberHostJuiceAmount int64  `json:"memberHostJuiceAmount"` // 真實玩家成為房主時的返水金額
	Profit                int64  `json:"profit"`                // memberIncome - memberOutcome
	PayoutAmount          int64  `json:"payoutAmount"`          // betAmount + profit
	Score                 int64  `json:"score"`                 // 無視以小博大的總積分
	TagType               int    `json:"tagType"`               // 標籤類型
}

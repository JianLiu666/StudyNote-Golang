package main

var startTimestamp int64 = 1690300800000
var dayAfterTimestamp int64 = 1690387199999
var weekAfterTimestamp int64 = 1690905599999

type job struct {
	Name  string
	AQL   string
	Binds map[string]interface{}
}

var aqlMemberRecordStatistics string = `
	FOR wr IN WagerRecords 
	FILTER wr.endTimestamp >= @startTimestamp && wr.endTimestamp <= @endTimestamp
	FILTER wr.ecSiteId IN ["1", "21729934", "21729963", "21729991", "21730024", "21730063", "21730086", "21730112", "21730139", "21730184", "21730208", "21730232", "3319009", "5044073", "63638"]
	FILTER wr.recordType IN ["Game", "Deposit", "Draw", "Club", "CompetitionOrder", "CompetitionReward"]

	COLLECT AGGREGATE
		totalCount = COUNT(wr),
		totalBetAmount = SUM(wr.betAmount),
		totalValidBetAmount =  SUM(wr.validBetAmount),
		totalPlayerProfit = SUM(wr.profit),
		totalHostJuice = SUM(wr.memberHostJuiceAmount)
		
	RETURN {
		totalCount: totalCount,
		totalBetAmount: totalBetAmount,
		totalValidBetAmount: totalValidBetAmount,
		totalPlayerProfit: totalPlayerProfit,
		totalHostJuice: totalHostJuice
	}`

var aqlMemberRecordFirstPage string = `
	FOR wr IN WagerRecords
	FILTER wr.endTimestamp >= @startTimestamp && wr.endTimestamp <= @endTimestamp
	FILTER wr.ecSiteId IN ["1", "21729934", "21729963", "21729991", "21730024", "21730063", "21730086", "21730112", "21730139", "21730184", "21730208", "21730232", "3319009", "5044073", "63638"]
	FILTER wr.recordType IN ["Game", "Deposit", "Draw", "Club", "CompetitionOrder", "CompetitionReward"]

	LIMIT 0, 200

	RETURN {
		"ecSiteId":         wr.ecSiteId,
		"gameId":           wr.gameId,
		"createdAt":        wr.endTimestamp,
		"username":         wr.ecUserId,
		"srIdSeq":          wr.sessionRecordId,
		"playerId":         wr.playerId,
		"themeId":          wr.themeId,
		"initAmount":       wr.beforeBetBalance,
		"betAmount":        wr.betAmount,
		"validateAmount":   wr.validBetAmount,
		"payoutAmount":     wr.payoutAmount,
		"fee":              wr.memberJuiceAmount,
		"profit":           -1*wr.profit,
		"finalAmount":      wr.betBalance,
		"wrId":             wr._key,
		"gameType":         TO_STRING(wr.gameType),
		"hostJuiceAmount":  wr.memberHostJuiceAmount,
		"kkBeforBalance":   wr.kkBeforBalance,
		"kkAfterBalance":   wr.kkBeforBalance,
		"createTimestamp":  TO_STRING(wr.endTimestamp),
	}`

var jobs []*job = []*job{
	{
		Name: "DailyMemberRecordStatistics",
		AQL:  aqlMemberRecordStatistics,
		Binds: map[string]interface{}{
			"startTimestamp": startTimestamp,
			"endTimestamp":   dayAfterTimestamp,
		},
	},
	{
		Name: "WeeklyMemberRecordStatistics",
		AQL:  aqlMemberRecordStatistics,
		Binds: map[string]interface{}{
			"startTimestamp": startTimestamp,
			"endTimestamp":   weekAfterTimestamp,
		},
	},
	{
		Name: "DailyMemberRecordFirstPage",
		AQL:  aqlMemberRecordFirstPage,
		Binds: map[string]interface{}{
			"startTimestamp": startTimestamp,
			"endTimestamp":   dayAfterTimestamp,
		},
	},
	{
		Name: "WeeklyMemberRecordFirstPage",
		AQL:  aqlMemberRecordFirstPage,
		Binds: map[string]interface{}{
			"startTimestamp": startTimestamp,
			"endTimestamp":   weekAfterTimestamp,
		},
	},
}

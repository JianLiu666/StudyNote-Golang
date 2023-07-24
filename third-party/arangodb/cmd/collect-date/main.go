package main

import (
	"arangodb/src"
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	ctx := context.Background()

	db := src.GetDBInstance(ctx, "Database")

	start := time.Date(2023, 4, 1, 0, 0, 0, 0, time.Local)
	end := time.Date(2023, 5, 1, 0, 0, 0, 0, time.Local)

	month_records := map[string]int{}
	mau := map[string]map[string]bool{}

	aql1 := `
	RETURN LENGTH(
		FOR doc IN @@col 
			FILTER doc.memberCount > 0
			RETURN 1
	)`
	aql2 := `
	LET arr = (
		FOR doc IN @@col
		LET ids = (
			FOR member IN doc.memberData
			RETURN member.playerId
		)
		RETURN ids
	)
	
	RETURN UNIQUE(FLATTEN(arr,2))`

	for ; !start.Equal(end); start = start.AddDate(0, 0, 1) {
		// 統計注單量
		cursor, err := db.Query(ctx, aql1,
			map[string]interface{}{
				"@col": fmt.Sprintf("GameRecords_%v", start.Format("20060102")),
			})
		if err != nil {
			log.Panic(err)
		}

		var num int
		cursor.ReadDocument(ctx, &num)
		cursor.Close()

		// 累計注單量到相同月份
		month_records[start.Format("200601")] += num

		// 統計 DAU
		cursor, err = db.Query(ctx, aql2,
			map[string]interface{}{
				"@col": fmt.Sprintf("GameRecords_%v", start.Format("20060102")),
			})
		if err != nil {
			log.Panic(err)
		}

		var names []string
		cursor.ReadDocument(ctx, &names)
		cursor.Close()

		// lazy init
		if mau[start.Format("200601")] == nil {
			mau[start.Format("200601")] = make(map[string]bool)
		}

		// 過濾重複玩家
		for _, name := range names {
			mau[start.Format("200601")][name] = true
		}
	}

	// 輸出每月分的注單數量與活躍玩家
	for key, sub := range mau {
		log.Printf("%v records:%v, mau:%v", key, month_records[key], len(sub))
	}
}

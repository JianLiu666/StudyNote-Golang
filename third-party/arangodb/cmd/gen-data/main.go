package main

import (
	"arangodb/src"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"
)

func main() {
	ctx := context.Background()

	db := src.GetDBInstance(ctx, "Test")

	grCol, _ := db.Collection(ctx, "GameRecords")
	wrCol, _ := db.Collection(ctx, "WagerRecords")

	t := time.Now()
	serviceId := 23
	incr := 1

	for i := 1; i <= 2500; i++ {
		var grs []*GameRecord
		var wrs []*WagerRecord

		for j := 1; j <= 40000; j++ {
			var gr GameRecord
			err := json.Unmarshal([]byte(baseGameRecord), &gr)
			if err != nil {
				panic(err)
			}

			gr.StartTimestamp = uint64(t.UnixMilli())
			gr.EndTimestamp = uint64(t.UnixMilli())
			uuid := fmt.Sprintf("%v-%v-%v", gr.EndTimestamp, serviceId, incr)
			gr.Key = fmt.Sprintf("GR-%v", uuid)
			gr.MemberData[0].WagerId = fmt.Sprintf("WGR-%v-%v", uuid, gr.MemberData[0].PlayerId)

			grs = append(grs, &gr)

			var wr WagerRecord
			err = json.Unmarshal([]byte(baseWagerRecord), &wr)
			if err != nil {
				panic(err)
			}

			wr.StartTimestamp = gr.StartTimestamp
			wr.EndTimestamp = gr.EndTimestamp
			wr.GameRecordID = gr.Key
			wr.SessionRecordID = gr.MemberData[0].WagerId

			wrs = append(wrs, &wr)

			t = t.Add(100 * time.Millisecond)
			incr++
			if incr > math.MaxUint16 {
				incr = 1
			}
		}

		_, _, err := grCol.CreateDocuments(ctx, grs)
		if err != nil {
			fmt.Println(err)
		}

		_, _, err = wrCol.CreateDocuments(ctx, wrs)
		if err != nil {
			fmt.Println(err)
		}

		log.Println(i)
	}
}

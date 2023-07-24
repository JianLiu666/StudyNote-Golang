package main

import (
	"arangodb/src"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"
)

var dbName string = "Test"

var batchSize int = 1000
var iteration int = 100000 // batchSize x iteration = total documents
// var batchSize int = 1
// var iteration int = 1 // batchSize x iteration = total documents

func main() {
	ctx := context.Background()

	db := src.GetDBInstance(ctx, dbName)

	grCol, err := db.Collection(ctx, "GameRecords")
	if err != nil {
		panic(err)
	}
	wrCol, err := db.Collection(ctx, "WagerRecords")
	if err != nil {
		panic(err)
	}

	t := time.Now()
	serviceId := 23
	incr := 1

	elasped := time.Now()
	for i := 1; i <= iteration; i++ {
		var grs []*GameRecord
		var wrs []*WagerRecord

		_elasped := time.Now()
		for batch := 1; batch <= batchSize; batch++ {
			//prepare game record
			var gr GameRecord
			err := json.Unmarshal([]byte(baseGameRecord), &gr)
			if err != nil {
				panic(err)
			}

			gr.StartTimestamp = uint64(t.UnixMilli())
			gr.EndTimestamp = uint64(t.UnixMilli())
			uuid := fmt.Sprintf("%v-%v-%v", gr.EndTimestamp, serviceId, incr)
			gr.RecordId = fmt.Sprintf("GR-%v", uuid)
			gr.MemberData[0].WagerId = fmt.Sprintf("WGR-%v-%v", uuid, gr.MemberData[0].PlayerId)
			gr.ShardKey = int(t.Day())

			grs = append(grs, &gr)

			//prepare wager record
			var wr WagerRecord
			err = json.Unmarshal([]byte(baseWagerRecord), &wr)
			if err != nil {
				panic(err)
			}

			wr.StartTimestamp = gr.StartTimestamp
			wr.EndTimestamp = gr.EndTimestamp
			wr.GameRecordID = gr.RecordId
			wr.SessionRecordID = gr.MemberData[0].WagerId
			wr.ShardKey = int(t.Day())

			wrs = append(wrs, &wr)

			t = t.Add(33 * time.Millisecond)
			incr++
			if incr > math.MaxUint16 {
				incr = 1
			}
		}

		_, _, err := grCol.CreateDocuments(ctx, grs)
		if err != nil {
			panic(err)
		}

		_, _, err = wrCol.CreateDocuments(ctx, wrs)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%v: %v\n", i, time.Now().Sub(_elasped).String())
	}

	fmt.Printf("total time: %v\n", time.Now().Sub(elasped).String())
}

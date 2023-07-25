package cmd

import (
	"encoding/json"
	"fmt"
	"natspractice/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var benchNatsMultiTopicsConsumerCmd = &cobra.Command{
	Use:   "bench_nats_multi_topics_sub",
	Short: "",
	Long:  ``,
	RunE:  RunBenchNatsMultiTopicsConsumerCmd,
}

func init() {
	rootCmd.AddCommand(benchNatsMultiTopicsConsumerCmd)
}

func RunBenchNatsMultiTopicsConsumerCmd(cmd *cobra.Command, args []string) error {
	type payload struct {
		Timestamp  int64       `json:"timestamp"`
		AnteType   int         `json:"anteType"`
		FishIds    []string    `json:"fishIds"`
		Payload    string      `json:"payload"`
		Positions  [][]float64 `json:"positions"`
		RoundToken string      `json:"roundToken"`
		SkillId    string      `json:"skillId"`
		SkillType  int         `json:"skillType"`
	}

	// 儲存資料用
	xaxis := []int{}
	items := []opts.LineData{}
	buffer := make(chan int64, 100000)

	// 開一條 goroutine 在背景匯聚資料
	go func(ch chan int64) {
		var xidx int = 1
		var batch, sum int64 = 0, 0
		for data := range buffer {
			sum += data
			batch++
			if batch == 1000 {
				xaxis = append(xaxis, xidx)
				items = append(items, opts.LineData{Value: sum / batch})
				batch = 0
				sum = 0
				xidx++
			}
		}
	}(buffer)

	// 對 NATS streaming server 建立連線
	nc, err := nats.Connect(
		config.Nats.Addr,
		nats.UserInfo(config.Nats.Username, config.Nats.Password),
	)
	if err != nil {
		return err
	}
	defer nc.Close()

	// 訂閱指定主題
	for i := 1; i <= config.Nats.BenchNumTopics; i++ {
		_, err = nc.Subscribe(fmt.Sprintf("Test%v", i), func(msg *nats.Msg) {
			var data payload
			if err := json.Unmarshal(msg.Data, &data); err != nil {
				panic(err)
			}
			eplased := time.Now().Sub(time.UnixMilli(int64(data.Timestamp))).Milliseconds()
			buffer <- eplased
		})
		if err != nil {
			return err
		}
	}

	// set graceful shutdown method
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	<-stopSignal

	// 繪製圖表
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title: "NATS Server Benchmark: Mutiple Topics",
			Subtitle: fmt.Sprintf("Number of topics: %v, Number of producers: %v, Time of each producer: %v, Sleep of each producer: %vms",
				config.Nats.BenchNumTopics,
				config.Nats.BenchNumProducers,
				config.Nats.BenchProducerEachTimes,
				config.Nats.BenchProducerSleepTime,
			),
		}))

	line.SetXAxis(xaxis).
		AddSeries("flight time (ms)", items).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	f, _ := os.Create(fmt.Sprintf("bench_nats_multi_topics_%vx%vx%v.html",
		config.Nats.BenchNumTopics,
		config.Nats.BenchNumProducers,
		config.Nats.BenchProducerEachTimes,
	))
	line.Render(f)

	return nil
}

package cmd

import (
	"fmt"
	"natspractice/config"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var benchJetStreamConsumerCmd = &cobra.Command{
	Use:   "bench_js_sub",
	Short: "",
	Long:  ``,
	RunE:  RunBenchJetStreamConsumerCmd,
}

func init() {
	rootCmd.AddCommand(benchJetStreamConsumerCmd)
}

func RunBenchJetStreamConsumerCmd(cmd *cobra.Command, args []string) error {
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

	// 對 NATS server 建立連線
	nc, err := nats.Connect(config.Nats.Addr)
	if err != nil {
		return err
	}
	defer nc.Close()

	// 建立 JetStream context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return err
	}

	// Check if the stream already exists
	_, err = js.StreamInfo("Test")
	if err != nil {
		if err == nats.ErrStreamNotFound {
			// If the stream does not exist, return an error
			return fmt.Errorf("Stream 'Test' does not exist")
		} else {
			// If there was another error, return it
			return err
		}
	}

	// 訂閱指定主題，從最新的訊息開始消費
	sub, err := js.Subscribe("Test", func(msg *nats.Msg) {
		data, _ := strconv.Atoi(string(msg.Data))
		eplased := time.Now().Sub(time.UnixMilli(int64(data))).Milliseconds()
		buffer <- eplased
	},
		nats.DeliverLast(),
		nats.ReplayInstant(),
		nats.AckNone(),
		nats.ConsumerMemoryStorage(),
	)
	if err != nil {
		return err
	}

	// set graceful shutdown method
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, os.Interrupt, syscall.SIGTERM)
	<-stopSignal

	// Unsubscribe
	if err := sub.Unsubscribe(); err != nil {
		return err
	}

	// 繪製圖表
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title: "NATS JetStream Server Benchmark",
			Subtitle: fmt.Sprintf("Number of producers: %v Time of each producer: %v",
				config.Nats.BenchNumProducers,
				config.Nats.BenchProducerEachTimes,
			),
		}))

	line.SetXAxis(xaxis).
		AddSeries("flight time (ms)", items).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	f, _ := os.Create(fmt.Sprintf("bench_js_%vx%v.html",
		config.Nats.BenchNumProducers,
		config.Nats.BenchProducerEachTimes,
	))
	line.Render(f)

	return nil
}

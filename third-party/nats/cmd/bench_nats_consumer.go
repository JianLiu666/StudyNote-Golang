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

var benchNatsConsumerCmd = &cobra.Command{
	Use:   "bench_nats_sub",
	Short: "",
	Long:  ``,
	RunE:  RunBenchNatsConsumerCmd,
}

func init() {
	rootCmd.AddCommand(benchNatsConsumerCmd)
}

func RunBenchNatsConsumerCmd(cmd *cobra.Command, args []string) error {
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
	)
	if err != nil {
		return err
	}
	defer nc.Close()

	// 訂閱指定主題
	_, err = nc.Subscribe("Test", func(msg *nats.Msg) {
		data, _ := strconv.Atoi(string(msg.Data))
		eplased := time.Now().Sub(time.UnixMilli(int64(data))).Milliseconds()
		buffer <- eplased
	})
	if err != nil {
		return err
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
			Title: "NATS Server Benchmark",
			Subtitle: fmt.Sprintf("Number of producers: %v Time of each producer: %v",
				config.Nats.BenchNumProducers,
				config.Nats.BenchProducerEachTimes,
			),
		}))

	line.SetXAxis(xaxis).
		AddSeries("flight time (ms)", items).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	f, _ := os.Create(fmt.Sprintf("bench_nats_%vx%v.html",
		config.Nats.BenchNumProducers,
		config.Nats.BenchProducerEachTimes,
	))
	line.Render(f)

	return nil
}

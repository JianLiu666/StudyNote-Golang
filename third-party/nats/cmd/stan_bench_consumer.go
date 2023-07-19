package cmd

import (
	"fmt"
	"natspractice/config"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/nats-io/stan.go"
	"github.com/spf13/cobra"
)

var stanConsumerCmd = &cobra.Command{
	Use:   "stan_sub",
	Short: "",
	Long:  ``,
	RunE:  RunStanConsumerCmd,
}

func init() {
	rootCmd.AddCommand(stanConsumerCmd)
}

func RunStanConsumerCmd(cmd *cobra.Command, args []string) error {
	// 儲存資料用
	xaxis := []int{}
	items := []opts.LineData{}
	buffer := make(chan int64, 100000)

	// 開一條 goroutine 在背景匯聚資料
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(ch chan int64, _wg *sync.WaitGroup) {
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
		wg.Done()
	}(buffer, &wg)

	// 對 NATS streaming server 建立連線
	sc, err := stan.Connect(
		config.Nats.ClusterId,
		fmt.Sprintf("stan-%v", time.Now().UnixNano()),
		stan.NatsURL(config.Nats.Addr),
	)
	if err != nil {
		return err
	}
	defer sc.Close()

	// 訂閱指定主題
	_, err = sc.Subscribe("Test", func(msg *stan.Msg) {
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

	close(buffer)
	wg.Wait()

	// 繪製圖表
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "Line example in Westeros theme",
			Subtitle: "Line chart rendered by the http server this time",
		}))

	line.SetXAxis(xaxis).
		AddSeries("Category A", items).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))

	f, _ := os.Create("line.html")
	line.Render(f)

	return nil
}

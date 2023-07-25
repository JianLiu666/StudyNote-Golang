package cmd

import (
	"fmt"
	"natspractice/config"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

var benchNatsMultiTopicsPublisherCmd = &cobra.Command{
	Use:   "bench_nats_multi_topics_pub",
	Short: "",
	Long:  ``,
	RunE:  RunBenchNatsMultiTopicsPublisherCmd,
}

func init() {
	rootCmd.AddCommand(benchNatsMultiTopicsPublisherCmd)
}

func RunBenchNatsMultiTopicsPublisherCmd(cmd *cobra.Command, args []string) error {
	nc, err := nats.Connect(
		config.Nats.Addr,
		nats.UserInfo(config.Nats.Username, config.Nats.Password),
	)
	if err != nil {
		return err
	}
	defer nc.Close()

	var wg sync.WaitGroup

	// 併發測試 NATS streaming publish 效能
	for i := 1; i <= config.Nats.BenchNumProducers; i++ {
		wg.Add(1)
		go func(_wg *sync.WaitGroup, idx int) {
			defer wg.Done()
			for i := 0; i < config.Nats.BenchProducerEachTimes; i++ {
				eplased := time.Now()
				err = nc.Publish(fmt.Sprintf("Test%v", idx), []byte(fmt.Sprintf("%v", time.Now().UnixMilli())))
				fmt.Println(i, time.Now().Sub(eplased))

				if err != nil {
					fmt.Println(err)
				}

				time.Sleep(time.Duration(config.Nats.BenchProducerSleepTime) * time.Millisecond)
			}
		}(&wg, i)
	}

	wg.Wait()

	return nil
}

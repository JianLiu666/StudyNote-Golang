package cmd

import (
	"fmt"
	"jian6/nats/config"
	"math/rand"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/spf13/cobra"
)

var durableQueueGroupCmd = &cobra.Command{
	Use:   "durable_queue",
	Short: "Run NATS streaming durable queue group test case.",
	Long:  `This case is aim to understand how max in flight will present in durable queue group subscription.`,
	RunE:  RunDurableQueueGroupCmd,
}

func init() {
	rootCmd.AddCommand(durableQueueGroupCmd)
}

func RunDurableQueueGroupCmd(cmd *cobra.Command, args []string) error {
	nc, err := nats.Connect(config.Nats.Addr)
	if err != nil {
		panic(err)
	}

	sc, err := stan.Connect(
		config.Nats.ClusterId,
		"Durable_Queue_Group_Client",
		stan.NatsConn(nc),
		// stan.MaxPubAcksInflight(1),
	)
	if err != nil {
		panic(err)
	}

	subjectName := "aasabaa"
	wg := sync.WaitGroup{}

	num := 1
	wg.Add(num)
	for i := 1; i <= num; i++ {
		go enableWorker(sc, subjectName, i, &wg)
	}
	wg.Wait()

	num = 100
	wg.Add(num)
	go func() {
		for i := 1; i <= num; i++ {
			if err := sc.Publish(subjectName, []byte(fmt.Sprintf("%v", i))); err != nil {
				fmt.Println("publish error:", err)
			}
		}
	}()
	wg.Wait()

	return nil
}

func enableWorker(sc stan.Conn, subjectName string, instanceId int, wg *sync.WaitGroup) {
	defer wg.Done()
	defer fmt.Printf("instance %d subscribed.\n", instanceId)

	callback := func(msg *stan.Msg) {
		time.Sleep(time.Duration(rand.Intn(2)) * time.Second)
		if err := msg.Ack(); err != nil {
			fmt.Printf("receiver error:%v\n", err)
		}
		fmt.Printf("[%d] %05d|Data:%v\n", instanceId, msg.Sequence, string(msg.Data))
		wg.Done()
	}

	_, err := sc.QueueSubscribe(
		subjectName,
		"durable_queue_group",
		callback,
		stan.SetManualAckMode(),
		stan.AckWait(10*time.Second),
		stan.MaxInflight(1),
	)
	if err != nil {
		panic(err)
	}
}

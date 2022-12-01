package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/cobra"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "mqttpubbot",
	Short: "Publishes random mqtt messages on specified queues",
	Long: `Publishes random mqtt on provided queues
at specified interval`,
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var ip string
var topicsArg string
var interval int
var port int

func init() {
	rootCmd.Flags().StringVarP(&ip, "ip", "i", "127.0.0.1", "ip address")
	rootCmd.Flags().IntVarP(&port, "port", "p", 1883, "mqtt port")
	rootCmd.Flags().StringVarP(&topicsArg, "topics", "t", "house/bulb1", "topics list seperated with a ,")
	rootCmd.Flags().IntVar(&interval, "interval", 3, "interval in seconds between ")
	rootCmd.MarkFlagRequired("topics")
}

func main() {
	rootCmd.Execute()
}

func start() {

	rand.Seed(time.Now().UnixNano())
	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt, syscall.SIGTERM)

	topics := strings.Split(topicsArg, ",")
	broker := fmt.Sprintf("tcp://%s:%d", ip, port)

	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("mqttpubbot")

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Cannot connect to ", broker)
	}
	go func() {
		<-ctrlC
		client.Disconnect(250)
		fmt.Println("Done publishing.")
	}()
	fmt.Println("Start publishing...")
	for i := 1; i > 0; i++ {
		//msg := fmt.Sprintf("Message %s", time.Now().Format("15:04:05"))
		randomNb := (1 + rand.Intn(30)) * 273
		msg := fmt.Sprintf("%d", randomNb)
		n := rand.Intn(len(topics))
		fmt.Println(n, "post", msg, "to", topics[n])
		if tok := client.Publish(topics[n], byte(0), true, msg); tok.Wait() && tok.Error() != nil {
			log.Fatalf("Cannot publish", msg, tok.Error())
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}

}

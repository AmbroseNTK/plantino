package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/ambrosentk/plantino/agent/domain"
	sensorRepo "github.com/ambrosentk/plantino/agent/internal/sensor/repo"
	sensorUcase "github.com/ambrosentk/plantino/agent/internal/sensor/ucase"
	"github.com/eclipse/paho.golang/autopaho"
	paho "github.com/eclipse/paho.golang/paho"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// read agent.yaml file with viper

	viper.SetConfigName("agent")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	// read the config file
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println(viper.GetString("serialNumber"))

	// init mqtt client
	u, err := url.Parse("mqtt://broker.emqx.io:1883")
	if err != nil {
		panic(err)
	}

	// connect postgres
	db, err := gorm.Open(postgres.Open(viper.GetString("db.dsn")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sensorDataRepo := sensorRepo.NewSensorDataRepository(db)
	sensorDataUseCase := sensorUcase.NewSensorDataUseCase(sensorDataRepo)

	wateringEnabled := false
	fanEnabled := false

	fanRelayAddress := viper.GetString("fan.relay.address")
	wateringRelayAddress := viper.GetString("watering.relay.address")

	// received data
	currentData := &domain.SensorData{
		Temperature:     0,
		Humidity:        0,
	}

	cliCfg := autopaho.ClientConfig{
		ServerUrls: []*url.URL{u},
		KeepAlive:  20, 
		CleanStartOnInitialConnection: false,
		SessionExpiryInterval: 60,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			fmt.Println("mqtt connection up")
			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{
					{Topic: fmt.Sprintf("plantino/%v/data", viper.GetString("serialNumber")), QoS: 0},
				},
			}); err != nil {
				fmt.Printf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
			}
			fmt.Println("mqtt subscription made")
		},
		OnConnectError: func(err error) { fmt.Printf("error whilst attempting connection: %s\n", err) },
		ClientConfig: paho.ClientConfig{
			ClientID: "plantino",
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				func(pr paho.PublishReceived) (bool, error) {
					data, err := domain.NewSensorDataFromRawData(pr.Packet.Payload)
					if err != nil {
						fmt.Printf("error parsing sensor data: %s\n", err)
						return true, nil
					}
					log.Printf("received sensor data: %+v\n", data)
					currentData = data
					err = sensorDataUseCase.Create(context.Background(), data)
					if err != nil {
						fmt.Printf("error saving sensor data: %s\n", err)
					}
					return true, nil
				}},
			OnClientError: func(err error) { fmt.Printf("client error: %s\n", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					fmt.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					fmt.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}

	ctx := context.Background()

	c, err := autopaho.NewConnection(ctx, cliCfg) // starts process; will reconnect until context cancelled
	if err != nil {
		panic(err)
	}
	// Wait for the connection to come up
	if err = c.AwaitConnection(ctx); err != nil {
		panic(err)
	}

	publishTopic := fmt.Sprintf("plantino/%v/command", viper.GetString("serialNumber"))

	// create ticker
	timer := time.NewTicker(5 * time.Second)
	deltaT := time.Now().UnixMilli()

	onWateringMode := false
	wateringRetainTime := int64(0)

	fanOffTemporarily := false
	fanOffRetainTime := int64(0)

	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			deltaT = time.Now().UnixMilli() - deltaT

			if onWateringMode {
				wateringEnabled = true
				fanEnabled = true
			} else {
				wateringEnabled = false
				fanEnabled = true
			}

			if fanOffTemporarily {
				fanOffRetainTime += deltaT
				fanEnabled = false
				if fanOffRetainTime > 5*60*1000 {
					fanOffTemporarily = false
					fanEnabled = true
					fanOffRetainTime = 0
				}

			}

			if currentData.Humidity < viper.GetFloat64("thresholds.humidity.low")*100 {
				onWateringMode = true
				wateringRetainTime = 0
			}
			if currentData.Humidity > viper.GetFloat64("thresholds.humidity.high")*100 {
				wateringRetainTime += deltaT
				if wateringRetainTime > viper.GetInt64("thresholds.humidity.duration")*1000 {
					onWateringMode = false
					fanOffTemporarily = true
				}
			}
			// publish data
			fanData := "off"
			if fanEnabled {
				fanData = "on"
			}
			wateringData := "off"
			if wateringEnabled {
				wateringData = "on"
			}

			fanCommand := map[string]interface{}{
				"command": "relay",
				"payload": map[string]interface{}{
					"relay_id": fanRelayAddress,
					"action":   fanData,
				},
			}

			wateringCommand := map[string]interface{}{
				"command": "relay",
				"payload": map[string]interface{}{
					"relay_id": wateringRelayAddress,
					"action":   wateringData,
				},
			}

			fanCommandBytes, _ := json.Marshal(fanCommand)
			wateringCommandBytes, _ := json.Marshal(wateringCommand)

			if _, err = c.Publish(ctx, &paho.Publish{
				QoS:     1,
				Topic:   publishTopic,
				Payload: fanCommandBytes,
			}); err != nil {
				if ctx.Err() == nil {
					panic(err) // Publish will exit when context cancelled or if something went wrong
				}
			}
			if _, err = c.Publish(ctx, &paho.Publish{
				QoS:     1,
				Topic:   publishTopic,
				Payload: wateringCommandBytes,
			}); err != nil {
				if ctx.Err() == nil {
					panic(err) // Publish will exit when context cancelled or if something went wrong
				}
			}
			deltaT = time.Now().UnixMilli()
		case <-ctx.Done():
		}
	}

}

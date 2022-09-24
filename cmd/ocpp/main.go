package main

import (
	"fmt"
	"log"

	ocpp16 "github.com/lorenzodonini/ocpp-go/ocpp1.6"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/remotetrigger"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
)

const (
	chargePointId = "cp0001"
	url           = "ws://localhost:8887"
)

func main() {
	chargePoint := ocpp16.NewChargePoint(chargePointId, nil, nil)

	// Set a handler for all callback functions
	handler := &ChargePointHandler{triggerC: make(chan remotetrigger.MessageTrigger, 1)}
	chargePoint.SetCoreHandler(handler)
	chargePoint.SetRemoteTriggerHandler(handler)

	go func() {
		for err := range chargePoint.Errors() {
			fmt.Println(err)
		}
	}()

	go func() {
		for msg := range handler.triggerC {
			fmt.Println("msg:", msg)
			switch msg {
			case core.BootNotificationFeatureName:
				if res, err := chargePoint.BootNotification("demo", "evcc"); err != nil {
					log.Println("BootNotification:", err)
				} else {
					log.Println("BootNotification:", res)
				}

			case core.StatusNotificationFeatureName:
				if res, err := chargePoint.StatusNotification(1, core.NoError, core.ChargePointStatusAvailable); err != nil {
					log.Println("StatusNotification:", err)
				} else {
					log.Println("StatusNotification:", res)
				}

			case core.MeterValuesFeatureName:
				if res, err := chargePoint.MeterValues(1, []types.MeterValue{
					{SampledValue: []types.SampledValue{
						{Measurand: types.MeasurandPowerActiveImport, Value: "1000"},
					}},
				}); err != nil {
					log.Println("MeterValues:", err)
				} else {
					log.Println("MeterValues:", res)
				}
			}
		}
	}()

	// Connects to central system
	if err := chargePoint.Start(url); err != nil {
		log.Fatal(err)
	}

	log.Printf("connected to central system at %v", url)

	select {}
}

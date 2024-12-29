package scanner

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
	"github.com/mikuta0407/switchbot-contact-sensor_exporter/internal/btle"
	"github.com/pkg/errors"
)

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("Done\n")
	case context.Canceled:
		fmt.Printf("Canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}

func Scan(device string, duration time.Duration) {
	flag.Parse()

	d, err := dev.NewDevice(device)
	if err != nil {
		log.Fatalf("Can't create device: %s", err)
	}
	ble.SetDefaultDevice(d)

	fmt.Printf("Scanning for %s...\n", duration)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), duration))
	chkErr(ble.Scan(ctx, true, btle.Handler, nil))

	fmt.Printf("Found SwitchBot devices:\n")
	for mac, data := range btle.BTDevice {
		isOpen := (data.ServiceData[3] & 0b00000010) >> 1      // 開閉状態取得
		isLeaveOpen := (data.ServiceData[3] & 0b00000100) >> 2 // 開けっ放し状態取得
		time := data.ServiceData[7]                            // 開けっ放しの時間取得
		fmt.Printf("MAC: %s, isOpen: %b, isLeaveOpen: %b, time: %d, Battery: %f\n", mac, isOpen, isLeaveOpen, time, data.Battery)
	}
}

package main

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"
)

func handleBackgroundGatherer(ctx context.Context, logger *log.Logger, errc chan error) {

	go func() {
		logger.Printf("[ Scraping thread ] -> Started sucessfully with period %s seconds", strconv.Itoa(10))
		go func() {
			for {
				logger.Printf("[ Scraping thread ] -> New definition has been loaded from OCP")
				time.Sleep(time.Duration(10) * time.Second)
			}
		}()

		errc <- errors.New("scraping thread is dead baby")

	}()
}

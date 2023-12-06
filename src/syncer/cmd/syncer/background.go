package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	ocp4cli "bitbucket.org/veldrane/golibs/ocp4cli"
	l "github.com/synclib"
)

func handleBackgroundGatherer(ctx context.Context, pods *l.NginxInstancies, config *l.Config, logger *log.Logger, errc chan error) {

	go func() {
		logger.Printf("[ Scraping thread ] -> Started sucessfully with period %s seconds", strconv.Itoa(10))
		ocpSession := ocp4cli.Session()
		ctx := context.Background()
		go func() {
			for {
				runningPods, _ := ocpSession.GetPods(&ctx, config)
				//logger.Printf("[ Scraping thread ] -> Waking up, checking ocp configuration....")
				if l.IsChanged(runningPods, pods.Pods, logger) {
					pods.Update(runningPods, logger)
					var pl string
					for k := range runningPods {
						pl = pl + k + " "
					}
					logger.Printf("[ Scraping thread ] -> Pods updated: %s", pl)
				}
				time.Sleep(time.Duration(10) * time.Second)
			}
		}()

		errc <- fmt.Errorf("%s", "[ Scraping thread ] -> scraping thread is dead baby")

	}()
}

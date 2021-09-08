package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	nsq "github.com/nsqio/go-nsq"
)

const updateDuration = 1 * time.Second

var fatalErr error
func fatal(e error) {
	fmt.Println(e)
	flag.PrintDefaults()
	fatalErr = e
}

var (
	countsLock sync.Mutex
	counts map[string]int
)
func main() {
	defer func() {
		if fatalErr != nil {
			os.Exit(1)
		}
	}()
	log.Println("Connecting to Database...")
	
	db, err := mgo.Dial("mongo")
	if err != nil {
		fatal(err)
		return
	}
	defer func() {
		log.Println("Closing connection to Database")
		db.Close()
	}()
	pollData := db.DB("ballots").C("polls")

	log.Println("Connecting to NSQ")
	q, err := nsq.NewConsumer("votes", "counter", nsq.NewConfig())
	if err != nil {
		fatal(err)
		return
	}
	q.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		countsLock.Lock()
		defer countsLock.Unlock()
		if counts == nil {
			counts = make(map[string]int)
		}
		vote := string(m.Body)
		counts[vote]++
		return nil
	}))
	if err := q.ConnectToNSQLookupd("nsqlookupd:4161"); err != nil {
		fatal(err)
		return
	}

	log.Println("Waiting for votes on NSQ...")
	var updater *time.Timer
	updater = time.AfterFunc(updateDuration, func() {
		countsLock.Lock()
		defer countsLock.Unlock()
		if len(counts) == 0 {
			log.Println("There is no new vote. Skipping update for database.")
		} else {
			log.Println("Updating database...")
			log.Println(counts)
			ok := true
			for option, count := range counts {
				sel := bson.M{"options": bson.M{"$in": []string{option}}}
				up := bson.M{"$inc": bson.M{"results." + option: count}}
				if _, err := pollData.UpdateAll(sel, up); err != nil {
					log.Println("Failed to update: ", err)
					ok = false
					continue
				}
				counts[option] = 0
			}
			if ok {
				log.Println("Database update is completed.")
				counts = nil // reset votes' result
			}
		}
		updater.Reset(updateDuration)
	})

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for {
		select {
		case <-termChan:
			updater.Stop()
			q.Stop()
		case <-q.StopChan:
			// finish
			return
		}
	}
}

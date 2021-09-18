package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/globalsign/mgo"
	nsq "github.com/nsqio/go-nsq"
)

var db *mgo.Session

func dialdb() error {
	var err error
	log.Println("Dial to MongoDB: mongo")
	db, err = mgo.Dial("mongo")
	return err
}

func closedb() {
	db.Close()
	log.Println("Close Connection to MongoDB")
}

type poll struct {
	Options []string
}

func loadOptions() ([]string, error) {
	var options []string
	iter := db.DB("ballots").C("polls").Find(nil).Iter()
	var p poll
	for iter.Next(&p){
		log.Println("loadOptions: ", p.Options)
		options = append(options, p.Options...)
	}
	iter.Close()
	log.Println(options)
	return options, iter.Err()
}

func publishVotes(votes <-chan string) <-chan struct{} {
	stopchan := make(chan struct{}, 1)
	pub, _ := nsq.NewProducer("nsqd:4150", nsq.NewConfig())
	go func() {
		for vote := range votes {
			pub.Publish("votes", []byte(vote)) // Publish the result of votes
		}
		log.Println("Publisher: Stopping")
		pub.Stop()
		log.Println("Publisher: Stopped")
		stopchan <- struct{}{}
	}()
	return stopchan
}

func main() {
	var stoplock sync.Mutex
	stop := false
	stopchan := make(chan struct{}, 1)
	signalchan := make(chan os.Signal, 1)
	go func() {
		<- signalchan
		stoplock.Lock()
		stop = true
		stoplock.Unlock()
		log.Println("Stopping...")
		stopchan <- struct{}{}
		closeConn()
	}()
	signal.Notify(signalchan, syscall.SIGINT, syscall.SIGTERM)
	if err := dialdb(); err != nil {
		log.Fatalln("Failed to connect MangoDB: ", err)
	}
	defer closedb()
	// start processing vote
	votes := make(chan string) // channell to process votes
	publisherStoppedChan := publishVotes(votes)
	twitterStoppedChan := startTwitterStream(stopchan, votes)
	go func() {
		for {
			time.Sleep(1*time.Minute)
			closeConn()
			stoplock.Lock()
			if stop {
				stoplock.Unlock()
				break
			}
			stoplock.Unlock()
		}
	}()
	<- twitterStoppedChan
	close(votes)
	<- publisherStoppedChan
}

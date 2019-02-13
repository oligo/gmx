package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/boltdb/bolt"
	"log"
	"time"
	"fmt"
	"encoding/json"
)

const (
	bucketKeys = "metrics-keys"
	bucketMetrics = "metrics"
)

type DataStore struct {
	*bolt.DB
}

func NewDataStore(file string) *DataStore {
	db, err := bolt.Open(file, 0600, &bolt.Options{
		Timeout: 3 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &DataStore{db,}
}


func (db *DataStore) saveKeys(keys []string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketKeys))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for _, key := range keys {
			id, err := b.NextSequence()
			if err != nil { return err }
			err = b.Put(itob(int(id)), []byte(key))
			if err != nil { return err }
		}

		return nil
	})

	return err
}

func (db *DataStore) getKeys() ([]string, error) {
	var keys []string

	err := db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucketKeys))
		if bk == nil {
			return errors.New("failed to get bucket")
		}
		c := bk.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			keys = append(keys, string(v))
		}

		return nil
	})

	return keys, err
}

func (db *DataStore) saveMetrics(metrics map[string]interface{}) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketMetrics))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		for key, value := range metrics {
			log.Printf("key: %s, value: %v\n", key, value)

			v, err := json.Marshal(metrics[key])
			if err != nil {
				log.Println("Marshalling data failed")
				return err
			}

			k := dbKey(key)
			err = b.Put([]byte(k), v)
			if err != nil { return err }
		}

		return nil
	})

	return err
}

func (db *DataStore) getMetrics(key string) ([]interface{}, error) {
	var metrics []interface{}

	err := db.View(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(bucketMetrics))
		if bk == nil {
			return errors.New("failed to get bucket")
		}
		c := bk.Cursor()
		prefix := []byte(key)
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			metrics = append(metrics, string(v))
		}

		return nil
	})

	return metrics, err
}


type MetricsSyncer struct {
	conn *GMXConn
	interval time.Duration	// scheduling interval
	db *DataStore
	stopChan chan bool
}


func NewMetricsSyncer(conn *GMXConn, interval time.Duration, db *DataStore) MetricsSyncer {
	return MetricsSyncer{
		conn: conn, 
		interval:interval, 
		db: db,
	}
}

func (s MetricsSyncer) retrieveKeys() error {
	keys := s.conn.FetchKeys()

	return s.db.saveKeys(keys)
}

func (s MetricsSyncer) queryAll() error {
	keys, err := s.db.getKeys()

	if err != nil { return err }

	values := s.conn.GetValues(keys)

	return s.db.saveMetrics(values)
}

func (s MetricsSyncer) Run() {
	ticker := time.NewTicker(s.interval)
	s.stopChan = make(chan bool)

	if err := s.retrieveKeys(); err != nil {
		log.Fatalln("Retrieving keys from GMX server failed!")
	}

	go func() {
		for {
			select {
			case <- ticker.C:
				err := s.queryAll()
				if err != nil {
					log.Println(err)
				}
			case <- s.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (s MetricsSyncer) Stop() {
	s.stopChan <- true
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// dbKey returns a key name concatenation by key and timestamp
func dbKey(key string) string {
	return fmt.Sprintf("%s:%d", key, time.Now().Unix())
}
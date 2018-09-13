package Wrench

import (
	bolt "go.etcd.io/bbolt"
	"log"
)

func (w Wrench) Insert(key string, val []byte) bool{
	if !w.Exists() {
		return false
	}
	
	db, err := bolt.Open(w.file, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			// create array to store references to buckets
			allBuckets := []*bolt.Bucket{}
			
			// set first bucket to root bucket
			allBuckets = append(allBuckets,tx.Bucket([]byte(w.path[1])))
			
			// burrow into bottom bucket of path
			for i:=2;i<len(w.path);i++{
				allBuckets = append(allBuckets,allBuckets[i-2].Bucket([]byte(w.path[i])))
			}
			
			allBuckets[len(allBuckets)-1].Put([]byte(key), val)
			
			return nil
		})
	})
	
	return true
}
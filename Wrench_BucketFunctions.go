package Wrench

import (
	bolt "go.etcd.io/bbolt"
	"errors"
	"fmt"
	"log"
)

// This file contains functions performed on the Wrench struct that affect the current bucket

// This function tests that the current bucket exists and is a bucket (not a value)
func (w Wrench) Exists() bool{
	bp := w.path
	if len(bp)==1{
		if bp[0]=="~"{
			return true
		} else {
			return false
		}
	}
	bucketExists := true
	db, err := bolt.Open(w.file, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rootBuckets := []string{}
	db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			rootBuckets=append(rootBuckets,string(name))
			return nil
		})
	})
	if contains(rootBuckets,bp[1]) {
		if len(bp) > 2 {
			db.View(func(tx *bolt.Tx) error {
				allBuckets := []*bolt.Bucket{}
				allBuckets = append(allBuckets, tx.Bucket([]byte(bp[1])))
				for i := 2; i < len(bp); i++ {
					keys := []string{}
					// for each key in parent bucket
					allBuckets[i - 2].ForEach(func(k, v []byte) error {
						keys = append(keys, string(k))
						return nil
					})
					if contains(keys, bp[i]) {
						allBuckets = append(allBuckets, allBuckets[i - 2].Bucket([]byte(bp[i])))
					} else {
						bucketExists = false
						break
					}
				}
				return nil
			})
		}
	} else {
		bucketExists = false
	}
	return bucketExists
}

// This returns an array of all key/value pairs in the current bucket
func (w Wrench) GetAll() []DBVal {
	// create array we'll return later
	r := []DBVal{}

	// verify that specified bucket path exists
	if w.Exists(){

		// open db for reading
		db, err := bolt.Open(w.file, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		if len(w.path)==1{

			// if path is root, dump root buckets
			db.View(func(tx *bolt.Tx) error {
				return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
					t := DBVal{}
					t.w = w
					t.path = w.path
					t.v = nil
					t.k=cpyBytes(name)
					r = append(r, t)
					return nil
				})
			})
		} else {
			// if path isn't the root, recurse into path
			db.View(func(tx *bolt.Tx) error {

				// create array to store references to buckets
				allBuckets := []*bolt.Bucket{}

				// set first bucket to root bucket
				allBuckets = append(allBuckets,tx.Bucket([]byte(w.path[1])))

				// burrow into bottom bucket of path
				for i:=2;i<len(w.path);i++{
					allBuckets = append(allBuckets,allBuckets[i-2].Bucket([]byte(w.path[i])))
				}

				// for all in last bucket, copy values to r
				allBuckets[len(allBuckets)-1].ForEach(func(k, v []byte) error {
					t := DBVal{}
					t.w = w
					t.path = w.path
					t.k=cpyBytes(k)
					if v!=nil {
						t.v = cpyBytes(v)
					} else {
						t.v = nil
					}
					r = append(r,t)
					return nil
				})
				return nil
			})
		}
	} else {
		return nil
	}

	// return the results
	return r
}

// This returns an individual key from the bucket, as well as whether the requested key exists
func (w Wrench) Get(key string) (DBVal,bool){
	vals := w.GetAll()
	for i:=0;i<len(vals);i++{
		if string(vals[i].k) == key{
			return vals[i], true
		}
	}
	return w.newDBVal(), false
}

// return just the values, not the buckets
func (w Wrench) GetValues() []DBVal {
	vals := []DBVal{}
	for _,val := range w.GetAll(){
		if !val.IsBucket() {
			vals = append(vals, val)
		}
	}
	return vals
}

// return just the buckets, not the values
func (w Wrench) GetBuckets() []DBVal {
	vals := []DBVal{}
	for _,val := range w.GetAll(){
		if val.IsBucket() {
			vals = append(vals, val)
		}
	}
	return vals
}

// returns 2 separate arrays, one for buckets and one for values. More efficient if getting both
func (w Wrench) GetBoth() ([]DBVal,[]DBVal){
	bs := []DBVal{}
	vs := []DBVal{}
	for _,val := range w.GetAll(){
		if val.IsBucket() {
			bs = append(bs, val)
		} else {
			vs = append(vs, val)
		}
	}
	return bs,vs
}

// return the total number of entries in the bucket
func (w Wrench) Count() int{
	return len(w.GetAll())
}

// return the number of non-bucket values stored in this bucket
func (w Wrench) CountValues() int {
	count := 0
	for _,val := range w.GetAll(){
		if !val.IsBucket(){
			count++
		}
	}
	return count
}

// return the number of buckets nested inside this bucket
func (w Wrench) CountBuckets() int {
	count := 0
	for _,val := range w.GetAll(){
		if val.IsBucket(){
			count++
		}
	}
	return count
}

func (w Wrench) CountBoth() (int,int){
	bc := 0
	vc := 0
	for _,val := range w.GetAll(){
		if val.IsBucket(){
			bc++
		} else {
			vc++
		}
	}
	return bc,vc
}

func (w Wrench) CreateBucket(name string) error{
	if !w.Exists() {
		err := errors.New("invalid bucket path")
		return err
	}

	key := []byte(name)

	db, err := bolt.Open(w.file, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		if len(w.path)>1 {
			return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
				// create array to store references to buckets
				allBuckets := []*bolt.Bucket{}

				// set first bucket to root bucket
				allBuckets = append(allBuckets, tx.Bucket([]byte(w.path[1])))

				// burrow into bottom bucket of path
				for i := 2; i < len(w.path); i++ {
					allBuckets = append(allBuckets, allBuckets[i-2].Bucket([]byte(w.path[i])))
				}

				allBuckets[len(allBuckets)-1].CreateBucketIfNotExists(key)

				return nil
			})
		} else {
			_,err := tx.CreateBucket(key)
			if err!=nil {
				return fmt.Errorf("[Error] Failed to create bucket: %s", err)
			}
			return nil
		}
	})

	return nil
}

// delete the value or bucket at the specified key
func (w Wrench) Delete(key string) bool{
	
	if !w.Exists() {
		return false
	}
	
	db, err := bolt.Open(w.file, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	
	db.Update(func(tx *bolt.Tx) error {
		if len(w.path)>1 {
			return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
				// create array to store references to buckets
				allBuckets := []*bolt.Bucket{}
				
				// set first bucket to root bucket
				allBuckets = append(allBuckets, tx.Bucket([]byte(w.path[1])))
				
				// burrow into bottom bucket of path
				for i := 2; i < len(w.path); i++ {
					allBuckets = append(allBuckets, allBuckets[i-2].Bucket([]byte(w.path[i])))
				}
				
				if allBuckets[len(allBuckets)-1].Get([]byte(key)) != nil {
					allBuckets[len(allBuckets)-1].Delete([]byte(key))
				} else {
					allBuckets[len(allBuckets)-1].DeleteBucket([]byte(key))
				}
				
				return nil
			})
		} else {
			tx.DeleteBucket([]byte(key))
			return nil
		}
	})
	
	return true
}

// delete all values and buckets located within this bucket
func (w Wrench) Empty() bool{
	
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
			
			allBuckets[len(allBuckets)-1].ForEach(func(k, v []byte) error {
				if v!=nil{
					allBuckets[len(allBuckets)-1].Delete(k)
				} else {
					allBuckets[len(allBuckets)-1].DeleteBucket(k)
				}
				return nil
			})
			
			return nil
		})
	})
	
	return true
}

func (w *Wrench) GoTo(path []string) bool{
	if len(path)==1 && path[0]=="~"{
		w.Reset()
		return true
	}
	oldPath := w.path
	w.path = path
	if !w.Exists() {
		w.path = oldPath
		return false
	}
	return true
}

func (w *Wrench) GoToString(path string) bool{
	return w.GoTo(w.StringToPath(path))
}

package Wrench

import (
	bolt "go.etcd.io/bbolt"
	"fmt"
	"os"
	"strings"
)

// This file is the class definition for the Wrench struct, and basic functions using this struct

// The class wrench is the main tool used to manipulate the bolt database
// Wrench treats the hierarchy of the buckets like navigating a filesystem,
// with all actions happening relative to your current position.
// The base class has two values, "path," which is the current location in the database
// and "file," which is the file being used as a database

type Wrench struct{
	path []string
	file string
}

// Wrench.openDB checks that the file exists, tests opening the db, sets Wrench.db
// However, to prevent the file from locking, it does not keep the db open.
// The db is opened when it needs to be read from or written to, but otherwise is closed.

func (w *Wrench) OpenDB(file string) error {
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	w.file = file
	w.path = []string{"~"}
	defer db.Close()
	return nil
}

// When bolt is told to open a file, and the file doesn't exist, it will create a new file to open at that path.
// The "strict" method prohibits this behaviour; It will not create any files if one doesn't exist.

func (w Wrench) OpenDBStrict(file string) error {
	_, err := os.Stat(w.file)
	if os.IsNotExist(err) {
		return err
	}
	return w.OpenDB(file)
}

// This method returns the current location in the db as a string, based on the "path" variable

func (w Wrench) PathString() string {
	s := ""
	for i:=0;i<len(w.path);i++{
		s=s+ w.path[i]+"/"
	}
	return s
}

// This returns the location as an array of string
// doing it like this makes w.path only set-able from the GoTo function, to prevent errors

func (w Wrench) Path() []string {
	return w.path
}

// This method resets the path variable to the base of the database
// "~" represents the lowest level of the database, based on the home directory in Linux
// All bucket paths will begin with ~ in the first position in the array.

func (w *Wrench) Reset() {
	w.path = []string{"~"}
}

// These functions are for getting the name of the file being modified.
// File returns the whole path, while FileName returns just the name.

func (w Wrench) File() string {
	return w.file
}

func (w Wrench) FileName() string {
	if strings.Contains(w.file,"/"){
		str := strings.Split(w.file,"/")
		return str[len(str)-1]
	}
	if strings.Contains(w.file,"\\"){
		str := strings.Split(w.file,"\\")
		return str[len(str)-1]
	}
	return w.file
}

// These functions are for getting the current location in the database
// The first is for getting the path as it's used in Wrench
// The second makes it into a single string for easier use

func (w Wrench) CurrentBucket() []string {
	return w.path
}

func (w Wrench) CurrentBucketString() string {
	s := ""
	for i:=0;i<len(w.path);i++{
		s=s+w.path[i]+"/"
	}
	return s
}

// This function creates a simple way to check if Wrench is currently
// positioned within the root of the database ("~")

func (w Wrench) IsRoot() bool{
	if len(w.path)==1{
		if w.path[0]=="~" {
			return true
		}
	}
	return false
}

// This function creates a new DBVal with the w value preset as this Wrench object

func (w Wrench) newDBVal() DBVal{
	d := DBVal{}
	d.w = w
	return d
}
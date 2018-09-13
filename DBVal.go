package Wrench

import (
	"fmt"
)

// The DBVal class is used to make working with values in the database easier
// It stores the path to the value, the key of the value, and the value itself

// Values can be either key/value pairs OR buckets. Functions exist for both.
// However, calling functions for one on a DBVal that is the other will cause errors.

type DBVal struct {
	path []string
	k []byte
	v []byte
	w Wrench
}

// Get the key as a string
func (d DBVal) Key() string{
	return string(d.k)
}

// Get the value as a string
func (d DBVal) Val() string{
	return string(d.v)
}

func (d DBVal) Path() []string {
	return d.path
}

func (d DBVal) K() []byte {
	return d.k
}

func (d DBVal) V() []byte {
	return d.v
}

func (d DBVal) ToString(verbose bool) string{
	r := ""
	if verbose && d.IsBucket() {
		w := d.AsWrench()
		bc,vc:=w.CountBoth()
		r = fmt.Sprintf("[Bucket] %s%s\n- Contains %d %s and %d %s\n",d.BucketString(),d.Key(),vc,valPlural(vc),bc,bcktPlural(bc))
	} else if d.IsBucket() {
		r = fmt.Sprintf("%s",d.Key())
	} else if verbose {
		r = fmt.Sprintf("[Key] %s%s\n- Value ([]Byte): %v\n",d.BucketString(),d.Key(),d.v)
	} else {
		r = d.Key()
	}
	return r
}

func plurality(singular string, plural string, count int) string{
	if count==1{
		return singular
	}
	return plural
}

func valPlural(count int) string{
	return plurality("value","values",count)
}

func bcktPlural(count int) string{
	return plurality("bucket","buckets",count)
}

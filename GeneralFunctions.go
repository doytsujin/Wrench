package Wrench

import (
	"math/rand"
	"strings"
)

// Does array contain given string?
func contains(haystack []string, needle string) bool {
	for _, str := range haystack {
		if str == needle {
			return true
		}
	}
	return false
}

// generate a random string of n length
func randString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// split a string for bucket path, but allow escaping backslashes
func escapedSplit(s string) []string{
	if !strings.Contains(s,"\\/"){
		return strings.Split(s,"/")
	}
	rndm := randString(10)
	for i:=0;strings.Contains(s,rndm);i++{
		rndm = randString(10)
	}
	s = strings.Replace(s,"\\/",rndm,-1)
	tmp := strings.Split(s,"/")
	for i:=0;i<len(tmp);i++{
		tmp[i]=strings.Replace(tmp[i],rndm,"/",-1)
	}
	return tmp
}

// copies an array of bytes to prevent it disappearing when the db is closed
func cpyBytes(s []byte) []byte{
	r := []byte{}
	for i:=0;i<len(s);i++ {
		r = append(r,s[i])
	}
	return r
}
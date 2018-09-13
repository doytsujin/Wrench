package Wrench

// This file is for functions performed on DBVal structs that represent buckets in the database

// This returns the path to the bucket this value is in
func (d DBVal) BucketString() string{
	s := ""
	for i:=0;i<len(d.path);i++{
		s=s+d.path[i]+"/"
	}
	return s
}

func (d DBVal) IsBucket() bool{
	return d.v==nil
}

func (d DBVal) AsWrench() Wrench{
	w := Wrench{}
	w.path = append(d.path,d.Key())
	w.file = d.w.file
	return w
}

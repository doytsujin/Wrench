package Wrench

import (
	"strings"
)

func (w Wrench) StringToPath(target string) []string{
	bp := []string{}
	if target=="~"{
		bp = []string{"~"}
	} else if strings.HasPrefix(target,"./"){
		
		// if the path is relative
		bp = w.path
		target = strings.Replace(target,"./","",1)
		tmp := escapedSplit(target)
		for i:=0;i<len(tmp);i++{
			bp = append(bp,tmp[i] )
		}
	} else if strings.HasPrefix(target,"/")||strings.HasPrefix(target,"~/"){
		
		// if the path is absolute
		if strings.HasPrefix(target,"/"){
			target = "~"+target
		}
		tmp := escapedSplit(target)
		for i:=0;i<len(tmp);i++{
			bp = append(bp,tmp[i] )
		}
	} else {
		if w.Exists(){
			// assume the path is relative
			bp = w.path
			tmp := escapedSplit(target)
			for i:=0;i<len(tmp);i++{
				bp = append(bp,tmp[i] )
			}
		} else {
			// assume the path is absolute
			tmp := escapedSplit(target)
			bp = append(bp,"~")
			for i:=0;i<len(tmp);i++{
				bp = append(bp,tmp[i] )
			}
		}
	}
	return bp
}
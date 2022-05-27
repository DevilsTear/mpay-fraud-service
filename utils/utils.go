package utils

import (
	"encoding/json"
	"fmt"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func Struct2Map(obj interface{}) map[string]interface{} {
	var mappedObj map[string]interface{}
	inRec, _ := json.Marshal(obj)
	json.Unmarshal(inRec, &mappedObj)

	return mappedObj
}

func Intersection(s1, s2 []string) (inter []string) {
	hash := make(map[string]bool)
	for _, e := range s1 {
		hash[e] = true
	}
	for _, e := range s2 {
		// If elements present in the hashmap then append intersection list.
		if hash[e] {
			inter = append(inter, e)
		}
	}
	//Remove dups from slice.
	inter = removeDups(inter)
	return
}

//Remove dups from slice.
func removeDups(elements []string) (nodups []string) {
	encountered := make(map[string]bool)
	for _, element := range elements {
		if !encountered[element] {
			nodups = append(nodups, element)
			encountered[element] = true
		}
	}
	return
}

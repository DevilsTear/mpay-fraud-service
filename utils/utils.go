package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// CheckError checks if error is nil and trigger panic
func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

// Struct2Map transform struct interface to map
func Struct2Map(obj interface{}) (map[string]interface{}, error) {
	var mappedObj map[string]interface{}
	inRec, _ := json.Marshal(obj)
	if err := json.Unmarshal(inRec, &mappedObj); err != nil {
		return nil, err
	}

	return mappedObj, nil
}

// Intersection returns intersection of two slices
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
	//Remove duplicates from slice.
	inter = RemoveDuplicates(inter)
	return
}

// RemoveDuplicates Remove duplicates from slice.
func RemoveDuplicates(elements []string) (purgedEls []string) {
	encountered := make(map[string]bool)
	for _, element := range elements {
		if !encountered[element] {
			purgedEls = append(purgedEls, element)
			encountered[element] = true
		}
	}
	return
}

// GetMD5Hash creates and returns the md5 encoded form of a string
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

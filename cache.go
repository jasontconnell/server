package server


import (
	"fmt"
	"crypto/md5"
	"time"
	"net/http"
	"os"
)


func MD5(content []byte) string {
	sum := md5.Sum(content)
	//final := base64.URLEncoding.EncodeToString(sum)
	final := fmt.Sprintf("%x", sum)
	return final
}


// returns true to continue processing (send content), false to not send content (304)

func checkDate(file string, w http.ResponseWriter, req *http.Request) bool {
	retVal := true
	
	if statinfo, err := os.Stat(file); err == nil {
		if t, err := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); err == nil && statinfo.ModTime().Unix() <= t.Unix() {
			w.WriteHeader(http.StatusNotModified)
			retVal = false
		} else if err == nil {
			w.Header().Add("Last-Modified", statinfo.ModTime().Format(http.TimeFormat))
			w.Header().Add("Cache-Control", "max-age=120")
			retVal = true
		}
	}
	return retVal
}

func checkETag(content []byte, w http.ResponseWriter, req *http.Request) bool {
	retVal := true
	md5 := MD5(content)

	if etag := req.Header.Get("If-None-Match"); etag == md5 {
		w.WriteHeader(http.StatusNotModified)
		retVal = false
	} else {
		w.Header().Add("ETag", md5)
		retVal = true
	}

	return retVal
}
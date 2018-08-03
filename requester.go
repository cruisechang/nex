package nex

import (
	"errors"
	"net/http"
	"time"
	"strings"
	"io/ioutil"
)


type Requester interface{
	URL()string
	SetPostURI(path string ,queryPair map[string]string)error
	PostURI()string
	PostQuery()string
	Post()(string,error)
}

type requester struct{
	address string
	port string
	url string
	netClient *http.Client
	postURI string
	postQuery string
}
func NewRequester(address,port string,connectTimeout,handshakeTimeout,requestTimeout int) (Requester,error){
	re:= &requester{
		address:address,
		port:port,
		url:"http://" + address + ":" + port + "/",

	}


	var netTransport = &http.Transport{
		//Dial: (&net.Dialer{
		//	Timeout: connectTimeout * time.Second,
		//}).Dial,
		TLSHandshakeTimeout: time.Duration(handshakeTimeout) * time.Second,
	}

	re.netClient = &http.Client{
		Timeout: time.Duration(requestTimeout)*time.Second ,
		Transport: netTransport,
	}

	return re,nil
}

func (re *requester)URL()string{
	return re.url
}

//SetPostURI
//queryPair is a map , key is query key, value is query value. eg.  ["data"]{"xxxx"} => data=xxxx
//query => data=""&age=33....
func (re *requester)SetPostURI(path string, queryPair map[string]string)error{
	if len(queryPair)<=0{
		return errors.New("query pair is len=0")
	}
	var query string
	for k,v:=range queryPair{
		query=k+"="+v+"&"
	}

	query=strings.TrimRight(query, "&")

	re.postQuery=query
	re.postURI=re.url+path
	return nil
}

//Get after set post url
func (re *requester)PostURI()string{
	return re.postURI
}

func (re *requester)PostQuery()string{
	return re.postQuery
}
func (re *requester)Post()(string,error){
	resp,err:=re.netClient.Post(re.postURI,"application/x-www-form-urlencoded",strings.NewReader(re.postQuery))

	if err != nil {
		return "", errors.New("HTTPPost error")
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("HTTPPost readAll error")
	}


	return string(body), nil

}




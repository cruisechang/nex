package http

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//Client represents client strut interface
type Client interface {
	URL() string
	SetPostURI(path string, queryPair map[string]string) error
	PostURI() string
	PostQuery() string
	Post() (string, error)

	//20190102
	SetHost(scheme, address, port string) (host string)
	Host() (host string)
	Do(req *http.Request) (*http.Response, error)
}

type client struct {
	address   string
	port      string
	url       string
	client    *http.Client
	postURI   string
	postQuery string
	//20190102
	host string
}

//NewClient returns Client interface
func NewClient(address, port string, connectTimeout, handshakeTimeout, requestTimeout int) (Client, error) {
	re := &client{
		address: address,
		port:    port,
		url:     "http://" + address + ":" + port + "/",
	}

	var netTransport = &http.Transport{
		//Dial: (&net.Dialer{
		//	Timeout: connectTimeout * time.Second,
		//}).Dial,
		TLSHandshakeTimeout: time.Duration(handshakeTimeout) * time.Second,
	}

	re.client = &http.Client{
		Timeout:   time.Duration(requestTimeout) * time.Second,
		Transport: netTransport,
	}

	return re, nil
}

func (c *client) URL() string {
	return c.url
}

//SetPostURI
//queryPair is a map , key is query key, value is query value. eg.  ["data"]{"xxxx"} => data=xxxx
//query => data=""&age=33....
func (c *client) SetPostURI(path string, queryPair map[string]string) error {
	if len(queryPair) <= 0 {
		return errors.New("query pair is len=0")
	}
	var query string
	for k, v := range queryPair {
		query = k + "=" + v + "&"
	}

	query = strings.TrimRight(query, "&")

	c.postQuery = query
	c.postURI = c.url + path
	return nil
}

//Get after set post url
func (c *client) PostURI() string {
	return c.postURI
}

func (c *client) PostQuery() string {
	return c.postQuery
}
func (c *client) Post() (string, error) {
	resp, err := c.client.Post(c.postURI, "application/x-www-form-urlencoded", strings.NewReader(c.postQuery))

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

//SetHost sets target host scheme , address and port
func (c *client) SetHost(scheme, address, port string) string {
	if strings.Contains(scheme, "://") {
		c.host = scheme + address + ":" + port
	} else {
		c.host = scheme + "://" + address + ":" + port
	}

	return c.host
}

//Host returns host url
func (c *client) Host() string {
	return c.host
}

//Do do the request
func (c *client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

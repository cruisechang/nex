package http

import (
	"context"
	"errors"
	goHTTP "net/http"
	"time"

	"net"
	"strings"
)

//Service is a httpApi service for http request operation.
type Server interface {
	Start() error
	Stop()
	GetRealAddr(r *goHTTP.Request) string
	RegisterHandler(pattern string, handler goHTTP.Handler)
}

type ServerParameters struct {
	Address        string
	Port           string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
}
type server struct {
	server *goHTTP.Server
	apiMux *goHTTP.ServeMux
}

//NewService make a new SocketManager
func NewServer(params *ServerParameters) (s Server, err error) {

	defer func() {
		if r := recover(); r != nil {

			err = errors.New("httpAPI NewService panic")
		}
	}()

	if params.Port == "" {
		return nil, errors.New("address error ")
	}

	sv := &server{
		apiMux: goHTTP.NewServeMux(),
		server: &goHTTP.Server{
			Addr:           params.Address + ":" + params.Port,
			ReadTimeout:    params.ReadTimeout,
			WriteTimeout:   params.WriteTimeout,
			IdleTimeout:    params.IdleTimeout,
			MaxHeaderBytes: params.MaxHeaderBytes,
			//MaxHeaderBytes: 1 << 20,
		},
	}

	sv.server.Handler = sv.apiMux

	return sv, err

}

func newAPIServeMux() *goHTTP.ServeMux {
	sm := &goHTTP.ServeMux{}
	sm.HandleFunc("/", requestHandler)
	return sm
}

//Start starts service.
func (s *server) Start() error {
	go s.server.ListenAndServe()
	return nil
}

//Stop is used to remove lobbyClient
func (s *server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
	}
}

func (s *server) RegisterHandler(pattern string, handler goHTTP.Handler) {
	s.apiMux.Handle(pattern, handler)
}

func requestHandler(w goHTTP.ResponseWriter, r *goHTTP.Request) {
	print("not foiun")
	w.Write([]byte("not found"))

}

/*
func requestHandler(w goHTTP.ResponseWriter, r *goHTTP.Request) {

	if r := recover(); r != nil {
	}

	var params = ""
	var path = ""

	if r.Method == goHTTP.MethodGet {
		params = r.URL.Query().Get("params")
	} else {
		params = strings.TrimSpace(r.FormValue("params"))
	}

	fmt.Printf("Method : %s %s \n", r.Method, params)

	if len(params) == 0 {
		goHTTP.Error(w, "httpAPI params len==0", goHTTP.StatusBadRequest)
		return
	}

	//b, err := encode.Base64Decode(params)
	b, err := base64.StdEncoding.DecodeString(params)
	if err != nil {
		goHTTP.Error(w, "httpAPI params base64 decode error", goHTTP.StatusBadRequest)
		return
	}

	// fmt.Printf("params : %v \n", string(b))

	path = r.URL.Path
	path = strings.TrimPrefix(r.URL.Path, "/")

	re := Response{}
	if err := json.Unmarshal(b, &re);err!=nil{

	}

	//encode.JSONDecodeToStruct(b, &re)

	fmt.Printf("path :%s Response : %+v \n", path, re)

	//pass path(grabRedEnv) and Response
	//ex: response=map["grabRedEnv"]
	writeParamsChan(map[string]Response{path: re})

	// w.Write(b)

}
*/

func (s *server) GetRealAddr(r *goHTTP.Request) string {

	remoteIP := ""
	// the default is the originating ip. but we try to find better options because this is almost
	// never the right IP
	if parts := strings.Split(r.RemoteAddr, ":"); len(parts) == 2 {
		remoteIP = parts[0]
	}
	// If we have a forwarded-for header, take the address from there
	if xff := strings.Trim(r.Header.Get("X-Forwarded-For"), ","); len(xff) > 0 {
		addrs := strings.Split(xff, ",")
		lastFwd := addrs[len(addrs)-1]
		if ip := net.ParseIP(lastFwd); ip != nil {
			remoteIP = ip.String()
		}
		// parse X-Real-Ip header
	} else if xri := r.Header.Get("X-Real-Ip"); len(xri) > 0 {
		if ip := net.ParseIP(xri); ip != nil {
			remoteIP = ip.String()
		}
	}

	return remoteIP
}

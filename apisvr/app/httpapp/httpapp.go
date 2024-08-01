package httpapp

import (
	"apisvr/app/am"
	"apisvr/app/dbapp"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdygd/goglib"

	"github.com/gorilla/mux"
)

type mapHandler struct {
	staticPath string
}

type spaHandler struct {
	staticPath string
	indexPath  string
}

var redirecturl []string = []string{""}

// ------------------------------------------------------------------------------
// HttpAppHandler
// ------------------------------------------------------------------------------
type HttpAppHandler struct {
	http.Handler
	dbHnd     dbapp.DBHandler
	tlsConfig *tls.Config
}

type rootHandler struct {
	staticPath string
	indexPath  string
}

func (h mapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	am.Applog.Print(1, "Map service..")

	reqPath := r.URL.Path
	ext := filepath.Ext(reqPath)
	am.Applog.Print(1, "ServeHTTP[maphandler1]:[%s][%s]", reqPath, ext)

	// check whether a file exists at the given path
	path := reqPath[1:]
	am.Applog.Print(1, "ServeHTTP[maphandler2] path :[%s]", path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		am.Applog.Print(1, "ServeHTTP[maphandler3] file does not exist :[%s][%v]", path, err)
		// file does not exist, serve index.html
		//http.ServeFile(w, r, path)
		//return
	} else {
		http.ServeFile(w, r, path)
	}
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	am.Applog.Print(2, "WebPage service")

	reqPath := r.URL.Path
	ext := filepath.Ext(reqPath)
	am.Applog.Print(6, "ServeHTTP:[%s][%s]", reqPath, ext)

	if reqPath == "/" {
		if ext != "" {
			am.Applog.Print(6, "ServeHTTP(1):[%s][%s]", reqPath, ext)
			w.WriteHeader(404)
			return
		}

	} else {
		if !(ext == ".css" || ext == ".js" || ext == ".json" || ext == ".ico" || ext == ".png" || ext == ".geojson" || ext == ".svg" || ext == ".otf" || ext == ".ttf" || ext == ".eot" || ext == ".woff" || ext == ".avi" || ext == ".mp4") {

			reqPath := fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
			strurl := fmt.Sprintf("%s", r.URL)
			am.Applog.Warn(" UnExist APIURL(2).. : %v", reqPath)
			//check redirec url
			// set,log redirecto
			for _, redirecurl := range redirecturl {
				if redirecurl == strurl {
					am.Applog.Print(2, "redirecto url.. %s", r.URL)
					var redurl string = ""
					if am.AppVar.Https == "yes" {
						redurl = fmt.Sprintf("https://%s:%d", am.AppVar.Domain, am.AppVar.HttpPort)
					} else {
						redurl = fmt.Sprintf("http://%s:%d", am.AppVar.Domain, am.AppVar.HttpPort)
					}
					am.Applog.Print(2, "Redirect url : %s", redurl)
					http.Redirect(w, r, redurl, http.StatusSeeOther)

					return
				}
			}

			w.WriteHeader(404)
			return
		}
	}

	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

// ------------------------------------------------------------------------------
// AuthMiddleware
// ------------------------------------------------------------------------------
func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqPath := fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		if !strings.Contains(r.URL.String(), "vworld_uw") {
			am.Applog.Always("REQUEST URL : [%s]", reqPath)
		}

		next.ServeHTTP(w, r)
		return
	})
}

// ------------------------------------------------------------------------------
// initHttpTLSconfig
// ------------------------------------------------------------------------------
func initHttpTLSconfig() *tls.Config {
	insecure := flag.Bool("insecure-ssl", false, "Accept/Ignore all server SSL certificates")
	flag.Parse()

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	// Read in the cert file
	certs, err := os.ReadFile(am.AppVar.Sslcertpem)
	if err != nil {
		log.Fatalf("Failed to append %q to RootCAs: %v", am.AppVar.Sslcertpem, err)
	}

	// Append our cert to the system pool
	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		log.Println("No certs appended, using system certs only %q", am.AppVar.Sslcertpem)
	}

	// Trust the augmented cert pool in our client
	config := &tls.Config{
		InsecureSkipVerify: *insecure,
		RootCAs:            rootCAs,
	}

	return config
}

func sendSse(data goglib.EventData) {
	am.Applog.Print(1, "Active sse session : %v", ActivesseSessionList)
	for _, actSession := range ActivesseSessionList {
		CheckSSEMsgChannel(actSession.Key)

		SseMsgChan[actSession.Key] <- data
	}

}

// ------------------------------------------------------------------------------
// processEventMsg
// ------------------------------------------------------------------------------
func ProcessEventMsg() {

	for {
		select {
		case event := <-goglib.ChEvent:
			am.Applog.Print(1, "Get Event message [%s]", event.Msgtype)

			if len(event.Msgtype) > 0 {
				msg := &event
				sendSse(*msg)
			} else {
				am.Applog.Error("undefined sse..[%s](%d)", event.Msgtype, event.Id)
			}

		}
	}
}

// ------------------------------------------------------------------------------
// handleSSE
// ------------------------------------------------------------------------------
func handleSSE() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// get sse session key
		sessionKey := GetSSeSessionKey()
		defer func() {
			am.Applog.Print(3, "Close sse.. [%d]", sessionKey)
			ClearSSeSessionKey(sessionKey)
		}()

		if sessionKey == 0 {
			// invalid key...
			am.Applog.Error("Access handleSSE invalid key.. [%d]", sessionKey)
			<-r.Context().Done()
			return
		}

		// prepare the header
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// prepare the flusher
		flusher, _ := w.(http.Flusher)

		// trap the request under loop forever
		for {
			select {

			case <-r.Context().Done():
				return
			default:
				sseMsg, ok := PopSSEMsgChannel(sessionKey)
				if ok {
					btData := sseMsg.PrepareMessage()
					//am.Applog.Print(2, "PRepareMessage : %v", string(btData[:]))
					fmt.Fprintf(w, "%s\n", btData)

					if sseMsg.Id == "3" {
						am.Applog.Print(1, "SSE SYSINFO (%v)", sseMsg)
					}

					flusher.Flush()
				}
			}
			time.Sleep(time.Millisecond * 5)
		}
	}
}

// ------------------------------------------------------------------------------
// MakeHandler
// ------------------------------------------------------------------------------
func MakeHandler(dbHandler dbapp.DBHandler) *HttpAppHandler {

	r := mux.NewRouter().StrictSlash(true)
	a := &HttpAppHandler{
		Handler: r,
		dbHnd:   dbHandler,
		//tlsConfig: initHttpTLSconfig(),
	}

	// Init API
	// test..

	r.HandleFunc("/gettest", a.GetTest).Methods("GET")
	r.HandleFunc("/posttest", a.PostTest).Methods("GET")
	r.HandleFunc("/deltest", a.DeleteTest).Methods("GET")

	// sse
	r.HandleFunc("/events", handleSSE())

	// webpage
	spa := spaHandler{staticPath: "wwwroot", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	// middleware
	r.Use(AuthMiddleware)

	// sse msg routine
	go ProcessEventMsg()

	return a
}

func GetRefTokenCookieInfo(name, value string, expTm time.Time, maxage int) http.Cookie {
	// http
	//  >> Secure:   false,
	//  >> SameSite: http.SameSiteDefaultMode,

	// https
	//  >> Secure:   true,
	//  >> SameSite: http.SameSiteStrictMode,

	var ck http.Cookie

	if am.AppVar.Https == "yes" {
		ck = http.Cookie{
			Path:     "/",
			Name:     name,
			Value:    value,
			SameSite: http.SameSiteStrictMode,
			Secure:   true,
			Domain:   am.AppVar.Domain,
			Expires:  expTm,

			MaxAge: maxage,
		}

	} else {
		ck = http.Cookie{
			Path:     "/",
			Name:     name,
			Value:    value,
			SameSite: http.SameSiteStrictMode,
			Secure:   false,
			Domain:   am.AppVar.Domain,
			Expires:  expTm,

			MaxAge: maxage,
		}
	}

	return ck

}

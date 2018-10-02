package golang

import (
"bytes"
"crypto/hmac"
"crypto/sha256"
"encoding/hex"
"io/ioutil"
"net/http"
"time"
"log"
	"io"
	"encoding/base64"
)

const (
	SignDateLayout   = "2006-01-02 15:04:05"
	AllowedTimeFrame = time.Minute * 15
	SignatureHeader  = "X-Application-Sign"
	DateHeader       = "X-Application-Sign-Date"
)

// SignRequest sign request with hmac
func SignRequest(r *http.Request, secret []byte) {
	requestURI := r.URL.Path + "?" + r.URL.RawQuery
	dateStr := time.Now().UTC().Format(SignDateLayout)
	var (
		requestBody []byte
		err         error
	)
	if r.ContentLength > 0 {
		// Read the body and let others down the chain to do the same
		requestBody, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error("Error while reading body of %s: %s", requestURI, err.Error())
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(requestURI))
	mac.Write(requestBody)
	mac.Write([]byte(dateStr))
	signature := mac.Sum(nil)
	r.Header.Set(SignatureHeader, hex.EncodeToString(signature))
	r.Header.Set(DateHeader, dateStr)
}

// SignatureMiddleware creates a http middleware which ensures that requests are signed with hmac.
func SignatureMiddleware(secret []byte, enabled bool) func(http.Handler) http.Handler {
	if enabled == false {
		return func(next http.Handler) http.Handler {
			return next
		}
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			signStr := r.Header.Get(SignatureHeader)
			dateStr := r.Header.Get(DateHeader)
			if signStr == "" || dateStr == "" {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("\"Signature is required\""))
				return
			}
			signDate, err := time.Parse(SignDateLayout, dateStr)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("\"Invalid sign date\""))
				return
			}
			signData, err := hex.DecodeString(signStr)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("\"Invalid signature\""))
				return
			}
			d := time.Since(signDate)
			if d < 0 {
				d = -d
			}
			// postpone the check in order not to time leak
			requestURI := r.URL.Path + "?" + r.URL.RawQuery
			var requestBody []byte
			if r.ContentLength > 0 {
				// Read the body and let others down the chain to do the same
				requestBody, err = ioutil.ReadAll(r.Body)
				if err != nil {
					log.Error("Error while reading body of %s: %s", requestURI, err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				r.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
			}
			mac := hmac.New(sha256.New, secret)
			mac.Write([]byte(requestURI))
			mac.Write(requestBody)
			mac.Write([]byte(dateStr))
			expectedMAC := mac.Sum(nil)
			if !hmac.Equal(signData, expectedMAC) || d > AllowedTimeFrame {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("\"Invalid signature\""))
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

var (
	httpClient = new(http.Client)
	APISecret = []byte{}
)


// Returns http.Request with Http basic Auth header
func newAuthenticatedRequest(method, url string, data io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("user", "passwd")
	return req, nil
}

func newSignedRequest(method, url string, data io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	SignRequest(req, APISecret)
	return req, nil
}


func GetSecretKey(key string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(key)
}

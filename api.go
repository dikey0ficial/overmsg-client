package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/sqweek/dialog"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	// TCPServerURL _
	TCPServerURL string
	// HTTPServerURL _
	HTTPServerURL string
	// TCPConn _
	TCPConn net.Conn
)

func initAPI() {
	err := tryGetURLs()
	if err != nil {
		errl.Println(err)
		dialog.Message("Error finding servers").Title("Error!!1").Error()
		os.Exit(1)
	}
	if conf.Name != "" && conf.Token != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var ech chan error
		TCPConn, ech, err = connectTCP(ctx, conf.Token)
		if err != nil {
			errl.Println(err)
			dialog.Message("Error connecting to server").Title("Error!!1").Error()
			os.Exit(1)
		}
		go func(ech chan error) {
			for {
				if conf.Token == "" {
					return
				}
				err := <-ech
				if err != nil {
					_, err := TCPConn.Read(make([]byte, 1))
					if TCPConn == nil || err != nil {
						errl.Println(err)
						if _, err := ping(ctx, HTTPServerURL); err != nil {
							errl.Println(err)
							TCPConn, ech, err = connectTCP(ctx, conf.Token)
							if err != nil {
								errl.Println(err)
								dialog.Message("Error: can't connect to server").Title("Error!!1").Error()
							}
						}
					}
				}
			}
		}(ech)
		go func() {
			in := bufio.NewScanner(TCPConn)
			for in.Scan() {
				if in.Err() != nil {
					errl.Println(err)
					dialog.Message("Error getting info from server (more info in logs) :(")
					continue
				}
				var got map[string]interface{}
				if err := json.Unmarshal(in.Bytes(), &got); err != nil {
					if strings.TrimSpace(in.Text()) == "success" {
						continue
					}
					errl.Printf("%s (%v)\n", in.Text(), err)
					dialog.Message("Error getting info from server (more info in logs) :(")
					continue
				}
				var t string
				if el, ok := got["type"]; !ok {
					continue
				} else if t, ok = el.(string); !ok {
					continue
				}
				switch t {
				case "message":
					var mess message
					if err := json.Unmarshal(in.Bytes(), &mess); err != nil {
						errl.Println(err)
						dialog.Message("Error getting info from server (more info in logs) :(")
						continue
					}
					if mess.Error != "" {
						errl.Println(mess.Error)
						dialog.Message("Error getting info from server (more info in logs) :(")
						continue
					}
					messCh <- mess
				default:
					continue
				}
			}
		}()
	}
}

func tryGetURLs() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var (
		bestDur time.Duration = 11 * time.Second
		bestURL string
	)
	for _, url := range conf.ServerURLs {
		time, err := ping(ctx, url)
		if err != nil {
			continue
		}
		if time < bestDur {
			bestDur = time
			bestURL = url
		}
	}
	if bestURL == "" {
		return errors.New("Found no aviable servers in list of servers")
	}
	HTTPServerURL, TCPServerURL = "http://"+bestURL+":4422", bestURL+":4242"
	return nil
}

func connectTCP(ctx context.Context, token string) (net.Conn, chan error, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", TCPServerURL)
	if err != nil {
		return nil, nil, err
	}
	conn.Write([]byte(token + "\n"))
	errCh := make(chan error, 0)
	go heartbeat(conn, errCh)
	return conn, errCh, nil
}

func heartbeat(conn net.Conn, ech chan error) {
	t := time.NewTicker(30 * time.Second)
	for {
		<-t.C
		_, err := http.Post(HTTPServerURL+"/heartbeat", "text/plain", strings.NewReader(conf.Token))
		ech <- err
	}
}

func ping(ctx context.Context, url string) (time.Duration, error) {
	start := time.Now()
	req, err := http.NewRequest("GET", "http://"+url+":4422", nil)
	if err != nil {
		return time.Duration(0), err
	}
	req = req.WithContext(ctx)

	client := http.DefaultClient
	_, err = client.Do(req)
	if err != nil {
		return time.Duration(0), err
	}
	return start.Sub(time.Now()), nil
}

var client = &http.Client{}

func reg(name, pass string) (string, error) {
	req, err := http.NewRequest("POST", HTTPServerURL+"/reg",
		strings.NewReader(`{"name":"`+name+`","pass":"`+pass+`"}`),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var ans answer
	if err := json.Unmarshal(dat, &ans); err != nil {
		return "", err
	}
	if !ans.Success {
		return "err", errors.New(ans.Error)
	}
	var (
		tokinf interface{}
		token  string
		ok     bool
	)
	if tokinf, ok = ans.Res["token"]; !ok {
		return "", errors.New("got no token")
	} else if token, ok = tokinf.(string); !ok {
		return "", errors.New("got not-string token")
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", errors.New("got empty token")
	}
	return token, nil
}

func getToken(name, pass string) (string, error) {
	req, err := http.NewRequest("POST", HTTPServerURL+"/get_token",
		strings.NewReader(`{"name":"`+name+`","pass":"`+pass+`"}`),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var ans answer
	if err := json.Unmarshal(dat, &ans); err != nil {
		return "", err
	}
	if !ans.Success {
		return "err", errors.New(ans.Error)
	}
	var (
		tokinf interface{}
		token  string
		ok     bool
	)
	if tokinf, ok = ans.Res["token"]; !ok {
		return "", errors.New("got no token")
	} else if token, ok = tokinf.(string); !ok {
		return "", errors.New("got not-string token")
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return "", errors.New("got empty token")
	}
	return token, nil
}

func goOffline(token string) error {
	req, err := http.NewRequest("POST", HTTPServerURL+"/go_offline", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Auth-Token", token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var ans answer
	if err := json.Unmarshal(dat, &ans); err != nil {
		return err
	}
	if !ans.Success {
		return errors.New(ans.Error)
	}
	return nil
}

func sendMessage(token, msg, to string) error {
	body, err := json.Marshal(sendMessageReq{
		PeerName: to,
		Message:  msg,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", HTTPServerURL+"/send_message",
		bytes.NewReader(body),
	)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Auth-Token", token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var ans answer
	if err := json.Unmarshal(dat, &ans); err != nil {
		return err
	}
	if !ans.Success {
		return errors.New(ans.Error)
	}
	return nil
}

func isOnline(nick string) (bool, bool, error) {
	body, err := json.Marshal(isOnlineReq{
		Name: nick,
	})
	if err != nil {
		return false, false, err
	}
	resp, err := http.Post(HTTPServerURL+"/is_online", "application/json", bytes.NewReader(body))
	if err != nil {
		return false, false, err
	}
	defer resp.Body.Close()
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, false, err
	}
	var ans answer
	if err := json.Unmarshal(dat, &ans); err != nil {
		return false, false, err
	}
	if !ans.Success {
		return false, false, errors.New(ans.Error)
	}
	var (
		is, exists bool
	)
	if elem, ok := ans.Res["is"]; !ok {
		return false, false, errors.New("got no is")
	} else if is, ok = elem.(bool); !ok {
		return false, false, errors.New("got not-bool is")
	}
	if elem, ok := ans.Res["exists"]; !ok {
		return false, false, errors.New("got no exists")
	} else if exists, ok = elem.(bool); !ok {
		return false, false, errors.New("got not-bool exists")
	}
	return is, exists, nil
}

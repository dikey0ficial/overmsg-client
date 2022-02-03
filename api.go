package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	// TCPServerURL _
	TCPServerURL string
	// HTTPServerURL _
	HTTPServerURL string
)

func connect() (net.Conn, error) {
	conn, err := net.Dial("tcp", TCPServerURL)
	if err != nil {
		return nil, err
	}
	if b, err := bufio.NewReader(conn).ReadBytes(7); err != nil {
		return nil, err
	} else if string(b) != "success" {
		debl.Println(b)
		return nil, errors.New("got wrong answer")
	}
	return conn, nil
}

func ping(conn net.Conn) (time.Duration, error) {
	var dur time.Duration
	if conn == nil {
		return dur, errors.New("nil connection")
	}
	start := time.Now()
	conn, err := net.Dial("tcp", TCPServerURL)
	if err != nil {
		return dur, err
	}
	if b, err := bufio.NewReader(conn).ReadBytes(4); err != nil {
		return dur, err
	} else if string(b) != "pong" {
		return dur, errors.New("got wrong answer")
	}
	return start.Sub(time.Now()), nil
}

var client = &http.Client{}

func auth(name string) (string, error) {
	req, err := http.NewRequest("POST", HTTPServerURL+"/auth",
		strings.NewReader(`{"name":"`+name+`"`),
	)
	if err != nil {
		return "", err
	}
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
		return "", errors.New(ans.Error)
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

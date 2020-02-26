package main

/*
 * Read ALB log lines as CSV format and output JSON
 */

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// ALBLog represents a log line from an ALB Log
type ALBLog struct {
	RequestType            string `json:"requestType"`
	Timestamp              string `json:"timestamp"`
	ELBResourceID          string `json:"elbResourceID"`
	ClientPort             string `json:"clientPort"`
	TargetPort             string `json:"targetPort"`
	RequestProcessingTime  string `json:"requestProcessingTime"`
	TargetProcessingTime   string `json:"targetProcessingTime"`
	ResponseProcessingTime string `json:"responseProcessingTime"`
	ELBStatusCode          string `json:"elbStatusCode"`
	TargetStatusCode       string `json:"targetStatusCode"`
	ReceivedBytes          string `json:"receivedBytes"`
	SentBytes              string `json:"sentBytes"`
	Request                string `json:"request"`
	UserAgent              string `json:"userAgent"`
	SSLCipher              string `json:"sslCipher"`
	SSLProtocol            string `json:"sslProtocol"`
	TargetGroupARN         string `json:"targetGroupARN"`
	TraceID                string `json:"traceID"`
	DomainName             string `json:"domainName"`
	ChosenCertARN          string `json:"chosenCertARN"`
	MatchedRulePriority    string `json:"matchedRulePriority"`
	RequestCreationTime    string `json:"requestCreationTime"`
	ActionsExecuted        string `json:"actionsExecuted"`
	RedirectURL            string `json:"redirectURL"`
	ErrorReason            string `json:"errorReason"`
}

// NewALBLog returns a new ALBLog object
func NewALBLog(record []string) ALBLog {

	logLine := ALBLog{
		RequestType:            record[0],
		Timestamp:              record[1],
		ELBResourceID:          record[2],
		ClientPort:             record[3],
		TargetPort:             record[4],
		RequestProcessingTime:  record[5],
		TargetProcessingTime:   record[6],
		ResponseProcessingTime: record[7],
		ELBStatusCode:          record[8],
		TargetStatusCode:       record[9],
		ReceivedBytes:          record[10],
		SentBytes:              record[11],
		Request:                record[12],
		UserAgent:              record[13],
		SSLCipher:              record[14],
		SSLProtocol:            record[15],
		TargetGroupARN:         record[16],
		TraceID:                record[17],
		DomainName:             record[18],
		ChosenCertARN:          record[19],
		MatchedRulePriority:    record[20],
		RequestCreationTime:    record[21],
		ActionsExecuted:        record[22],
		RedirectURL:            record[23],
		ErrorReason:            record[24],
	}
	return logLine
}

func main() {
	r := csv.NewReader(bufio.NewReader(os.Stdin))
	r.Comma = ' '
	w := json.NewEncoder(os.Stdout)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		logLine := NewALBLog(record)

		err = w.Encode(logLine)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println()
}

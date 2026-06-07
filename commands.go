package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

// PerformPong response pack with "PONG", or optionally a passed in argument
func PerformPong(args []string) string {
	if len(args) > 0 {
		return stringMsg(args[0])
	}

	return stringMsg("PONG")
}

type StoreItem struct {
	Value  string
	Expiry time.Time
	Mutex  sync.Mutex
}

var store = make(map[string]*StoreItem)

func PerformSet(args []string) string {
	if len(args) < 2 {
		return errorMsg("invalid syntax provided to SET")
	}
	key := &args[0]
	val := &args[1]
	var exp time.Time

	if len(args) > 2 {
		position := 2
		for position < len(args) {
			switch strings.ToLower(args[position]) {
			case "px":
				// set expiry in ms
				if len(args) < position+1 {
					return errorMsg("no time provided to 'PX'")
				}
				expMillis, err := strconv.Atoi(string(args[position+1]))
				if err != nil {
					return errorMsg("invalid time provided to 'PX'")
				}
				exp = time.Now().Add(time.Duration(expMillis) * time.Millisecond)
				// + 2 because we need to skip the "px" and the value
				position += 2
			case "ex":
				// set expiry in seconds
				if len(args) < position+1 {
					return errorMsg("no time provided to 'EX'")
				}
				expSeconds, err := strconv.Atoi(string(args[position+1]))
				if err != nil {
					return errorMsg("invalid time provided to 'EX'")
				}
				exp = time.Now().Add(time.Duration(expSeconds) * time.Second)
				// + 2 because we need to skip the "ex" and the value
				position += 2
			default:
				return errorMsg(fmt.Sprintf("invalid argument '%s'", args[position]))
			}
		}
	}

	// Default to an hour
	if exp.Equal((time.Time{})) {
		exp = time.Now().Add(time.Hour)
	}

	store[*key] = &StoreItem{
		Value:  *val,
		Expiry: exp,
	}

	return stringMsg("OK")
}

func PerformGet(args []string) string {
	if len(args) == 0 {
		return errorMsg("no value provided to 'GET'")
	}

	item := store[args[0]]
	if item == nil {
		return nilBulkStringMsg()
	}

	// Lock because many clients may be trying to access the same item
	item.Mutex.Lock()
	defer item.Mutex.Unlock()
	

	now := time.Now()
	if item.Expiry.Before(now) {
		store[args[0]] = nil
	}

	return stringMsg(item.Value)
}

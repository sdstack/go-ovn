/**
 * Copyright (c) 2017 eBay Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 **/

package goovn

import (
	"fmt"
	"sync"

	"github.com/socketplane/libovsdb"
)

const (
	insert string = "insert"
	mutate string = "mutate"
	del    string = "delete"
	list   string = "select"
	update string = "update"
)

const (
	NBDB string = "OVN_Northbound"
)

const (
	LSWITCH     string = "Logical_Switch"
	LPORT       string = "Logical_Switch_Port"
	ACLS        string = "ACL"
	Address_Set string = "Address_Set"
)

const (
	UNIX string = "unix"
	TCP  string = "tcp"
)

const (
	//random seed.
	MAX_TRANSACTION = 1000
)

type ovnDBClient struct {
	socket   string
	server   string
	port     int
	protocol string
	dbclient *libovsdb.OvsdbClient
}

type ovnDBImp struct {
	client     *ovnDBClient
	cache      map[string]map[string]libovsdb.Row
	cachemutex sync.Mutex
	tranmutex  sync.Mutex
	callback   OVNSignal
}

type OVNDB struct {
	imp *ovnDBImp
}

func New(socketfile string, proto string, server string, port int, callback OVNSignal) (OVNDBApi, error) {
	var dbapi *OVNDB
	var err error

	switch proto {
	case UNIX:
		dbapi, err = newNBBySocket(socketfile, callback)
	case TCP:
		dbapi, err = newNBByServer(server, port, callback)
	default:
		err = fmt.Errorf("The protocol [%s] is not supported", proto)
	}

	if err != nil {
		return nil, err
	}

	return dbapi, nil
}

func SetCallBack(c OVNDBApi, callback OVNSignal) {
	if c != nil {
		c.SetCallBack(callback)
	}
}

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

package ovn

import (
	"fmt"
	"sync"

	"github.com/socketplane/libovsdb"
)

const (
	OpInsert string = "insert"
	OpMutate string = "mutate"
	OpDelete string = "delete"
	OpList   string = "select"
	OpUpdate string = "update"
)

const (
	NBDB string = "OVN_Northbound"
)

const (
	TableNBGlobal       string = "NB_Global"
	TableLS             string = "Logical_Switch"
	TableLSP            string = "Logical_Switch_Port"
	TableACL            string = "ACL"
	TableAS             string = "Address_Set"
	TableLB             string = "Load_Balancer"
	TableLR             string = "Logical_Router"
	TableQoS            string = "QoS"
	TableLRP            string = "Logical_Router_Port"
	TableLRSR           string = "Logical_Router_Static_Route"
	TableNAT            string = "NAT"
	TableDHCPOptions    string = "DHCP_Options"
	TableConnection     string = "Connection"
	TableDNS            string = "DNS"
	TableSSL            string = "SSL"
	TableGatewayChassis string = "Gateway_Chassis"
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

var mu sync.Mutex
var ovnDBApi OVNDBApi

func GetInstance(socketfile string, protocol string, server string, port int, callback OVNSignal) (OVNDBApi, error) {
	var dbapi *OVNDB
	var err error

	mu.Lock()
	defer mu.Unlock()
	if ovnDBApi != nil {
		return ovnDBApi, nil
	}

	if protocol == UNIX {
		dbapi, err = newNBBySocket(socketfile, callback)
	} else if protocol == TCP {
		dbapi, err = newNBByServer(server, port, callback)
	} else {
		err = fmt.Errorf("The protocol [%s] is not supported", protocol)
	}

	if err != nil {
		return nil, err
	}

	ovnDBApi = dbapi

	return ovnDBApi, nil
}

func SetCallBack(callback OVNSignal) {
	if ovnDBApi != nil {
		ovnDBApi.SetCallBack(callback)
	}
}

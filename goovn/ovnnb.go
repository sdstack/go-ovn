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
	"errors"

	"github.com/socketplane/libovsdb"
)

func newNBClient(socketfile string, protocol string, server string, port int) (*ovnDBClient, error) {
	client := &ovnDBClient{
		socket:   socketfile,
		server:   server,
		port:     port,
		protocol: UNIX,
	}

	if protocol == UNIX {
		clt, err := libovsdb.ConnectWithUnixSocket(socketfile)
		if err != nil {
			//glog.Fatalf("OVN DB initial failed: (%v) with socket file %s.", err, socketfile)
			return nil, err
		}
		client.dbclient = clt
		return client, nil

	} else if protocol == TCP {
		clt, err := libovsdb.Connect(server, port)
		if err != nil {
			//glog.Fatalf("OVN DB initial failed: (%v) on %s:%d", err, server, port)
			return nil, err
		}
		client.dbclient = clt
		return client, nil
	}
	return nil, errors.New("OVN DB initial failed: (unsupported protocol)")
}

func newNBBySocket(socketfile string, callback OVNSignal) (*OVNDB, error) {
	odb, err := newNBClient(socketfile, UNIX, "", 0)
	if err == nil {
		nb, err := newNBImp(odb, callback)
		if err != nil {
			return nil, err
		}
		return &OVNDB{nb}, nil
	} else {
		return nil, err
	}
}

func newNBByServer(server string, port int, callback OVNSignal) (*OVNDB, error) {
	odb, err := newNBClient("", TCP, server, port)
	if err != nil {
		nb, err := newNBImp(odb, callback)
		if err != nil {
			return nil, err
		}
		return &OVNDB{nb}, nil
	} else {
		return nil, err
	}
}

func (odb *OVNDB) LSWAdd(lsw string) (*OvnCommand, error) {
	return odb.imp.lswAddImp(lsw)
}

func (odb *OVNDB) LSWDel(lsw string) (*OvnCommand, error) {
	return odb.imp.lswDelImp(lsw)
}

func (odb *OVNDB) LSWList() (*OvnCommand, error) {
	return odb.imp.lswListImp()
}

func (odb *OVNDB) LSPAdd(lsw string, lsp string) (*OvnCommand, error) {
	return odb.imp.lspAddImp(lsw, lsp)
}

func (odb *OVNDB) LSPDel(lsp string) (*OvnCommand, error) {
	return odb.imp.lspDelImp(lsp)
}

func (odb *OVNDB) LSPSetAddress(lsp string, addresses ...string) (*OvnCommand, error) {
	return odb.imp.lspSetAddressImp(lsp, addresses...)
}

func (odb *OVNDB) LSPSetPortSecurity(lsp string, security ...string) (*OvnCommand, error) {
	return odb.imp.lspSetPortSecurityImp(lsp, security...)
}

func (odb *OVNDB) LSPSetDHCPv4Options(lsp string, cidr string, options map[string]string, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.lspSetDHCPv4OptionsImp(lsp, cidr, options, external_ids)
}

func (odb *OVNDB) ACLAdd(lsw, direct, match, action string, priority int, external_ids map[string]string, logflag bool) (*OvnCommand, error) {
	return odb.imp.aclAddImp(lsw, direct, match, action, priority, external_ids, logflag)
}

func (odb *OVNDB) ACLDel(lsw, direct, match string, priority int, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.aclDelImp(lsw, direct, match, priority, external_ids)
}

func (odb *OVNDB) ASAdd(name string, addrs []string, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.ASAdd(name, addrs, external_ids)
}

func (odb *OVNDB) ASDel(name string) (*OvnCommand, error) {
	return odb.imp.ASDel(name)
}

func (odb *OVNDB) ASUpdate(name string, addrs []string, external_ids map[string]string) (*OvnCommand, error) {
	return odb.imp.ASUpdate(name, addrs, external_ids)
}

func (odb *OVNDB) LSSetOpt(lsp string, options map[string]string) (*OvnCommand, error) {
	return odb.imp.LSSetOpt(lsp, options)
}

func (odb *OVNDB) Execute(cmds ...*OvnCommand) error {
	return odb.imp.Execute(cmds...)
}

func (odb *OVNDB) GetLogicPortsBySwitch(lsw string) []*LogcalPort {
	return odb.imp.GetLogicPortsBySwitch(lsw)
}

func (odb *OVNDB) GetACLsBySwitch(lsw string) []*ACL {
	return odb.imp.GetACLsBySwitch(lsw)
}

func (odb *OVNDB) GetAddressSets() []*AddressSet {
	return odb.imp.GetAddressSets()
}

func (odb *OVNDB) GetASByName(name string) *AddressSet {
	return odb.imp.GetASByName(name)
}

func (odb *OVNDB) SetCallBack(callback OVNSignal) {
	odb.imp.callback = callback
}

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
	"github.com/socketplane/libovsdb"
)

type LogicalRouter struct {
	UUID    string
	Name    string
	Enabled bool

	Ports        []*LogicalRouterPort
	StaticRoutes []*LogicalRouterStaticRoute
	NAT          []*NAT
	LoadBalancer []*LoadBalancer

	Chassis         string
	DNATForceSNATIP string
	LBForceSNATIP   string
	ExternalID      map[interface{}]interface{}
}

func (odbi *ovnDBImp) lrAddImp(name string) (*OvnCommand, error) {
	namedUUID, err := newUUID()
	if err != nil {
		return nil, err
	}

	//row to insert
	lrouter := make(OVNRow)
	lrouter["name"] = name

	if uuid := odbi.getRowUUID(tableLogicalRouter, lrouter); len(uuid) > 0 {
		return nil, ErrorExist
	}

	insertOp := libovsdb.Operation{
		Op:       opInsert,
		Table:    tableLogicalRouter,
		Row:      lrouter,
		UUIDName: namedUUID,
	}

	operations := []libovsdb.Operation{insertOp}

	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

func (odbi *ovnDBImp) lrDelImp(name string) (*OvnCommand, error) {
	condition := libovsdb.NewCondition("name", "==", name)

	deleteOp := libovsdb.Operation{
		Op:    opDelete,
		Table: tableLogicalRouter,
		Where: []interface{}{condition},
	}

	operations := []libovsdb.Operation{deleteOp}

	return &OvnCommand{operations, odbi, make([][]map[string]interface{}, len(operations))}, nil
}

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

func (odbi *ovnDBImp) GetLR(name string) []*LogicalRouter {
	var lrList []*LogicalRouter
	odbi.cachemutex.Lock()
	defer odbi.cachemutex.Unlock()

	for uuid, drows := range odbi.cache[tableLogicalRouter] {
		if lrName, ok := drows.Fields["name"].(string); ok && lrName == name {
			lr := odbi.RowToLR(uuid)
			lrList = append(lrList, lr)
		}
	}
	return lrList
}

func (odbi *ovnDBImp) RowToLR(uuid string) *LogicalRouter {
	return &LogicalRouter{
		UUID:       uuid,
		Name:       odbi.cache[tableLoadBalancer][uuid].Fields["name"].(string),
		ExternalID: odbi.cache[tableLoadBalancer][uuid].Fields["external_ids"].(libovsdb.OvsMap).GoMap,
	}
}

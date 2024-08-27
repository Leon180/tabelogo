package entitymodel

import "tabelog-spider/model/enum"

type TabelogElementInfo struct {
	Link                  enum.URL
	ElementCollectorSlice ElementCollectorSlice
}

type TabelogElementInfoSlice []TabelogElementInfo

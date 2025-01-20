package responsemodel

import (
	"encoding/json"
	"google-map/grpc/proto"

	"google.golang.org/protobuf/types/known/anypb"
)

func ConvertInterfaceResponse(res interface{}) *proto.InterfaceResponse {
	jsonBytes, _ := json.Marshal(res)
	anyValue, _ := anypb.New(&anypb.Any{Value: jsonBytes})
	return &proto.InterfaceResponse{
		Data: anyValue,
	}
}

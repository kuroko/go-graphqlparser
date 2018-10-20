package server

import (
	"encoding/json"
	"testing"
)

func BenchmarkResponse_MarshalGraphQL(b *testing.B) {
	var res []byte
	var err error

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		response := Response{
			Data: ResponseValue{
				Kind: ResponseValueKindObject,
				ObjectValue: []ResponseObjectField{
					{Name: "f1", Value: ResponseValue{
						Kind: ResponseValueKindObject,
						ObjectValue: []ResponseObjectField{
							{Name: "episodeID", Value: ResponseValue{
								Kind:     ResponseValueKindInt,
								IntValue: 4,
							}},
						},
					}},
					{Name: "f2", Value: ResponseValue{
						Kind: ResponseValueKindObject,
						ObjectValue: []ResponseObjectField{
							{Name: "title", Value: ResponseValue{
								Kind:        ResponseValueKindString,
								StringValue: "Attack of the Clones",
							}},
						},
					}},
				},
			},
		}

		res, err = response.MarshalGraphQL()
		if err != nil {
			b.Error(err)
		}
	}

	_ = res
}

func BenchmarkResponse_MapToJSON(b *testing.B) {
	var mapBS []byte
	var err error

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		mapRes := make(map[string]interface{})
		mapRes["data"] = map[string]interface{}{
			"f1": map[string]interface{}{
				"episodeID": 4,
			},
			"f2": map[string]interface{}{
				"title": "Attack of the Clones",
			},
		}

		mapBS, err = json.Marshal(mapRes)
		if err != nil {
			b.Error(err)
		}
	}

	_ = mapBS
}

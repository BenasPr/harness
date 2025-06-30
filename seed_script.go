package main

import (
	"fmt"
	"os"

	pb_tuning "github.com/VU-ASE/rovercom/packages/go/tuning"
	"google.golang.org/protobuf/proto"
)

func main() {
	os.Mkdir("seed_bytes3", 0755)

	seeds := []*pb_tuning.TuningState{
		{
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("speed", 0.6),
			},
		},
		{
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("kp", 0.005),
			},
		},
		{
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("kd", 0.00002),
			},
		},
		{
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("ki", 1.0),
			},
		},
		{
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("servo-trim", -0.10),
			},
		},
		{
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("threshold-value", 0),
			},
		},
		{
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("servo-scaler", 0.8),
			},
		},
	}

	for i, seed := range seeds {
		data, err := proto.Marshal(seed)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(fmt.Sprintf("seed_bytes3/seed%d", i), data, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func numParam(key string, value float64) *pb_tuning.TuningState_Parameter {
	return &pb_tuning.TuningState_Parameter{
		Parameter: &pb_tuning.TuningState_Parameter_Number{
			Number: &pb_tuning.TuningState_Parameter_NumberParameter{
				Key:   key,
				Value: float32(value),
			},
		},
	}
}

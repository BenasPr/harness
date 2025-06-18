package main

import (
	"fmt"
	"os"

	pb_tuning "github.com/VU-ASE/rovercom/packages/go/tuning"
	"google.golang.org/protobuf/proto"
)

func main() {
	os.Mkdir("in", 0755)

	seeds := []*pb_tuning.TuningState{
		{
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("speed", 1.1),
			},
		},
		{ // 1: Realistic single parameter
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("kp", 0.003),
			},
		},
		{ // 2: Realistic two parameters
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("kp", 0.003),
				numParam("kd", 1e-5),
			},
		},
		{ // 3: Realistic three parameters
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("kp", 0.003),
				numParam("kd", 1e-5),
				numParam("ki", 0.0),
			},
		},
		{ // 4: Edge case - empty key
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("", 123.45),
			},
		},
		{ // 5: Unknown key
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("unknown-key", 42.0),
			},
		},
		{ // 6: Large float
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("servo-trim", 1e38),
			},
		},
		{ // 7: Negative float
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("servo-trim", -0.15),
			},
		},
		{ // 8: Zero value
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("threshold-value", 0.0),
			},
		},
		{ // 9: Valid but borderline float
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("servo-scaler", 0.9),
			},
		},
		{ // 10: Mix of known and unknown
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				numParam("kd", 1e-5),
				numParam("mystery", -12345.67),
			},
		},
	}

	for i, seed := range seeds {
		data, err := proto.Marshal(seed)
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(fmt.Sprintf("in/seed%d", i), data, 0644)
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

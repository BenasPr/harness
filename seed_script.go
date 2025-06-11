package main

import (
	"fmt"
	"os"
	"time"

	pb_tuning "github.com/VU-ASE/rovercom/packages/go/tuning"
	"google.golang.org/protobuf/proto"
)

func main() {
	os.Mkdir("in", 0755)

	seeds := []*pb_tuning.TuningState{
		{ // Valid float, known key
			Timestamp: uint64(time.Now().UnixMilli()),
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				{
					Parameter: &pb_tuning.TuningState_Parameter_Number{
						Number: &pb_tuning.TuningState_Parameter_NumberParameter{
							Key:   "kp",
							Value: 1.23,
						},
					},
				},
			},
		},
		{ // Invalid key
			Timestamp: uint64(time.Now().UnixMilli()),
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				{
					Parameter: &pb_tuning.TuningState_Parameter_Number{
						Number: &pb_tuning.TuningState_Parameter_NumberParameter{
							Key:   "speed",
							Value: 3.14,
						},
					},
				},
			},
		},
		{ // Negative value
			Timestamp: uint64(time.Now().UnixMilli()),
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				{
					Parameter: &pb_tuning.TuningState_Parameter_Number{
						Number: &pb_tuning.TuningState_Parameter_NumberParameter{
							Key:   "kd",
							Value: -999.99,
						},
					},
				},
			},
		},
		{ // Empty key
			Timestamp: uint64(time.Now().UnixMilli()),
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				{
					Parameter: &pb_tuning.TuningState_Parameter_Number{
						Number: &pb_tuning.TuningState_Parameter_NumberParameter{
							Key:   "",
							Value: 0.0,
						},
					},
				},
			},
		},
		{ // Large float
			Timestamp: uint64(time.Now().UnixMilli()),
			DynamicParameters: []*pb_tuning.TuningState_Parameter{
				{
					Parameter: &pb_tuning.TuningState_Parameter_Number{
						Number: &pb_tuning.TuningState_Parameter_NumberParameter{
							Key:   "ki",
							Value: 1e38,
						},
					},
				},
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

package main

import (
	"log"
	"net"
	"time"

	//"vu/ase/transceiver/src/serverconnection"
	//"vu/ase/transceiver/src/state"
	pb_tuning "github.com/VU-ASE/rovercom/packages/go/tuning"
	"google.golang.org/protobuf/proto"
)

func main() {
	tuning := &pb_tuning.TuningState{
		Timestamp: uint64(time.Now().UnixMilli()),
		DynamicParameters: []*pb_tuning.TuningState_Parameter{
			{
				Parameter: &pb_tuning.TuningState_Parameter_Number{
					Number: &pb_tuning.TuningState_Parameter_NumberParameter{
						Key:   "kp",
						Value: 0.002,
					},
				},
			},
		},
	}

	data, err := proto.Marshal(tuning)
	if err != nil {
		log.Fatalf("Failed to marshal tuning: %v", err)
	}

	conn, err := net.Dial("tcp", "192.168.0.121:9000")
	if err != nil {
		log.Fatalf("Failed to connect to transceiver: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		log.Fatalf("Failed to send data: %v", err)
	}

	log.Printf("TuningState sent to transceiver successfully.")
}

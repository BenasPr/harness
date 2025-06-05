package main

import (
	"os"
	"log"
	"net"
	"time"
	"strconv"

	//"vu/ase/transceiver/src/serverconnection"
	//"vu/ase/transceiver/src/state"
	pb_tuning "github.com/VU-ASE/rovercom/packages/go/tuning"
	"google.golang.org/protobuf/proto"
)

func main() {
	// fuzzArg := os.Args[1]

	// value, err := strconv.ParseFloat(fuzzArg, 32)
	// if err != nil {
	// 	log.Fatalf("Error converting argument to float: %v", err)
	// }

	// // Convert float64 to float32
	// fuzzValue := float32(value)


	fuzzArg := ""
	if len(os.Args) > 1 {
		fuzzArg = os.Args[1]
	} else {
		// If no argument is provided, read from stdin
		data := make([]byte, 256) // Adjust size as needed
		n, err := os.Stdin.Read(data)
		if err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
		fuzzArg = string(data[:n - 1])
	}

	value, err := strconv.ParseFloat(fuzzArg, 32)
	if err != nil {
		log.Fatalf("Error converting argument to float: %v", err)
	}

	// Convert float64 to float32
	fuzzValue := float32(value)



	tuning := &pb_tuning.TuningState{
		Timestamp: uint64(time.Now().UnixMilli()),
		DynamicParameters: []*pb_tuning.TuningState_Parameter{
			{
				Parameter: &pb_tuning.TuningState_Parameter_Number{
					Number: &pb_tuning.TuningState_Parameter_NumberParameter{
						Key:   "kp",
						Value: fuzzValue,
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

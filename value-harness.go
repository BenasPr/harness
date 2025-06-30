// package main

// import (
// 	"log"
// 	"net"
// 	"os"
// 	"strconv"

// 	"syscall"
// 	"time"

// 	//"vu/ase/transceiver/src/serverconnection"
// 	//"vu/ase/transceiver/src/state"
// 	pb_tuning "github.com/VU-ASE/rovercom/packages/go/tuning"
// 	"google.golang.org/protobuf/proto"
// )

// func main() {
// 	// fuzzArg := os.Args[1]

// 	// value, err := strconv.ParseFloat(fuzzArg, 32)
// 	// if err != nil {
// 	// 	log.Fatalf("Error converting argument to float: %v", err)
// 	// }

// 	// // Convert float64 to float32
// 	// fuzzValue := float32(value)

// 	fuzzArg := ""
// 	if len(os.Args) > 1 {
// 		fuzzArg = os.Args[1]
// 	} else {
// 		// If no argument is provided, read from stdin
// 		data := make([]byte, 256) // Adjust size as needed
// 		n, err := os.Stdin.Read(data)
// 		if err != nil {
// 			log.Fatalf("Error reading from stdin: %v", err)
// 		}
// 		fuzzArg = string(data[:n-1])
// 	}

// 	value, err := strconv.ParseFloat(fuzzArg, 32)
// 	if err != nil {
// 		log.Fatalf("Error converting argument to float: %v", err)
// 	}

// 	// Convert float64 to float32
// 	fuzzValue := float32(value)

// 	tuning := &pb_tuning.TuningState{
// 		Timestamp: uint64(time.Now().UnixMilli()),
// 		DynamicParameters: []*pb_tuning.TuningState_Parameter{
// 			{
// 				Parameter: &pb_tuning.TuningState_Parameter_Number{
// 					Number: &pb_tuning.TuningState_Parameter_NumberParameter{
// 						Key:   "speed",
// 						Value: fuzzValue,
// 					},
// 				},
// 			},
// 		},
// 	}
// 	syscall.Kill(syscall.Getpid(), syscall.SIGSEGV)

// 	data, err := proto.Marshal(tuning)
// 	if err != nil {
// 		log.Fatalf("Failed to marshal tuning: %v", err)
// 	}

// 	conn, err := net.Dial("tcp", "192.168.0.146:9000")
// 	if err != nil {
// 		log.Fatalf("Failed to connect to transceiver: %v", err)
// 	}
// 	defer conn.Close()

// 	_, err = conn.Write(data)
// 	if err != nil {
// 		log.Fatalf("Failed to send data: %v", err)
// 	}

// 	log.Printf("TuningState sent to transceiver successfully.")
// }


package main

import (
    "encoding/binary"
    "fmt"
    "io"
    "io/ioutil"
    "net"
    "os"
    "path/filepath"
    "strconv"
    "time"
	"math"

    // Import your generated protobuf package here
    pb "github.com/VU-ASE/rovercom/packages/go/tuning"
    "google.golang.org/protobuf/proto"
)

func main() {
    // Read 56 bytes from stdin
    inputBuf := make([]byte, 28)
    n, err := io.ReadFull(os.Stdin, inputBuf)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to read input: %v\n", err)
        os.Exit(1)
    }
    if n != 28 {
        fmt.Fprintf(os.Stderr, "Expected 56 bytes (7 doubles), got %d\n", n)
        os.Exit(1)
    }

    // Parse 7 float64 values (little endian)
    values := make([]float32, 7)
    for i := 0; i < 7; i++ {
        values[i] = math.Float32frombits(binary.LittleEndian.Uint32(inputBuf[i*4 : (i+1)*4]))
	}
    

    keys := []string{
        "speed",
        "kp",
        "kd",
        "ki",
        "servo-trim",
        "threshhold-value",
        "servo-scaler",
    }

    // Construct TuningState protobuf
	dynamicParams := make([]*pb.TuningState_Parameter, 7)
	for i := 0; i < 7; i++ {
	dynamicParams[i] = &pb.TuningState_Parameter{
			Parameter: &pb.TuningState_Parameter_Number{
				Number: &pb.TuningState_Parameter_NumberParameter{
					Key:   keys[i],
					Value: values[i],
				},
			},
		}
	}		
    tuning := &pb.TuningState{
		DynamicParameters: dynamicParams,
	}
    // for i := 0; i < 7; i++ {
    //     param := &pb.TuningState_Parameter{}
    //     numberParam := &pb.TuningState_Parameter_NumberParameter{
    //         Key:   keys[i],
    //         Value: values[i],
    //     }
    //     param.Param = &pb.TuningState_Parameter_Number{Number: numberParam}
    //     tuning.DynamicParameters = append(tuning.DynamicParameters, param)
    // }

    // Serialize protobuf
    serialized, err := proto.Marshal(tuning)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to serialize protobuf message: %v\n", err)
        os.Exit(1)
    }

    // Handle connection state
    stateFile := "pipeline_status.flag"
    wasConnectedBefore := false
    if data, err := ioutil.ReadFile(stateFile); err == nil && len(data) > 0 {
        wasConnectedBefore = (data[0] == '1')
    }

    // Connect to server
    conn, err := net.Dial("tcp", "192.168.0.146:9000")
    if err != nil {
        _ = ioutil.WriteFile(stateFile, []byte("0"), 0644)
        // _ = startPipeline("./start_pipeline.sh")

        if wasConnectedBefore {
            _ = os.MkdirAll("crashes", 0755)
            crashFile := filepath.Join("crashes", "crash_"+strconv.FormatInt(time.Now().Unix(), 10)+".bin")
            if err := os.Rename("last_input.bin", crashFile); err != nil {
                fmt.Fprintf(os.Stderr, "Failed to save crash input: %v\n", err)
            }
            // Here you would pause/resume AFL++ as needed
            panic(42)
        } else {
            os.Exit(1)
        }
    }
    defer conn.Close()
    _ = ioutil.WriteFile(stateFile, []byte("1"), 0644)

    // Send serialized protobuf message
    totalSent := 0
    toSend := len(serialized)
    for totalSent < toSend {
        n, err := conn.Write(serialized[totalSent:])
        if err != nil {
            fmt.Fprintf(os.Stderr, "send: %v\n", err)
            os.Exit(1)
        }
        totalSent += n
    }

    // Save input to last_input.bin
    if err := ioutil.WriteFile("last_input.bin", inputBuf, 0644); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to write last_input.bin: %v\n", err)
        os.Exit(1)
    }
}

// startPipeline runs a shell script.
// func startPipeline(path string) error {
//     return execShell(path)
// }

// func execShell(path string) error {
//     // Simple implementation; for more control use os/exec
//     return os.system("/bin/sh", "-c", path)
// }

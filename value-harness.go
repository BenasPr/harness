package main

import (
    "encoding/binary"
    "io"
    "io/ioutil"
    "net"
    "os"
    "path/filepath"
    "strconv"
    "time"
	"math"
    "os/exec"

    // Import your generated protobuf package here
    pb "github.com/VU-ASE/rovercom/packages/go/tuning"
    "google.golang.org/protobuf/proto"
)

func main() {
    inputBuf := make([]byte, 28)
    n, err := io.ReadFull(os.Stdin, inputBuf)
    if err != nil {
        os.Exit(1)
    }
    if n != 28 {
        os.Exit(1)
    }
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

    serialized, err := proto.Marshal(tuning)
    if err != nil {
        os.Exit(1)
    }

    stateFile := "pipeline_status.flag"
    wasConnectedBefore := false
    if data, err := ioutil.ReadFile(stateFile); err == nil && len(data) > 0 {
        wasConnectedBefore = (data[0] == '1')
    }

    conn, err := net.Dial("tcp", "192.168.0.146:9000")
    if err != nil {
        _ = ioutil.WriteFile(stateFile, []byte("0"), 0644)
        exec.Command("./start_pipeline.sh")
        if wasConnectedBefore {
            _ = os.MkdirAll("crashes", 0755)
            crashFile := filepath.Join("crashes", "crash_"+strconv.FormatInt(time.Now().Unix(), 10)+".bin")
            os.Rename("last_input.bin", crashFile)
            panic(42)
        } else {
            exec.Command("./start_pipeline.sh")
        }
        return
    }
    defer conn.Close()
    _ = ioutil.WriteFile(stateFile, []byte("1"), 0644)

    totalSent := 0
    toSend := len(serialized)
    for totalSent < toSend {
        n, err := conn.Write(serialized[totalSent:])
        if err != nil {
            os.Exit(1)
        }
        totalSent += n
    }

    if err := ioutil.WriteFile("last_input.bin", inputBuf, 0644); err != nil {
        os.Exit(1)
    }
}


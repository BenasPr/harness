package main

import (
    "fmt"
    "io/ioutil"
    "net"
    "os"
    "os/exec"
    "path/filepath"
    "time"
)

func main() {
    inputBuf, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        return
    }

    if len(inputBuf) == 0 {
        return
    }

    stateFile := "pipeline_status.flag"
    wasConnectedBefore := false

    if _, err := os.Stat(stateFile); err == nil {
        data, err := ioutil.ReadFile(stateFile)
        if err == nil && len(data) > 0 {
            wasConnectedBefore = data[0] == '1'
        }
    }

    conn, err := net.Dial("tcp", "192.168.0.146:9000")
    if err != nil {
        ioutil.WriteFile(stateFile, []byte("0"), 0644)

        exec.Command("./start_pipeline.sh")

        if wasConnectedBefore {
            os.MkdirAll("crashes", os.ModePerm);
            crashFile := filepath.Join("crashes", fmt.Sprintf("crash_%d.bin", time.Now().Unix()))
            os.Rename("last_input.bin", crashFile)
        } else {
            exec.Command("./start_pipeline.sh")
        }
        return
    }
    // defer conn.Close()

    ioutil.WriteFile(stateFile, []byte("1"), 0644)

    totalSent := 0
    toSend := len(inputBuf)
    for totalSent < toSend {
        sent, err := conn.Write(inputBuf[totalSent:])
        if err != nil {
            return
        }
        totalSent += sent
    }
    ioutil.WriteFile("last_input.bin", inputBuf, 0644)
}


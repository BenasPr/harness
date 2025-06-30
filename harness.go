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
    // Read input from stdin
    inputBuf, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Error reading input:", err)
        return
    }

    if len(inputBuf) == 0 {
        fmt.Fprintln(os.Stderr, "Empty input")
        return
    }

    stateFile := "pipeline_status.flag"
    wasConnectedBefore := false

    // Check if the state file exists and read its content
    if _, err := os.Stat(stateFile); err == nil {
        data, err := ioutil.ReadFile(stateFile)
        if err == nil && len(data) > 0 {
            wasConnectedBefore = data[0] == '1'
        }
    }

    // Create a TCP connection
    conn, err := net.Dial("tcp", "192.168.0.146:9000")
    if err != nil {
        // If connection fails, write "0" to state file and run the script
        ioutil.WriteFile(stateFile, []byte("0"), 0644)

        // Execute the start_pipeline.sh script
        cmd := exec.Command("./start_pipeline.sh")
        if err := cmd.Run(); err != nil {
            fmt.Fprintln(os.Stderr, "Failed to start pipeline:", err)
        }

        if wasConnectedBefore {
            // Create crashes directory if it doesn't exist
            if err := os.MkdirAll("crashes", os.ModePerm); err != nil {
                fmt.Fprintln(os.Stderr, "Failed to create crashes directory:", err)
            }

            crashFile := filepath.Join("crashes", fmt.Sprintf("crash_%d.bin", time.Now().Unix()))
            if err := os.Rename("last_input.bin", crashFile); err != nil {
                fmt.Fprintln(os.Stderr, "Failed to save crash input:", err)
            }

            // panic("Crash") // Simulate throwing an exception
			x := 1
    		y := 0
    		fmt.Println(x / y)
        } else {
            return
        }
    }
    // defer conn.Close()

    // If connected successfully, write "1" to state file
    ioutil.WriteFile(stateFile, []byte("1"), 0644)

    // Send the input buffer over the connection
    totalSent := 0
    toSend := len(inputBuf)
    for totalSent < toSend {
        sent, err := conn.Write(inputBuf[totalSent:])
        if err != nil {
            fmt.Fprintln(os.Stderr, "Error sending data:", err)
            return
        }
        totalSent += sent
    }

    // Write the input buffer to last_input.bin
    ioutil.WriteFile("last_input.bin", inputBuf, 0644)
}


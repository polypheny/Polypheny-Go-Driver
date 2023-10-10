package main

import (
        os "os"
        fmt "fmt"
        strings "strings"
        log "log"
        exec "os/exec"
)

const (
        GO_PACKAGE_NAME = "polypheny.com/protos" // the package name of generated proto code
        PROTOS_PATH = "protos" // where the proto files are
        OUTPUT_PATH = "." // where the generated files will go
	COMMAND = "protoc" // the protoc command to compile proto files
)

func main() {
	rawSwitches := fmt.Sprintf("--proto_path=%s --go_out=%s --go-grpc_out=%s", PROTOS_PATH, OUTPUT_PATH, OUTPUT_PATH) // command line args starting with -
	var args []string // proto file names
	entries, err := os.ReadDir(PROTOS_PATH)

        if err != nil {
                log.Fatal(err)
        }

        for _, e := range entries {
		if strings.Contains(e.Name(), ".proto") {
			temp := fmt.Sprintf(" --go_opt=M%s=%s --go-grpc_opt=M%s=%s ", e.Name(), GO_PACKAGE_NAME, e.Name(), GO_PACKAGE_NAME)
                        rawSwitches = rawSwitches + temp
                        args = append(args, PROTOS_PATH + "/" + e.Name())
		}
	}
	fields := strings.Fields(rawSwitches)
	fields = append(fields, args[:]...)
	err = exec.Command(COMMAND, fields[:]...).Run()
	if err != nil {
                log.Fatalf("Failed to complie the proto files. Please try to manually execute\n%s", COMMAND + strings.Join(fields, " "))
        }
}

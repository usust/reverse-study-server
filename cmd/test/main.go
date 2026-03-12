package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CompileRequest struct {
	SourceCode      string                `json:"source_code"`
	CompilerOptions CompileRequestOptions `json:"compiler_options"`
}

type CompileRequestOptions struct {
	OptLevel         string `json:"opt_level"`
	NoInline         bool   `json:"no_inline"`
	KeepFramePointer bool   `json:"keep_frame_pointer"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	reqBody := CompileRequest{
		SourceCode: "#include <stdio.h>\nint main(void){printf(\"hello\\n\");return 0;}",
		CompilerOptions: CompileRequestOptions{
			OptLevel:         "O2",
			NoInline:         true,
			KeepFramePointer: true,
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://192.168.33.118:10000/compile", bytes.NewReader(payload))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var er ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
			panic(fmt.Errorf("request failed with status %d", resp.StatusCode))
		}
		panic(fmt.Errorf("request failed with status %d: %s", resp.StatusCode, er.Error))
	}

	zipData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("/Users/lyu/Downloads/2/build-artifact.zip", zipData, 0o644); err != nil {
		panic(err)
	}

	fmt.Printf("saved build-artifact.zip (%d bytes)\n", len(zipData))
	fmt.Printf("X-Compile-Return-Code: %s\n", resp.Header.Get("X-Compile-Return-Code"))
	fmt.Printf("X-Compile-Timed-Out: %s\n", resp.Header.Get("X-Compile-Timed-Out"))
}

package compiler

import (
	"context"
	"os"
	"testing"
)

func TestCompileByMSVC(t *testing.T) {
	reqBody := CompileRequest{
		SourceCode: "#include <stdio.h>\nint main(void){printf(\"hello\\n\");return 0;}",
		CompilerOptions: CompileRequestOptions{
			OptLevel:         "O2",
			NoInline:         true,
			KeepFramePointer: true,
		},
	}
	artifact, err := CompileByMSVC(context.Background(), reqBody)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile("/Users/lyu/Downloads/2/sample.exe", artifact.ExecutableData, 0o755)
	if err != nil {
		t.Fatal(err)
	}
}

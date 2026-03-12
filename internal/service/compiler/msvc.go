package compiler

import (
	"context"
	"fmt"
	"strings"

	compilercore "reverse-study-server/compiler"
)

// CompileCRequest 是“调用 MSVC 编译服务”的输入。
type CompileCRequest struct {
	SourceCode      string                             `json:"source_code"`
	CompilerOptions compilercore.CompileRequestOptions `json:"compiler_options"`
}

// CompileCResponse 是“调用 MSVC 编译服务”的输出。
type CompileCResponse struct {
	Artifact compilercore.Artifact `json:"artifact"`
}

// CompileCByMSVC 调用 compiler.CompileByMSVC 执行 C 源码编译。
func CompileCByMSVC(ctx context.Context, req CompileCRequest) (CompileCResponse, error) {
	sourceCode := strings.TrimSpace(req.SourceCode)
	if sourceCode == "" {
		return CompileCResponse{}, fmt.Errorf("source_code is required")
	}

	artifact, err := compilercore.CompileByMSVC(ctx, compilercore.CompileRequest{
		SourceCode:      sourceCode,
		CompilerOptions: req.CompilerOptions,
	})
	if err != nil {
		return CompileCResponse{}, err
	}
	if len(artifact.ExecutableData) == 0 {
		return CompileCResponse{}, fmt.Errorf("compiled executable is empty")
	}

	return CompileCResponse{Artifact: artifact}, nil
}

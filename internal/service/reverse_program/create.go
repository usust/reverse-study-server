package reverse_program

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reverse-study-server/compiler"
	"reverse-study-server/internal/bootstrap"
	dbmodel "reverse-study-server/internal/model"
	rp "reverse-study-server/internal/repository/reverse_program"
	"reverse-study-server/internal/service/config/prompt"
	chatapi "reverse-study-server/volcengine/chat"
	"strconv"
)

func GenerateNewReverseProgram(ctx context.Context,
	req chatapi.PromptChatRequest,
	compileOption compiler.CompileRequestOptions,
	programInfo dbmodel.ReverseProgram,
) error {
	promptChatResponse, err := prompt.GenerateCCode(ctx, req)
	if err != nil {
		return err
	}
	artifact, err := compiler.CompileByMSVC(ctx, compiler.CompileRequest{
		SourceCode:      promptChatResponse.Response.Content,
		CompilerOptions: compileOption,
	})
	if err != nil {
		return err
	}

	// 计算程序文件的MD5
	if len(artifact.ExecutableData) > 0 {
		executableHash := md5.Sum(artifact.ExecutableData)
		programInfo.ProgramFileMD5 = hex.EncodeToString(executableHash[:])
	}

	// 创建编译记录
	if err = rp.Create(ctx, bootstrap.GormDB, &programInfo); err != nil {
		return err
	}
	
	programDir := filepath.Join(programInfo.BaseDir, "programs")
	if err = os.MkdirAll(programDir, 0o755); err != nil {
		return err
	}
	rootDir := filepath.Join(programDir, strconv.FormatUint(programInfo.ID, 10))
	if err = os.MkdirAll(rootDir, 0o755); err != nil {
		return err
	}
	// 1) 写入源码
	if err = os.WriteFile(
		filepath.Join(rootDir, programInfo.SourceFileName), // 保存路径
		[]byte(promptChatResponse.Response.Content),        // 由AI生成的源码
		0o644); err != nil {
		return err
	}

	// 2) 文件
	if err = os.WriteFile(
		filepath.Join(rootDir, programInfo.ProgramFileName), // 保存路径
		artifact.ExecutableData,                             // 编译后的文件
		0o644); err != nil {
		return err
	}

	// 3) meta文件
	metaJSON, err := json.MarshalIndent(artifact.Meta, "", "  ")
	if err = os.WriteFile(
		filepath.Join(rootDir, programInfo.MetaFileName), // 保存路径
		metaJSON,                                         // meta文件
		0o644); err != nil {
		return err
	}

	return nil
}

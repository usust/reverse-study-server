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

	/*
		// ID 是程序记录主键，使用数据库自增数字。
		ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

		// Score 是题目分值。
		Score int `gorm:"not null;default:100" json:"score"`
		// Difficulty 是题目难度（如 easy/medium/hard）。
		Difficulty string `gorm:"size:32;not null;default:medium" json:"difficulty"`

		// CompletedCount 是完成数量。
		CompletedCount int `gorm:"not null;default:0" json:"completedCount"`

		// BaseDir 是下载路径或保存路径标识。
		BaseDir string `gorm:"size:512;not null" json:"baseDir"`
		// SourceFileName 是源码文件名，默认 main.c。
		SourceFileName string `gorm:"size:128;not null;default:main.c" json:"sourceFileName"`
		// ProgramFileName 编译后文件名
		ProgramFileName string `gorm:"size:128;not null;default:main.c" json:"programFileName"`
		// ProgramFileMD5 是程序源码文件的 MD5。
		ProgramFileMD5 string `gorm:"size:64;not null;index" json:"programFileMd5"`
	*/

	if err = os.MkdirAll(programInfo.BaseDir, 0o755); err != nil {
		return err
	}
	rootDir := filepath.Join(programInfo.BaseDir, strconv.FormatUint(programInfo.ID, 10))
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
		metaJSON, // meta文件
		0o644); err != nil {
		return err
	}

	return nil
}

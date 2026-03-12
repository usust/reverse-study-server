package compiler

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

// 使用 http://exsi-win10-ltsc:10000/compile 的API来编译c源码

const (
	defaultCompileAPIURL  = "http://192.168.33.118:10000/compile"
	defaultCompileTimeout = 90 * time.Second
)

// CompileRequest 对应 POST /compile 的 JSON 请求结构。
// 这里只接受源码字符串和少量高层编译开关，不接受任何原始命令行参数。
type CompileRequest struct {
	SourceCode      string                `json:"source_code" binding:"required"`
	CompilerOptions CompileRequestOptions `json:"compiler_options" binding:"required"`
}

// CompileRequestOptions 是暴露给调用方的编译选项。
// 这些字段会在 compiler 包内部映射成严格白名单中的 MSVC 参数，
// HTTP 层不会把用户输入直接拼接成 cl.exe 参数。
type CompileRequestOptions struct {
	OptLevel         string `json:"opt_level" binding:"required"`
	NoInline         bool   `json:"no_inline"`
	KeepFramePointer bool   `json:"keep_frame_pointer"`
}

// Meta 会写入 meta.json，同时也会返回给 HTTP 层，
// 供响应头输出部分状态信息。
type Meta struct {
	CompiledAt      string   `json:"compiled_at"`
	CompilerOptions []string `json:"compiler_options"`
	ReturnCode      int      `json:"return_code"`
	TimedOut        bool     `json:"timed_out"`
	WorkDir         string   `json:"work_dir"`
}

// Artifact 是单次编译完成后在内存中的结果对象。
// ExecutableData 在编译成功且 sample.exe 存在时才有值。
type Artifact struct {
	ExecutableData []byte
	BuildLog       string
	Meta           Meta
}

// CompileResponse 是编译服务的通用响应结构。
// 不同版本编译服务字段可能有差异，这里预留多种常见字段。
type CompileResponse struct {
	Success          *bool    `json:"success"`
	Message          string   `json:"message"`
	Error            string   `json:"error"`
	BuildLog         string   `json:"build_log"`
	Meta             Meta     `json:"meta"`
	CompiledAt       string   `json:"compiled_at"`
	CompilerOptions  []string `json:"compiler_options"`
	ReturnCode       *int     `json:"return_code"`
	TimedOut         *bool    `json:"timed_out"`
	WorkDir          string   `json:"work_dir"`
	ExecutableData   string   `json:"executable_data"`
	ExecutableBase64 string   `json:"executable_base64"`
	ZipBase64        string   `json:"zip_base64"`
	ZipData          string   `json:"zip_data"`
	BinaryBase64     string   `json:"binary_base64"`
	Binary           string   `json:"binary"`
	OutputBase64     string   `json:"output_base64"`
	Stdout           string   `json:"stdout"`
	Stderr           string   `json:"stderr"`
}

// CompileByMSVC 调用远端编译服务编译 C 源码，并返回可执行文件、日志和元信息。
func CompileByMSVC(ctx context.Context, req CompileRequest) (Artifact, error) {
	req.SourceCode = strings.TrimSpace(req.SourceCode)
	if req.SourceCode == "" {
		return Artifact{}, fmt.Errorf("source_code is required")
	}

	req.CompilerOptions.OptLevel = strings.TrimSpace(req.CompilerOptions.OptLevel)
	if req.CompilerOptions.OptLevel == "" {
		req.CompilerOptions.OptLevel = "O0"
	}

	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, defaultCompileTimeout)
		defer cancel()
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return Artifact{}, fmt.Errorf("marshal compile request failed: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, defaultCompileAPIURL, bytes.NewReader(payload))
	if err != nil {
		return Artifact{}, fmt.Errorf("create compile request failed: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return Artifact{}, fmt.Errorf("call compile api failed: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return Artifact{}, fmt.Errorf("read compile response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		msg := strings.TrimSpace(string(bodyBytes))
		if msg == "" {
			msg = fmt.Sprintf("request failed with status %d", resp.StatusCode)
		}
		return Artifact{}, fmt.Errorf("%s", msg)
	}

	artifact := Artifact{
		BuildLog: "",
		Meta: Meta{
			CompiledAt:      resp.Header.Get("X-Compile-Compiled-At"),
			CompilerOptions: nil,
			ReturnCode:      parseIntHeader(resp.Header.Get("X-Compile-Return-Code")),
			TimedOut:        parseBoolHeader(resp.Header.Get("X-Compile-Timed-Out")),
			WorkDir:         resp.Header.Get("X-Compile-Work-Dir"),
		},
	}

	contentType := strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Type")))
	extractedMultipart, multipartErr := fillArtifactFromMultipartBody(&artifact, bodyBytes, contentType)
	if multipartErr != nil {
		return Artifact{}, multipartErr
	}
	if extractedMultipart {
		return artifact, nil
	}

	trimmedBody := strings.TrimSpace(string(bodyBytes))
	if strings.Contains(contentType, "application/json") || (len(trimmedBody) > 0 && strings.HasPrefix(trimmedBody, "{")) {
		var parsed CompileResponse
		if err := json.Unmarshal(bodyBytes, &parsed); err == nil {
			meta := parsed.Meta
			if strings.TrimSpace(meta.CompiledAt) == "" {
				meta.CompiledAt = parsed.CompiledAt
			}
			if len(meta.CompilerOptions) == 0 {
				meta.CompilerOptions = parsed.CompilerOptions
			}
			if parsed.ReturnCode != nil {
				meta.ReturnCode = *parsed.ReturnCode
			}
			if parsed.TimedOut != nil {
				meta.TimedOut = *parsed.TimedOut
			}
			if strings.TrimSpace(meta.WorkDir) == "" {
				meta.WorkDir = parsed.WorkDir
			}
			if strings.TrimSpace(meta.CompiledAt) == "" {
				meta.CompiledAt = artifact.Meta.CompiledAt
			}
			if len(meta.CompilerOptions) == 0 {
				meta.CompilerOptions = artifact.Meta.CompilerOptions
			}
			if meta.ReturnCode == 0 {
				meta.ReturnCode = artifact.Meta.ReturnCode
			}
			if strings.TrimSpace(meta.WorkDir) == "" {
				meta.WorkDir = artifact.Meta.WorkDir
			}
			artifact.Meta = meta

			buildLog := strings.TrimSpace(parsed.BuildLog)
			if buildLog == "" {
				buildLog = strings.TrimSpace(parsed.Stderr)
			}
			if buildLog == "" {
				buildLog = strings.TrimSpace(parsed.Stdout)
			}
			artifact.BuildLog = buildLog

			zipBase64 := strings.TrimSpace(parsed.ZipData)
			if zipBase64 == "" {
				zipBase64 = strings.TrimSpace(parsed.ZipBase64)
			}
			if zipBase64 != "" {
				zipBytes, decodeErr := base64.StdEncoding.DecodeString(zipBase64)
				if decodeErr != nil {
					return Artifact{}, fmt.Errorf("decode zip data failed: %w", decodeErr)
				}
				extracted, extractErr := fillArtifactFromZipBytes(&artifact, zipBytes)
				if extractErr != nil {
					return Artifact{}, extractErr
				}
				if extracted {
					return artifact, nil
				}
			}

			exeBase64 := strings.TrimSpace(parsed.ExecutableData)
			if exeBase64 == "" {
				exeBase64 = strings.TrimSpace(parsed.ExecutableBase64)
			}
			if exeBase64 == "" {
				exeBase64 = strings.TrimSpace(parsed.BinaryBase64)
			}
			if exeBase64 == "" {
				exeBase64 = strings.TrimSpace(parsed.OutputBase64)
			}
			if exeBase64 == "" {
				exeBase64 = strings.TrimSpace(parsed.Binary)
			}
			if exeBase64 != "" {
				decoded, decodeErr := base64.StdEncoding.DecodeString(exeBase64)
				if decodeErr != nil {
					return Artifact{}, fmt.Errorf("decode executable data failed: %w", decodeErr)
				}
				artifact.ExecutableData = decoded
				return artifact, nil
			}

			return artifact, nil
		}
	}

	extracted, extractErr := fillArtifactFromZipBytes(&artifact, bodyBytes)
	if extractErr != nil {
		return Artifact{}, extractErr
	}
	if extracted {
		return artifact, nil
	}

	artifact.ExecutableData = bodyBytes
	return artifact, nil
}

func fillArtifactFromMultipartBody(artifact *Artifact, bodyBytes []byte, contentType string) (bool, error) {
	if artifact == nil {
		return false, fmt.Errorf("artifact is nil")
	}
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false, nil
	}
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(mediaType)), "multipart/") {
		return false, nil
	}
	boundary := strings.TrimSpace(params["boundary"])
	if boundary == "" {
		return false, nil
	}

	reader := multipart.NewReader(bytes.NewReader(bodyBytes), boundary)
	var foundAny bool
	for {
		part, nextErr := reader.NextPart()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			return false, fmt.Errorf("read multipart part failed: %w", nextErr)
		}

		data, readErr := io.ReadAll(part)
		part.Close()
		if readErr != nil {
			return false, fmt.Errorf("read multipart body failed: %w", readErr)
		}

		formName := strings.ToLower(strings.TrimSpace(part.FormName()))
		fileName := strings.ToLower(strings.TrimSpace(part.FileName()))
		switch {
		case formName == "artifact" || strings.HasSuffix(fileName, ".exe"):
			artifact.ExecutableData = data
			foundAny = true
		case formName == "build_log":
			artifact.BuildLog = strings.TrimSpace(string(data))
			foundAny = true
		case formName == "meta":
			var meta Meta
			if err := json.Unmarshal(data, &meta); err == nil {
				artifact.Meta = mergeMeta(artifact.Meta, meta)
			}
			foundAny = true
		}
	}

	return foundAny, nil
}

func fillArtifactFromZipBytes(artifact *Artifact, zipBytes []byte) (bool, error) {
	if artifact == nil {
		return false, fmt.Errorf("artifact is nil")
	}
	if len(zipBytes) == 0 {
		return false, nil
	}

	reader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return false, nil
	}

	var foundAny bool
	for _, file := range reader.File {
		name := strings.ToLower(path.Base(file.Name))
		if name == "" {
			continue
		}

		data, readErr := readZipFile(file)
		if readErr != nil {
			return false, fmt.Errorf("read zip file %s failed: %w", file.Name, readErr)
		}

		switch {
		case strings.HasSuffix(name, ".exe"):
			if len(artifact.ExecutableData) == 0 {
				artifact.ExecutableData = data
			}
			foundAny = true
		case name == "meta.json":
			var meta Meta
			if err := json.Unmarshal(data, &meta); err == nil {
				artifact.Meta = mergeMeta(artifact.Meta, meta)
			}
			foundAny = true
		case name == "build.log" || strings.HasSuffix(name, ".log") || strings.HasSuffix(name, ".txt"):
			text := strings.TrimSpace(string(data))
			if text != "" {
				if strings.TrimSpace(artifact.BuildLog) == "" {
					artifact.BuildLog = text
				} else {
					artifact.BuildLog = artifact.BuildLog + "\n" + text
				}
			}
			foundAny = true
		}
	}

	return foundAny, nil
}

func readZipFile(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

func mergeMeta(base Meta, incoming Meta) Meta {
	if strings.TrimSpace(base.CompiledAt) == "" {
		base.CompiledAt = incoming.CompiledAt
	}
	if len(base.CompilerOptions) == 0 {
		base.CompilerOptions = incoming.CompilerOptions
	}
	if base.ReturnCode == 0 {
		base.ReturnCode = incoming.ReturnCode
	}
	if !base.TimedOut {
		base.TimedOut = incoming.TimedOut
	}
	if strings.TrimSpace(base.WorkDir) == "" {
		base.WorkDir = incoming.WorkDir
	}
	return base
}

func parseIntHeader(value string) int {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return 0
	}
	n, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0
	}
	return n
}

func parseBoolHeader(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return false
	}
	v, err := strconv.ParseBool(trimmed)
	if err != nil {
		return false
	}
	return v
}

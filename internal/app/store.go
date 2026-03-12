package app

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"
)

type Store struct {
	mu      sync.RWMutex
	sets    []TrainingSet
	history []ScoreResult
}

func NewStore() *Store {
	return &Store{
		sets: []TrainingSet{
			{
				ID:            "set-001",
				Title:         "login_gate.exe 静态分析",
				Description:   "每次训练只提供一个 exe，题目全部围绕同一二进制的字符串、导入表和成功分支。",
				Difficulty:    "初级",
				CreatedAt:     "2026-03-02T09:00:00.000Z",
				QuestionCount: 3,
				Tags:          []string{"PE", "Strings", "IAT"},
				Artifact: TrainingArtifact{
					FileName:     "login_gate.exe",
					DownloadName: "login_gate.exe",
					Focus:        "定位关键字符串、导入函数和成功提示。",
					Content:      "MZ\nentry: sub_401000\nimports: strcmp, MessageBoxA\nstrings: ADMIN_OVERRIDE, ACCESS_GRANTED",
				},
				Questions: []Question{
					{ID: "q-001", Prompt: "这个 exe 校验的关键字符串是什么？", ExpectedAnswer: "ADMIN_OVERRIDE"},
					{ID: "q-002", Prompt: "它依赖的关键比较函数是什么？", ExpectedAnswer: "strcmp"},
					{ID: "q-003", Prompt: "成功提示字符串是什么？", ExpectedAnswer: "ACCESS_GRANTED"},
				},
			},
			{
				ID:            "set-002",
				Title:         "net_probe.exe 行为定位",
				Description:   "围绕同一个网络探测样本，识别导入函数、端口常量和分支行为。",
				Difficulty:    "中级",
				CreatedAt:     "2026-03-01T14:30:00.000Z",
				QuestionCount: 3,
				Tags:          []string{"CFG", "Branch"},
				Artifact: TrainingArtifact{
					FileName:     "net_probe.exe",
					DownloadName: "net_probe.exe",
					Focus:        "关注 WinSock 调用链和目标端口。",
					Content:      "MZ\nimports: socket, connect, send, recv\nmov ecx, 4444h\ncmp eax, 0\njz short connected",
				},
				Questions: []Question{
					{ID: "q-101", Prompt: "该样本最关键的连接函数是什么？", ExpectedAnswer: "connect"},
					{ID: "q-102", Prompt: "目标端口常量是多少（十六进制）？", ExpectedAnswer: "4444h"},
					{ID: "q-103", Prompt: "跳转到成功分支前要求 eax 等于什么？", ExpectedAnswer: "0"},
				},
			},
		},
		history: []ScoreResult{},
	}
}

func (s *Store) ListTrainingSets(sortBy string) []TrainingSet {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := cloneTrainingSets(s.sets)
	switch sortBy {
	case "title":
		slices.SortFunc(items, func(a, b TrainingSet) int {
			return strings.Compare(a.Title, b.Title)
		})
	default:
		slices.SortFunc(items, func(a, b TrainingSet) int {
			return strings.Compare(b.CreatedAt, a.CreatedAt)
		})
	}
	return items
}

func (s *Store) GetTrainingSet(id string) (TrainingSet, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, item := range s.sets {
		if item.ID == id {
			return cloneTrainingSet(item), true
		}
	}
	return TrainingSet{}, false
}

func (s *Store) CreateTrainingSet(input CreateTrainingInput) TrainingSet {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	id := fmt.Sprintf("set-%d", now.UnixMilli())
	title := strings.TrimSpace(input.Title)
	description := strings.TrimSpace(input.Description)
	if description == "" {
		description = "自定义训练题集"
	}
	normalized := strings.ToUpper(strings.ReplaceAll(title, " ", "_"))
	artifact := TrainingArtifact{
		FileName:     fmt.Sprintf("%s-challenge.exe", title),
		DownloadName: fmt.Sprintf("%s-challenge.exe", title),
		Focus:        "围绕单个 exe 完成字符串、导入表和控制流定位。",
		Content:      fmt.Sprintf("MZ\nentry: sub_401000\nimports: CreateFileA, ReadFile, CloseHandle\nstrings: %s_KEY, %s_SUCCESS\ncmp eax, 7", normalized, normalized),
	}
	questions := []Question{
		{ID: id + "-q1", Prompt: "识别这个 exe 中的关键校验字符串。", ExpectedAnswer: normalized + "_KEY"},
		{ID: id + "-q2", Prompt: "识别主要文件读取 API。", ExpectedAnswer: "CreateFileA"},
		{ID: id + "-q3", Prompt: "分支比较要求 eax 等于多少？", ExpectedAnswer: "7"},
	}

	created := TrainingSet{
		ID:            id,
		Title:         title,
		Description:   description,
		Difficulty:    normalizeDifficulty(input.Difficulty),
		CreatedAt:     now.Format(time.RFC3339),
		QuestionCount: len(questions),
		Tags:          []string{"Custom", "Practice"},
		Artifact:      artifact,
		Questions:     questions,
	}

	s.sets = append([]TrainingSet{created}, s.sets...)
	return cloneTrainingSet(created)
}

func (s *Store) GenerateArtifact(id string) (TrainingSet, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for idx := range s.sets {
		if s.sets[idx].ID != id {
			continue
		}
		base := strings.TrimPrefix(s.sets[idx].Artifact.Content, "MZ\n")
		s.sets[idx].Artifact.DownloadName = s.sets[idx].Artifact.FileName
		s.sets[idx].Artifact.Content = fmt.Sprintf("MZ\ngeneratedAt: %s\n%s", time.Now().UTC().Format(time.RFC3339), base)
		return cloneTrainingSet(s.sets[idx]), true
	}

	return TrainingSet{}, false
}

func (s *Store) ScoreQuestion(id string, questionID string, answer string) (QuestionScoreResult, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var dataset *TrainingSet
	for idx := range s.sets {
		if s.sets[idx].ID == id {
			dataset = &s.sets[idx]
			break
		}
	}
	if dataset == nil {
		return QuestionScoreResult{}, false
	}

	var current *Question
	for idx := range dataset.Questions {
		if dataset.Questions[idx].ID == questionID {
			current = &dataset.Questions[idx]
			break
		}
	}
	if current == nil {
		return QuestionScoreResult{}, false
	}

	submitted := strings.TrimSpace(answer)
	correct := submitted == current.ExpectedAnswer
	result := QuestionScoreResult{
		DatasetID:       dataset.ID,
		DatasetTitle:    dataset.Title,
		QuestionID:      current.ID,
		Prompt:          current.Prompt,
		SubmittedAt:     time.Now().UTC().Format(time.RFC3339),
		SubmittedAnswer: submitted,
		ExpectedAnswer:  current.ExpectedAnswer,
		Correct:         correct,
	}

	historyItem := ScoreResult{
		DatasetID:    dataset.ID,
		DatasetTitle: dataset.Title,
		SubmittedAt:  result.SubmittedAt,
		Score:        0,
		Total:        1,
		Passed:       correct,
		CorrectCount: 0,
		WrongCount:   0,
		Mistakes:     []MistakeDetail{},
	}
	if correct {
		historyItem.Score = 1
		historyItem.CorrectCount = 1
	} else {
		historyItem.WrongCount = 1
		historyItem.Mistakes = []MistakeDetail{{
			QuestionID:      current.ID,
			Prompt:          current.Prompt,
			ExpectedAnswer:  current.ExpectedAnswer,
			SubmittedAnswer: submitted,
		}}
	}

	s.history = append([]ScoreResult{historyItem}, s.history...)
	return result, true
}

func (s *Store) ListHistory() []ScoreResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]ScoreResult, 0, len(s.history))
	for _, item := range s.history {
		items = append(items, cloneScoreResult(item))
	}
	return items
}

func normalizeDifficulty(value string) string {
	switch value {
	case "中级", "高级":
		return value
	default:
		return "初级"
	}
}

func cloneTrainingSets(input []TrainingSet) []TrainingSet {
	items := make([]TrainingSet, 0, len(input))
	for _, item := range input {
		items = append(items, cloneTrainingSet(item))
	}
	return items
}

func cloneTrainingSet(item TrainingSet) TrainingSet {
	cloned := item
	cloned.Tags = append([]string(nil), item.Tags...)
	cloned.Artifact = item.Artifact
	cloned.Questions = append([]Question(nil), item.Questions...)
	return cloned
}

func cloneScoreResult(item ScoreResult) ScoreResult {
	cloned := item
	cloned.Mistakes = append([]MistakeDetail(nil), item.Mistakes...)
	return cloned
}

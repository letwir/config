package prompt

import (
	"fmt"
	"strings"
)

// DefaultSystemPrompt はデフォルトのシステムプロンプトですわ。
// 旦那様が後で手軽に微調整できるように、この変数の冒頭に配置いたしました。
// %s には対象のプログラミング言語（例: "python3.13"）が入ります。
var DefaultSystemPrompt = `You are an expert %s programmer.
Your task is to generate ONLY valid, runnable %s code based on the user instruction.
Do NOT output any markdown code blocks, do NOT wrap the code in Triple Backticks (like ` + "```" + `), do NOT write any introduction or explanation.
Start your response immediately with the source code, and end it as soon as the code is finished.
Every character in your response must be a part of the executable source code.`

// BuildSystemPrompt はシステムプロンプトを言語名に合わせて構築しますわ。
func BuildSystemPrompt(language string) string {
	return fmt.Sprintf(DefaultSystemPrompt, language, language)
}

// BuildUserPrompt はユーザーの入力を <instruction> XMLタグでカプセル化しますわ。
func BuildUserPrompt(promptText string) string {
	return fmt.Sprintf("<instruction>\n%s\n</instruction>", strings.TrimSpace(promptText))
}

// CleanCode はLLMの出力からマークダウンコードフェンス（例: ```python ... ```）を除去し、
// 純粋なコードのみを抽出するフィルタリング処理を行いますわ。
func CleanCode(rawOutput string) string {
	trimmed := strings.TrimSpace(rawOutput)
	lines := strings.Split(trimmed, "\n")
	if len(lines) == 0 {
		return ""
	}

	// 開始フェンス（例: ```python）のチェック
	hasStartFence := false
	startIdx := 0
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "```") {
			hasStartFence = true
			startIdx = i + 1
			break
		}
	}

	if !hasStartFence {
		return trimmed
	}

	// 終了フェンス（```）のチェック
	endIdx := len(lines)
	for i := len(lines) - 1; i >= startIdx; i-- {
		trimmedLine := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trimmedLine, "```") {
			endIdx = i
			break
		}
	}

	// フェンスの中身だけを結合して返しますわ
	if startIdx < endIdx {
		return strings.Join(lines[startIdx:endIdx], "\n")
	}

	return trimmed
}

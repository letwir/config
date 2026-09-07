package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"llama2coder/client"
	"llama2coder/prompt"

	"github.com/spf13/cobra"
)

var (
	promptText string
	language   string
	baseURL    string
	modelName  string
	stream     bool
	temp       float64
	maxTokens  int
	sysPrompt  string
)

// StreamFilter はストリーミング出力の過程で、
// コードフェンス（```や```pythonなど）を検出してStdoutに出力しないようにするラッパーですわ。
type StreamFilter struct {
	w            io.Writer
	buf          strings.Builder
	firstLineRaw string
	skippedStart bool
	inFence      bool
}

func NewStreamFilter(w io.Writer) *StreamFilter {
	return &StreamFilter{w: w}
}

func (sf *StreamFilter) Write(p []byte) (n int, err error) {
	// 簡易的なフィルタリング:
	// まだ開始フェンスの処理が終わっていない場合
	if !sf.skippedStart {
		sf.buf.Write(p)
		content := sf.buf.String()

		// 改行が含まれるか、または十分長いかチェックしますわ
		if strings.Contains(content, "\n") {
			lines := strings.Split(content, "\n")
			sf.skippedStart = true
			sf.buf.Reset()

			firstLine := strings.TrimSpace(lines[0])
			startIdx := 0

			// 最初の行が ``` で始まる場合はスキップしますわ
			if strings.HasPrefix(firstLine, "```") {
				sf.inFence = true
				startIdx = 1
			}

			// 残りの行を出力バッファまたは書き込み先に流しますわ
			for i := startIdx; i < len(lines); i++ {
				line := lines[i]
				// 最終要素で改行が終わっていない場合はバッファに戻しますわ
				if i == len(lines)-1 && !strings.HasSuffix(content, "\n") {
					sf.buf.WriteString(line)
				} else {
					// 終了フェンスの簡易チェック
					trimmed := strings.TrimSpace(line)
					if sf.inFence && strings.HasPrefix(trimmed, "```") {
						sf.inFence = false
						continue
					}
					_, err = io.WriteString(sf.w, line+"\n")
					if err != nil {
						return len(p), err
					}
				}
			}
		}
		return len(p), nil
	}

	// 既に最初の行の処理が完了している場合
	// 終了フェンス（```）が混ざるのを防ぐため、バッファリングしつつ流しますわ
	text := string(p)
	// ``` が含まれている場合は、その周辺をフィルタリングします
	if strings.Contains(text, "```") {
		// ストリーム中の ``` を消去しますわ（最も安全で単純なアプローチです）
		lines := strings.Split(text, "```")
		for _, segment := range lines {
			// プログラミング言語識別子（例: python, go）が単独で残るのを防ぐため、
			// フェンス直後の改行コードまでの単語を簡易的に除去します
			cleanedSegment := segment
			if sf.inFence {
				// フェンス開始直後の可能性を考慮し、言語名＋改行を削りますわ
				for _, lang := range []string{"python", "go", "rust", "javascript", "typescript", "cpp", "c", "bash", "sh", "html", "css"} {
					if strings.HasPrefix(strings.ToLower(cleanedSegment), lang) {
						cleanedSegment = cleanedSegment[len(lang):]
						break
					}
				}
				sf.inFence = false
			} else {
				sf.inFence = true
			}
			_, err = io.WriteString(sf.w, cleanedSegment)
			if err != nil {
				return len(p), err
			}
		}
	} else {
		_, err = io.WriteString(sf.w, text)
		if err != nil {
			return len(p), err
		}
	}

	return len(p), nil
}

// Flush はバッファに残ったデータを出力しますわ。
func (sf *StreamFilter) Flush() {
	rem := sf.buf.String()
	if rem != "" {
		trimmed := strings.TrimSpace(rem)
		if !strings.HasPrefix(trimmed, "```") {
			_, _ = io.WriteString(sf.w, rem)
		}
	}
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "coder",
		Short: "llama2coder は llama.cpp を用いてコードのみを出力する推論ツールですわ",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			// システムプロンプトの決定
			// --system が指定されていればそれを使用し、空であればデフォルトから自動生成しますわ
			var finalSysPrompt string
			if sysPrompt != "" {
				finalSysPrompt = sysPrompt
			} else {
				finalSysPrompt = prompt.BuildSystemPrompt(language)
			}

			// ユーザープロンプトの構築 (<instruction> タグで包む)
			userPrompt := prompt.BuildUserPrompt(promptText)

			cli := client.NewLlamaClient(baseURL)

			if stream {
				// ストリーミング出力
				filter := NewStreamFilter(os.Stdout)
				err := cli.Stream(ctx, modelName, finalSysPrompt, userPrompt, temp, maxTokens, filter)
				filter.Flush()
				if err != nil {
					return fmt.Errorf("streaming error: %w", err)
				}
			} else {
				// 一括取得してクレンジングしてから出力
				output, err := cli.Generate(ctx, modelName, finalSysPrompt, userPrompt, temp, maxTokens)
				if err != nil {
					return fmt.Errorf("generation error: %w", err)
				}

				cleaned := prompt.CleanCode(output)
				fmt.Print(cleaned)
			}

			return nil
		},
	}

	// コマンドライン引数（フラグ）の設定ですわ
	rootCmd.Flags().StringVarP(&promptText, "prompt", "p", "", "コード生成の命令文 (必須ですわ)")
	rootCmd.Flags().StringVarP(&language, "language", "l", "python3", "対象のプログラミング言語")
	rootCmd.Flags().StringVarP(&baseURL, "url", "u", "http://localhost:8080", "llama-server のベース URL")
	rootCmd.Flags().StringVar(&modelName, "model", "default", "使用するモデル名")
	rootCmd.Flags().BoolVarP(&stream, "stream", "s", false, "ストリーミング出力を行うかどうか")
	rootCmd.Flags().Float64Var(&temp, "temperature", 0.2, "生成時の温度パラメータ")
	rootCmd.Flags().IntVar(&maxTokens, "max-tokens", -1, "最大生成トークン数 (-1でサーバーデフォルト)")
	rootCmd.Flags().StringVar(&sysPrompt, "system", "", "カスタムシステムプロンプト (指定するとデフォルトを上書きしますわ)")

	_ = rootCmd.MarkFlagRequired("prompt")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

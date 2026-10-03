# Walkthrough - Resolve documentation LRF effect mismatch

## 概要

Level 2の効果不一致を、読み取り専用の停止規則とVCS操作の承認規則に分けて解消しました。

## 成果物一覧

- `rules/documentation.lrf`: `doc.secret-found` を停止・報告に限定し、履歴改変・追跡ファイル削除の承認条件を `VCS_WRITE` の別レコードに分離。

## 検証結果

- harness-lint Level 1: エラー0、警告0。
- harness-lint Level 2: エラー0、警告0。
- `git diff --check -- rules/documentation.lrf`: 成功。

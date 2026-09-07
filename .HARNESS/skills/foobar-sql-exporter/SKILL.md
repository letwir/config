---
name: foobar-sql-exporter
description: Execute SQL queries against the foobar2000 PostgreSQL database to generate M3U playlists, extract metadata, or inspect configurations.
---
[SPR/PIDGIN::ρ→max] 🎵m3u 📊json 📄text 🗄️schema
🎵 プレイリスト m3u ⊢ Query⊗Output ⇒ 高評価/直近再生M3U生成 | `& $bin --db $env:FOOBAR_DB_URL -f m3u -o <out.m3u> -q "<SQL>"`
📊 データ json ⊢ Query ⇒ トラック・再生統計・メタデータJSON出力 | `& $bin --db $env:FOOBAR_DB_URL -f json -q "<SQL>"`
📄 テキスト text ⊢ Query ⇒ 設定値・タグ・統計テキスト出力 | `& $bin --db $env:FOOBAR_DB_URL -f text -q "<SQL>"`
🗄️ スキーマ schema ⊢ ∅ ⇒ raw.foobar2000_{tracks,stats,playlists,playlist_tracks,configs}

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=foobar-sql-exporter
R|bin.resolve|task:foobar-sql-exporter|MUST|RO_LOCAL|bin=first-existing($PSScriptRoot/bin/fb2k-sql.exe,$env:USERPROFILE/.harness/skills/foobar-sql-exporter/bin/fb2k-sql.exe); missing=>STOP
R|use.query|task:foobar-sql-exporter|MUST|RO_LOCAL|cmds=m3u,json,text; effect=read-only-inspection; default-format=text
R|use.export|task:foobar-sql-exporter&flag:output|MUST|LW_SCOPE|write=requested-output-file-only
R|guard.cred|task:foobar-sql-exporter|MUST_NOT|CRED|env=FOOBAR_DB_URL,DATABASE_URL; deny=credentials,raw-connection-uri; allow=task-scoped-queries
R|guard.truth|task:foobar-sql-exporter|MUST|RO_LOCAL|unverified=fact-forbidden; source=foobar2000-postgres-db

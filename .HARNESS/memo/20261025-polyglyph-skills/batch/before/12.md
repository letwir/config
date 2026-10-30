---
name: foobar-sql-exporter
description: Execute SQL queries against the foobar2000 PostgreSQL database to generate M3U playlists, extract metadata, or inspect configurations.
---
@meta=text/POLYGLYPH×SIGMA/1
@codec=nearest_ancestor(POLYGLYPH.md);read_before_body=true
@authority=SIGMA/1+routed-LRF-only;PG=representation-only
@decode_failure=STOP;unknown=preserve;action=blocked;ask=true

[SPR/PIDGIN::ρ→max] 🔎task-scoped foobar2000 PostgreSQL inspection ⊢ Resolve first-existing `$bin` only; missing executable⇒STOP. Supply SQL query; permitted formats={m3u,json,text}; default=text. Query operation is read-only; ¬implicit DB write. Output effect limited by format and request: m3u→requested output file only; json/text→emit data. Never expose credentials or raw connection URI. Unverified DB facts⇒do not assert.
🎵m3u: Query⊗requested output→playlist (high-rated/recently-played) | `$bin --db $FOOBAR_DB_URL -f m3u -o <out.m3u> -q "<SQL>"`
📊json: Query→track/play-count/statistics/metadata JSON | `$bin --db $FOOBAR_DB_URL -f json -q "<SQL>"`
📄text: Query→configuration values/tags/statistics text | `$bin --db $FOOBAR_DB_URL -f text -q "<SQL>"`
🗄️schema: visible schema boundary=raw.foobar2000_{tracks,stats,playlists,playlist_tracks,configs}; ∅≠permission to alter.

#SIGMA LRF/1
@lrf=1|aud=GPT-5.6|scope=foobar-sql-exporter
R|bin.resolve|task:foobar-sql-exporter|MUST|RO_LOCAL|bin=first-existing(~/.harness/skills/foobar-sql-exporter/bin/fb2k-sql,~/.harness/skills/foobar-sql-exporter/bin/fb2k-sql.exe,./bin/fb2k-sql,./bin/fb2k-sql.exe); missing=>STOP
R|use.query|task:foobar-sql-exporter|MUST|RO_LOCAL|cmds=m3u,json,text; effect=read-only-inspection; default-format=text
R|use.export|task:foobar-sql-exporter&flag:output|MUST|LW_SCOPE|write=requested-output-file-only
R|guard.cred|task:foobar-sql-exporter|MUST_NOT|CRED|env=FOOBAR_DB_URL,DATABASE_URL; deny=credentials,raw-connection-uri; allow=task-scoped-queries
R|guard.truth|task:foobar-sql-exporter|MUST|RO_LOCAL|unverified=fact-forbidden; source=foobar2000-postgres-db

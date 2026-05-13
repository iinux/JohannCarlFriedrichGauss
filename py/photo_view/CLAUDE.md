# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Run

```bash
python3 photo_view.py
```

Flask dev server on `127.0.0.1:5000`. Uncomment the `app.run(host=...)` line in `photo_view.py` to expose on LAN. There is no requirements file — the app needs `flask`, `werkzeug`, `pdf2image` (Python), and `ffmpeg` + `ffprobe` (system binaries) on PATH. `pdf2image` additionally needs `poppler` installed.

No tests, no lint config. The only smoke check is `python3 -m py_compile *.py`.

## Layout

The Flask app is split into blueprints registered in `photo_view.py`. Each `views_*.py` owns its routes:

- `views_index.py` — `/index`, `/index/image/<filename>`, `/index/random` (random video screenshot via ffmpeg)
- `views_words.py` — `/words` (random lines from `dict2.txt`)
- `views_rtsp.py` — `/rtsp/ch00`, `/rtsp/ch01`, `/rtsp/stream/<channel>`, `/rtsp/feed/<channel>`
- `views_upload.py` — `/upload`, `/upload/pdf`, `/upload/delete/<filename>`, `/upload/download/<filename>`
- `paths.py` — shared directory constants and `ALLOWED_EXTENSIONS`
- `my_config.py` — RTSP URLs (gitignored, must exist locally; expects `rtsp_url1`, `rtsp_url2`)
- `sf.py` / `sf_data.py` — standalone scraper, not part of the Flask app

Working directories (`mp4/`, `img/`, `flv/`, `upload/`, `screenshot/`, `download_tag/`) are read/written at relative path `.` — the server must be run from the project root. `upload/` and `screenshot/` are auto-created on first write; the others must exist.

Templates in `templates/` use **hardcoded URL paths** (no `url_for`). So renaming a Flask endpoint is safe, but changing a URL path (`/upload` → `/files` etc.) requires updating the corresponding `.html`. The index page also links to `/download/...` paths that are **not** Flask routes — those rely on an upstream reverse proxy (nginx) to serve filesystem directories.

## RTSP streaming architecture

`views_rtsp.py` is the most complex piece. Key invariants:

- **One ffmpeg per channel, on demand.** A background thread (`_stream_reader`) runs a single ffmpeg subprocess that emits MJPEG to stdout. The latest JPEG frame is parsed by scanning for SOI/EOI markers (`\xff\xd8` / `\xff\xd9`) and stored in `stream_cache[channel]`. All HTTP requests read from this cache.
- **Reference counting.** `generate_frames` (used by `/rtsp/stream/<channel>`) increments `stream_viewers[channel]` on entry and decrements in `finally`. `/rtsp/feed/<channel>` (single-image polling) updates `stream_last_feed[channel]` for a `STREAM_FEED_GRACE = 30s` grace window.
- **Idle shutdown.** The reader thread checks `_stream_active()` (1) before launching ffmpeg, (2) every 1s inside the read loop, and (3) after ffmpeg exits. When inactive, it atomically pops itself from `stream_threads` and `stream_cache` under `stream_lock` before returning. The next request triggers a fresh thread + ffmpeg.
- **stderr must be drained.** `_drain_stderr` runs in its own daemon thread to keep the 64KB pipe from blocking ffmpeg.

If touching this code, preserve the lock pattern: `_stream_active` check and `stream_threads.pop` happen in the same `with stream_lock:` block so a new viewer can't race with a thread that's about to exit.

The Flask dev server is single-threaded by default. Multiple concurrent MJPEG clients (or a stream client + upload at once) need `app.run(threaded=True)` or a real WSGI server.

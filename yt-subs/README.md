# yt-subs

Extract YouTube video subtitles as JSON with timestamps.

## Features

- ✅ Extract subtitles from any YouTube video
- ✅ Works with video IDs or full URLs
- ✅ Auto-fallback to auto-generated captions if manual unavailable
- ✅ JSON output with timestamps (text, start, duration)
- ✅ Pretty-print option for readable output
- ✅ Save to file or stdout
- ✅ User-friendly error messages
- ✅ Modern CLI with Typer

## Installation

```bash
# Navigate to the yt-subs directory
cd yt-subs

# Install dependencies with UV
uv sync
```

## Usage

### Basic Usage

```bash
# Extract and display subtitles
uv run yt-subs dQw4w9WgXcQ

# Extract and save to file
uv run yt-subs dQw4w9WgXcQ -o subtitles.json

# Extract with pretty-printed JSON
uv run yt-subs dQw4w9WgXcQ -p

# Extract and suppress output (quiet mode)
uv run yt-subs dQw4w9WgXcQ -q
```

### Using Full URLs

```bash
# Extract from full URL
uv run yt-subs "https://youtube.com/watch?v=dQw4w9WgXcQ"

# Extract from shortened URL
uv run yt-subs "https://youtu.be/dQw4w9WgXcQ"
```

## Output Format

The tool outputs JSON in this format:

```json
[
  {
    "text": "Welcome to this video",
    "start": 0.0,
    "duration": 3.5
  },
  {
    "text": "In this tutorial we will...",
    "start": 3.5,
    "duration": 4.2
  }
]
```

## Error Handling

- **Invalid video ID**: Shows clear error message
- **Video unavailable**: Indicates if video is private/unavailable
- **No subtitles**: Tells you if no captions are available (manual or auto)
- **Network errors**: Handles connection issues gracefully

## Requirements

- Python 3.11+
- UV (optional but recommended for dependency management)
- Internet connection for YouTube API access

## Dependencies

- `youtube-transcript-api` - Core subtitle extraction
- `typer` - Modern CLI framework

## Version

Current version: 0.1.0
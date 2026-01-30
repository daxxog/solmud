"""
Core YouTube subtitle extraction logic.

Handles video ID validation, subtitle extraction, and error handling.
"""

import re
from typing import List, Dict, Any
from pathlib import Path

from youtube_transcript_api import YouTubeTranscriptApi, NoTranscriptFound, VideoUnavailable, IpBlocked, TranscriptsDisabled


def validate_video_id(video_id: str) -> str:
    """
    Validate and normalize YouTube video ID.
    
    Accepts:
    - Full YouTube URL (https://youtube.com/watch?v=VIDEO_ID)
    - Shortened URL (https://youtu.be/VIDEO_ID)
    - Raw video ID
    
    Returns normalized video ID or raises ValueError.
    """
    # Pattern for full YouTube URL
    full_url_pattern = r"youtube\.com/watch\?v=([\w-]+)"
    # Pattern for shortened YouTube URL
    short_url_pattern = r"youtu\.be/([\w-]+)"
    # Pattern for raw video ID (11 characters)
    video_id_pattern = r"^[\w-]{11}$"
    
    # Check if it's a full URL
    full_url_match = re.search(full_url_pattern, video_id)
    if full_url_match:
        return full_url_match.group(1)
    
    # Check if it's a shortened URL
    short_url_match = re.search(short_url_pattern, video_id)
    if short_url_match:
        return short_url_match.group(1)
    
    # Check if it's a raw video ID
    if re.match(video_id_pattern, video_id):
        return video_id
    
    # If none match, raise error
    raise ValueError(
        f"Invalid video ID or URL format: {video_id}. "
        f"Expected 11-character ID or valid YouTube URL."
    )


def _snippets_to_dicts(snippets: List[Any]) -> List[Dict[str, Any]]:
    """Convert FetchedTranscriptSnippet objects to dictionaries."""
    return [
        {
            "text": snippet.text,
            "start": snippet.start,
            "duration": snippet.duration
        }
        for snippet in snippets
    ]


def extract_subtitles(video_id: str) -> List[Dict[str, Any]]:
    """
    Extract subtitles from YouTube video.
    
    Args:
        video_id: YouTube video ID or URL
    
    Returns:
        List of subtitle segments as dictionaries with text, start time, and duration
    
    Raises:
        ValueError: If video ID is invalid
        VideoUnavailable: If video is private/unavailable
        NoTranscriptFound: If no subtitles found
        IpBlocked: If YouTube blocks the request
        TranscriptsDisabled: If video has transcripts disabled
    """
    # Validate and normalize video ID
    normalized_id = validate_video_id(video_id)
    
    # Try to get manual subtitles first
    try:
        transcript = YouTubeTranscriptApi().fetch(normalized_id, ['en'])
        return _snippets_to_dicts(transcript)
    except NoTranscriptFound:
        # Fallback to auto-generated subtitles
        try:
            transcript_list = YouTubeTranscriptApi().list(normalized_id)
            # Get the first available transcript (usually English)
            transcript = transcript_list.find_transcript(['en'])
            if transcript is None:
                # Get first available if English not found
                transcript = list(transcript_list)[0]
            # Get translated transcript if needed
            transcript = transcript.translate('en').fetch()
            return _snippets_to_dicts(transcript)
        except Exception as e:
            raise NoTranscriptFound(
                f"No subtitles available for video {normalized_id}: {str(e)}"
            ) from e


def format_subtitles(transcript: List[Dict[str, Any]]) -> str:
    """
    Format subtitles as pretty-printed JSON.
    
    Args:
        transcript: List of subtitle segments
    
    Returns:
        JSON string representation of subtitles
    """
    import json
    return json.dumps(transcript, indent=2, ensure_ascii=False)


def save_subtitles(transcript: List[Dict[str, Any]], output_path: Path) -> None:
    """
    Save subtitles to file.
    
    Args:
        transcript: List of subtitle segments
        output_path: Path to save file
    """
    output_path = Path(output_path)
    output_path.parent.mkdir(parents=True, exist_ok=True)
    
    json_content = format_subtitles(transcript)
    
    with open(output_path, 'w', encoding='utf-8') as f:
        f.write(json_content)
    
    print(f"Subtitles saved to: {output_path}")
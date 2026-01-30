"""
Command-line interface for YouTube subtitle extraction.

Uses Typer to create a modern CLI with automatic help generation.
"""

import json
from pathlib import Path
from typing import Optional

import typer

from yt_subs.extractor import (
    extract_subtitles,
    save_subtitles,
    validate_video_id,
    NoTranscriptFound,
    VideoUnavailable,
    IpBlocked,
    TranscriptsDisabled,
)

app = typer.Typer(
    name="yt-subs",
    add_completion=False,
    help="Extract YouTube video subtitles as JSON with timestamps",
)


@app.command()
def extract(
    video_id: str,
    output: Optional[Path] = typer.Option(None, "-o", "--output", help="Save subtitles to file"),
    pretty: bool = typer.Option(False, "-p", "--pretty", help="Pretty-print JSON output"),
    quiet: bool = typer.Option(False, "-q", "--quiet", help="Suppress non-error output"),
) -> None:
    """
    Extract YouTube video subtitles as JSON.
    
    VIDEO_ID can be:
    - A YouTube video ID (11 characters)
    - A full YouTube URL
    - A shortened YouTube URL
    """
    try:
        # Validate video ID format
        normalized_id = validate_video_id(video_id)
        
        if not quiet:
            typer.echo(f"Extracting subtitles for: {normalized_id}")
        
        # Extract subtitles
        transcript = extract_subtitles(normalized_id)
        
        if not quiet:
            typer.echo(f"Found {len(transcript)} subtitle segments")
        
        # Format output
        json_output = json.dumps(transcript, indent=2 if pretty else None, ensure_ascii=False)
        
        # Save or print
        if output:
            save_subtitles(transcript, output)
            if not quiet:
                typer.echo(f"Subtitles saved to: {output}")
        else:
            typer.echo(json_output)
            if not quiet:
                typer.echo(f"\nTotal segments: {len(transcript)}")
                
    except ValueError as e:
        typer.echo(f"Error: {str(e)}", err=True)
        raise typer.Exit(1)
    
    except VideoUnavailable:
        typer.echo("Error: Video unavailable (may be private or deleted)", err=True)
        raise typer.Exit(1)
    
    except NoTranscriptFound:
        typer.echo("Error: No subtitles available for this video", err=True)
        raise typer.Exit(1)
    
    except IpBlocked:
        typer.echo(
            "Error: YouTube is blocking requests from your IP. "
            "This may be due to rate limiting or using a cloud provider IP.",
            err=True
        )
        raise typer.Exit(1)
    
    except TranscriptsDisabled:
        typer.echo("Error: This video has subtitles disabled", err=True)
        raise typer.Exit(1)
    
    except Exception as e:
        typer.echo(f"Unexpected error: {str(e)}", err=True)
        raise typer.Exit(1)


@app.command(name="version")
def version_command():
    """Show version information."""
    from yt_subs import __version__
    typer.echo(f"yt-subs version {__version__}")


if __name__ == "__main__":
    app()
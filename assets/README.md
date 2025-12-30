# Assets Directory

This directory contains game assets (sprites, sounds, fonts).

## Structure

```
assets/
├── sprites/    # PNG textures and sprite sheets
├── sounds/     # Audio files (WAV, OGG, MP3)
└── fonts/      # Font files (TTF)
```

## Asset Guidelines

### Sprites
- Use PNG format with transparency
- Consistent pixel art style (e.g., 16x16 tiles)
- Consider sprite atlases for performance
- Name descriptively: `player_idle.png`, `ore_gold.png`, etc.

### Sounds
- Use OGG or WAV format
- Normalize audio levels
- Keep file sizes reasonable
- Name by action: `drill_start.wav`, `ore_collect.wav`, etc.

### Fonts
- Use TTF format
- Include licensing information
- Keep a default fallback font

## Licensing

Make sure all assets are:
- Created by you, or
- Licensed for commercial use, or
- Public domain / CC0

Document sources and licenses for third-party assets.

## Placeholder Assets

During development, consider using:
- Colored rectangles for sprites
- Raylib's default font
- Simple beep/tone sounds

Replace with proper assets as development progresses.

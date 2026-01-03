# dmsg

Type-safe Discord message component builder for Components v2.

## Installation

```bash
go get github.com/thomasgtaylor/dmsg
```

## Requirements

- Go 1.21 or later
- [discordgo](https://github.com/bwmarrin/discordgo) - Discord API wrapper for Go

## Basic Usage

```go
import (
    "github.com/bwmarrin/discordgo"
    "github.com/thomasgtaylor/dmsg"
)

response := dmsg.Response(
    dmsg.Container(
        dmsg.AccentColor(5763719), // green
        dmsg.ContainerSection(
            dmsg.SectionText("## Hello, World!"),
        ),
    ),
)
```

## Features

- **Type-safe**: Options are typed per component - can't pass wrong option to wrong component
- **Composable**: Components nest naturally, reads top-to-bottom like the UI
- **Clean API**: Functional options pattern with intuitive naming
- **Direct**: Uses discordgo types directly, no conversion needed

## Quick Start

### Simple Error Response

```go
dmsg.Response(
    dmsg.Container(
        dmsg.AccentColor(16740978), // red
        dmsg.ContainerSection(
            dmsg.SectionText("## Error\n\nYou don't have enough coins"),
            dmsg.SectionAccessory(
                dmsg.Thumbnail(errorImageURL, "Error icon"),
            ),
        ),
        dmsg.ContainerSeparator(),
    ),
)
```

### Success with Buttons

```go
dmsg.Response(
    dmsg.Container(
        dmsg.AccentColor(5763719), // green
        dmsg.ContainerSection(
            dmsg.SectionText("## Success!\n\nYou won **1,000 coins**"),
            dmsg.SectionAccessory(
                dmsg.Thumbnail(trophyURL, trophyDesc),
            ),
        ),
        dmsg.ContainerSeparator(),
        dmsg.ContainerActionRow(
            dmsg.Button("Play Again", "play_again", dmsg.Style(dmsg.Success)),
            dmsg.Button("Quit", "quit", dmsg.Style(dmsg.Secondary)),
        ),
    ),
)
```

### Ephemeral Response

```go
dmsg.Ephemeral(
    dmsg.Container(
        dmsg.AccentColor(14197815), // gold
        dmsg.ContainerSection(
            dmsg.SectionText("This message is only visible to you"),
        ),
    ),
)
```

### Update Message

```go
dmsg.Update(
    dmsg.Container(
        dmsg.AccentColor(16777215), // white
        dmsg.ContainerSection(
            dmsg.SectionText("## Game Closed\n\nThanks for playing!"),
        ),
    ),
)
```

## Component Reference

### Entry Points

- `Response(components ...Component)` - Standard response
- `Ephemeral(components ...Component)` - Ephemeral response
- `Update(components ...Component)` - Update message response

### Top-Level Components

Components that can be passed directly to `Response`, `Ephemeral`, or `Update`:

- `Container()` - Container with optional accent color and nested components
- `ActionRow()` - Action row with buttons (top-level)
- `TextDisplay()` - Text display component (top-level)
- `Separator()` - Separator/divider component (top-level)

### Container Components

Components for use inside `Container()`:

**ContainerSection** - Section with text and optional accessory
```go
dmsg.ContainerSection(
    dmsg.SectionText("content"),
    dmsg.SectionAccessory(component),
)
```

**ContainerSeparator** - Separator/divider
```go
dmsg.ContainerSeparator(
    dmsg.WithDivider(true),
    dmsg.Spacing(discordgo.SeparatorSpacingSizeSmall),
)
```

**ContainerActionRow** - Action row with buttons
```go
dmsg.ContainerActionRow(
    dmsg.Button(...),
    dmsg.Button(...),
)
```

**ContainerFile** - File attachment
```go
dmsg.ContainerFile(url,
    dmsg.Spoiler(),
)
```

**ContainerGallery** - Media gallery
```go
dmsg.ContainerGallery(
    dmsg.Media(url, description, spoiler),
    dmsg.Media(url, description, spoiler),
)
```

### Section Components

Components for use inside `ContainerSection()`:

**SectionText** - Text display
```go
dmsg.SectionText("## Markdown content")
```

**SectionAccessory** - Accessory (thumbnail or button)
```go
dmsg.SectionAccessory(
    dmsg.Thumbnail(url, description),
)
```

### Standalone Components

**Thumbnail** - Can be used as accessory or standalone
```go
dmsg.Thumbnail(url, description,
    dmsg.Spoiler(),
)
```

**Media** - Helper for gallery items
```go
dmsg.Media(url, description, spoiler)
```

### Interactive Components

**Button**
```go
dmsg.Button(label, customID,
    dmsg.Style(dmsg.Primary),   // or Secondary, Success, Danger
    dmsg.Emoji(&discordgo.ComponentEmoji{Name: "🎉"}),
    dmsg.Disabled(),
)
```

**Link Button**
```go
dmsg.LinkButton(label, url,
    dmsg.Emoji(&discordgo.ComponentEmoji{Name: "❤️"}),
)
```

## Type Safety

Options are typed per component. The compiler prevents mistakes:

```go
// ✅ This compiles
dmsg.Container(
    dmsg.AccentColor(red),       // ContainerOption
    dmsg.ContainerSection(...),  // Also ContainerOption
)

// ❌ This won't compile
dmsg.Container(
    dmsg.Style(dmsg.Primary),  // ButtonOption - wrong type!
)

// ✅ Clear naming shows where components belong
dmsg.ContainerSection(        // For use in Container
    dmsg.SectionText(...),     // For use in Section
)
```

## Design Patterns

### Multiple Sections

```go
dmsg.Container(
    dmsg.AccentColor(gold),
    dmsg.ContainerSection(
        dmsg.SectionText("## Section 1"),
        dmsg.SectionText("Content here"),
    ),
    dmsg.ContainerSeparator(),
    dmsg.ContainerSection(
        dmsg.SectionText("## Section 2"),
        dmsg.SectionText("More content"),
    ),
)
```

### Section with Accessory

```go
dmsg.ContainerSection(
    dmsg.SectionText("Main content"),
    dmsg.SectionAccessory(
        dmsg.Thumbnail(imageURL, imageDesc),
    ),
)
```

### Multiple Buttons

```go
dmsg.ContainerActionRow(
    dmsg.Button("Confirm", "confirm", dmsg.Style(dmsg.Success)),
    dmsg.Button("Cancel", "cancel", dmsg.Style(dmsg.Danger)),
    dmsg.LinkButton("Help", "https://example.com"),
)
```

### Top-Level Components

```go
dmsg.Update(
    dmsg.TextDisplay("Simple text message"),
    dmsg.Separator(),
    dmsg.ActionRow(
        dmsg.Button("Action", "action_id"),
    ),
)
```

## Button Styles

- `dmsg.Primary` - Blue (default)
- `dmsg.Secondary` - Gray
- `dmsg.Success` - Green
- `dmsg.Danger` - Red

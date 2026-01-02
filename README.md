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
        dmsg.Section(
            dmsg.Text("## Hello, World!"),
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
        dmsg.Section(
            dmsg.Text("## Error\n\nYou don't have enough coins"),
            dmsg.Accessory(
                dmsg.Thumbnail(errorImageURL, "Error icon"),
            ),
        ),
        dmsg.Divider(),
    ),
)
```

### Success with Buttons

```go
dmsg.Response(
    dmsg.Container(
        dmsg.AccentColor(5763719), // green
        dmsg.Section(
            dmsg.Text("## Success!\n\nYou won **1,000 coins**"),
            dmsg.Accessory(
                dmsg.Thumbnail(trophyURL, trophyDesc),
            ),
        ),
        dmsg.Divider(),
        dmsg.Actions(
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
        dmsg.Section(
            dmsg.Text("This message is only visible to you"),
        ),
    ),
)
```

### Update Message

```go
dmsg.Update(
    dmsg.Container(
        dmsg.AccentColor(16777215), // white
        dmsg.Section(
            dmsg.Text("## Game Closed\n\nThanks for playing!"),
        ),
    ),
)
```

## Component Reference

### Entry Points

- `Response(components ...Component)` - Standard response
- `Ephemeral(components ...Component)` - Ephemeral response
- `Update(components ...Component)` - Update message response

**Note:** Use `Container()` and `ActionRow()` as top-level components. Use `Actions()` only inside containers.

### Layout Components

**Container**
```go
dmsg.Container(
    dmsg.AccentColor(color int),
    dmsg.Spoiler(),
    // child components...
)
```

**Section**
```go
dmsg.Section(
    dmsg.Text("content"),
    dmsg.Accessory(component),
)
```

**Divider**
```go
dmsg.Divider(
    dmsg.WithDivider(true),
    dmsg.Spacing(discordgo.SeparatorSpacingSizeSmall),
)
```

**Actions** (for use inside containers)
```go
dmsg.Actions(
    dmsg.Button(...),
    dmsg.Button(...),
)
```

**ActionRow** (for use at top-level)
```go
dmsg.ActionRow(
    dmsg.Button(...),
    dmsg.Button(...),
)
```

### Content Components

**Text**
```go
dmsg.Text("## Markdown content")
```

**Thumbnail**
```go
dmsg.Thumbnail(url, description,
    dmsg.Spoiler(),
)
```

**Gallery**
```go
dmsg.Gallery(
    dmsg.Media(url, description, spoiler),
    dmsg.Media(url, description, spoiler),
)
```

**File**
```go
dmsg.File("attachment://file.txt",
    dmsg.Spoiler(),
)
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
    dmsg.AccentColor(red),  // ContainerOption
    dmsg.Section(...),      // Also ContainerOption
)

// ❌ This won't compile
dmsg.Container(
    dmsg.Style(dmsg.Primary),  // ButtonOption - wrong type!
)
```

## Design Patterns

### Multiple Sections

```go
dmsg.Container(
    dmsg.AccentColor(gold),
    dmsg.Section(
        dmsg.Text("## Section 1"),
        dmsg.Text("Content here"),
    ),
    dmsg.Divider(),
    dmsg.Section(
        dmsg.Text("## Section 2"),
        dmsg.Text("More content"),
    ),
)
```

### Section with Accessory

```go
dmsg.Section(
    dmsg.Text("Main content"),
    dmsg.Accessory(
        dmsg.Thumbnail(imageURL, imageDesc),
    ),
)
```

### Multiple Buttons

```go
dmsg.Actions(
    dmsg.Button("Confirm", "confirm", dmsg.Style(dmsg.Success)),
    dmsg.Button("Cancel", "cancel", dmsg.Style(dmsg.Danger)),
    dmsg.LinkButton("Help", "https://example.com"),
)
```

## Button Styles

- `dmsg.Primary` - Blue (default)
- `dmsg.Secondary` - Gray
- `dmsg.Success` - Green
- `dmsg.Danger` - Red

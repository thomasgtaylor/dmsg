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
            dmsg.TextDisplay("## Hello, World!"),
        ),
    ),
)
```

## Features

- **Type-safe**: Options are typed per component - can't pass wrong option to wrong component
- **Composable**: Components nest naturally, reads top-to-bottom like the UI
- **Flexible**: Components work in multiple contexts (top-level, in containers, in sections)
- **Clean API**: Simple, intuitive naming without prefixes
- **Direct**: Uses discordgo types directly, no conversion needed

## Quick Start

### Simple Error Response

```go
dmsg.Response(
    dmsg.Container(
        dmsg.AccentColor(16740978), // red
        dmsg.Section(
            dmsg.TextDisplay("## Error\n\nYou don't have enough coins"),
            dmsg.Accessory(
                dmsg.Thumbnail(errorImageURL, "Error icon"),
            ),
        ),
        dmsg.Separator(),
    ),
)
```

### Success with Buttons

```go
dmsg.Response(
    dmsg.Container(
        dmsg.AccentColor(5763719), // green
        dmsg.Section(
            dmsg.TextDisplay("## Success!\n\nYou won **1,000 coins**"),
            dmsg.Accessory(
                dmsg.Thumbnail(trophyURL, trophyDesc),
            ),
        ),
        dmsg.Separator(),
        dmsg.ActionRow(
            dmsg.Button("Play Again", "play_again", dmsg.Style(dmsg.Success)),
            dmsg.Button("Quit", "quit", dmsg.Style(dmsg.Secondary)),
        ),
    ),
)
```

### Simple Top-Level Message

```go
dmsg.Update(
    dmsg.TextDisplay("It's dangerous to go alone!"),
    dmsg.Separator(),
    dmsg.TextDisplay("Take this."),
)
```

### Ephemeral Response

```go
dmsg.Ephemeral(
    dmsg.Container(
        dmsg.AccentColor(14197815), // gold
        dmsg.Section(
            dmsg.TextDisplay("This message is only visible to you"),
        ),
    ),
)
```

## Component Reference

### Entry Points

- `Response(components ...Component)` - Standard response
- `Ephemeral(components ...Component)` - Ephemeral response
- `Update(components ...Component)` - Update message response

### Layout Components

**Container**
```go
dmsg.Container(
    dmsg.AccentColor(color int),
    dmsg.Spoiler(),
    // child components...
)
```

Container can contain: `Section`, `TextDisplay`, `Separator`, `ActionRow`, `File`, `Gallery`

**Section**
```go
dmsg.Section(
    dmsg.TextDisplay("content"),
    dmsg.Accessory(component),
)
```

Section can be used:
- As a top-level component in `Response`, `Update`, `Ephemeral`
- Inside `Container`

Section can contain: `TextDisplay`
Section can have an accessory: `Button` or `Thumbnail`

**Separator**
```go
dmsg.Separator(
    dmsg.WithDivider(true),
    dmsg.Spacing(discordgo.SeparatorSpacingSizeSmall),
)
```

Separator can be used:
- As a top-level component
- Inside `Container`

**ActionRow**
```go
dmsg.ActionRow(
    dmsg.Button(...),
    dmsg.Button(...),
)
```

ActionRow can be used:
- As a top-level component
- Inside `Container`

ActionRow can contain: `Button`, `LinkButton`

### Content Components

**TextDisplay**
```go
dmsg.TextDisplay("## Markdown content")
```

TextDisplay can be used:
- As a top-level component
- Inside `Container`
- Inside `Section`

**Thumbnail**
```go
dmsg.Thumbnail(url, description,
    dmsg.Spoiler(),
)
```

Thumbnail can be used as a `Section` accessory.

**File**
```go
dmsg.File("attachment://file.txt",
    dmsg.Spoiler(),
)
```

File can be used inside `Container`.

**Gallery**
```go
dmsg.Gallery(
    dmsg.Media(url, description, spoiler),
    dmsg.Media(url, description, spoiler),
)
```

Gallery can be used inside `Container`.

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

## Flexible Component Usage

Components automatically work in multiple contexts:

```go
// TextDisplay works everywhere
dmsg.Response(
    dmsg.TextDisplay("Top level"),           // ✅ Top-level
    dmsg.Container(
        dmsg.TextDisplay("In container"),    // ✅ In Container
        dmsg.Section(
            dmsg.TextDisplay("In section"),  // ✅ In Section
        ),
    ),
)

// Section works at top-level and in containers
dmsg.Response(
    dmsg.Section(                            // ✅ Top-level
        dmsg.TextDisplay("Content"),
    ),
    dmsg.Container(
        dmsg.Section(                        // ✅ In Container
            dmsg.TextDisplay("Content"),
        ),
    ),
)
```

## Design Patterns

### Multiple Sections

```go
dmsg.Container(
    dmsg.AccentColor(gold),
    dmsg.Section(
        dmsg.TextDisplay("## Section 1"),
        dmsg.TextDisplay("Content here"),
    ),
    dmsg.Separator(),
    dmsg.Section(
        dmsg.TextDisplay("## Section 2"),
        dmsg.TextDisplay("More content"),
    ),
)
```

### Section with Accessory

```go
dmsg.Section(
    dmsg.TextDisplay("Main content"),
    dmsg.Accessory(
        dmsg.Thumbnail(imageURL, imageDesc),
    ),
)
```

### Multiple Buttons

```go
dmsg.ActionRow(
    dmsg.Button("Confirm", "confirm", dmsg.Style(dmsg.Success)),
    dmsg.Button("Cancel", "cancel", dmsg.Style(dmsg.Danger)),
    dmsg.LinkButton("Help", "https://example.com"),
)
```

### Mixing Top-Level Components

```go
dmsg.Response(
    dmsg.TextDisplay("Simple message"),
    dmsg.Separator(),
    dmsg.Section(
        dmsg.TextDisplay("Section content"),
    ),
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

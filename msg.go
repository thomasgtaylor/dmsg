// Package dmsg provides a type-safe builder for Discord Components v2 messages.
//
// It wraps discordgo types with a functional options pattern for clean,
// composable message construction. The API prevents common mistakes by using
// typed options that can only be applied to their corresponding components.
//
// Basic usage:
//
//	response := dmsg.Response(
//	    dmsg.Container(
//	        dmsg.AccentColor(5763719),
//	        dmsg.Section(
//	            dmsg.Text("## Hello World"),
//	        ),
//	    ),
//	)
//
// The package supports all Discord Components v2 features including containers,
// sections, buttons, thumbnails, galleries, and more.
package dmsg

import "github.com/bwmarrin/discordgo"

// Component is any Discord message component
type Component = discordgo.MessageComponent

// Response creates a standard interaction response
func Response(components ...Component) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsIsComponentsV2,
			Components: components,
		},
	}
}

// Ephemeral creates an ephemeral interaction response
func Ephemeral(components ...Component) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsIsComponentsV2 | discordgo.MessageFlagsEphemeral,
			Components: components,
		},
	}
}

// Update creates an update message response
func Update(components ...Component) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsIsComponentsV2,
			Components: components,
		},
	}
}

// ContainerOption configures a Container
type ContainerOption interface {
	applyToContainer(*discordgo.Container)
}

// Container creates a container component
func Container(opts ...ContainerOption) Component {
	container := &discordgo.Container{
		Components: []discordgo.MessageComponent{},
	}
	for _, opt := range opts {
		opt.applyToContainer(container)
	}
	return container
}

type accentColorOption struct {
	color int
}

func (o accentColorOption) applyToContainer(c *discordgo.Container) {
	c.AccentColor = &o.color
}

// AccentColor sets the container's accent color
func AccentColor(color int) ContainerOption {
	return accentColorOption{color}
}

type spoilerOption struct{}

func (o spoilerOption) applyToContainer(c *discordgo.Container) {
	c.Spoiler = true
}

func (o spoilerOption) applyToThumbnail(t *discordgo.Thumbnail) {
	t.Spoiler = true
}

func (o spoilerOption) applyToFile(f *discordgo.FileComponent) {
	f.Spoiler = true
}

// Spoiler marks a component as a spoiler (Container, Thumbnail, or File)
func Spoiler() interface {
	ContainerOption
	ThumbnailOption
	FileOption
} {
	return spoilerOption{}
}

type containerComponentOption struct {
	component discordgo.MessageComponent
}

func (o containerComponentOption) applyToContainer(c *discordgo.Container) {
	c.Components = append(c.Components, o.component)
}

// SectionOption configures a Section
type SectionOption interface {
	applyToSection(*discordgo.Section)
}

// ContainerSection creates a section component for use inside containers
func ContainerSection(opts ...SectionOption) ContainerOption {
	section := &discordgo.Section{
		Components: []discordgo.MessageComponent{},
	}
	for _, opt := range opts {
		opt.applyToSection(section)
	}
	return containerComponentOption{section}
}

type sectionComponentOption struct {
	component discordgo.MessageComponent
}

func (o sectionComponentOption) applyToSection(s *discordgo.Section) {
	s.Components = append(s.Components, o.component)
}

type accessoryOption struct {
	component discordgo.MessageComponent
}

func (o accessoryOption) applyToSection(s *discordgo.Section) {
	s.Accessory = o.component
}

// SectionAccessory sets the section's accessory (Button or Thumbnail)
func SectionAccessory(component Component) SectionOption {
	return accessoryOption{component}
}

// SectionText creates a text display component for use inside sections
func SectionText(content string) SectionOption {
	return sectionComponentOption{
		&discordgo.TextDisplay{
			Content: content,
		},
	}
}

// TextDisplay creates a text display component for use at top-level
func TextDisplay(content string) Component {
	return &discordgo.TextDisplay{
		Content: content,
	}
}

// ThumbnailOption configures a Thumbnail
type ThumbnailOption interface {
	applyToThumbnail(*discordgo.Thumbnail)
}

// Thumbnail creates a thumbnail component
func Thumbnail(url, description string, opts ...ThumbnailOption) Component {
	thumbnail := &discordgo.Thumbnail{
		Media: discordgo.UnfurledMediaItem{
			URL: url,
		},
		Description: &description,
	}
	for _, opt := range opts {
		opt.applyToThumbnail(thumbnail)
	}
	return thumbnail
}

// SeparatorOption configures a Separator
type SeparatorOption interface {
	applyToSeparator(*discordgo.Separator)
}

// ContainerSeparator creates a separator component for use inside containers
func ContainerSeparator(opts ...SeparatorOption) ContainerOption {
	truth := true
	spacing := discordgo.SeparatorSpacingSizeSmall
	separator := &discordgo.Separator{
		Divider: &truth,
		Spacing: &spacing,
	}
	for _, opt := range opts {
		opt.applyToSeparator(separator)
	}
	return containerComponentOption{separator}
}

// Separator creates a separator component for use at top-level
func Separator(opts ...SeparatorOption) Component {
	truth := true
	spacing := discordgo.SeparatorSpacingSizeSmall
	separator := &discordgo.Separator{
		Divider: &truth,
		Spacing: &spacing,
	}
	for _, opt := range opts {
		opt.applyToSeparator(separator)
	}
	return separator
}

type withDividerOption struct {
	show bool
}

func (o withDividerOption) applyToSeparator(s *discordgo.Separator) {
	s.Divider = &o.show
}

// WithDivider sets whether to show the divider line
func WithDivider(show bool) SeparatorOption {
	return withDividerOption{show}
}

type spacingOption struct {
	size discordgo.SeparatorSpacingSize
}

func (o spacingOption) applyToSeparator(s *discordgo.Separator) {
	s.Spacing = &o.size
}

// Spacing sets the separator spacing (1 = small, 2 = large)
func Spacing(size discordgo.SeparatorSpacingSize) SeparatorOption {
	return spacingOption{size}
}

// ContainerActionRow creates an action row with buttons for use inside containers
func ContainerActionRow(buttons ...Component) ContainerOption {
	return containerComponentOption{
		&discordgo.ActionsRow{
			Components: buttons,
		},
	}
}

// ActionRow creates an action row with buttons for use at top-level
func ActionRow(buttons ...Component) Component {
	return &discordgo.ActionsRow{
		Components: buttons,
	}
}

// ButtonStyle represents button visual styles
type ButtonStyle int

const (
	Primary   ButtonStyle = 1
	Secondary ButtonStyle = 2
	Success   ButtonStyle = 3
	Danger    ButtonStyle = 4
)

// ButtonOption configures a Button
type ButtonOption interface {
	applyToButton(*discordgo.Button)
}

// Button creates an action button
func Button(label, customID string, opts ...ButtonOption) Component {
	button := &discordgo.Button{
		Label:    label,
		CustomID: customID,
		Style:    discordgo.PrimaryButton,
	}
	for _, opt := range opts {
		opt.applyToButton(button)
	}
	return button
}

// LinkButton creates a link button
func LinkButton(label, url string, opts ...ButtonOption) Component {
	button := &discordgo.Button{
		Label: label,
		URL:   url,
		Style: discordgo.LinkButton,
	}
	for _, opt := range opts {
		opt.applyToButton(button)
	}
	return button
}

type styleOption struct {
	style ButtonStyle
}

func (o styleOption) applyToButton(b *discordgo.Button) {
	b.Style = discordgo.ButtonStyle(o.style)
}

// Style sets the button style
func Style(style ButtonStyle) ButtonOption {
	return styleOption{style}
}

type emojiOption struct {
	emoji *discordgo.ComponentEmoji
}

func (o emojiOption) applyToButton(b *discordgo.Button) {
	b.Emoji = o.emoji
}

// Emoji sets the button emoji
func Emoji(emoji *discordgo.ComponentEmoji) ButtonOption {
	return emojiOption{emoji}
}

type disabledOption struct{}

func (o disabledOption) applyToButton(b *discordgo.Button) {
	b.Disabled = true
}

// Disabled marks the button as disabled
func Disabled() ButtonOption {
	return disabledOption{}
}

// FileOption configures a File
type FileOption interface {
	applyToFile(*discordgo.FileComponent)
}

// ContainerFile creates a file component for use inside containers
func ContainerFile(url string, opts ...FileOption) ContainerOption {
	file := &discordgo.FileComponent{
		File: discordgo.UnfurledMediaItem{
			URL: url,
		},
	}
	for _, opt := range opts {
		opt.applyToFile(file)
	}
	return containerComponentOption{file}
}

// MediaItem represents a media gallery item
type MediaItem struct {
	URL         string
	Description string
	Spoiler     bool
}

// Media creates a media item for galleries
func Media(url, description string, spoiler bool) MediaItem {
	return MediaItem{
		URL:         url,
		Description: description,
		Spoiler:     spoiler,
	}
}

// ContainerGallery creates a media gallery component for use inside containers
func ContainerGallery(items ...MediaItem) ContainerOption {
	galleryItems := make([]discordgo.MediaGalleryItem, len(items))
	for i, item := range items {
		galleryItems[i] = discordgo.MediaGalleryItem{
			Media: discordgo.UnfurledMediaItem{
				URL: item.URL,
			},
			Description: &item.Description,
			Spoiler:     item.Spoiler,
		}
	}
	return containerComponentOption{
		&discordgo.MediaGallery{
			Items: galleryItems,
		},
	}
}

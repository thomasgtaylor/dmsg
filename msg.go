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

// unwrappable is implemented by wrapper types that need to be unwrapped
type unwrappable interface {
	unwrap() Component
}

func unwrapComponents(components []Component) []Component {
	unwrapped := make([]Component, len(components))
	for i, c := range components {
		if u, ok := c.(unwrappable); ok {
			unwrapped[i] = u.unwrap()
		} else {
			unwrapped[i] = c
		}
	}
	return unwrapped
}

// Response creates a standard interaction response
func Response(components ...Component) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsIsComponentsV2,
			Components: unwrapComponents(components),
		},
	}
}

// Ephemeral creates an ephemeral interaction response
func Ephemeral(components ...Component) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsIsComponentsV2 | discordgo.MessageFlagsEphemeral,
			Components: unwrapComponents(components),
		},
	}
}

// Update creates an update message response
func Update(components ...Component) *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Flags:      discordgo.MessageFlagsIsComponentsV2,
			Components: unwrapComponents(components),
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

type sectionComponent struct {
	*discordgo.Section
}

func (s sectionComponent) unwrap() Component {
	return s.Section
}

func (s sectionComponent) applyToContainer(c *discordgo.Container) {
	c.Components = append(c.Components, s.Section)
}

// Section creates a section component (can be used top-level or in containers)
func Section(opts ...SectionOption) interface {
	Component
	ContainerOption
} {
	section := &discordgo.Section{
		Components: []discordgo.MessageComponent{},
	}
	for _, opt := range opts {
		opt.applyToSection(section)
	}
	return sectionComponent{section}
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

// Accessory sets the section's accessory (Button or Thumbnail)
func Accessory(component Component) SectionOption {
	return accessoryOption{component}
}

type textDisplayComponent struct {
	*discordgo.TextDisplay
}

func (t textDisplayComponent) unwrap() Component {
	return t.TextDisplay
}

func (t textDisplayComponent) applyToSection(s *discordgo.Section) {
	s.Components = append(s.Components, t.TextDisplay)
}

func (t textDisplayComponent) applyToContainer(c *discordgo.Container) {
	c.Components = append(c.Components, t.TextDisplay)
}

// TextDisplay creates a text display component (can be used top-level, in containers, or in sections)
func TextDisplay(content string) interface {
	Component
	ContainerOption
	SectionOption
} {
	return textDisplayComponent{
		&discordgo.TextDisplay{
			Content: content,
		},
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

type separatorComponent struct {
	*discordgo.Separator
}

func (s separatorComponent) unwrap() Component {
	return s.Separator
}

func (s separatorComponent) applyToContainer(c *discordgo.Container) {
	c.Components = append(c.Components, s.Separator)
}

// Separator creates a separator component (can be used top-level or in containers)
func Separator(opts ...SeparatorOption) interface {
	Component
	ContainerOption
} {
	truth := true
	spacing := discordgo.SeparatorSpacingSizeSmall
	separator := &discordgo.Separator{
		Divider: &truth,
		Spacing: &spacing,
	}
	for _, opt := range opts {
		opt.applyToSeparator(separator)
	}
	return separatorComponent{separator}
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

type actionRowComponent struct {
	*discordgo.ActionsRow
}

func (a actionRowComponent) unwrap() Component {
	return a.ActionsRow
}

func (a actionRowComponent) applyToContainer(c *discordgo.Container) {
	c.Components = append(c.Components, a.ActionsRow)
}

// ActionRow creates an action row with buttons (can be used top-level or in containers)
func ActionRow(buttons ...Component) interface {
	Component
	ContainerOption
} {
	return actionRowComponent{
		&discordgo.ActionsRow{
			Components: buttons,
		},
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

type fileComponent struct {
	*discordgo.FileComponent
}

func (f fileComponent) applyToContainer(c *discordgo.Container) {
	c.Components = append(c.Components, f.FileComponent)
}

// File creates a file component (for use in containers)
func File(url string, opts ...FileOption) ContainerOption {
	file := &discordgo.FileComponent{
		File: discordgo.UnfurledMediaItem{
			URL: url,
		},
	}
	for _, opt := range opts {
		opt.applyToFile(file)
	}
	return fileComponent{file}
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

type mediaGalleryComponent struct {
	*discordgo.MediaGallery
}

func (m mediaGalleryComponent) applyToContainer(c *discordgo.Container) {
	c.Components = append(c.Components, m.MediaGallery)
}

// Gallery creates a media gallery component (for use in containers)
func Gallery(items ...MediaItem) ContainerOption {
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
	return mediaGalleryComponent{
		&discordgo.MediaGallery{
			Items: galleryItems,
		},
	}
}

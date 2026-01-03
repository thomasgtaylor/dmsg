package dmsg

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestResponse(t *testing.T) {
	t.Run("creates standard response", func(t *testing.T) {
		container := Container()
		response := Response(container)

		if response.Type != discordgo.InteractionResponseChannelMessageWithSource {
			t.Errorf("expected type %d, got %d", discordgo.InteractionResponseChannelMessageWithSource, response.Type)
		}

		if response.Data.Flags != discordgo.MessageFlagsIsComponentsV2 {
			t.Errorf("expected flags %d, got %d", discordgo.MessageFlagsIsComponentsV2, response.Data.Flags)
		}

		if len(response.Data.Components) != 1 {
			t.Errorf("expected 1 component, got %d", len(response.Data.Components))
		}
	})

	t.Run("handles multiple components", func(t *testing.T) {
		container1 := Container()
		container2 := Container()
		response := Response(container1, container2)

		if len(response.Data.Components) != 2 {
			t.Errorf("expected 2 components, got %d", len(response.Data.Components))
		}
	})

	t.Run("handles empty components", func(t *testing.T) {
		response := Response()

		if len(response.Data.Components) != 0 {
			t.Errorf("expected 0 components, got %d", len(response.Data.Components))
		}
	})
}

func TestEphemeral(t *testing.T) {
	t.Run("creates ephemeral response", func(t *testing.T) {
		container := Container()
		response := Ephemeral(container)

		if response.Type != discordgo.InteractionResponseChannelMessageWithSource {
			t.Errorf("expected type %d, got %d", discordgo.InteractionResponseChannelMessageWithSource, response.Type)
		}

		expectedFlags := discordgo.MessageFlagsIsComponentsV2 | discordgo.MessageFlagsEphemeral
		if response.Data.Flags != expectedFlags {
			t.Errorf("expected flags %d, got %d", expectedFlags, response.Data.Flags)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("creates update response", func(t *testing.T) {
		container := Container()
		response := Update(container)

		if response.Type != discordgo.InteractionResponseUpdateMessage {
			t.Errorf("expected type %d, got %d", discordgo.InteractionResponseUpdateMessage, response.Type)
		}

		if response.Data.Flags != discordgo.MessageFlagsIsComponentsV2 {
			t.Errorf("expected flags %d, got %d", discordgo.MessageFlagsIsComponentsV2, response.Data.Flags)
		}
	})
}

func TestContainer(t *testing.T) {
	t.Run("creates empty container", func(t *testing.T) {
		container := Container()

		c, ok := container.(*discordgo.Container)
		if !ok {
			t.Fatal("expected *discordgo.Container")
		}

		if c.AccentColor != nil {
			t.Error("expected nil accent color")
		}

		if c.Spoiler {
			t.Error("expected spoiler to be false")
		}

		if len(c.Components) != 0 {
			t.Errorf("expected 0 components, got %d", len(c.Components))
		}
	})

	t.Run("applies accent color", func(t *testing.T) {
		color := 16740978
		container := Container(AccentColor(color))

		c := container.(*discordgo.Container)
		if c.AccentColor == nil {
			t.Fatal("expected accent color to be set")
		}

		if *c.AccentColor != color {
			t.Errorf("expected color %d, got %d", color, *c.AccentColor)
		}
	})

	t.Run("applies spoiler", func(t *testing.T) {
		container := Container(Spoiler())

		c := container.(*discordgo.Container)
		if !c.Spoiler {
			t.Error("expected spoiler to be true")
		}
	})

	t.Run("adds section", func(t *testing.T) {
		container := Container(ContainerSection())

		c := container.(*discordgo.Container)
		if len(c.Components) != 1 {
			t.Errorf("expected 1 component, got %d", len(c.Components))
		}

		_, ok := c.Components[0].(*discordgo.Section)
		if !ok {
			t.Error("expected component to be *discordgo.Section")
		}
	})

	t.Run("adds multiple options", func(t *testing.T) {
		container := Container(
			AccentColor(123),
			Spoiler(),
			ContainerSection(),
			ContainerSeparator(),
		)

		c := container.(*discordgo.Container)
		if c.AccentColor == nil || *c.AccentColor != 123 {
			t.Error("accent color not applied correctly")
		}

		if !c.Spoiler {
			t.Error("spoiler not applied")
		}

		if len(c.Components) != 2 {
			t.Errorf("expected 2 components, got %d", len(c.Components))
		}
	})
}

func TestSection(t *testing.T) {
	t.Run("creates empty section", func(t *testing.T) {
		section := ContainerSection()

		opt, ok := section.(containerComponentOption)
		if !ok {
			t.Fatal("expected containerComponentOption")
		}

		s, ok := opt.component.(*discordgo.Section)
		if !ok {
			t.Fatal("expected *discordgo.Section")
		}

		if len(s.Components) != 0 {
			t.Errorf("expected 0 components, got %d", len(s.Components))
		}

		if s.Accessory != nil {
			t.Error("expected nil accessory")
		}
	})

	t.Run("adds text", func(t *testing.T) {
		section := ContainerSection(SectionText("test"))

		opt := section.(containerComponentOption)
		s := opt.component.(*discordgo.Section)

		if len(s.Components) != 1 {
			t.Fatalf("expected 1 component, got %d", len(s.Components))
		}

		td, ok := s.Components[0].(*discordgo.TextDisplay)
		if !ok {
			t.Fatal("expected *discordgo.TextDisplay")
		}

		if td.Content != "test" {
			t.Errorf("expected content 'test', got '%s'", td.Content)
		}
	})

	t.Run("adds multiple text components", func(t *testing.T) {
		section := ContainerSection(SectionText("first"), SectionText("second"))

		opt := section.(containerComponentOption)
		s := opt.component.(*discordgo.Section)

		if len(s.Components) != 2 {
			t.Errorf("expected 2 components, got %d", len(s.Components))
		}
	})

	t.Run("sets accessory", func(t *testing.T) {
		thumbnail := Thumbnail("http://example.com/image.png", "test image")
		section := ContainerSection(SectionAccessory(thumbnail))

		opt := section.(containerComponentOption)
		s := opt.component.(*discordgo.Section)

		if s.Accessory == nil {
			t.Fatal("expected accessory to be set")
		}

		_, ok := s.Accessory.(*discordgo.Thumbnail)
		if !ok {
			t.Error("expected accessory to be *discordgo.Thumbnail")
		}
	})
}

func TestText(t *testing.T) {
	t.Run("creates text display for sections", func(t *testing.T) {
		content := "## Hello\nWorld"
		text := SectionText(content)

		opt, ok := text.(sectionComponentOption)
		if !ok {
			t.Fatal("expected sectionComponentOption")
		}

		td, ok := opt.component.(*discordgo.TextDisplay)
		if !ok {
			t.Fatal("expected *discordgo.TextDisplay")
		}

		if td.Content != content {
			t.Errorf("expected content '%s', got '%s'", content, td.Content)
		}
	})

	t.Run("handles empty string", func(t *testing.T) {
		text := SectionText("")

		opt := text.(sectionComponentOption)
		td := opt.component.(*discordgo.TextDisplay)

		if td.Content != "" {
			t.Errorf("expected empty content, got '%s'", td.Content)
		}
	})

	t.Run("creates top-level text display", func(t *testing.T) {
		content := "## Top Level Text"
		text := TextDisplay(content)

		td, ok := text.(*discordgo.TextDisplay)
		if !ok {
			t.Fatal("expected *discordgo.TextDisplay")
		}

		if td.Content != content {
			t.Errorf("expected content '%s', got '%s'", content, td.Content)
		}
	})
}

func TestThumbnail(t *testing.T) {
	t.Run("creates thumbnail", func(t *testing.T) {
		url := "http://example.com/image.png"
		desc := "test description"
		thumbnail := Thumbnail(url, desc)

		th, ok := thumbnail.(*discordgo.Thumbnail)
		if !ok {
			t.Fatal("expected *discordgo.Thumbnail")
		}

		if th.Media.URL != url {
			t.Errorf("expected URL '%s', got '%s'", url, th.Media.URL)
		}

		if th.Description == nil {
			t.Fatal("expected description to be set")
		}

		if *th.Description != desc {
			t.Errorf("expected description '%s', got '%s'", desc, *th.Description)
		}

		if th.Spoiler {
			t.Error("expected spoiler to be false")
		}
	})

	t.Run("applies spoiler", func(t *testing.T) {
		thumbnail := Thumbnail("http://example.com/image.png", "test", Spoiler())

		th := thumbnail.(*discordgo.Thumbnail)
		if !th.Spoiler {
			t.Error("expected spoiler to be true")
		}
	})

	t.Run("handles empty description", func(t *testing.T) {
		thumbnail := Thumbnail("http://example.com/image.png", "")

		th := thumbnail.(*discordgo.Thumbnail)
		if th.Description == nil {
			t.Fatal("expected description to be set")
		}

		if *th.Description != "" {
			t.Errorf("expected empty description, got '%s'", *th.Description)
		}
	})
}

func TestSeparator(t *testing.T) {
	t.Run("creates default container separator", func(t *testing.T) {
		separator := ContainerSeparator()

		opt, ok := separator.(containerComponentOption)
		if !ok {
			t.Fatal("expected containerComponentOption")
		}

		sep, ok := opt.component.(*discordgo.Separator)
		if !ok {
			t.Fatal("expected *discordgo.Separator")
		}

		if sep.Divider == nil {
			t.Fatal("expected divider to be set")
		}

		if !*sep.Divider {
			t.Error("expected divider to be true")
		}

		if sep.Spacing == nil {
			t.Fatal("expected spacing to be set")
		}

		if *sep.Spacing != discordgo.SeparatorSpacingSizeSmall {
			t.Errorf("expected spacing %d, got %d", discordgo.SeparatorSpacingSizeSmall, *sep.Spacing)
		}
	})

	t.Run("applies with divider false", func(t *testing.T) {
		separator := ContainerSeparator(WithDivider(false))

		opt := separator.(containerComponentOption)
		sep := opt.component.(*discordgo.Separator)

		if sep.Divider == nil {
			t.Fatal("expected divider to be set")
		}

		if *sep.Divider {
			t.Error("expected divider to be false")
		}
	})

	t.Run("applies spacing", func(t *testing.T) {
		separator := ContainerSeparator(Spacing(discordgo.SeparatorSpacingSizeLarge))

		opt := separator.(containerComponentOption)
		sep := opt.component.(*discordgo.Separator)

		if sep.Spacing == nil {
			t.Fatal("expected spacing to be set")
		}

		if *sep.Spacing != discordgo.SeparatorSpacingSizeLarge {
			t.Errorf("expected spacing %d, got %d", discordgo.SeparatorSpacingSizeLarge, *sep.Spacing)
		}
	})

	t.Run("creates top-level separator", func(t *testing.T) {
		separator := Separator()

		sep, ok := separator.(*discordgo.Separator)
		if !ok {
			t.Fatal("expected *discordgo.Separator")
		}

		if sep.Divider == nil {
			t.Fatal("expected divider to be set")
		}

		if !*sep.Divider {
			t.Error("expected divider to be true")
		}
	})
}

func TestActions(t *testing.T) {
	t.Run("creates container actions row", func(t *testing.T) {
		button := Button("Test", "test_id")
		actions := ContainerActionRow(button)

		opt, ok := actions.(containerComponentOption)
		if !ok {
			t.Fatal("expected containerComponentOption")
		}

		row, ok := opt.component.(*discordgo.ActionsRow)
		if !ok {
			t.Fatal("expected *discordgo.ActionsRow")
		}

		if len(row.Components) != 1 {
			t.Errorf("expected 1 component, got %d", len(row.Components))
		}
	})

	t.Run("handles multiple buttons", func(t *testing.T) {
		button1 := Button("Button 1", "id1")
		button2 := Button("Button 2", "id2")
		actions := ContainerActionRow(button1, button2)

		opt := actions.(containerComponentOption)
		row := opt.component.(*discordgo.ActionsRow)

		if len(row.Components) != 2 {
			t.Errorf("expected 2 components, got %d", len(row.Components))
		}
	})

	t.Run("handles empty buttons", func(t *testing.T) {
		actions := ContainerActionRow()

		opt := actions.(containerComponentOption)
		row := opt.component.(*discordgo.ActionsRow)

		if len(row.Components) != 0 {
			t.Errorf("expected 0 components, got %d", len(row.Components))
		}
	})
}

func TestActionRow(t *testing.T) {
	t.Run("creates action row", func(t *testing.T) {
		button := Button("Test", "test_id")
		row := ActionRow(button)

		ar, ok := row.(*discordgo.ActionsRow)
		if !ok {
			t.Fatal("expected *discordgo.ActionsRow")
		}

		if len(ar.Components) != 1 {
			t.Errorf("expected 1 component, got %d", len(ar.Components))
		}
	})
}

func TestButton(t *testing.T) {
	t.Run("creates default button", func(t *testing.T) {
		label := "Click Me"
		customID := "click_button"
		button := Button(label, customID)

		btn, ok := button.(*discordgo.Button)
		if !ok {
			t.Fatal("expected *discordgo.Button")
		}

		if btn.Label != label {
			t.Errorf("expected label '%s', got '%s'", label, btn.Label)
		}

		if btn.CustomID != customID {
			t.Errorf("expected customID '%s', got '%s'", customID, btn.CustomID)
		}

		if btn.Style != discordgo.PrimaryButton {
			t.Errorf("expected style %d, got %d", discordgo.PrimaryButton, btn.Style)
		}

		if btn.Disabled {
			t.Error("expected button to be enabled")
		}

		if btn.Emoji != nil {
			t.Error("expected no emoji")
		}
	})

	t.Run("applies style", func(t *testing.T) {
		button := Button("Test", "test", Style(Secondary))

		btn := button.(*discordgo.Button)
		if btn.Style != discordgo.SecondaryButton {
			t.Errorf("expected style %d, got %d", discordgo.SecondaryButton, btn.Style)
		}
	})

	t.Run("applies all styles", func(t *testing.T) {
		tests := []struct {
			style    ButtonStyle
			expected discordgo.ButtonStyle
		}{
			{Primary, discordgo.PrimaryButton},
			{Secondary, discordgo.SecondaryButton},
			{Success, discordgo.SuccessButton},
			{Danger, discordgo.DangerButton},
		}

		for _, tt := range tests {
			button := Button("Test", "test", Style(tt.style))
			btn := button.(*discordgo.Button)

			if btn.Style != tt.expected {
				t.Errorf("expected style %d, got %d", tt.expected, btn.Style)
			}
		}
	})

	t.Run("applies emoji", func(t *testing.T) {
		emoji := &discordgo.ComponentEmoji{Name: "🎉"}
		button := Button("Test", "test", Emoji(emoji))

		btn := button.(*discordgo.Button)
		if btn.Emoji == nil {
			t.Fatal("expected emoji to be set")
		}

		if btn.Emoji.Name != emoji.Name {
			t.Errorf("expected emoji name '%s', got '%s'", emoji.Name, btn.Emoji.Name)
		}
	})

	t.Run("applies disabled", func(t *testing.T) {
		button := Button("Test", "test", Disabled())

		btn := button.(*discordgo.Button)
		if !btn.Disabled {
			t.Error("expected button to be disabled")
		}
	})

	t.Run("applies multiple options", func(t *testing.T) {
		emoji := &discordgo.ComponentEmoji{Name: "✅"}
		button := Button("Confirm", "confirm",
			Style(Success),
			Emoji(emoji),
			Disabled(),
		)

		btn := button.(*discordgo.Button)
		if btn.Style != discordgo.SuccessButton {
			t.Error("style not applied")
		}

		if btn.Emoji == nil || btn.Emoji.Name != emoji.Name {
			t.Error("emoji not applied")
		}

		if !btn.Disabled {
			t.Error("disabled not applied")
		}
	})

	t.Run("handles empty strings", func(t *testing.T) {
		button := Button("", "")

		btn := button.(*discordgo.Button)
		if btn.Label != "" {
			t.Errorf("expected empty label, got '%s'", btn.Label)
		}

		if btn.CustomID != "" {
			t.Errorf("expected empty customID, got '%s'", btn.CustomID)
		}
	})
}

func TestLinkButton(t *testing.T) {
	t.Run("creates link button", func(t *testing.T) {
		label := "Visit Site"
		url := "https://example.com"
		button := LinkButton(label, url)

		btn, ok := button.(*discordgo.Button)
		if !ok {
			t.Fatal("expected *discordgo.Button")
		}

		if btn.Label != label {
			t.Errorf("expected label '%s', got '%s'", label, btn.Label)
		}

		if btn.URL != url {
			t.Errorf("expected URL '%s', got '%s'", url, btn.URL)
		}

		if btn.Style != discordgo.LinkButton {
			t.Errorf("expected style %d, got %d", discordgo.LinkButton, btn.Style)
		}

		if btn.CustomID != "" {
			t.Error("expected empty customID for link button")
		}
	})

	t.Run("applies emoji", func(t *testing.T) {
		emoji := &discordgo.ComponentEmoji{Name: "🔗"}
		button := LinkButton("Link", "https://example.com", Emoji(emoji))

		btn := button.(*discordgo.Button)
		if btn.Emoji == nil {
			t.Fatal("expected emoji to be set")
		}

		if btn.Emoji.Name != emoji.Name {
			t.Errorf("expected emoji name '%s', got '%s'", emoji.Name, btn.Emoji.Name)
		}
	})
}

func TestFile(t *testing.T) {
	t.Run("creates file component", func(t *testing.T) {
		url := "attachment://file.txt"
		file := ContainerFile(url)

		opt, ok := file.(containerComponentOption)
		if !ok {
			t.Fatal("expected containerComponentOption")
		}

		fc, ok := opt.component.(*discordgo.FileComponent)
		if !ok {
			t.Fatal("expected *discordgo.FileComponent")
		}

		if fc.File.URL != url {
			t.Errorf("expected URL '%s', got '%s'", url, fc.File.URL)
		}

		if fc.Spoiler {
			t.Error("expected spoiler to be false")
		}
	})

	t.Run("applies spoiler", func(t *testing.T) {
		file := ContainerFile("attachment://file.txt", Spoiler())

		opt := file.(containerComponentOption)
		fc := opt.component.(*discordgo.FileComponent)

		if !fc.Spoiler {
			t.Error("expected spoiler to be true")
		}
	})
}

func TestMedia(t *testing.T) {
	t.Run("creates media item", func(t *testing.T) {
		url := "http://example.com/image.png"
		desc := "test image"
		spoiler := false

		media := Media(url, desc, spoiler)

		if media.URL != url {
			t.Errorf("expected URL '%s', got '%s'", url, media.URL)
		}

		if media.Description != desc {
			t.Errorf("expected description '%s', got '%s'", desc, media.Description)
		}

		if media.Spoiler != spoiler {
			t.Errorf("expected spoiler %t, got %t", spoiler, media.Spoiler)
		}
	})

	t.Run("creates spoiler media item", func(t *testing.T) {
		media := Media("http://example.com/image.png", "desc", true)

		if !media.Spoiler {
			t.Error("expected spoiler to be true")
		}
	})
}

func TestGallery(t *testing.T) {
	t.Run("creates gallery", func(t *testing.T) {
		item1 := Media("http://example.com/image1.png", "Image 1", false)
		item2 := Media("http://example.com/image2.png", "Image 2", true)

		gallery := ContainerGallery(item1, item2)

		opt, ok := gallery.(containerComponentOption)
		if !ok {
			t.Fatal("expected containerComponentOption")
		}

		gal, ok := opt.component.(*discordgo.MediaGallery)
		if !ok {
			t.Fatal("expected *discordgo.MediaGallery")
		}

		if len(gal.Items) != 2 {
			t.Fatalf("expected 2 items, got %d", len(gal.Items))
		}

		if gal.Items[0].Media.URL != item1.URL {
			t.Errorf("expected URL '%s', got '%s'", item1.URL, gal.Items[0].Media.URL)
		}

		if gal.Items[0].Description == nil {
			t.Fatal("expected description to be set")
		}

		if *gal.Items[0].Description != item1.Description {
			t.Errorf("expected description '%s', got '%s'", item1.Description, *gal.Items[0].Description)
		}

		if gal.Items[0].Spoiler != item1.Spoiler {
			t.Errorf("expected spoiler %t, got %t", item1.Spoiler, gal.Items[0].Spoiler)
		}

		if !gal.Items[1].Spoiler {
			t.Error("expected second item to be spoiler")
		}
	})

	t.Run("creates empty gallery", func(t *testing.T) {
		gallery := ContainerGallery()

		opt := gallery.(containerComponentOption)
		gal := opt.component.(*discordgo.MediaGallery)

		if len(gal.Items) != 0 {
			t.Errorf("expected 0 items, got %d", len(gal.Items))
		}
	})
}

func TestAccessory(t *testing.T) {
	t.Run("sets thumbnail accessory", func(t *testing.T) {
		thumbnail := Thumbnail("http://example.com/image.png", "test")
		section := ContainerSection(SectionAccessory(thumbnail))

		opt := section.(containerComponentOption)
		s := opt.component.(*discordgo.Section)

		if s.Accessory == nil {
			t.Fatal("expected accessory to be set")
		}

		_, ok := s.Accessory.(*discordgo.Thumbnail)
		if !ok {
			t.Error("expected accessory to be *discordgo.Thumbnail")
		}
	})

	t.Run("sets button accessory", func(t *testing.T) {
		button := Button("Click", "click_id")
		section := ContainerSection(SectionAccessory(button))

		opt := section.(containerComponentOption)
		s := opt.component.(*discordgo.Section)

		if s.Accessory == nil {
			t.Fatal("expected accessory to be set")
		}

		_, ok := s.Accessory.(*discordgo.Button)
		if !ok {
			t.Error("expected accessory to be *discordgo.Button")
		}
	})
}

func TestComplexMessage(t *testing.T) {
	t.Run("builds complex message structure", func(t *testing.T) {
		response := Response(
			Container(
				AccentColor(5763719),
				Spoiler(),
				ContainerSection(
					SectionText("## Title"),
					SectionText("Description text"),
					SectionAccessory(Thumbnail("http://example.com/img.png", "image")),
				),
				ContainerSeparator(),
				ContainerSection(
					SectionText("Another section"),
				),
				ContainerActionRow(
					Button("Action 1", "action1", Style(Primary)),
					Button("Action 2", "action2", Style(Secondary), Disabled()),
				),
			),
			ActionRow(
				LinkButton("External Link", "https://example.com"),
			),
		)

		if response.Type != discordgo.InteractionResponseChannelMessageWithSource {
			t.Error("wrong response type")
		}

		if len(response.Data.Components) != 2 {
			t.Fatalf("expected 2 top-level components, got %d", len(response.Data.Components))
		}

		container, ok := response.Data.Components[0].(*discordgo.Container)
		if !ok {
			t.Fatal("first component should be container")
		}

		if !container.Spoiler {
			t.Error("container should be spoiler")
		}

		if container.AccentColor == nil || *container.AccentColor != 5763719 {
			t.Error("wrong accent color")
		}

		if len(container.Components) != 4 {
			t.Errorf("expected 4 container components, got %d", len(container.Components))
		}

		actionRow, ok := response.Data.Components[1].(*discordgo.ActionsRow)
		if !ok {
			t.Fatal("second component should be action row")
		}

		if len(actionRow.Components) != 1 {
			t.Error("expected 1 button in action row")
		}
	})
}

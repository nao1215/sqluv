package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"gopkg.in/yaml.v3"
)

// ColorScheme represents a complete color scheme for the TUI
type ColorScheme struct {
	Name            string `yaml:"name"`
	Background      string `yaml:"background"`
	Foreground      string `yaml:"foreground"`
	Border          string `yaml:"border"`
	BorderFocus     string `yaml:"border_focus"`
	Selection       string `yaml:"selection"`
	SelectionText   string `yaml:"selection_text"`
	Header          string `yaml:"header"`
	Button          string `yaml:"button"`
	ButtonFocus     string `yaml:"button_focus"`
	ButtonText      string `yaml:"button_text"`
	ButtonTextFocus string `yaml:"button_text_focus"`
}

// GetTcellColor converts a color string to a tcell.Color
func GetTcellColor(colorName string) tcell.Color {
	switch colorName {
	case "black":
		return tcell.ColorBlack
	case "red":
		return tcell.ColorRed
	case "green":
		return tcell.ColorGreen
	case "yellow":
		return tcell.ColorYellow
	case "blue":
		return tcell.ColorBlue
	case "magenta":
		return tcell.ColorDarkMagenta
	case "cyan":
		return tcell.ColorDarkCyan
	case "white":
		return tcell.ColorWhite
	case "gray", "grey":
		return tcell.ColorGray
	default:
		// Try to parse the color as a hex value
		return tcell.GetColor(colorName)
	}
}

// ColorConfig manages the application color schemes
type ColorConfig struct {
	CurrentScheme *ColorScheme
	Schemes       map[string]*ColorScheme
	configPath    string
}

// DefaultColorSchemes returns the default color schemes
func DefaultColorSchemes() map[string]*ColorScheme {
	return map[string]*ColorScheme{
		"default": {
			Name:            "Default",
			Background:      "black",
			Foreground:      "white",
			Border:          "white",
			BorderFocus:     "green",
			Selection:       "green",
			SelectionText:   "black",
			Header:          "yellow",
			Button:          "white",
			ButtonFocus:     "lightgray",
			ButtonText:      "black",
			ButtonTextFocus: "black",
		},
		"dark": {
			Name:            "Dark",
			Background:      "#222222",
			Foreground:      "#e0e0e0",
			Border:          "#555555",
			BorderFocus:     "#00cc00",
			Selection:       "#004400",
			SelectionText:   "#00ff00",
			Header:          "#ffcc00",
			Button:          "#555555",
			ButtonFocus:     "#888888",
			ButtonText:      "#ffffff",
			ButtonTextFocus: "#ffffff",
		},
		"light": {
			Name:            "Light",
			Background:      "#f0f0f0",
			Foreground:      "#222222",
			Border:          "#999999",
			BorderFocus:     "#0066cc",
			Selection:       "#ccddff",
			SelectionText:   "#000000",
			Header:          "#333399",
			Button:          "#dddddd",
			ButtonFocus:     "#0066cc",
			ButtonText:      "#222222",
			ButtonTextFocus: "#ffffff",
		},
		"solarized": {
			Name:            "Solarized",
			Background:      "#002b36",
			Foreground:      "#839496",
			Border:          "#586e75",
			BorderFocus:     "#cb4b16",
			Selection:       "#073642",
			SelectionText:   "#93a1a1",
			Header:          "#b58900",
			Button:          "#073642",
			ButtonFocus:     "#586e75",
			ButtonText:      "#839496",
			ButtonTextFocus: "#fdf6e3",
		},
		"monokai": {
			Name:            "Monokai",
			Background:      "#272822",
			Foreground:      "#f8f8f2",
			Border:          "#75715e",
			BorderFocus:     "#f92672",
			Selection:       "#49483e",
			SelectionText:   "#f8f8f2",
			Header:          "#66d9ef",
			Button:          "#75715e",
			ButtonFocus:     "#f92672",
			ButtonText:      "#f8f8f2",
			ButtonTextFocus: "#ffffff",
		},
		"dracula": {
			Name:            "Dracula",
			Background:      "#282a36",
			Foreground:      "#f8f8f2",
			Border:          "#6272a4",
			BorderFocus:     "#ff79c6",
			Selection:       "#44475a",
			SelectionText:   "#f8f8f2",
			Header:          "#8be9fd",
			Button:          "#6272a4",
			ButtonFocus:     "#ff79c6",
			ButtonText:      "#f8f8f2",
			ButtonTextFocus: "#f8f8f2",
		},
		"nord": {
			Name:            "Nord",
			Background:      "#2e3440",
			Foreground:      "#d8dee9",
			Border:          "#4c566a",
			BorderFocus:     "#88c0d0",
			Selection:       "#3b4252",
			SelectionText:   "#eceff4",
			Header:          "#5e81ac",
			Button:          "#4c566a",
			ButtonFocus:     "#88c0d0",
			ButtonText:      "#e5e9f0",
			ButtonTextFocus: "#2e3440",
		},
		"gruvbox": {
			Name:            "Gruvbox",
			Background:      "#282828",
			Foreground:      "#ebdbb2",
			Border:          "#665c54",
			BorderFocus:     "#fe8019",
			Selection:       "#504945",
			SelectionText:   "#ebdbb2",
			Header:          "#b8bb26",
			Button:          "#665c54",
			ButtonFocus:     "#fe8019",
			ButtonText:      "#ebdbb2",
			ButtonTextFocus: "#fbf1c7",
		},
		"tokyo-night": {
			Name:            "Tokyo Night",
			Background:      "#1a1b26",
			Foreground:      "#a9b1d6",
			Border:          "#414868",
			BorderFocus:     "#7aa2f7",
			Selection:       "#24283b",
			SelectionText:   "#c0caf5",
			Header:          "#bb9af7",
			Button:          "#414868",
			ButtonFocus:     "#7aa2f7",
			ButtonText:      "#c0caf5",
			ButtonTextFocus: "#1a1b26",
		},
		"catppuccin": {
			Name:            "Catppuccin",
			Background:      "#1e1e2e",
			Foreground:      "#cdd6f4",
			Border:          "#585b70",
			BorderFocus:     "#f5c2e7",
			Selection:       "#313244",
			SelectionText:   "#cdd6f4",
			Header:          "#89b4fa",
			Button:          "#45475a",
			ButtonFocus:     "#f5c2e7",
			ButtonText:      "#cdd6f4",
			ButtonTextFocus: "#1e1e2e",
		},
		"vscode": {
			Name:            "VS Code",
			Background:      "#1e1e1e",
			Foreground:      "#d4d4d4",
			Border:          "#3c3c3c",
			BorderFocus:     "#007acc",
			Selection:       "#264f78",
			SelectionText:   "#ffffff",
			Header:          "#569cd6",
			Button:          "#3c3c3c",
			ButtonFocus:     "#007acc",
			ButtonText:      "#d4d4d4",
			ButtonTextFocus: "#ffffff",
		},
		"atom": {
			Name:            "Atom",
			Background:      "#282c34",
			Foreground:      "#abb2bf",
			Border:          "#3b4048",
			BorderFocus:     "#528bff",
			Selection:       "#3e4451",
			SelectionText:   "#abb2bf",
			Header:          "#61afef",
			Button:          "#3b4048",
			ButtonFocus:     "#528bff",
			ButtonText:      "#abb2bf",
			ButtonTextFocus: "#ffffff",
		},
		"sublime": {
			Name:            "Sublime Text",
			Background:      "#272822",
			Foreground:      "#f8f8f2",
			Border:          "#75715e",
			BorderFocus:     "#a6e22e",
			Selection:       "#49483e",
			SelectionText:   "#f8f8f2",
			Header:          "#66d9ef",
			Button:          "#75715e",
			ButtonFocus:     "#a6e22e",
			ButtonText:      "#f8f8f2",
			ButtonTextFocus: "#272822",
		},
	}
}

// NewColorConfig creates a new color manager
func NewColorConfig() (*ColorConfig, error) {
	configDir, err := ensureConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, "color_scheme.yml")
	manager := &ColorConfig{
		Schemes:    DefaultColorSchemes(),
		configPath: configPath,
	}

	// Try to load saved scheme, fall back to default if not found
	savedScheme, err := manager.loadSavedScheme()
	if err == nil {
		manager.CurrentScheme = savedScheme
	} else {
		manager.CurrentScheme = manager.Schemes["default"]
	}

	return manager, nil
}

// loadSavedScheme loads the color scheme from the configuration file
func (cm *ColorConfig) loadSavedScheme() (*ColorScheme, error) {
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return nil, err
	}

	var scheme ColorScheme
	err = yaml.Unmarshal(data, &scheme)
	if err != nil {
		return nil, err
	}

	return &scheme, nil
}

// SaveCurrentScheme saves the current color scheme to the configuration file
func (cm *ColorConfig) SaveCurrentScheme() error {
	if cm.CurrentScheme == nil {
		return errors.New("no current color scheme to save")
	}

	data, err := yaml.Marshal(cm.CurrentScheme)
	if err != nil {
		return err
	}

	return os.WriteFile(cm.configPath, data, 0600)
}

// SetScheme sets the current color scheme
func (cm *ColorConfig) SetScheme(name string) error {
	scheme, exists := cm.Schemes[name]
	if !exists {
		return errors.New("color scheme not found")
	}

	cm.CurrentScheme = scheme
	return cm.SaveCurrentScheme()
}

// GetSchemeNames returns the names of all available color schemes
func (cm *ColorConfig) GetSchemeNames() []string {
	names := make([]string, 0, len(cm.Schemes))
	for name := range cm.Schemes {
		names = append(names, name)
	}
	return names
}

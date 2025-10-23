# Pel 🎨

[![Go Reference](https://pkg.go.dev/badge/github.com/carlomunguia/pel.svg)](https://pkg.go.dev/github.com/carlomunguia/pel)
[![Go Report Card](https://goreportcard.com/badge/github.com/carlomunguia/pel)](https://goreportcard.com/report/github.com/carlomunguia/pel)
[![License](https://img.shields.io/github/license/carlomunguia/pel)](LICENSE)

**Pel** is a lightweight, cross-platform pixel art editor written in Go using the [Fyne](https://fyne.io) UI toolkit. Create beautiful pixel art with an intuitive interface, powerful color management, and real-time canvas manipulation.

![Pel Screenshot](docs/screenshot.png) <!-- Add a screenshot to docs/ folder -->

## ✨ Features

- 🖼️ **Flexible Canvas** - Create custom-sized pixel art canvases (up to 1024x1024)
- 🎨 **Advanced Color Picker** - HSV/RGBA color selection with visual preview
- 🎯 **Color Palette Management** - Save and manage up to 64 color swatches
- 🔍 **Zoom & Pan** - Smooth canvas navigation with mouse wheel zoom
- 💾 **File Operations** - Open, save, and export PNG images
- 🖱️ **Multiple Brush Tools** - Pencil, eraser, fill, and more (coming soon)
- 🌈 **Automatic Color Extraction** - Import colors from existing images
- ⚡ **Performance** - Built with Go for speed and efficiency

## 📦 Installation

### Prerequisites

- Go 1.25 or later
- Platform-specific dependencies for Fyne:
  - **macOS**: Xcode command line tools
  - **Linux**: `libgl1-mesa-dev xorg-dev` (Debian/Ubuntu) or equivalent
  - **Windows**: No additional dependencies

### Install from source

```bash
go install github.com/carlomunguia/pel/pel@latest
```

### Clone and build

```bash
git clone https://github.com/carlomunguia/pel.git
cd pel
go build -v ./pel
```

## 🚀 Quick Start

### Running Pel

```bash
# If installed via go install
pel

# Or run from source
go run ./pel
```

### Creating Your First Artwork

1. **Create a new canvas**: `File → New`

   - Enter your desired width and height (e.g., 32x32)
   - Click "Create"

2. **Pick a color**: Use the color picker on the right

   - Drag the hue slider to select a color
   - Click in the color square to choose your shade
   - The color appears in the preview box

3. **Add to palette**: Click an empty swatch at the bottom to save your color

4. **Start drawing**: Click on the canvas to paint pixels!

5. **Save your work**: `File → Save` or `File → Save As`

## 🎨 Usage Guide

### Canvas Navigation

| Action          | Method               |
| --------------- | -------------------- |
| **Zoom In/Out** | Mouse scroll wheel   |
| **Pan Canvas**  | Middle-click + drag  |
| **Draw Pixel**  | Left-click on canvas |

### Color Management

#### Using the Color Picker

1. Move the vertical slider to select hue
2. Click in the color square for saturation/value
3. Selected color appears in preview
4. Click canvas to draw with selected color

#### Using Color Swatches

1. Select a color with the picker
2. Click an empty swatch slot to save
3. Click saved swatches to quickly switch colors
4. Up to 64 swatches available

#### Import Colors from Images

1. `File → Open` to load an image
2. Pel automatically extracts unique colors
3. Colors populate available swatch slots

### File Operations

| Operation  | Menu Path        | Shortcut    |
| ---------- | ---------------- | ----------- |
| New Canvas | `File → New`     | -           |
| Open Image | `File → Open`    | -           |
| Save       | `File → Save`    | -           |
| Save As    | `File → Save As` | -           |
| Export     | `File → Export`  | Coming soon |
| Quit       | `File → Quit`    | -           |

## 🏗️ Architecture

Pel follows a clean, modular architecture:

```
pel/
├── pel/           # Main application entry point
├── apptype/       # Core types and interfaces
├── pelcanvas/     # Canvas widget and rendering
│   └── brush/     # Brush tools implementation
├── swatch/        # Color swatch widgets
├── ui/            # User interface components
│   ├── layout.go  # Main layout
│   ├── menus.go   # File menu handlers
│   ├── picker.go  # Color picker setup
│   └── swatches.go # Swatch panel
└── util/          # Utility functions
```

## 🛠️ Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/carlomunguia/pel.git
cd pel

# Download dependencies
go mod download

# Build
go build -v ./pel

# Run tests
go test ./...
```

### Project Requirements

- **Go**: 1.25+
- **Fyne**: v2.7.0+

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Fyne](https://fyne.io) - A beautiful cross-platform GUI toolkit
- Color picker by [lusingander/colorpicker](https://github.com/lusingander/colorpicker)

## 📧 Contact

Carlo Munguia - [@carlomunguia](https://github.com/carlomunguia)

Project Link: [https://github.com/carlomunguia/pel](https://github.com/carlomunguia/pel)

---

Made with ❤️ and Go

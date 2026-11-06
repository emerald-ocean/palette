# Emerald Ocean palette

The single source of truth for the Emerald Ocean colours: a dark palette taken from a breaking emerald wave, where green dominates the blue.

`palette.json` holds 15 colours in four groups (abyss, foam, kelp, reef). Each colour is defined in OKLCH, with its sRGB hex and what it is used for. The editor themes (`emerald-ocean.nvim`, `emerald-ocean-vscode`) are built from this file.

| Group | Colours |
|-------|---------|
| abyss | `#01130b` `#001c12` `#02291c` `#405d50` |
| foam  | `#eceff1` `#cbdbdf` `#8ca59a` |
| kelp  | `#83eab5` `#97c37f` `#7bd3cf` `#76dffb` `#5b92bf` |
| reef  | `#f7a792` `#df676b` `#baa4e2` |

## Rules

The `rules` section of `palette.json` sets the readability rules every colour must meet: a minimum contrast against the backgrounds, and a minimum perceptual distance (OKLab ΔE) between colours that commonly sit next to each other in code.

## Generator

`emerald-gen` checks the palette and renders editor colour files from templates. The palette is embedded in the binary, so the version you run is the palette version you get.

```sh
go run github.com/emerald-ocean/palette/cmd/emerald-gen@<version> check
go run github.com/emerald-ocean/palette/cmd/emerald-gen@<version> generate -template palette.lua.tmpl -out lua/emerald-ocean/palette.lua
```

`check` fails if a colour is outside the sRGB gamut, its hex does not match its OKLCH value, or it breaks a contrast or neighbour rule. `generate` renders a [text/template](https://pkg.go.dev/text/template) with:

- `{{.emerald}}`: a colour's hex value, e.g. `#83eab5`
- `{{alpha .glow 0.16}}`: a colour with an alpha channel, e.g. `#76dffb29`
- `{{version}}`: the palette version

`report` prints the APCA lightness contrast (Lc) of every colour on each background, and the ΔE of every neighbour pair as seen with normal vision, protanopia, deuteranopia and tritanopia (Machado et al. 2009). It informs the rules but does not fail.

All commands take `-palette path/to/palette.json` to use a local file instead of the embedded palette.

## Development

```sh
make build vet test test-race lint
```

## Preview

`preview/index.html` shows the palette on a full Go file in a Neovim/tmux mock-up. Open it in a browser.

## Images

- `assets/emerald-wave.png`: a barrel wave recoloured so every pixel sits on the palette.
- `assets/emerald-wave-avatar.jpg`: 800×800 crop of the same image, used as the organisation profile picture.

Wave photo: ["Vortex" by Matt Hardy](https://unsplash.com/photos/zJ-Cj3nT6EU) on Unsplash, used under the Unsplash License.

## License

MIT

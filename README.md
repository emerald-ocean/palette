# Emerald Ocean palette

The single source of truth for the Emerald Ocean colours: a dark palette taken from a breaking emerald wave, where green dominates the blue.

`pkg/palette/palette.json` holds 44 colours in nine families. Each family is one hue (bioluminescence groups the light sources) with shades numbered from dark to light: accent families have two dark tints for backgrounds and highlights, then a muted, a core and a bright shade for text. Each colour is defined in OKLCH, with its sRGB hex, what it is used for, and whether it is the family's core. The editor themes (`emerald-ocean.nvim`, `emerald-ocean-vscode`) are built from this file.

| Family | Shades 1 → 5 (core in bold) |
|--------|------|
| abyss | **`#01130b`** `#001c12` `#02291c` `#163b2c` `#2f5042` |
| foam | `#6d7e7d` `#9ab1ab` `#aeb4b4` `#c6d4d5` **`#dbeaeb`** |
| emerald | `#052f1a` `#074628` `#35b977` **`#54f9a3`** `#bafed3` |
| kelp | `#142e11` `#1f451b` `#61a659` **`#79cb70`** `#98e98f` |
| lagoon | `#012d2f` `#0a4346` `#40b1b6` **`#47d4db`** `#87f2f8` |
| sea | `#052941` `#083e61` `#4a93ca` **`#61aee9`** `#84cef9` |
| reef | `#3f1918` `#612422` `#f66c6d` **`#f5b1aa`** `#fed2cb` |
| anemone | `#2b1f3f` `#41305e` `#aa8dde` **`#c6b1f0`** `#e3d7fb` |
| biolume | `#06343e` `#431f3c` `#f59ae1` **`#76dffb`** |

## Rules

The `rules` section of `pkg/palette/palette.json` and each colour's optional `contrast` range set what every colour must meet:

- **APCA lightness contrast:** colours with a `contrast` range must reach its `min` on the `background` (and `min` minus `selectionAllowance` on the `selection`), and stay under its `max`: 60 for code, 45 for comments, punctuation and packages, 30 for line numbers, at most 92 for the brightest text to limit halation.
- **Neighbour separation:** colour pairs that sit next to each other in code must differ by at least `minDeltaE` (OKLab ΔE), and by `minDeltaEColourBlind` as seen with protanopia and deuteranopia.

## Generator

`emerald-gen` checks the palette and renders editor colour files from templates. The palette is embedded in the binary, so the version you run is the palette version you get.

```sh
go run github.com/emerald-ocean/palette/cmd/emerald-gen@<version> check
go run github.com/emerald-ocean/palette/cmd/emerald-gen@<version> generate -template palette.lua.tmpl -out lua/emerald-ocean/palette.lua
```

`check` fails if a colour is outside the sRGB gamut, its hex does not match its OKLCH value, or it breaks a contrast or neighbour rule. `generate` renders a [text/template](https://pkg.go.dev/text/template) with:

- `{{.emerald4}}`: a colour's hex value, e.g. `#54f9a3`
- `{{alpha .biolume4 0.16}}`: a colour with an alpha channel, e.g. `#76dffb29`
- `{{version}}`: the palette version

`report` prints the APCA lightness contrast (Lc) of every colour on each background, and the ΔE of every neighbour pair as seen with normal vision, protanopia, deuteranopia and tritanopia (Machado et al. 2009). It informs the rules but does not fail.

All commands take `-palette path/to/palette.json` to use a local file instead of the embedded palette.

## Development

```sh
make build vet test test-race lint check
```

## Releases

Releases are automated with [release-please](https://github.com/googleapis/release-please). Every push to `main` updates a release PR with the next version, worked out from the conventional commit messages. Merging that PR tags `vX.Y.Z` and publishes a GitHub Release with the notes. The tag is the only place the version lives: `emerald-gen` reads it from its Go build info, so `go run …/emerald-gen@vX.Y.Z` reports exactly that version.

## Preview

`preview/index.html` shows the palette on a full Go file in a Neovim/tmux mock-up, and `preview/scales.html` shows every family with its shades, cores, contrast and colour-blind separation, and how the tests work. Open them in a browser.

## Images

- `assets/emerald-wave.png`: a barrel wave recoloured so every pixel sits on the palette.
- `assets/emerald-wave-avatar.jpg`: 800×800 crop of the same image, used as the organisation profile picture.

Wave photo: ["Vortex" by Matt Hardy](https://unsplash.com/photos/zJ-Cj3nT6EU) on Unsplash, used under the Unsplash License.

## License

MIT

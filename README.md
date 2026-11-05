# Emerald Ocean palette

The single source of truth for the Emerald Ocean colours: a dark palette taken from a breaking emerald wave, where green dominates the blue.

`palette.json` holds 15 colours in four groups (abyss, foam, kelp, reef). Each colour is defined in OKLCH, with its sRGB hex and what it is used for. The editor themes (`emerald-ocean.nvim`, `emerald-ocean-vscode`) are built from this file.

| Group | Colours |
|-------|---------|
| abyss | `#01130b` `#001c12` `#02291c` `#405d50` |
| foam  | `#eceff1` `#cbdbdf` `#8ca59a` |
| kelp  | `#83eab5` `#97c37f` `#7bd3cf` `#76dffb` `#5b92bf` |
| reef  | `#f7a792` `#df676b` `#baa4e2` |

## Preview

`preview/index.html` shows the palette on a full Go file in a Neovim/tmux mock-up. Open it in a browser.

## Images

- `assets/emerald-wave.png`: a barrel wave recoloured so every pixel sits on the palette.
- `assets/emerald-wave-avatar.jpg`: 800×800 crop of the same image, used as the organisation profile picture.

Wave photo: ["Vortex" by Matt Hardy](https://unsplash.com/photos/zJ-Cj3nT6EU) on Unsplash, used under the Unsplash License.

## License

MIT

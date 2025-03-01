# rlbot gui

A GUI for [RLBot](https://rlbot.org) v5 written in go and powered by [wails](https://wails.io) (v3 alpha).

## Info

This project is built using a justfile. Run `just build` to build, `just dev` to run in dev mode, `just lint` to lint
and `just format` to format. Remember to [install dependencies](#Installing-build-dependencies) first.

## Installing build dependencies

1. Install the [wails v3-alpha cli](https://v3alpha.wails.io/getting-started/installation/)
2. Install [node](https://nodejs.org/en) and [pnpm](https://pnpm.io/installation)
3. Run `cd frontend` and then `pnpm i` (you probably want to `cd ../` after)

## Linux runtime dependencies

* Either [zenity, matedialog or garma](https://github.com/ncruces/zenity?tab=readme-ov-file#benefits-of-the-go-package)

---
name: uhd-cli
description: 'Use uhd CLI to hex dump files with Unicode/multi-byte encoding support. Use when: inspecting binary files, checking text encodings (UTF-8, Shift-JIS, EUC-JP, Big5, etc.), comparing hex output, analyzing file contents, running uhd commands, choosing layout options (jhd/hexdump/bytes).'
argument-hint: 'filename, encoding, layout option, width, etc.'
---

# uhd CLI Guide

`uhd` is a hexdump tool with Unicode and multi-byte encoding support.
It serves as an alternative to `xxd`, `hexdump`, and `jhd`, rendering CJK characters and emoji correctly in the printable column.

## Installation

```sh
# Build from source
go install github.com/wtnb75/uhd@latest

# Pre-built binaries (GitHub Releases)
# https://github.com/wtnb75/uhd/releases

# Whalebrew (container)
whalebrew install ghcr.io/wtnb75/uhd
```

## Basic Usage

### Dump a file

```sh
uhd <filename>
```

### Dump from stdin

```sh
echo hello world | uhd
cat data.bin | uhd
```

### Explicit stdin

```sh
uhd -
```

## Options

| Option | Short | Default | Description |
|---|---|---|---|
| `--encoding` | | `utf-8` | Input text encoding |
| `--width` | | `16` | Bytes per line |
| `--sep` | | `8` | Separator interval (bytes) |
| `--layout` | | `jhd` | Output format (`jhd` / `hexdump` / `bytes`) |
| `--no-color` | | false | Disable color output |
| `--verbose` | `-v` | false | Enable debug logging |
| `--list-codes` | `-l` | | Print supported encodings and exit |

## Layout Options

### `jhd` (default)

Uppercase hex + Unicode-aware printable column. CJK characters and emoji render correctly.

```
00000000  E3 81 93 E3 82 93 E3 81  AB E3 81 A1 E3 82 8F E3    こ_ん_に_ち_わ_、
```

### `hexdump`

Lowercase hex + ASCII-only printable column, surrounded by pipes. Equivalent to `hexdump -C`.

```
00000000  68 65 6c 6c 6f 20 77 6f  72 6c 64 0a               |hello world.|
```

### `bytes`

Each byte displayed as `0xXX`. Useful for C array initialization.

```
00000000  0x68 0x65 0x6c 0x6c 0x6f 0x20 0x77 0x6f 0x72 0x6c 0x64 0x0a    hello world.
```

## Encoding

### List supported encodings

```sh
uhd -l
```

### Common encoding examples

```sh
# Shift-JIS
uhd --encoding shift-jis file.txt

# EUC-JP
uhd --encoding euc-jp file.txt

# Big5 (Traditional Chinese)
uhd --encoding big5 file.txt

# GB18030 (Simplified Chinese)
uhd --encoding gb18030 file.txt

# EUC-KR (Korean)
uhd --encoding euc-kr file.txt

# UTF-16 LE
uhd --encoding utf-16le file.bin
```

### Combine with iconv

```sh
# Convert UTF-8 to Shift-JIS, then dump
echo 'ﾊﾛｰﾜｰﾙﾄﾞ' | iconv -f utf-8 -t shift-jis | uhd --encoding shift-jis
```

## Customizing Output Width

```sh
# 8 bytes per line
uhd --width 8 file.bin

# 32 bytes per line, no separators
uhd --width 32 --sep 32 file.bin

# 24 bytes per line, separator every 8 bytes
uhd --width 24 --sep 8 file.bin
```

## Disabling Color

For non-TTY environments or piped output:

```sh
uhd --no-color file.bin | less
```

## Installing This Skill

The `uhd` binary embeds this `SKILL.md` and can install it to your personal skill directory.

```sh
# Install to ~/.copilot/skills/uhd-cli/ (default)
uhd --install-skill

# Install to ~/.agents/skills/uhd-cli/
uhd --install-skill --skill-target agents

# Install to ~/.claude/skills/uhd-cli/
uhd --install-skill --skill-target claude
```

Once installed, the skill is available as `/uhd-cli` in any VS Code workspace.

## Procedure: Inspecting a Binary File

1. Run `uhd <file>` to see the default output.
2. If text appears garbled, run `uhd -l` to find the correct encoding name.
3. Re-run with `uhd --encoding <enc> <file>`.
4. Switch layout with `--layout hexdump` or `--layout bytes` as needed.
5. Adjust line width with `--width` and `--sep`.

## Procedure: Extending the Codebase

1. Check the `option` struct and flags in `main.go`.
2. Review layout definitions in `get_layout()`.
3. Review the Writer dispatch in `do_uhd()` (e.g. `NewHexdump`, `NewPrintable`).
4. To add a new layout or Writer, update both `get_layout()` and `do_uhd()`.

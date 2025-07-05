# uhd: hexdump supports text encodings

## Install

- pre-built binary: <https://github.com/wtnb75/uhd/releases>
- build from source
    - `go install github.com/wtnb75/uhd@latest`
- container/whalebrew
    - `whalebrew install ghcr.io/wtnb75/uhd`

## Usage

help

```plaintext
# uhd --help
Usage:
  uhd [OPTIONS]

Application Options:
  -v, --verbose                    Enable verbose logging
      --encoding=
      --width=
      --sep=
      --layout=[hexdump|jhd|bytes]
  -l, --list-codes                 list encoding

Help Options:
  -h, --help                       Show this help message
```

simple dump

```plaintext
# echo hello world | uhd
00000000  68 65 6C 6C 6F 20 77 6F  72 6C 64 0A                hello world.
```

utf-8

```plaintext
# echo こんにちわ、🍺はいかがですか | uhd
00000000  E3 81 93 E3 82 93 E3 81  AB E3 81 A1 E3 82 8F E3    こ_ん_に_ち_わ_、
00000010  80 81 F0 9F 8D BA E3 81  AF E3 81 84 E3 81 8B E3    __🍺__は_い_か_が
00000020  81 8C E3 81 A7 E3 81 99  E3 81 8B 0A                __で_す_か_.
# echo 'ﾊﾛｰﾜｰﾙﾄﾞ' | uhd
00000000  EF BE 8A EF BE 9B EF BD  B0 EF BE 9C EF BD B0 EF    ﾊ__ﾛ__ｰ__ﾜ__ｰ__ﾙ
00000010  BE 99 EF BE 84 EF BE 9E  0A                         __ﾄ__ﾞ__.
```

other encodings

```plaintext
# echo 你好 | iconv -f utf-8 -t big5 | uhd --encoding big5
00000000  A7 41 A6 6E 0A                                      你好.
# echo 'ﾊﾛｰﾜｰﾙﾄﾞ' | iconv -f utf-8 -t shift-jis | uhd --encoding shift-jis
00000000  CA DB B0 DC B0 D9 C4 DE  0A                         ﾊﾛｰﾜｰﾙﾄﾞ.
```

# see also

- jhd
- hexdump
- xxd
- od
- [bvi/bview](https://bvi.sourceforge.net/)

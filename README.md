# `qr`

A command-line utility to render scannable QR codes in the terminal.

## Installation

Grab a binary from the `bin` directory, or compile it yourself by running:

```
# Grab a binary from `bin`, or:

make           # just for your architecture
make build-all # for all architectures
```

## Usage

Just pipe data through stdin to render a text-based QR code:

```
echo "hello world" | qr
echo "hello world" | qr -render text # equivalent
```

If you want a prettier one and your terminal supports sixel (most do), run with the `-render sixel` flag:

```
echo "hello world" | qr -render sixel
```

This has the benefit of being more compact and easier to scan.

### Flags

* `-render <"text"|"sixel">`: set the rendering mode, as shown above.
        ```
        qr -render text  # print characters, the default
        qr -render sixel # print an image using sixel
        ```

* `-white <char or color>`: set the white (background) character/color for text/sixel rendering.
        ```
        qr -white ' '           # the default for text
        qr -white '255;255;255' # the default for sixel
        ```

* `-black <char or color>`: set the black (foreground) character/color for text/sixel rendering.
        ```
        qr -black '█'     # the default for text
        qr -black '0;0;0' # the default for sixel
        ```


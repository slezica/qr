package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/qpliu/qrencode-go/qrencode"
)

const MAX_INPUT_BYTES = int64(8 * 1024)

type Config struct {
	black  string     // Black character or sixel RGB string (e.g. '#' or '0;0;0')
	white  string     // White character or sixel RGB string (e.g. ' ' or '255;255;255')
	render RenderFunc // Rendering function, either renderText or renderSixel
}

type RenderName string
type RenderFunc = func(w io.Writer, grid Grid, black string, white string)

const (
	TEXT  RenderName = "text"
	SIXEL RenderName = "sixel"
)

type Grid = *qrencode.BitGrid


// -------------------------------------------------------------------------------------------------

func parseArgs() (*Config, error) {
	config := &Config{}
	// help := flagSet.Bool("help", false, "Show help")

	blackArg := flag.String("black", "", "Black character/color for text/sixel renderer")
	whiteArg := flag.String("white", "", "White character/color for text/sixel renderer")
	renderArg := flag.String("render", "text", "Render method ('text' or 'sixel')")

	flag.Parse()

	// TODO validate

	if *renderArg == "text" {
		config.white = " "
		config.black = "█"
		config.render = renderText

	} else if *renderArg == "sixel" {
		config.white = "255;255;255"
		config.black = "0;0;0"

	} else {
		return nil, fmt.Errorf("invalid renderer name") // TODO
	}

	if *whiteArg != "" {
		config.white = *whiteArg
	}

	if *blackArg != "" {
		config.black = *blackArg
	}

	return config, nil
}


// -------------------------------------------------------------------------------------------------

func renderText(w io.Writer, grid Grid, black string, white string) {
	width := grid.Width()
	height := grid.Height()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if grid.Get(x, y) {
				fmt.Fprintf(w, "%s%s", black, black)
			} else {
				fmt.Fprintf(w, "%s%s", white, white)
			}
		}
		fmt.Fprint(w, "\n")
	}

	fmt.Fprint(w, "\n")
}

func renderSixel(w io.Writer, grid Grid, black string, white string) {
	fmt.Fprint(w, "\n\033Pq") // start sixel mode

	fmt.Fprintf(w, "#0;2;%s", black) // assign #0 to RGB string from `black`
	fmt.Fprintf(w, "#1;2;%s", white) // assign #1 to RGB string from `white`

	width := grid.Width()
	height := grid.Height()
	scale := 2
	padding := 1

	unit := 6 * scale
	size := unit * (width + 2*padding)

	for y := 0; y < padding*scale; y++ {
		fmt.Fprintf(w, "#1!%d~-", size) // in white (#1), repeat (!) %d times a full block (~)
	}

	line := &strings.Builder{}

	for y := 0; y < height; y++ {
		line.Reset()

		fmt.Fprintf(line, "#1!%d~", unit*padding)

		for x := 0; x < width; x++ {
			if grid.Get(x, y) {
				fmt.Fprintf(line, "#0!%d~", unit) // in black (#0), repeat (!) %d times a full block (~)
			} else {
				fmt.Fprintf(line, "#1!%d~", unit) // same, in white (#1)
			}
		}

		fmt.Fprintf(line, "#1!%d~", unit)

		for i := 0; i < scale; i++ {
			fmt.Fprint(w, line.String(), "-") // buffered line, then down (-)
		}
	}

	for y := 0; y < padding*scale; y++ {
		fmt.Fprintf(w, "#1!%d~-", size) // in white (#1), repeat (!) %d times a full block (~)
	}

	fmt.Fprint(w, "\033\\") // end sixel mode
	fmt.Println()
}


// -------------------------------------------------------------------------------------------------

func readLimitOrFail(r io.Reader, n int64) ([]byte, error) {
  reader := &io.LimitedReader{R: r, N: n+1} // extra byte to detect excess data

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

  if reader.N > n {
    return nil, fmt.Errorf("Input longer than the maximum of %d bytes\n", n)
  }

  return data, nil
}


// -------------------------------------------------------------------------------------------------

func main() {
	config, err := parseArgs()
	if err != nil {
    panic(err)
	}

  data, err := readLimitOrFail(os.Stdin, MAX_INPUT_BYTES)
  if err != nil {
    panic(err)
	}

	grid, err := qrencode.EncodeBytes(data, qrencode.ECLevelL)
	if err != nil {
		panic(err)
	}

	config.render(os.Stdout, grid, config.black, config.white)
}

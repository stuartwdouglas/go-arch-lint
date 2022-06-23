package code

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"

	"github.com/fe3dback/go-arch-lint/internal/models"
)

type Render struct {
	printer ColorPrinter
}

type annotateOpts struct {
	code              []byte
	ref               models.CodeReference
	region            codeRegion
	showColumnPointer bool
}

func NewRender(printer ColorPrinter) *Render {
	return &Render{
		printer: printer,
	}
}

func (r *Render) SourceCode(ref models.CodeReference, highlight bool) []byte {
	code, region := r.readCode(ref, highlight)
	return r.annotate(annotateOpts{
		code:              code,
		ref:               ref,
		region:            region,
		showColumnPointer: true,
	})
}

func (r *Render) SourceCodeWithoutOffset(ref models.CodeReference, highlight bool) []byte {
	code, region := r.readCode(ref, highlight)
	return r.annotate(annotateOpts{
		code:              code,
		ref:               ref,
		region:            region,
		showColumnPointer: false,
	})
}

func (r *Render) readCode(ref models.CodeReference, highlight bool) ([]byte, codeRegion) {
	if !ref.Pointer.Valid {
		return []byte{}, codeRegion{}
	}

	rawCode, region := r.readRaw(ref)
	if !highlight {
		return rawCode, region
	}

	return highlightRawCode(ref, rawCode), region
}

func (r *Render) readRaw(ref models.CodeReference) ([]byte, codeRegion) {
	if !ref.Pointer.Valid {
		return []byte{}, codeRegion{}
	}

	file, err := os.Open(ref.Pointer.File)
	if err != nil {
		return []byte{}, codeRegion{}
	}

	linesCount, err := lineCounter(file)
	if err != nil {
		return []byte{}, codeRegion{}
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return []byte{}, codeRegion{}
	}

	region := calculateCodeRegion(ref, linesCount)
	return readLines(file, region), region
}

func highlightRawCode(ref models.CodeReference, code []byte) []byte {
	lexer := lexers.Match(ref.Pointer.File)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Trac
	formatter := formatters.TTY8

	iterator, err := lexer.Tokenise(nil, string(code))
	if err != nil {
		return []byte{}
	}

	var buf bytes.Buffer
	err = formatter.Format(&buf, style, iterator)
	if err != nil {
		return []byte{}
	}

	return buf.Bytes()
}

func readLines(r io.Reader, region codeRegion) []byte {
	sc := bufio.NewScanner(r)
	currentLine := 0
	var buffer bytes.Buffer

	for sc.Scan() {
		currentLine++

		if currentLine >= region.lineFirst && currentLine <= region.lineLast {
			buffer.Write(sc.Bytes())

			if currentLine != region.lineLast {
				buffer.WriteByte('\n')
			}
		}
	}

	return buffer.Bytes()
}

func (r *Render) annotate(opt annotateOpts) []byte {
	buf := bytes.NewBuffer(opt.code)
	sc := bufio.NewScanner(buf)
	currentLine := opt.region.lineFirst

	var resultBuffer bytes.Buffer
	for sc.Scan() {
		prefixLine := r.printer.Gray(fmt.Sprintf("%4d |", currentLine))
		prefixEmpty := r.printer.Gray("        ")

		// add line pointer
		if currentLine == opt.region.lineMain {
			prefixLine = fmt.Sprintf("> %s", prefixLine)
		} else {
			prefixLine = fmt.Sprintf("  %s", prefixLine)
		}

		// draw line
		resultBuffer.WriteString(fmt.Sprintf("%s %s\n",
			prefixLine,
			r.replaceTabsToSpaces(sc.Bytes()),
		))

		// add offset pointer
		if opt.showColumnPointer {
			if currentLine == opt.region.lineMain && opt.ref.Pointer.Valid {
				spaces := strings.Repeat(" ", int(math.Max(0, float64(opt.ref.Pointer.Offset-1))))
				resultBuffer.WriteString(fmt.Sprintf("%s %s^\n", prefixEmpty, spaces))
			}
		}

		currentLine++
	}

	return resultBuffer.Bytes()
}

func (r *Render) replaceTabsToSpaces(src []byte) []byte {
	return []byte(strings.ReplaceAll(string(src), "\t", "  "))
}
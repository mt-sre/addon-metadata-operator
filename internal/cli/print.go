package cli

import (
	"errors"
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/fatih/color"
)

var ErrInvalidTableStyle = errors.New("invalid table style")

type TableStyle string

func (ts TableStyle) ToSimpleTableStyle() (*simpletable.Style, error) {
	switch ts {
	case TableStyleCompactLite:
		return simpletable.StyleCompactLite, nil
	default:
		return nil, ErrInvalidTableStyle
	}
}

const (
	TableStyleCompactLite TableStyle = "compactLite"
)

type TableHeader []string

func (th TableHeader) ToSimpleTableHeader() *simpletable.Header {
	var res simpletable.Header

	for _, col := range th {
		res.Cells = append(res.Cells,
			&simpletable.Cell{Align: simpletable.AlignCenter, Text: col},
		)
	}

	return &res
}

func NewTable(opts ...TableOption) (*Table, error) {
	var cfg TableConfig

	if err := cfg.Option(opts...); err != nil {
		return nil, fmt.Errorf("applying options: %w", err)
	}

	cfg.Default()

	table := simpletable.New()
	table.SetStyle(cfg.Style)
	table.Header = cfg.Header

	return &Table{
		cfg: cfg,
		t:   table,
	}, nil
}

type Table struct {
	cfg TableConfig
	t   *simpletable.Table
}

func (t *Table) String() string {
	return t.t.String()
}

func (t *Table) WriteRow(row TableRow) {
	t.t.Body.Cells = append(t.t.Body.Cells, row.ToCells())
}

type TableRow []Field

func (r TableRow) ToCells() []*simpletable.Cell {
	res := make([]*simpletable.Cell, 0, len(r))

	for _, f := range r {
		res = append(res, f.ToCell())
	}

	return res
}

type Field struct {
	Value string
	Color FieldColor
}

func (f Field) ToCell() *simpletable.Cell {
	return &simpletable.Cell{
		Align: simpletable.AlignLeft,
		Text:  f.Color.Apply(f.Value),
	}
}

type FieldColor string

func (fc FieldColor) Apply(s string) string {
	switch fc {
	case FieldColorGreen:
		return green(s)
	case FieldColorRed:
		return red(s)
	case FieldColorIntenselyBoldRed:
		return intenselyBoldRed(s)
	default:
		return s
	}
}

const (
	FieldColorGreen            FieldColor = "green"
	FieldColorRed              FieldColor = "red"
	FieldColorIntenselyBoldRed FieldColor = "intenselyBoldRed"
)

var (
	green            = color.New(color.FgGreen).SprintFunc()
	red              = color.New(color.FgRed).SprintFunc()
	intenselyBoldRed = color.New(color.Bold, color.FgHiRed).SprintFunc()
)

type TableConfig struct {
	Header *simpletable.Header
	Style  *simpletable.Style
}

func (c *TableConfig) Option(opts ...TableOption) error {
	for _, opt := range opts {
		if err := opt.ConfigureTable(c); err != nil {
			return fmt.Errorf("configuring table: %w", err)
		}
	}

	return nil
}

func (c *TableConfig) Default() {
	if c.Style == nil {
		c.Style = simpletable.StyleCompactLite
	}
}

type TableOption interface {
	ConfigureTable(*TableConfig) error
}

type WithHeaders TableHeader

func (h WithHeaders) ConfigureTable(c *TableConfig) error {
	c.Header = TableHeader(h).ToSimpleTableHeader()

	return nil
}

type WithStyle TableStyle

func (s WithStyle) ConfigureTable(c *TableConfig) error {
	style, err := TableStyle(s).ToSimpleTableStyle()
	if err != nil {
		return fmt.Errorf("parsing table style %q: %w", s, err)
	}

	c.Style = style

	return nil
}

// PrintValidationErrors - helper to pretty print validationErrors
func PrintValidationErrors(errs []error) {
	fmt.Printf("\n%s\n", red("Failed with the following errors:"))
	for _, err := range errs {
		fmt.Printf("%s\n", err.Error())
	}
}

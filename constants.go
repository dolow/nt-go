package ntgo

const (
	DictionaryKeySeparator = ':'
	ListSymbol             = '-'
	TextSymbol             = '>'
	CommentSymbol          = '#'
	IndentChar             = ' '
	Space                  = ' '
	Tab                    = '\t'
	CR                     = '\r'
	LF                     = '\n'
	Quote                  = '\''
	DoubleQuote            = '"'

	EmptyChar     byte = 0x00
	NotFoundIndex int  = -1

	MarshallerTag = "nt"
	MarshallerTagSeparator = ","
	MarshallerTagMultilineText = "multilinetext"
	UnmarshalDefaultIndentSize = 2
)

type MultiLineText []string

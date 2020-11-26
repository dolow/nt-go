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

	MarshallerTag              = "nt"
	MarshallerTagSeparator     = ","
	MarshallerTagOmitEmpty     = "omitempty"
	MarshallerTagMultilineText = "multilinetext"
	UnmarshalDefaultIndentSize = 2

	MarshallerTagFlagOmitEmpty     = 1
	MarshallerTagFlagMultilineText = 1 << 1
)

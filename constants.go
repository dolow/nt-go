package ntgo

const (
	DictionaryKeySeparator = ':'
	ListToken              = '-'
	TextToken              = '>'
	CommentToken           = '#'
	IndentChar             = ' '
	Space                  = ' '
	Tab                    = '\t'
	CR                     = '\r'
	LF                     = '\n'
	Quote                  = '\''
	DoubleQuote            = '"'

	EmptyChar     byte = 0x00
	NotFoundIndex int  = -1

	MarshallerTag                 = "nt"
	MarshallerTagSeparator        = ","
	MarshallerTagOmitEmpty        = "omitempty"
	MarshallerTagMultilineStrings = "multilinestrings"
	UnmarshalDefaultIndentSize    = 2

	MarshallerTagFlagOmitEmpty = 1 << iota
	MarshallerTagFlagMultilineStrings
)

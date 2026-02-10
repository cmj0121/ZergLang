package evaluator

import (
	"io"
	"os"
	"runtime"
	"strings"
	"unicode"
)

// BuiltinModules is the map of all builtin modules.
var BuiltinModules = map[string]*Module{
	"_sys": SysModule,
	"_io":  IoModule,
	"str":  StrModule,
	"char": CharModule,
}

// SysModule provides system-related functions.
var SysModule = &Module{
	Name: "sys",
	Methods: map[string]*Builtin{
		"os":   {Name: "os", Fn: sysOs},
		"arch": {Name: "arch", Fn: sysArch},
		"args": {Name: "args", Fn: sysArgs},
		"exit": {Name: "exit", Fn: sysExit},
		"env":  {Name: "env", Fn: sysEnv},
	},
}

// IoModule provides I/O functions.
// File methods (read, write, seek, tell, close) are handled by GetFileMethod.
var IoModule = &Module{
	Name: "_io",
	Methods: map[string]*Builtin{
		"open":   {Name: "open", Fn: ioOpen},
		"exists": {Name: "exists", Fn: ioExists},
	},
}

// StrModule provides string manipulation functions.
var StrModule = &Module{
	Name: "str",
	Methods: map[string]*Builtin{
		"split":       {Name: "split", Fn: strSplit},
		"join":        {Name: "join", Fn: strJoin},
		"trim":        {Name: "trim", Fn: strTrim},
		"find":        {Name: "find", Fn: strFind},
		"replace":     {Name: "replace", Fn: strReplace},
		"substring":   {Name: "substring", Fn: strSubstring},
		"starts_with": {Name: "starts_with", Fn: strStartsWith},
		"ends_with":   {Name: "ends_with", Fn: strEndsWith},
		"upper":       {Name: "upper", Fn: strUpper},
		"lower":       {Name: "lower", Fn: strLower},
		"contains":    {Name: "contains", Fn: strContains},
	},
}

// CharModule provides character manipulation functions.
var CharModule = &Module{
	Name: "char",
	Methods: map[string]*Builtin{
		"ord":      {Name: "ord", Fn: charOrd},
		"chr":      {Name: "chr", Fn: charChr},
		"is_digit": {Name: "is_digit", Fn: charIsDigit},
		"is_alpha": {Name: "is_alpha", Fn: charIsAlpha},
		"is_space": {Name: "is_space", Fn: charIsSpace},
		"is_alnum": {Name: "is_alnum", Fn: charIsAlnum},
	},
}

// ============================================================================
// sys module implementations
// ============================================================================

func sysOs(args ...Object) Object {
	if len(args) != 0 {
		return newError("sys.os() takes no arguments (%d given)", len(args))
	}
	return &String{Value: runtime.GOOS}
}

func sysArch(args ...Object) Object {
	if len(args) != 0 {
		return newError("sys.arch() takes no arguments (%d given)", len(args))
	}
	return &String{Value: runtime.GOARCH}
}

// programArgs stores the command-line arguments for sys.args()
var programArgs []string

// SetProgramArgs sets the command-line arguments for sys.args()
func SetProgramArgs(args []string) {
	programArgs = args
}

func sysArgs(args ...Object) Object {
	if len(args) != 0 {
		return newError("sys.args() takes no arguments (%d given)", len(args))
	}

	// Use stored program args, or fall back to os.Args
	argsToUse := programArgs
	if argsToUse == nil {
		argsToUse = os.Args
	}

	elements := make([]Object, len(argsToUse))
	for i, arg := range argsToUse {
		elements[i] = &String{Value: arg}
	}
	return &List{Elements: elements}
}

func sysExit(args ...Object) Object {
	if len(args) != 1 {
		return newError("sys.exit() takes exactly 1 argument (%d given)", len(args))
	}

	code, ok := args[0].(*Integer)
	if !ok {
		return newError("sys.exit() argument must be an integer, not %s", args[0].Type())
	}

	os.Exit(int(code.Value))
	return NULL // never reached
}

func sysEnv(args ...Object) Object {
	if len(args) != 1 {
		return newError("sys.env() takes exactly 1 argument (%d given)", len(args))
	}

	name, ok := args[0].(*String)
	if !ok {
		return newError("sys.env() argument must be a string, not %s", args[0].Type())
	}

	return &String{Value: os.Getenv(name.Value)}
}

// ============================================================================
// io module implementations
// ============================================================================

func ioOpen(args ...Object) Object {
	if len(args) < 1 || len(args) > 2 {
		return newError("io.open() takes 1 or 2 arguments (%d given)", len(args))
	}

	path, ok := args[0].(*String)
	if !ok {
		return newError("io.open() path must be a string, not %s", args[0].Type())
	}

	mode := "r"
	if len(args) == 2 {
		modeStr, ok := args[1].(*String)
		if !ok {
			return newError("io.open() mode must be a string, not %s", args[1].Type())
		}
		mode = modeStr.Value
	}

	var flag int
	var perm os.FileMode = 0644

	switch mode {
	case "r":
		flag = os.O_RDONLY
	case "w":
		flag = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	case "a":
		flag = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	case "rw", "wr":
		flag = os.O_RDWR | os.O_CREATE
	default:
		return newError("io.open() invalid mode: %s (use 'r', 'w', 'a', or 'rw')", mode)
	}

	file, err := os.OpenFile(path.Value, flag, perm)
	if err != nil {
		return newError("io.open() failed: %s", err.Error())
	}

	return &File{Path: path.Value, Mode: mode, Handle: file}
}

func ioRead(args ...Object) Object {
	if len(args) != 1 {
		return newError("io.read() takes exactly 1 argument (%d given)", len(args))
	}

	fh, ok := args[0].(*File)
	if !ok {
		return newError("io.read() argument must be a file handle, not %s", args[0].Type())
	}

	file, ok := fh.Handle.(*os.File)
	if !ok {
		return newError("io.read() invalid file handle")
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return newError("io.read() failed: %s", err.Error())
	}

	return &String{Value: string(content)}
}

func ioReadLines(args ...Object) Object {
	if len(args) != 1 {
		return newError("io.read_lines() takes exactly 1 argument (%d given)", len(args))
	}

	fh, ok := args[0].(*File)
	if !ok {
		return newError("io.read_lines() argument must be a file handle, not %s", args[0].Type())
	}

	file, ok := fh.Handle.(*os.File)
	if !ok {
		return newError("io.read_lines() invalid file handle")
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return newError("io.read_lines() failed: %s", err.Error())
	}

	lines := strings.Split(string(content), "\n")
	elements := make([]Object, len(lines))
	for i, line := range lines {
		elements[i] = &String{Value: line}
	}

	return &List{Elements: elements}
}

func ioWrite(args ...Object) Object {
	if len(args) != 2 {
		return newError("io.write() takes exactly 2 arguments (%d given)", len(args))
	}

	fh, ok := args[0].(*File)
	if !ok {
		return newError("io.write() first argument must be a file handle, not %s", args[0].Type())
	}

	data, ok := args[1].(*String)
	if !ok {
		return newError("io.write() second argument must be a string, not %s", args[1].Type())
	}

	file, ok := fh.Handle.(*os.File)
	if !ok {
		return newError("io.write() invalid file handle")
	}

	n, err := file.WriteString(data.Value)
	if err != nil {
		return newError("io.write() failed: %s", err.Error())
	}

	return &Integer{Value: int64(n)}
}

func ioClose(args ...Object) Object {
	if len(args) != 1 {
		return newError("io.close() takes exactly 1 argument (%d given)", len(args))
	}

	fh, ok := args[0].(*File)
	if !ok {
		return newError("io.close() argument must be a file handle, not %s", args[0].Type())
	}

	file, ok := fh.Handle.(*os.File)
	if !ok {
		return newError("io.close() invalid file handle")
	}

	if err := file.Close(); err != nil {
		return newError("io.close() failed: %s", err.Error())
	}

	return NULL
}

func ioExists(args ...Object) Object {
	if len(args) != 1 {
		return newError("io.exists() takes exactly 1 argument (%d given)", len(args))
	}

	path, ok := args[0].(*String)
	if !ok {
		return newError("io.exists() argument must be a string, not %s", args[0].Type())
	}

	_, err := os.Stat(path.Value)
	return nativeBoolToBooleanObject(err == nil)
}

func ioReadFile(args ...Object) Object {
	if len(args) != 1 {
		return newError("io.read_file() takes exactly 1 argument (%d given)", len(args))
	}

	path, ok := args[0].(*String)
	if !ok {
		return newError("io.read_file() argument must be a string, not %s", args[0].Type())
	}

	content, err := os.ReadFile(path.Value)
	if err != nil {
		return newError("io.read_file() failed: %s", err.Error())
	}

	return &String{Value: string(content)}
}

func ioWriteFile(args ...Object) Object {
	if len(args) != 2 {
		return newError("io.write_file() takes exactly 2 arguments (%d given)", len(args))
	}

	path, ok := args[0].(*String)
	if !ok {
		return newError("io.write_file() first argument must be a string, not %s", args[0].Type())
	}

	data, ok := args[1].(*String)
	if !ok {
		return newError("io.write_file() second argument must be a string, not %s", args[1].Type())
	}

	if err := os.WriteFile(path.Value, []byte(data.Value), 0644); err != nil {
		return newError("io.write_file() failed: %s", err.Error())
	}

	return NULL
}

// ============================================================================
// Additional asm functions for self-hosted evaluator
// ============================================================================

// asmFileReadN reads exactly n bytes from a file handle, returns string.
func asmFileReadN(args ...Object) Object {
	if len(args) < 2 {
		return newError("file_read_n requires 2 arguments (handle, n)")
	}
	fh, ok := args[0].(*File)
	if !ok {
		return newError("file_read_n first argument must be a file handle")
	}
	file, ok := fh.Handle.(*os.File)
	if !ok {
		return newError("file_read_n invalid file handle")
	}
	n, ok := args[1].(*Integer)
	if !ok {
		return newError("file_read_n second argument must be an integer")
	}
	data := make([]byte, n.Value)
	bytesRead, err := file.Read(data)
	if err != nil && err != io.EOF {
		return newError("file_read_n failed: %s", err.Error())
	}
	return &String{Value: string(data[:bytesRead])}
}

// asmFileSeek seeks to a position in a file handle.
func asmFileSeek(args ...Object) Object {
	if len(args) < 2 {
		return newError("file_seek requires at least 2 arguments (handle, offset)")
	}
	fh, ok := args[0].(*File)
	if !ok {
		return newError("file_seek first argument must be a file handle")
	}
	file, ok := fh.Handle.(*os.File)
	if !ok {
		return newError("file_seek invalid file handle")
	}
	offset, ok := args[1].(*Integer)
	if !ok {
		return newError("file_seek offset must be an integer")
	}
	whence := 0
	if len(args) > 2 {
		w, ok := args[2].(*Integer)
		if !ok {
			return newError("file_seek whence must be an integer")
		}
		whence = int(w.Value)
	}
	pos, err := file.Seek(offset.Value, whence)
	if err != nil {
		return newError("file_seek failed: %s", err.Error())
	}
	return &Integer{Value: pos}
}

// asmFileTell returns the current position in a file handle.
func asmFileTell(args ...Object) Object {
	if len(args) < 1 {
		return newError("file_tell requires 1 argument (handle)")
	}
	fh, ok := args[0].(*File)
	if !ok {
		return newError("file_tell argument must be a file handle")
	}
	file, ok := fh.Handle.(*os.File)
	if !ok {
		return newError("file_tell invalid file handle")
	}
	pos, err := file.Seek(0, 1)
	if err != nil {
		return newError("file_tell failed: %s", err.Error())
	}
	return &Integer{Value: pos}
}

// asmByteChr converts an integer (0-255) to a single raw byte string.
// Unlike char_chr which does UTF-8 encoding, this produces exactly 1 byte.
func asmByteChr(args ...Object) Object {
	if len(args) != 1 {
		return newError("byte_chr takes exactly 1 argument (%d given)", len(args))
	}
	code, ok := args[0].(*Integer)
	if !ok {
		return newError("byte_chr argument must be an integer, not %s", args[0].Type())
	}
	return &String{Value: string([]byte{byte(code.Value)})}
}

// ============================================================================
// str module implementations
// ============================================================================

func strSplit(args ...Object) Object {
	if len(args) != 2 {
		return newError("str.split() takes exactly 2 arguments (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("str.split() first argument must be a string, not %s", args[0].Type())
	}

	sep, ok := args[1].(*String)
	if !ok {
		return newError("str.split() second argument must be a string, not %s", args[1].Type())
	}

	parts := strings.Split(s.Value, sep.Value)
	elements := make([]Object, len(parts))
	for i, part := range parts {
		elements[i] = &String{Value: part}
	}

	return &List{Elements: elements}
}

func strJoin(args ...Object) Object {
	if len(args) != 2 {
		return newError("str.join() takes exactly 2 arguments (%d given)", len(args))
	}

	list, ok := args[0].(*List)
	if !ok {
		return newError("str.join() first argument must be a list, not %s", args[0].Type())
	}

	sep, ok := args[1].(*String)
	if !ok {
		return newError("str.join() second argument must be a string, not %s", args[1].Type())
	}

	strs := make([]string, len(list.Elements))
	for i, el := range list.Elements {
		strs[i] = el.Inspect()
	}

	return &String{Value: strings.Join(strs, sep.Value)}
}

func strTrim(args ...Object) Object {
	if len(args) != 1 {
		return newError("str.trim() takes exactly 1 argument (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("str.trim() argument must be a string, not %s", args[0].Type())
	}

	return &String{Value: strings.TrimSpace(s.Value)}
}

func strFind(args ...Object) Object {
	if len(args) != 2 {
		return newError("str.find() takes exactly 2 arguments (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("str.find() first argument must be a string, not %s", args[0].Type())
	}

	sub, ok := args[1].(*String)
	if !ok {
		return newError("str.find() second argument must be a string, not %s", args[1].Type())
	}

	return &Integer{Value: int64(strings.Index(s.Value, sub.Value))}
}

func strReplace(args ...Object) Object {
	if len(args) != 3 {
		return newError("str.replace() takes exactly 3 arguments (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("str.replace() first argument must be a string, not %s", args[0].Type())
	}

	old, ok := args[1].(*String)
	if !ok {
		return newError("str.replace() second argument must be a string, not %s", args[1].Type())
	}

	newStr, ok := args[2].(*String)
	if !ok {
		return newError("str.replace() third argument must be a string, not %s", args[2].Type())
	}

	return &String{Value: strings.ReplaceAll(s.Value, old.Value, newStr.Value)}
}

func strSubstring(args ...Object) Object {
	if len(args) != 3 {
		return newError("str.substring() takes exactly 3 arguments (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("str.substring() first argument must be a string, not %s", args[0].Type())
	}

	start, ok := args[1].(*Integer)
	if !ok {
		return newError("str.substring() second argument must be an integer, not %s", args[1].Type())
	}

	end, ok := args[2].(*Integer)
	if !ok {
		return newError("str.substring() third argument must be an integer, not %s", args[2].Type())
	}

	str := s.Value
	startIdx := int(start.Value)
	endIdx := int(end.Value)

	if startIdx < 0 {
		startIdx = 0
	}
	if endIdx > len(str) {
		endIdx = len(str)
	}
	if startIdx > endIdx {
		return &String{Value: ""}
	}

	return &String{Value: str[startIdx:endIdx]}
}

func strStartsWith(args ...Object) Object {
	if len(args) != 2 {
		return newError("str.starts_with() takes exactly 2 arguments (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("str.starts_with() first argument must be a string, not %s", args[0].Type())
	}

	prefix, ok := args[1].(*String)
	if !ok {
		return newError("str.starts_with() second argument must be a string, not %s", args[1].Type())
	}

	return nativeBoolToBooleanObject(strings.HasPrefix(s.Value, prefix.Value))
}

func strEndsWith(args ...Object) Object {
	if len(args) != 2 {
		return newError("str.ends_with() takes exactly 2 arguments (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("str.ends_with() first argument must be a string, not %s", args[0].Type())
	}

	suffix, ok := args[1].(*String)
	if !ok {
		return newError("str.ends_with() second argument must be a string, not %s", args[1].Type())
	}

	return nativeBoolToBooleanObject(strings.HasSuffix(s.Value, suffix.Value))
}

func strUpper(args ...Object) Object {
	if len(args) != 1 {
		return newError("str.upper() takes exactly 1 argument (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("str.upper() argument must be a string, not %s", args[0].Type())
	}

	return &String{Value: strings.ToUpper(s.Value)}
}

func strLower(args ...Object) Object {
	if len(args) != 1 {
		return newError("str.lower() takes exactly 1 argument (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("str.lower() argument must be a string, not %s", args[0].Type())
	}

	return &String{Value: strings.ToLower(s.Value)}
}

func strContains(args ...Object) Object {
	if len(args) != 2 {
		return newError("str.contains() takes exactly 2 arguments (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("str.contains() first argument must be a string, not %s", args[0].Type())
	}

	sub, ok := args[1].(*String)
	if !ok {
		return newError("str.contains() second argument must be a string, not %s", args[1].Type())
	}

	return nativeBoolToBooleanObject(strings.Contains(s.Value, sub.Value))
}

// ============================================================================
// char module implementations
// ============================================================================

func charOrd(args ...Object) Object {
	if len(args) != 1 {
		return newError("char.ord() takes exactly 1 argument (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("char.ord() argument must be a string, not %s", args[0].Type())
	}

	if len(s.Value) == 0 {
		return newError("char.ord() argument must be a non-empty string")
	}

	return &Integer{Value: int64(s.Value[0])}
}

func charChr(args ...Object) Object {
	if len(args) != 1 {
		return newError("char.chr() takes exactly 1 argument (%d given)", len(args))
	}

	code, ok := args[0].(*Integer)
	if !ok {
		return newError("char.chr() argument must be an integer, not %s", args[0].Type())
	}

	return &String{Value: string(rune(code.Value))}
}

func charIsDigit(args ...Object) Object {
	if len(args) != 1 {
		return newError("char.is_digit() takes exactly 1 argument (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("char.is_digit() argument must be a string, not %s", args[0].Type())
	}

	if len(s.Value) == 0 {
		return FALSE
	}

	return nativeBoolToBooleanObject(unicode.IsDigit(rune(s.Value[0])))
}

func charIsAlpha(args ...Object) Object {
	if len(args) != 1 {
		return newError("char.is_alpha() takes exactly 1 argument (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("char.is_alpha() argument must be a string, not %s", args[0].Type())
	}

	if len(s.Value) == 0 {
		return FALSE
	}

	return nativeBoolToBooleanObject(unicode.IsLetter(rune(s.Value[0])))
}

func charIsSpace(args ...Object) Object {
	if len(args) != 1 {
		return newError("char.is_space() takes exactly 1 argument (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("char.is_space() argument must be a string, not %s", args[0].Type())
	}

	if len(s.Value) == 0 {
		return FALSE
	}

	return nativeBoolToBooleanObject(unicode.IsSpace(rune(s.Value[0])))
}

func charIsAlnum(args ...Object) Object {
	if len(args) != 1 {
		return newError("char.is_alnum() takes exactly 1 argument (%d given)", len(args))
	}

	s, ok := args[0].(*String)
	if !ok {
		return newError("char.is_alnum() argument must be a string, not %s", args[0].Type())
	}

	if len(s.Value) == 0 {
		return FALSE
	}

	r := rune(s.Value[0])
	return nativeBoolToBooleanObject(unicode.IsLetter(r) || unicode.IsDigit(r))
}

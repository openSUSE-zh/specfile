package specfile

import (
	"strings"
	"unicode"
)

var (
	dependencyMap         = map[string]struct{}{"BuildRequires": {}, "Requires": {}, "Obsoletes": {}, "Provides": {}, "Conflicts": {}, "Recommends": {}, "Supplements": {}, "Suggests": {}, "Enhances": {}}
	sectionMap            = map[string]struct{}{"%package": {}, "%prep": {}, "generate_buildrequires": {}, "%build": {}, "%install": {}, "%check": {}, "%clean": {}, "%preun": {}, "%postun": {}, "%pretrans": {}, "%posttrans": {}, "%pre": {}, "%post": {}, "%files": {}, "%changelog": {}, "%description": {}, "%triggerpostun": {}, "%triggerprein": {}, "%triggerun": {}, "%triggerin": {}, "%trigger": {}, "%verifyscript": {}, "%sepolicy": {}, "%filetriggerin": {}, "%filetrigger": {}, "%filetriggerun": {}, "%filetriggerpostun": {}, "%transfiletriggerin": {}, "%transfiletrigger": {}, "%transfiletriggerun": {}, "%transfiletriggerpostun": {}, "%patchlist": {}, "%sourcelist": {}}
	conditionalIndicators = []string{"%if", "%else", "%elif", "%end"}
)

// Line information for a syntactically valid line
type Line struct {
	Lines  []string
	Last   string
	Len    int
	Offset int64
}

// NewLine initialize a new Line
func NewLine(offset int64, lines ...string) Line {
	last := ""
	if len(lines) > 0 {
		last = lines[len(lines)-1]
	}
	return Line{lines, last, len(lines), offset}
}

// isConditional if a Line is a conditional line
func (line Line) isConditional() bool {
	for _, c := range conditionalIndicators {
		if strings.HasPrefix(line.Last, c) {
			return true
		}
	}
	return false
}

// isSection if the line is a specfile section like %build, %install
func (line Line) isSection() bool {
	for section := range sectionMap {
		// section indicator must be itself like "%files" or with whitespaces like "%files -n"
		// or "%install" will match "%install_info"
		if strings.HasPrefix(line.Last, section+"\n") || strings.HasPrefix(line.Last, section+" ") {
			return true
		}
	}
	return false
}

// isMacro if the line is a rpm macro like "%define fcitx5_version 5.0.1"
func (line Line) isMacro() bool {
	if strings.HasPrefix(line.Last, "%define") || strings.HasPrefix(line.Last, "%global") {
		return true
	}
	return false
}

// isDependency if the line is a dependency tag
func (line Line) isDependency() bool {
	for k := range dependencyMap {
		if strings.HasPrefix(line.Last, k) {
			return true
		}
	}
	return false
}

// isTag if the line is a rpm tag like "Name: xz"
func (line Line) isTag() bool {
	r := []rune(line.Last)
	if len(r) == 0 {
		return false
	}
	return unicode.IsUpper([]rune(line.Last)[0])
}

// Concat prepend or append lines of string to Line struct
func (line *Line) Concat(prepend bool, lines ...string) {
	old := line.Lines
	new := lines

	if prepend {
		old = lines
		new = line.Lines
	}

	line.Len = line.Len + len(lines)

	for _, v := range new {
		old = append(old, v)
	}

	line.Lines = old
	line.Last = line.Lines[len(line.Lines)-1]
}

package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var rgx = regexp.MustCompile(`\((.*?)\)`)
var pre = regexp.MustCompile(`.*?\((.*?)\)`)
var sig = regexp.MustCompile(`.*?\:\((.*?)\)\:\((.*?)\)`)
var noRetSig = regexp.MustCompile(`.*?\:\((.*?)\)`)

// Signature is a parsed signature
type Signature struct {
	Method string
	Args   []string
	Return []string
}

// Parser -
type Parser struct{}

// NewParser -
func NewParser() *Parser {
	return &Parser{}
}

// Parse -
func (p *Parser) Parse(name, signature string, hasReturns bool) string {
	args := extractArgs(signature)
	var rts []string
	if hasReturns {
		rts = extractReturnTypes(signature)
	}
	svcSig := formatDSL(name, args, rts)
	return svcSig
}

// Unmarshal -
func (p *Parser) Unmarshal(signature string) Signature {
	return parseSignature(signature)
}

func formatDSL(method string, args, returns []string) string {
	argsStr := strings.Join(removeMainPrefix(args[1:]), ", ")
	if len(returns) == 0 {
		return fmt.Sprintf("%s:(%s)", method, argsStr)
	}
	// srvName := strings.TrimPrefix(args[0], "*main.") // Don't need this, yet.
	returnsStr := strings.Join(removeMainPrefix(returns), ", ")
	return fmt.Sprintf("%s:(%s):(%s)", method, argsStr, returnsStr)
}

func extractArgs(args string) []string {
	rs := rgx.FindAllString(args, -1)[0]
	rs = strings.Replace(rs, "main.", "", -1)
	rs = strings.TrimPrefix(rs, "(")
	rs = strings.TrimSuffix(rs, ")")
	return strings.Split(rs, ", ")
}

func extractReturnTypes(args string) []string {
	res := pre.FindAllString(args, -1)[0]
	res = strings.TrimPrefix(args, res+" ")
	res = strings.Replace(res, "main.", "", -1)
	if strings.HasPrefix(res, "(") {
		return extractArgs(res)
	}
	return []string{res}
}

// parseSignature takes a signature and returns a map
// of the method and arguments
func parseSignature(signature string) Signature {
	rs := strings.Split(signature, ":")
	count := len(rs)
	hasReturns := true
	argGroup := rs[len(rs)-2]
	if count <= 2 {
		hasReturns = false
	}

	// If this signature has no return type,
	// args will be in a different position
	if hasReturns == false {
		argGroup = rs[len(rs)-1]
	}

	argGroup = strings.TrimPrefix(argGroup, "(")
	argGroup = strings.TrimSuffix(argGroup, ")")
	args := strings.Split(argGroup, ", ")

	var rets []string
	if hasReturns {
		retGroup := rs[len(rs)-1:][0]
		ret := strings.TrimPrefix(retGroup, "(")
		ret = strings.TrimSuffix(ret, ")")
		rets = strings.Split(ret, ", ")
	}

	return Signature{
		Method: rs[0],
		Args:   args,
		Return: rets,
	}
}

func removeMainPrefix(types []string) []string {
	res := []string{}
	for _, t := range types {
		if strings.HasPrefix(t, "main.") {
			t = strings.TrimPrefix(t, "main.")
		}
		res = append(res, t)
	}
	return res
}

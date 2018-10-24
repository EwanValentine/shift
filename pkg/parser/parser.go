package parser

import (
	"fmt"
	"regexp"
	"strings"
)

var rgx = regexp.MustCompile(`\((.*?)\)`)
var pre = regexp.MustCompile(`.*?\((.*?)\)`)
var sig = regexp.MustCompile(`.*?\:\((.*?)\)\:\((.*?)\)`)

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
func (p *Parser) Parse(name, signature string) string {
	args := extractArgs(signature)
	rts := extractReturnTypes(signature)
	svcSig := formatDSL(name, args, rts)
	return svcSig
}

// Unmarshal -
func (p *Parser) Unmarshal(signature string) Signature {
	return parseSignature(signature)
}

func formatDSL(method string, args, returns []string) string {
	argsStr := strings.Join(removeMainPrefix(args[1:]), ", ")
	returnsStr := strings.Join(removeMainPrefix(returns), ", ")
	// srvName := strings.TrimPrefix(args[0], "*main.")
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
	argGroup := rs[len(rs)-2]
	argGroup = strings.TrimPrefix(argGroup, "(")
	argGroup = strings.TrimSuffix(argGroup, ")")
	args := strings.Split(argGroup, ", ")
	retGroup := rs[len(rs)-1:][0]
	ret := strings.TrimPrefix(retGroup, "(")
	ret = strings.TrimSuffix(ret, ")")
	rets := strings.Split(ret, ", ")
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

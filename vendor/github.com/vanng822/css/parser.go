package css

import (
	"fmt"
	"github.com/gorilla/css/scanner"
	"strings"
)

/*
	stylesheet  : [ CDO | CDC | S | statement ]*;
	statement   : ruleset | at-rule;
	at-rule     : ATKEYWORD S* any* [ block | ';' S* ];
	block       : '{' S* [ any | block | ATKEYWORD S* | ';' S* ]* '}' S*;
	ruleset     : selector? '{' S* declaration? [ ';' S* declaration? ]* '}' S*;
	selector    : any+;
	declaration : property S* ':' S* value;
	property    : IDENT;
	value       : [ any | block | ATKEYWORD S* ]+;
	any         : [ IDENT | NUMBER | PERCENTAGE | DIMENSION | STRING
	              | DELIM | URI | HASH | UNICODE-RANGE | INCLUDES
	              | DASHMATCH | ':' | FUNCTION S* [any|unused]* ')'
	              | '(' S* [any|unused]* ')' | '[' S* [any|unused]* ']'
	              ] S*;
	unused      : block | ATKEYWORD S* | ';' S* | CDO S* | CDC S*;
*/

type State int

const (
	STATE_NONE State = iota
	STATE_SELECTOR
	STATE_PROPERTY
	STATE_VALUE
)

type parserContext struct {
	State            State
	NowSelectorText  string
	NowRuleType      RuleType
	CurrentRule      *CSSRule
	CurrentMediaRule *CSSRule
}

func resetContextStyleRule(context *parserContext) {
	context.CurrentRule = nil
	context.NowSelectorText = ""
	context.NowRuleType = STYLE_RULE
	context.State = STATE_NONE
}

func parseRule(context *parserContext, s *scanner.Scanner, css *CSSStyleSheet) {
	context.CurrentRule = NewRule(context.NowRuleType)
	context.NowSelectorText += parseSelector(s)
	context.CurrentRule.Style.SelectorText = strings.TrimSpace(context.NowSelectorText)
	context.CurrentRule.Style.Styles = parseBlock(s)
	if context.CurrentMediaRule != nil {
		context.CurrentMediaRule.Rules = append(context.CurrentMediaRule.Rules, context.CurrentRule)
	} else {
		css.CssRuleList = append(css.CssRuleList, context.CurrentRule)
	}
}

// Parse takes a string of valid css rules, stylesheet,
// and parses it. Be aware this function has poor error handling
// so you should have valid syntax in your css
func Parse(csstext string) *CSSStyleSheet {
	context := &parserContext{
		State:            STATE_NONE,
		NowSelectorText:  "",
		NowRuleType:      STYLE_RULE,
		CurrentMediaRule: nil,
	}

	css := &CSSStyleSheet{}
	css.CssRuleList = make([]*CSSRule, 0)
	s := scanner.New(csstext)

	for {
		token := s.Next()

		//fmt.Printf("Parse(%d): %s:'%s'\n", context.State, token.Type.String(), token.Value)

		if token.Type == scanner.TokenEOF || token.Type == scanner.TokenError {
			break
		}

		switch token.Type {
		case scanner.TokenCDO:
			break
		case scanner.TokenCDC:
			break
		case scanner.TokenComment:
			break
		case scanner.TokenS:
			break
		case scanner.TokenAtKeyword:
			switch token.Value {
			case "@media":
				context.NowRuleType = MEDIA_RULE
			case "@font-face":
				// Parse as normal rule, would be nice to parse according to syntax
				// https://developer.mozilla.org/en-US/docs/Web/CSS/@font-face
				context.NowRuleType = FONT_FACE_RULE
				parseRule(context, s, css)
				resetContextStyleRule(context)
			case "@import":
				// No validation
				// https://developer.mozilla.org/en-US/docs/Web/CSS/@import
				rule := parseImport(s)
				if rule != nil {
					css.CssRuleList = append(css.CssRuleList, rule)
				}
				resetContextStyleRule(context)
			case "@charset":
				// No validation
				// https://developer.mozilla.org/en-US/docs/Web/CSS/@charset
				rule := parseCharset(s)
				if rule != nil {
					css.CssRuleList = append(css.CssRuleList, rule)
				}
				resetContextStyleRule(context)

			case "@page":
				context.NowRuleType = PAGE_RULE
				parseRule(context, s, css)
				resetContextStyleRule(context)
			default:
				panic(fmt.Sprintf("At rule '%s' is not supported", token.Value))
			}
		default:
			if context.State == STATE_NONE {
				if token.Value == "}" && context.CurrentMediaRule != nil {
					// close media rule
					css.CssRuleList = append(css.CssRuleList, context.CurrentMediaRule)
					context.CurrentMediaRule = nil
					break
				}
			}

			if context.NowRuleType == MEDIA_RULE {
				context.CurrentMediaRule = NewRule(context.NowRuleType)
				context.CurrentMediaRule.Style.SelectorText = strings.TrimSpace(token.Value + parseSelector(s))
				resetContextStyleRule(context)
				break
			} else {
				context.NowSelectorText += token.Value
				parseRule(context, s, css)
				resetContextStyleRule(context)
				break
			}
		}
	}
	return css
}

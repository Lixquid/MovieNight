package common

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type EmotesMap map[string]string

var Emotes EmotesMap
var WrappedEmotesOnly bool = false

var (
	reStripStatic   = regexp.MustCompile(`^(\\|/)?static`)
	reWrappedEmotes = regexp.MustCompile(`[:\[][^\s:\/\\\?=#\]\[]+[:\]]`)
	reImg           = regexp.MustCompile(`^<a href="(https?://i\.imgur\.com/\w+\.png)" target="_blank">https?://i\.imgur\.com/\w+\.png</a>$`)
)

func init() {
	Emotes = NewEmotesMap()
}

func NewEmotesMap() EmotesMap {
	return map[string]string{}
}

func (em EmotesMap) Add(fullpath string) EmotesMap {
	fullpath = reStripStatic.ReplaceAllLiteralString(fullpath, "")

	base := filepath.Base(fullpath)
	code := base[0 : len(base)-len(filepath.Ext(base))]

	_, exists := em[code]

	num := 0
	for exists {
		num += 1
		_, exists = em[fmt.Sprintf("%s-%d", code, num)]
	}

	if num > 0 {
		code = fmt.Sprintf("%s-%d", code, num)
	}

	em[code] = fullpath
	return em
}

func EmoteToHtml(file, title string) string {
	return fmt.Sprintf(`<img src="%s" class="emote" title="%s" />`, file, title)
}

var valid_modifiers = map[string]string{
	"blur":   "blur",
	"grey":   "grayscale",
	"gray":   "grayscale",
	"invert": "invert",
	"flip":   "horizontalflip",
	"upside": "upsidedown",
	"squish": "squished",
	"lg":     "large",
}

func ParseEmotesArray(words []string) []string {
	newWords := []string{}
	for _, word := range words {
		emote_name := strings.Trim(word, ":")
		emote_modifier := ""

		if mod := strings.SplitN(emote_name, "~", 2); len(mod) == 2 {
			if class, ok := valid_modifiers[mod[1]]; ok {
				emote_name = mod[0]
				emote_modifier = class
			}
		}

		if filename, ok := Emotes[emote_name]; ok {
			newWords = append(newWords, fmt.Sprintf(`<img src="%s" class="emote emote_%s" title="%s" />`, filename, emote_modifier, emote_name))
		} else {
			newWords = append(newWords, word)
		}
	}

	return newWords
}

func ParseEmotes(msg string) string {
	if reImg.MatchString(msg) {
		href := reImg.FindStringSubmatch(msg)[1]
		return `<a href="` + href + `" target="_blank"><img src="` + href + `" class="chat-image" /></a>`
	} else {
		words := ParseEmotesArray(strings.Split(msg, " "))
		return strings.Join(words, " ")
	}
}

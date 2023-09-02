package main

import (
	"fmt"
	"regexp"
	"strings"
)

var maxMeaningLine int = 5

type meaningAndSentence struct {
	Meaning  string   `json:"meaning"`
	Sentence []string `json:"sentence"`
}

type wordMeaning struct {
	WordToSearch string               `json:"word"`
	KK           string               `json:"kk"`
	POS          string               `json:"pos"`
	ResultList   []meaningAndSentence `json:"result"`
}

func formatCrawlerResult(result string) string {
	space := regexp.MustCompile(`\s+`)
	removeSpace := space.ReplaceAllString(result, " ")
	// remove case of [ C ] or [ T ]
	corT := regexp.MustCompile(`\[\s+.\s+\]|:`)
	removeCorT := corT.ReplaceAllString(removeSpace, "")
	// remove leading space
	noSpaceDuplicate := strings.TrimSpace(removeCorT)
	// replace the needed escape character
	// escape := regexp.MustCompile(`\.|\'|\*|\[|\]|\(|\)|\~|\>|\#|\+|\-|\=|\||\{|\}|\.|\!`)
	// removeEscape := escape.ReplaceAllString(noSpaceDuplicate, `\$0`)

	s := noSpaceDuplicate
	return s
}

func PreprocessingJSONToString(preOutput wordMeaning) string {
	output := ""
	// title
	output += fmt.Sprintf(`*%s*  (_%s_)`, preOutput.WordToSearch, preOutput.POS) + "\n"
	output += preOutput.KK + "\n"

	for i, result := range preOutput.ResultList {
		if i+1 > maxMeaningLine {
			break
		}
		output += fmt.Sprintf("%d", i+1) + ". *" + result.Meaning + "*\n"
		if len(result.Sentence) > 0 {
			output += `\* _` + result.Sentence[0] + "_\n"
		}
	}

	return output
}

package main

import (
	"github.com/gocolly/colly"
)

var (
	c             *colly.Collector
	parsingResult wordMeaning
)

func init() {
	c = colly.NewCollector(
		colly.AllowedDomains("dictionary.cambridge.org"),
	)
	c.AllowURLRevisit = true
	settingSearch()
}

func settingSearch() {
	c.OnHTML(".pos-header.dpos-h", func(e *colly.HTMLElement) {
		// KK
		e.ForEach(".us.dpron-i .pron.dpron", func(i int, m *colly.HTMLElement) {
			parsingResult.KK = m.Text
		})
		// part of speech
		e.ForEach(".posgram.dpos-g.hdib.lmr-5", func(i int, m *colly.HTMLElement) {
			parsingResult.POS = m.Text
		})
	})
	// On every a element which has href attribute call callback
	c.OnHTML(".def-block.ddef_block", func(e *colly.HTMLElement) {
		var newMeaningAndSentence meaningAndSentence
		// meaning
		e.ForEach(".def.ddef_d.db", func(i int, m *colly.HTMLElement) {
			newMeaningAndSentence.Meaning = formatCrawlerResult(m.Text)
		})
		// sentence
		e.ForEach(".def-body.ddef_b .examp.dexamp", func(i int, m *colly.HTMLElement) {
			newMeaningAndSentence.Sentence = append(newMeaningAndSentence.Sentence, formatCrawlerResult(m.Text))
		})
		parsingResult.ResultList = append(parsingResult.ResultList, newMeaningAndSentence)
	})
}

func Search(wordToSearch string) wordMeaning {
	parsingResult = wordMeaning{}
	parsingResult.WordToSearch = wordToSearch
	err := c.Visit("https://dictionary.cambridge.org/dictionary/english/" + wordToSearch)
	if err != nil {
		panic(err)
	}
	return parsingResult
}

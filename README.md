# Cambridge Dictionary API

 Let's Build a fully functional  Cambridge Dictionary API by crawling it's data based on word search
 We will use  colly to crawl the result data and  Regex to format our result and encode and Encode the data using JSON Marshalling.

Create a folder a dictionary-api in your computer and create a `main.go` file and run `go mod init dictionary-api` and then `go mod tidy` when ever you add a new dependencies run `go mod tidy` in terminal

## Step - 1 Installing Dependencies

We will install Colly a Excellent Web Scraping framework Written in GO

```go
go get github.com/gocolly/colly
```

### Writing Our Scraping logic

1. **Initialization:**
   - The code starts by initializing two important variables: `c` (Collector) and `parsingResult`.
   - `c` is an instance of a web scraping tool called Colly, which is used to navigate and extract data from web pages.
   - `parsingResult` is a structure that will hold the results of the web scraping.

2. **Initialization Function (`init`):**
   - In the `init` function, the Colly collector `c` is set up.
   - It is configured to only visit web pages from the domain "dictionary.cambridge.org".
   - `AllowURLRevisit` is set to true, which means the collector can revisit the same URL if needed.
   - The `settingSearch` function is called, which sets up the rules for data extraction.

3. **Setting Up Rules for Data Extraction (`settingSearch` Function):**
   - This function defines the rules for extracting specific data from the Cambridge Dictionary web pages.
   - It uses Colly's `OnHTML` method to specify what to do when certain HTML elements are found on the page.
   - For example, when an HTML element with the class "pos-header" is found, it extracts the pronunciation (KK) and part of speech (POS).
   - When an HTML element with the class "def-block" is found, it extracts the meanings and example sentences.

4. **Search Function (`Search` Function):**
   - The `Search` function is used to perform a word search on the Cambridge Dictionary website.
   - It takes a `wordToSearch` parameter, which is the word you want to look up.
   - It resets `parsingResult` to empty before starting a new search.
   - It uses Colly to visit the Cambridge Dictionary page for the specified word.
   - After scraping the page, it returns the `parsingResult`, which now contains the word's pronunciation, part of speech, meanings, and example sentences.

```go

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
```

## Step - 2 Formatting our Scraping result Using Regex

This data structures and functions are used to organize and format information about words and their meanings. They help make the data more human-readable and suitable for presentation, such as in a user interface or chatbot response.

1. **Data Structures**:
   - `meaningAndSentence`: This structure represents a word's meaning and related example sentences. It contains two fields:
     - `Meaning`: A string that stores the meaning of the word.
     - `Sentence`: A list of strings that stores example sentences related to the word's meaning.

   - `wordMeaning`: This structure represents information about a word, including its pronunciation, part of speech, and a list of meanings and sentences. It contains the following fields:
     - `WordToSearch`: A string storing the word being searched.
     - `KK`: A string storing the word's pronunciation.
     - `POS`: A string storing the word's part of speech.
     - `ResultList`: A list of `meaningAndSentence` structures, representing different meanings and sentences associated with the word.

2. **`formatCrawlerResult` Function**:
   - This function takes a string (`result`) as input and processes it to remove unnecessary spaces and characters.
   - It uses regular expressions to remove extra spaces, colons, and other unwanted characters from the input string.
   - The processed string is then returned.

3. **`PreprocessingJSONToString` Function**:
   - This function takes a `wordMeaning` structure (`preOutput`) as input and prepares a formatted string for display.
   - It constructs a string by combining various pieces of information, including the word, its part of speech, pronunciation, meanings, and example sentences.
   - It limits the number of meanings displayed to a maximum of 5 (as specified by `maxMeaningLine`).
   - The formatted string is returned, which can be used for displaying word information in a more readable format.

```go
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
```

### Step - 3 Writing a HTTP Handler with EndPoints

The `SearchQueryHandler` function handles HTTP requests that include a word to search for in the URL path (e.g., "/search/word"). It extracts the word, retrieves its meaning, and sends a JSON response with the meaning information. If there are errors during this process, it may panic and terminate the program.

```go
func SearchQueryHandler() {
 http.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
  wordToSearch := strings.TrimPrefix(r.URL.Path, "/search/")
  outputJSON := getMeaning(wordToSearch)
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  encodeError := json.NewEncoder(w).Encode(outputJSON)
  if encodeError != nil {
   panic(encodeError)
  }
 })
}
```

### Listen and Serve

```go
func main() {
 // Standard JSON REQUEST
 SearchQueryHandler()
 // Start the HTTP server
 // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
 if err := http.ListenAndServe(":8080", nil); err != nil {
  log.Fatal(err)
 }
}
```

### Run

```go
// it will find all the global functions and run it. It is because it did not step up go env GOPATH
// So for you - you can try  go run main.go
go run .
// OR
go run main.go
```

## Demo

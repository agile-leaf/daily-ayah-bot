package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
)

const TOTAL_AYAHS = 6236
const MY_USERNAME = "@daily_ayah_bot"

func main() {
	currentConfig := loadConfig()

	randomAyah, err := getRandomAyah()
	if err != nil {
		fmt.Printf("Unable to get random ayah. Error: %s", err.Error())
		os.Exit(1)
	}

	ayahParts := splitAyahForTweet(randomAyah)
	tweets := getTweetsFromParts(randomAyah, ayahParts)

	api := getTwitterApiForConfig(currentConfig)
	firstTweet, err := api.PostTweet(tweets[0], nil)
	if err != nil {
		fmt.Printf("Unable to post tweet. Error: %s", err.Error())
		os.Exit(1)
	}

	tweetPartOptions := url.Values{}
	tweetPartOptions.Set("in_reply_to_status_id", fmt.Sprint(firstTweet.Id))

	if len(tweets) > 1 {
		for _, tweet := range tweets[1:] {
			_, err = api.PostTweet(tweet, tweetPartOptions)
			if err != nil {
				fmt.Printf("Unable to post tweet. Tweet: %s\nError: %s", tweet, err.Error())
				os.Exit(1)
			}
		}
	}

	fmt.Println("Success")
}

func formatTweet(ayahText string, footer string) string {
	return fmt.Sprintf("%s\n\n%s", ayahText, footer)
}

func formatTweetWithCount(ayahText string, footer string, partNumber int, totalParts int) string {
	return fmt.Sprintf("[%d/%d] %s", partNumber, totalParts, formatTweet(ayahText, footer))
}

func getTweetsFromParts(a *ayah, parts []string) []string {
	tweets := []string{}
	footer := a.getFooter()

	if len(parts) == 1 {
		tweets = append(tweets, formatTweet(parts[0], footer))
	} else {
		// The first tweet is special because it doesn't include the username
		tweets = append(tweets, formatTweetWithCount(parts[0], footer, 1, len(parts)))

		for i, part := range parts[1:] {
			currentTweet := formatTweetWithCount(part, footer, i+2, len(parts))
			tweets = append(tweets, currentTweet)
		}
	}

	return tweets
}

func getTwitterApiForConfig(cfg *config) *anaconda.TwitterApi {
	anaconda.SetConsumerKey(cfg.ConsumerKey)
	anaconda.SetConsumerSecret(cfg.ConsumerSecret)
	return anaconda.NewTwitterApi(cfg.AccessToken, cfg.AccessSecret)
}

func splitAyahForTweet(ayahToSplit *ayah) []string {
	words := strings.Fields(ayahToSplit.ayahText)

	footer := ayahToSplit.getFooter()
	// <AYAH>\n\n<FOOTER>
	spaceLeftForAyah := 140 - len(footer) - 2

	if spaceLeftForAyah > len(ayahToSplit.ayahText) {
		return []string{ayahToSplit.ayahText}
	}

	numberOfParts := int(len(ayahToSplit.ayahText)/spaceLeftForAyah) + 1
	tweetCountLen := 1
	if numberOfParts > 10 {
		tweetCountLen = 2
	}

	// The count of chars needed to display meta data for the tweet part. This includes:
	// <OPENING BRACE><PART NUMBER><SLASH><TOTAL PARTS><CLOSING BRACE><SPACE>
	spaceNeededForTweetCount := 1 + tweetCountLen + 1 + tweetCountLen + 1 + 1
	spaceLeftForAyah -= spaceNeededForTweetCount

	parts := []string{}
	currentPart := words[0]
	for _, word := range words[1:] {
		if len(currentPart)+len(word)+1 > spaceLeftForAyah {
			parts = append(parts, currentPart)
			currentPart = word
		} else {
			currentPart += " " + word
		}
	}

	if len(currentPart) > 0 {
		parts = append(parts, currentPart)
	}
	return parts
}

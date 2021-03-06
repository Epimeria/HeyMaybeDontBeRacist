// Version 1.04
package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

const copyPasta = `Haha, those sure are some kooky crime statistics you got there, dad! Where are those from, the FBI? Well, did you know that FBI crime statistics only track those that are successfully convicted of a crime, even though the vast majority of crimes committed aren't even reported, let alone lead to a successful arrest, prosecution, and conviction? Sounds like a big hole in the data if you ask me! And when you combine that with the fact that black neighborhoods are more heavily policed black people are more likely to get stopped by the police, and more likely to be convicted by a jury for the same crimes as a white person, it makes you think that maybe those FBI stats have a lot more to do with systemic racial profiling than some sort of innate behavior in black people! How's that for a brain blast?

^^This ^^action ^^was ^^performed ^^automatically. ^^To ^^report ^^a ^^false ^^positive, ^^please ^^open ^^an ^^issue ^^on ^^my ^^[github](https://github.com/Epimeria/HeyMaybeDontBeRacist)
`

type announcer struct {
	bot reddit.Bot
}

// CommentHandler ...
type CommentHandler interface {
	Comment(comment *reddit.Comment) error
}

func (a *announcer) Comment(comment *reddit.Comment) error {
	matched, _ := regexp.MatchString(`((13\%?\:?\;?\/? ?5(0|2)\%?)( |$|\.|\,|\: |\;|\-))|(Despite$)|(Despite (only )?making up (only )?13\%?)`, comment.Body)
	negativeMatch, _ := regexp.MatchString(`((youtu\.be)|(youtube\.com))|(dogwhistle|dog whistle)|(\$1350|1350\$)`, comment.Body)
	if matched && !negativeMatch {
		fmt.Printf("%s was a little racist on %s: '%s'\n", comment.Author, comment.Subreddit, comment.Body)
		return a.bot.Reply(comment.Name, copyPasta)
	}

	return nil
}

func main() {
	subsToMonitor := []string{
		"ActualPublicFreakouts",
		"AskThe_Donald",
		"Conservative",
		"Conspiracy",
		"FightPorn",
		"FuckThesePeople",
		"Gaming",
		"Holup",
		"Iamatotalpieceofshit",
		"InstantKarma",
		"JusticeServed",
		"KotakuInAction",
		"MakeMeSuffer",
		"News",
		"PewdiePieSubmissions",
		"Politic",
		"PoliticalCompassMemes",
		"PoliticalHumor",
		"PublicFreakout",
		"PussyPassDenied",
		"SargonOfAkkad",
		"The_Donald",
		"Trashy",
		"TumblrInAction",
		"UnpopularOpinion",
		"WatchRedditDie",
		"WhitePeopleTwitter",
		"WinStupidPrizes",
		"WorldNews",
	}

	bot, err := reddit.NewBotFromAgentFile("bot.agent", 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	cfg := graw.Config{
		SubredditComments: subsToMonitor,
	}

	fmt.Println("Bot started")

	if _, wait, err := graw.Run(&announcer{bot: bot}, bot, cfg); err != nil {
		fmt.Println("Failed to start graw run: ", err)
	} else {
		fmt.Println("graw run failed: ", wait())
	}
}

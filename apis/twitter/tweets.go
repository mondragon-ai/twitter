package twitter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/twitter/apis/openai"
	"github.com/twitter/auth"
)


var (
	consumerKey = os.Getenv("CONSUMER_KEY")
	consumerSecret = os.Getenv("CONSUMER_SECRET")
	accessTokenKey = os.Getenv("ACCESS_TOKEN_KEY")
	accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
)


// List of topics
var topics = []string{
	"The stock market as a game of musical chairs.",
	"Cryptocurrency: the new gold rush, or just a digital illusion?",
	"Tech startups: where every idea is a billion-dollar unicorn until proven otherwise.",
	"Investing in NFTs: turning JPEGs into digital gold.",
	"The AI takeover: when algorithms become our overlords.",
	"Venture capitalists: the modern-day kings of Silicon Valley.",
	"Blockchain: decentralization or just another buzzword?",
	"Tech IPOs: where hype meets reality.",
	"The gig economy: where everyone's a freelancer until they need healthcare.",
	"Tech bro culture: where hoodies are the new power suits.",
	"Disrupting industries: because who doesn't love a good revolution?",
	"Elon Musk: our real-life Tony Stark or just a very ambitious meme lord?",
	"The sharing economy: where we rent everything except our opinions.",
	"FOMO investing: because missing out is the new FOMO.",
	"Tech billionaires: where wealth accumulates faster than you can say 'innovation.'",
	"The dark side of data: privacy concerns in the age of surveillance capitalism.",
	"Stock market volatility: because rollercoasters are so 20th century.",
	"Silicon Valley: where dreams are made, and startups go to die.",
	"The gig economy: where hustle culture meets burnout.",
	"Tech addiction: because dopamine hits are the new currency.",
	"Robo-advisors: because who needs human expertise when you have algorithms?",
	"Tech conferences: where networking is just a fancy word for 'awkward small talk.'",
	"The rise of remote work: because commuting is so last century.",
	"Venture capitalists: where throwing money at problems is the norm.",
	"The rise of AI influencers: because even algorithms need a following.",
	"Tech IPOs: where valuations defy gravity.",
	"The future of work: where robots are colleagues and bosses alike.",
	"Sustainable investing: because saving the planet is the new ROI.",
	"Tech unicorns: because nothing says success like a billion-dollar valuation.",
	"The gig economy: where side hustles are the new 9 to 5.",
	"Quantum computing: because sometimes binary just isn't enough.",
	"The sharing economy: where we rent out our spare rooms and our dignity.",
	"Tech bro culture: because disruptors gonna disrupt (and mansplain).",
	"The rise of fintech: because who needs banks when you have apps?",
	"The gig economy: where everyone's a CEO of their own hustle.",
	"Artificial intelligence: because who needs common sense when you have algorithms?",
	"Tech IPOs: where losses are just pre-profit.",
	"The sharing economy: where we monetize our excess stuff and our free time.",
	"Tech billionaires: because money can't buy humility (or taste).",
	"The future of transportation: where we trade gas guzzlers for electric dreams.",
	"The gig economy: where every side hustle is a potential IPO.",
	"Venture capitalists: because betting on the future is the ultimate gamble.",
	"The sharing economy: where trust is currency, and ratings are king.",
	"Tech bro culture: because disrupting industries is just another form of mansplaining.",
	"The rise of remote work: where pajamas are the new power suits.",
	"Artificial intelligence: because who needs human intelligence when you have algorithms?",
	"Tech IPOs: where speculation meets speculation.",
	"The sharing economy: where we monetize our trust issues.",
	"Tech billionaires: where philanthropy is just a tax write-off.",
	"The gig economy: where hustle culture meets gig fatigue.",
	"The future of AI: where ethics are just a software update away.",
	"Tech IPOs: where unicorns become donkeys overnight.",
	"The sharing economy: where we monetize our social lives.",
	"Tech bro culture: where brogrammers rule the world.",
	"The rise of remote work: where Zoom calls are the new water cooler chats.",
	"Artificial intelligence: because sometimes the future is smarter than we are.",
	"Tech IPOs: where dreams meet reality (and sometimes crash).",
	"The sharing economy: where we rent out everything but our dignity.",
	"Tech billionaires: where ego meets innovation.",
	"The gig economy: where everyone's a freelancer until the rent is due.",
	"The future of transportation: where flying cars are just a software update away.",
	"The rise of AI: where Siri is the new oracle.",
	"Tech IPOs: where optimism meets delusion.",
	"The sharing economy: where we trade privacy for convenience.",
	"Tech bro culture: where disruption is just another word for 'boys club.'",
	"The future of work: where robots are colleagues, and humans are resources.",
	"Artificial intelligence: because sometimes the future is scarier than fiction.",
	"Tech IPOs: where startups go to make it big (or go bust).",
	"The sharing economy: where strangers become customers, and customers become strangers.",
	"Tech billionaires: where wealth meets hubris.",
	"The gig economy: where everyone's a CEO until taxes are due.",
	"The future of AI: where algorithms are the new gods.",
	"Tech IPOs: where dreams are priced by the share.",
	"The sharing economy: where we monetize our possessions and our privacy.",
	"Tech bro culture: where hoodies are the new power suits.",
	"The rise of remote work: where Slack is the new office.",
	"Artificial intelligence: because sometimes the future is smarter than we are.",
	"Tech IPOs: where unicorns become unicorpses.",
	"The sharing economy: where we rent out our lives one Airbnb at a time.",
	"Tech billionaires: where wealth meets influence.",
	"The gig economy: where everyone's a freelancer until the rent is due.",
	"The future of transportation: where self-driving cars are just a software update away.",
	"The rise of AI: where Alexa is the new oracle.",
	"Tech IPOs: where optimism meets delusion.",
	"The sharing economy: where we trade privacy for convenience.",
	"Tech bro culture: where disruption is just another word for 'boys club.'",
	"The future of work: where robots are colleagues, and humans are resources.",
	"Artificial intelligence: because sometimes the future is scarier than fiction.",
	"Tech IPOs: where startups go to make it big (or go bust).",
	"The sharing economy: where strangers become customers, and customers become strangers.",
	"Tech billionaires: where wealth meets hubris.",
	"The gig economy: where everyone's a CEO until taxes are due.",
	"The future of AI: where algorithms are the new gods.",
	"Tech IPOs: where dreams are priced by the share.",
	"The sharing economy: where we monetize our possessions and our privacy.",
	"Tech bro culture: where hoodies are the new power suits.",
	"The rise of remote work: where Slack is the new office.",
	"Artificial intelligence: because sometimes the future is smarter than we are.",
	"Tech IPOs: where unicorns become unicorpses.",
	"The sharing economy: where we rent out our lives one Airbnb at a time.",
}
// getRandomTopic selects a random topic from the list
func getRandomTopic() string {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(topics))
	return topics[randomIndex]
}

// insertRandomTopicIntoPrompt inserts a randomly selected topic into the prompt
func insertRandomTopicIntoPrompt() string {
	randomTopic := getRandomTopic()
	return fmt.Sprintf("##Twitter Bot Prompt\n\n### Description: This Twitter bot generates tech and investment-oriented tweets with a focus on wit and satire. The tweets should be insightful yet predominantly parody or witty. \n\n### Bio:Tech enthusiast whos just as likely to be lost in a foreign city as debugging lines of code. Foodie pro navigating the world one misadventure at a time 🧳🏋🏼🍲\n\n### Prompt: Generate a tweet related to technology or investment, infused with wit or satire, within 160 characters. \n\nTopic: %s \n\n #INSTRUCTIONS: \n- **DO NOT** generate a tweet with hashtags. \n- **DO NOT** use quotations around the generated tweet. \n\n#EXAMPLES - Here are examples of popular tweets you can use to get inspiration. do not include any of the users:  Tech and marketing co founders are the deadliest combination for any startup. - @iuditg \n If your family doesnt think youre crazy, are you even a startup founder? - @dagorenouf \nIm a Pull Stack Developer.I just pull things off the Internet and put it into my code. - @TheJackForge \nDebugging is like an onion. There are multiple layers to it, and the more you peel them back, the more likely youre going to start crying at inappropriate times. - @iamdevloper\nif youre not happy single, you wont be happy in a relationship.true happiness comes from closing 100 chrome tabs after solving an obscure programming bug, not from someone else -@cszhu \n No meetings and uninterrupted coding is the best thing that can happen to a developer. -@ vittoStack \n Programming is 60 percent thinking, 5 percent coding, and 45 percent yelling the F*** word at your laptop.  -@ vittoStack\n Every startup has acquisition as a backup plan, except for the ones that win. -@naval \n A car with a driver is the new horse and buggy. -@naval \nThe narrative for this cycle is “this is the last cycle.” -@naval.", randomTopic)
}

func Tweet(writer http.ResponseWriter, request *http.Request) {

	// Retrieve environment variables
	consumerKey = os.Getenv("CONSUMER_KEY")
	consumerSecret = os.Getenv("CONSUMER_SECRET")
	accessTokenKey = os.Getenv("ACCESS_TOKEN_KEY")
	accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")

	// Generate signature
    signature := auth.PrepareOAuthSignature(consumerKey, accessTokenKey, consumerSecret, accessTokenSecret)

	prompt := insertRandomTopicIntoPrompt()
	completion, err := openai.OpenAIChatCompletion(prompt)
	if err != nil {
		// log.Fatal("Error creating request: ", err)
		http.Error(writer, "Error generating completion from OpenAI: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare tweet data
	tweetData := map[string]string{
		"text": strings.Replace(completion, "\"", "", -1),
	}
	tweetJSON, _ := json.Marshal(tweetData)

	// Send tweet request
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.twitter.com/2/tweets", bytes.NewBuffer(tweetJSON))
	if err != nil {
		log.Fatal("Error creating request: ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", signature)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request: ", err)
	}
	defer resp.Body.Close()

	// Print response
	fmt.Println("Response Status:", resp.Status)
}



func Post() {

	// Retrieve environment variables
	consumerKey = os.Getenv("CONSUMER_KEY")
	consumerSecret = os.Getenv("CONSUMER_SECRET")
	accessTokenKey = os.Getenv("ACCESS_TOKEN_KEY")
	accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")

	// Generate signature
    signature := auth.PrepareOAuthSignature(consumerKey, accessTokenKey, consumerSecret, accessTokenSecret)

	prompt := insertRandomTopicIntoPrompt()
	completion, err := openai.OpenAIChatCompletion(prompt)
	if err != nil {
		log.Fatal("Error creating request: ", err)
		return
	}

	// Prepare tweet data
	tweetData := map[string]string{
		"text": strings.Replace(completion, "\"", "", -1),
	}
	tweetJSON, _ := json.Marshal(tweetData)

	// Send tweet request
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.twitter.com/2/tweets", bytes.NewBuffer(tweetJSON))
	if err != nil {
		log.Fatal("Error creating request: ", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", signature)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request: ", err)
	}
	defer resp.Body.Close()

	// Print response
	fmt.Println("Response Status:", resp.Status)
}

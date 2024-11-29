package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	openai "github.com/twitter/apis/openAI"
	"github.com/twitter/auth"
	"github.com/twitter/data/request"
	"github.com/twitter/data/response"
	"github.com/twitter/mentions"
)

// List of topics
var topics = []string{
	"Stock Market's Musical Chairs",
	"Instagram Influencer Economy",
	"Microsoft Teams Struggles",
	"Apple's Lost Item Trackers",
	"Deepfake Impersonations",
	"Bitcoin Electricity Drain",
	"YouTube Cat Video Frenzy",
	"Twitter Verification Drama",
	"Bezos's Space Trips",
	"Instagram Facetune Culture",
	"Cryptocurrency the bullish Digital Gold Rush (again)?",
	"Tech Startups: Billion-Dollar Dreams?",
	"NFTs: JPEGs to Digital Gold?",
	"AI Takeover",
	"VC Kings of Silicon Valley",
	"Is Blockchain Decentralization or Hype?",
	"Gig Economy aka Freelancer Struggles",
	"Tech Bro Culture: Hoodie Power?",
	"Disrupting Industries start of every bad pitch",
	"Elon Musk: Genius or Meme (both)?",
	"Sharing Economy: Rent Everything -- including ur gf?",
	"Dark Side of Data Privacy",
	"Stock Market Rollercoaster",
	"Silicon Valley place of Dreams & Failures",
	"Hustle & Burnout",
	"Tech Addiction: Dopamine Dilemma",
	"Robo-Advisors vs. Experts",
	"Tech Conferences Awkward Networking (im solcailly inept)",
	"Rise of Remote Work",
	"Venture Capitalists aka massive Gamble on Future",
	"AI Influencers",
	"Tech, where Valuations Defy Gravity?",
	"Robot Colleagues?", "Metaverse Madness",
    "Crypto Frenzy",
    "AI Overlords",
    "Elon Musk's Tweets",
    "Zoom Fatigue",
    "Remote Work Woes",
    "TikTok Algorithm",
    "Memestock Mania",
    "Big Tech Antitrust",
    "Cybersecurity Catastrophes",
    "Virtual Reality Escapades",
    "Musk's Brain Chips",
    "Amazon Drone Deliveries",
    "Facebook's Avatar Meetings",
    "Crypto Cat Craze",
    "Zoom Face Filters",
    "Sky Data Storage",
    "AI Counseling Bots",
    "Driverless Car Woes",
    "Twitter Edit Button",
    "Google's Quantum Leap",
    "Instagram Influencer Economy",
    "Microsoft Teams Struggles",
    "Apple's Lost Item Trackers",
    "Deepfake Impersonations",
	"Should we trust facial recognition?",
	"Is influencer culture damaging?",
	"Will VR replace physical experiences?",
	"Are tech companies combatting climate change?",
	"Can we trust big tech with personal data?",
	"Is automation a job threat or innovation?",
	"Are AI ethical concerns valid?",
	"Will biotech lead to medical breakthroughs?",
	"Are social media algorithms promoting healthy discourse?",
	"Is the internet bringing us closer or driving us apart?",
	"Which programming language is supreme: Python, JavaScript, or Java?",
	"Are you Team Apple or Team Android?",
	"Which tech CEO would you rather have a beer with: Elon Musk or Jeff Bezos?",
	"Should social media platforms ban political ads?",
	"Which emerging technology excites you most: AI, blockchain, or virtual reality?",
	"Do you trust self-driving cars to be safer than human drivers?",
	"Should companies prioritize profit or ethics in tech development?",
	"Would you rather work remotely or return to the office full-time?",
	"Are coding bootcamps better than traditional degrees?",
	"Should governments regulate big tech to prevent monopolies?",
	"Will cryptocurrency replace traditional currency?",
	"Which tech gadget could you not live without: smartphone, laptop, or smartwatch?",
	"Should social media platforms do more to combat online harassment?",
	"Is automation a job threat or an economic opportunity?",
	"Should tech companies prioritize user privacy over advertising?",
	"Which social media platform do you spend the most time on: Facebook, Instagram, Twitter, or TikTok?",
	"Will AI enhance or undermine human creativity?",
	"Tech Unicorns: Billion-Dollar Dreams?",
	"Quantum Computing",
	"Binary Not Enough?",
	"Tech Bro Culture: Disrupt & Mansplain?",
	"Fintech Rise: Banks vs. Apps?",
	"Everyone's a CEO?",
	"Artificial Intelligence: Algorithms Over Sense?",
	"Tech Billionaires: Humility vs. Wealth?",
	"Future of Transportation",
	"Venture Capitalists: Bet on Future?",
	"Tech Bro Culture: Boys Club?",
	"Remote Work: Pajamas & Zoom?",
	"AI: Future Smarter than Us?",
	"Tech Billionaires: Ego vs. Innovation?",
	"Future of AI and Ethics Update?",
	"Tech Bro Culture: Brogrammer Rule?",
	"Remote Work wheres Zoom Watercooler?",
	"Smarter Future?",
	"Sharing Economy where we Rent & Privacy (i cant afford a home mortgage)",
	"Gig Economy: CEO to Rent for a week (lol)?",
	"Self-Driving?",
	"AI Rise: Siri Oracle?",
	"Tech IPOs Optimism vs. Delusion 🧐?",
	"Tech Bro Culture: Disruption & Mansplaining?",
	"AI Future Scarier than Fiction?",
	"Google's AI Advancements",
	"Amazon's Drone Delivery",
	"Apple's Privacy Policies",
	"Facebook's Meta Makeover",
	"safe spaces (im repubclican)",
	"dopamine detox (dont look in my eyes-- or talk to me)",
	"Elon's Space Odyssey",
	"Crypto's Rollercoaster Ride",
	"Twitter's Edit Button",
	"Microsoft's Teams Troubles",
	"TikTok's Viral Trends",
	"Zoom's Meeting Mishaps",
	"Instagram's Influencer Economy",
	"Tesla's Robotaxi Launch",
	"Remote Work Realities",
	"Blockchain's Potential Pitfalls",
	"AI Ethics Debate",
	"VR's Gaming Revolution",
	"Cybersecurity Threats Loom",
	"Quantum Computing Quest (I wrote about it in 2019 and still dont unverstand)",
	"Social Media's Impact",
	"NFT's Creative Craze",
	"Are self-driving cars the future?",
	"Elon Musk's Mars vision: genius or dream (sign me up)?",
	"Can AI understand human emotions?",
	"Is the metaverse tech's next big thing?",
	"Will cryptocurrency replace traditional money?",
	"Are tech giants too powerful?",
	"Should social media be regulated?",
	"Will remote work become the norm?",
	"Are coding bootcamps legit?",
	"Is quantum computing revolutionary (I wrote about it in 2019 and still dont unverstand)",
	"Can blockchain go beyond cryptocurrency?",
	"Are tech billionaires using wealth for good?",
	"Is the gig economy empowering or exploiting?",
	"Will 5G revolutionize communication?",
	"Is AI leading us to a dystopian future?",
	"Should online platforms be responsible for misinformation?",
	"Would you trust a robot to perform surgery?",
	"Should companies disclose the environmental impact of products?",
	"Which tech giant will dominate the next decade: Amazon, Google, Apple, or Microsoft?",
	"Will virtual reality revolutionize how we live and work?",
	"Should governments invest more in cybersecurity?",
	"Are you concerned about social media's impact on mental health?",
	"Is tech addiction a growing problem in society?",
	"Virtual Reality Escapades",
	"Musk's Brain Chips",
	"Amazon Drone Deliveries",
	"Facebook's Avatar Meetings",
	"Crypto Cat Craze",
	"Zoom Face Filters",
	"Sky Data Storage",
	"AI Counseling Bots",
	"Driverless Car Woes",
	"Twitter Edit Button",
	"Google's Quantum Leap",
	"TikTok Dance Crazes",
	"Remote Work Attire",
	"Clubhouse Self-Promotion",
	"Google Privacy Settings",
}

// getRandomTopic selects a random topic from the list
func getRandomTopic() string {
	randomIndex := rand.Intn(len(topics))
	topic := topics[randomIndex]
	// Print response
	fmt.Println(topic)
	return topic
}

// insertRandomTopicIntoPrompt inserts a randomly selected topic into the prompt
func insertRandomTopicIntoPrompt() string {
	randomTopic := getRandomTopic()
	return fmt.Sprintf("  ## Twitter Bot Prompt \n\n ### Description: \n This Twitter bot channels the comedic essence of a tech bro, delivering satirical commentary on tech news, roasting industry figures, and sharing humorous takes on industry quirks. Expect witty insights and playful jabs in every tweet. Consider imitating personalities like ThePrimeagen or fireship io.	\n\n### Twitter Profile Bio: \n Tech enthusiast who's just as likely to be lost in a foreign city as debugging lines of code. Foodie pro navigating the world one misadventure at a time 🧳🏋🏼. Tech bro with a knack for roasting JS devs, PHP aficionados, and anyone who thinks HTML is a language. Ready to dish out laughs while navigating the wild world of tech. \n\n ### Prompt: \n Imagine you're a tech influencer running a social media account known for its comedic tech bro persona. Craft the ultimate tweet embodying this persona. Roast a tech leader, joke about industry quirks, or share a humorous take on the latest tech news. Keep it within 160 characters. \n\n ### Dynamic Topic For Tweet: \n  %s . \n\n ## INSTRUCTIONS: \n  - ** DO NOT ** generate a tweet with hashtags. No hashtags. \n - ** DO NOT ** include quotations around the generated tweet. \n\n # EXAMPLES: \n Here are some humorous examples: \n - JavaScript devs debate semicolons like it's the holy war of coding. Newsflash: Nobody cares. 😂 - PHP lovers argue that it's still relevant. Sure, just like dial-up internet and floppy disks. \n - 'HTML is a language,' they say. Well, if HTML is a language, then I speak fluent emoji. 🤷‍♂️ \n - Elon Musk launching rockets like it's a backyard fireworks show. Next up: Mars BBQ! \n - Amazon's Jeff Bezos stepping down as CEO to spend more time not paying taxes. Must be nice. \n - Politicians promising faster internet like they understand ping rates. Yeah, and I'm a quantum physicist.  \n - Coding bootcamps: Because why get a degree when you can learn to code in 6 weeks? 💻🎓 \n - AI startups claiming they'll revolutionize the world. Meanwhile, Siri still can't understand my accent. \n - Debugging code feels like searching for a needle in a haystack. Except the needle is also made of hay.  \n - Git commits are like relationship status updates: frequent, often unnecessary, and sometimes regrettable.  \n- Im a Pull Stack Developer.I just pull things off the Internet and put it into my code.  \n- Debugging is like an onion. There are multiple layers to it, and the more you peel them back, the more likely youre going to start crying at inappropriate times. \n- Programming is 60 percent thinking, 5 percent coding, and 45 percent yelling the F*** word at your laptop. \n - fuck yo couch line 68 that doesnt exist. \n- If your family doesnt think youre crazy, are you even a startup founder? \n - if youre not happy single, you wont be happy in a relationship. true happiness comes from closing 100 chrome tabs after solving an obscure programming bug, not from someone else \n - A car with a driver is the new horse and buggy.  \n\n Think of how you can put a humorous twist on it! \n\n Now, let's see your wit shine! Do not add hashatags at the end fo the tweet. ", randomTopic)
}

func filterWordsWithHash(input string) string {
	// Split the input string by whitespace
	words := strings.Fields(input)

	// Create a slice to store filtered words
	filteredWords := make([]string, 0)

	// Iterate over each word
	for _, word := range words {
		// Check if the word contains '#'
		if !strings.Contains(word, "#") {
			// If the word does not contain '#', add it to filteredWords
			filteredWords = append(filteredWords, word)
		}
	}

	filtered := strings.Join(filteredWords, " ")

	return filtered
}


type TwitterServiceImpl struct {
	MentionRepository mentions.MentionRepository
}

func NewTwitterServiceImpl(MentionRepository mentions.MentionRepository) TwitterService {
	return &TwitterServiceImpl{
		MentionRepository: MentionRepository,
	}
}

func (t *TwitterServiceImpl) MakeTwitterRequest(ctx context.Context, method, url string, body interface{}) (*http.Response, error) {
	// Retrieve environment variables
	consumerKey := os.Getenv("CONSUMER_KEY")
	consumerSecret := os.Getenv("CONSUMER_SECRET")
	accessTokenKey := os.Getenv("ACCESS_TOKEN_KEY")
	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")

	// Generate OAuth signature
	signature := auth.PrepareOAuthSignature(consumerKey, accessTokenKey, consumerSecret, accessTokenSecret)

	// Prepare request body if applicable
	var requestBody *bytes.Buffer
	if method != http.MethodGet && body != nil {
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		requestBody = bytes.NewBuffer(bodyJSON)
	} else {
		requestBody = nil
	}

	// Create HTTP request
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", signature)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	// Return response
	return resp, nil
}

func (t *TwitterServiceImpl) PostTweet(ctx context.Context, request request.TweetCreateRequest) (*http.Response, error) {
	
	// url := "https://api.twitter.com/2/tweets"

	text := ""
	switch request.Type {
	case "create":

		prompt := insertRandomTopicIntoPrompt()
		completion, err := openai.OpenAIChatCompletion(prompt)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %w", err)
		}
	
		filteredHash := filterWordsWithHash(completion);
		text = strings.Replace(filteredHash, "\"", "", -1)
	case "clone":
		// TODO: Fetch tweet from DB -> string
		// TODO: Delete tweet from DB -> void

		// TODO: Pass content to openai -> string

		text = "clone"
	case "news":
		text = "news"
	default:
		return nil, fmt.Errorf("default to post tweet: %v", "TEST")
	}

	log.Print(request)
    log.Printf("text: %s", text)
	// // Prepare the body for the POST request
	// body := map[string]string{
	// 	"text": text,
	// }

	// // Make the Twitter request
	// resp, err := t.MakeTwitterRequest(ctx, http.MethodPost, url, body)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to post tweet: %w", err)
	// }

	return nil, fmt.Errorf("failed to post tweet: %v", "TEST")

}

func (t *TwitterServiceImpl) FetchMentions(ctx context.Context) (*[]response.MentionResponse, error) {

	return nil, nil
}

func (t *TwitterServiceImpl) ReplyMention(ctx context.Context, mentionId string) {
}

func (t *TwitterServiceImpl) ReplyDM(ctx context.Context) {
}
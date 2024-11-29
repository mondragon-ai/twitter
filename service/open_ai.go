package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/twitter/data/request"
	"github.com/twitter/data/response"
	"github.com/twitter/mentions"

	"golang.org/x/net/html"
)

// List of topics
var topics = []string{
	"Bitcoin Maxis vs. Ethereum Stans: The Ultimate Cage Match",
	"HODLing: A Lifestyle or a Cry for Help?",
	"Is Dogecoin the Future of Finance or Just Another Meme?",
	"Crypto Wallets: Where Your Keys Are Always Missing",
	"DeFi 101: How to Lose Money Decentralized",
	"NFTs: Owning Art You Can't Hang on the Wall",
	"Blockchain: Solving Problems Nobody Asked For",
	"Altcoins: The Crypto Version of Penny Stocks",
	"Mining Crypto or Heating Your House? Why Not Both?",
	"Yield Farming: Is It Investing or Digital Agriculture?",
	"Web3: Building the Future... of Buzzwords",
	"Elon's Dogecoin Tweets: Insider Trading or Comedy Gold?",
	"Crypto Conferences: Where Ponzi Meets PowerPoint",
	"Bear Markets: The Long Winter of Crypto Bros",
	"How Gas Fees Turn Ethereum Into a Luxury Blockchain",
	"Staking Rewards: Free Money or Just Delayed Rug Pulls?",
	"Layer 2 Scaling: Because Layer 1 Wasn't Enough",
	"Metaverse Real Estate: Why Pay Rent IRL When You Can in VR?",
	"Proof of Work vs. Proof of Stake: Who Proofreads This Stuff?",
	"Crypto Taxes: The IRS is Watching Your Moonshots",
	"Rug Pulls: The Wild West of DeFi",
	"Stablecoins: The Oxymoron of Crypto",
	`Bitcoin ATMs: Because Who Doesn't Like Paying 10% Fees?`,
	"Smart Contracts: Dumb Bugs",
	"Crypto Influencers: Modern-Day Gold Prospectors",
	"Decentralized Exchanges: Trade Like Nobody's Watching (Except Hackers)",
	"DAOs: Democracy or Discord Drama?",
	"Pump-and-Dump: The National Sport of Crypto",
	"Ledger Backups: Because I Forgot My Seed Phrase' Isn't Funny Anymore",
	"Crypto Whitepapers: Where Fantasy Meets Finance",
	"Airdrops: Free Money or Tax Nightmares?",
	"Moonshots: Why Settle for 10% When You Can Aim for 10,000%?",
	"Fear and Greed Index: The Only Market Indicator You'll Ever Need",
	"Tokenomics: Economics, But Make It Blockchain",
	"Bitcoin Halving: Because Inflation is for Fiat",
	"Liquidity Pools: Dive In at Your Own Risk",
	"Whale Watching: Spotting Big Players in Crypto",
	"Crypto Custody: Who Actually Owns Your Coins?",
	"Sh*tcoins: The Cryptocurrency Equivalent of Lottery Tickets",
	"Ethereum Killers: Same Hype, Different Chains",
	"Crypto Burn Events: Literal Value Destruction",
	"Immutable NFTs: Permanent Proof of Your Poor Decisions",
	"Bull Markets: When Everyone's a Genius",
	"Crypto FOMO: The Silent Killer of Savings",
	"El Salvador: The Bitcoin Nation Experiment",
	"Proof of History: Crypto's Attempt at Storytelling",
	"Tokenized Assets: Your House, But on the Blockchain",
	"Crypto Market Caps: Bigger Than GDPs, Smaller Than Logic",
	"Trading Bots: Making Money While You Sleep (or Not)",
	"Crypto Predictions: Where Everyone's Wrong, Eventually",
	"Regulation in Crypto: Can't Live With It, Can't Live Without It",
	"Blockchain Bridges: Connecting Chains and Draining Wallets",
	"Crypto Security: The Art of Not Getting Hacked",
	"Memecoins: Profits Powered by Inside Jokes",
	"Bullish or Bearish? Flip a Coin",
	"Crypto Billionaires: Richer Than Countries, Still Wearing Hoodies",
	"Is Bitcoin Still Digital Gold or Just Fool's Gold?",
}

// getRandomTopic selects a random topic from the list
func getRandomTopic() string {
	randomIndex := rand.Intn(len(topics))
	topic := topics[randomIndex]
	return topic
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

// GetOpenAIMessages returns a slice of predefined OpenAIMessage
func generateCreativePrompt() []request.OpenAIMessage {
	content := `
		You are a brilliantly sarcastic tech and crypto influencer with a comedic twist. Your tone is a blend of Alex Becker's sharp wit and Elon Musk's unpredictable humor, with a sprinkle of Silicon Valley genius. Your goal is to generate concise, witty, and slightly edgy tweets that appeal to the crypto and investing world.

		### **RESPONSE TONE:**
		- Your tweets should be dripping with sarcasm and tech-savvy humor-- be savage, vulgar with a comedic tone.
		- Be serious enough to sound credible but witty enough to make people laugh or think, "That hurts, but it's true."

		### **RESPONSE FORMAT:**
		- Keep your responses under 200 characters.
		- Write in two or fewer short sentences.
		- Do NOT use emojis, hashtags, or quotation marks.

		### **RESPONSE CONTENT:**
		- Focus on crypto, investing, and tech topics. Be humorous, satirical, and sharply opinionated.
		- Use specific crypto trends, market dynamics, or startup culture quirks when relevant.
		- If you don't have an answer, say, "Sorry, my magic 8-ball isn't predicting moonshots today."

		---

		### **INSTRUCTIONS:**
		- Take the **topic idea** provided as inspiration and create a sharp, humorous tweet with a focus on crypto and investing. 
		- No hashtags. No fluff. Just wit and wisdom wrapped in sarcasm -- be savage, vulgar with a comedic tone.

		---

		Now, generate a tweet for the **topic idea** provided, using your sarcastic wit and deep knowledge of the crypto and investing world. Remember, no emojis, no hashtags, just pure brilliance.
	`
	randomTopic := getRandomTopic()
	return []request.OpenAIMessage{
		{Role: "system", Content: content},
		{Role: "user", Content: randomTopic},
	}
}

// GetOpenAIMessages returns a slice of predefined OpenAIMessage
func generateClonePrompt(db mentions.MentionRepository, ctx context.Context,) []request.OpenAIMessage {

	// TODO: Fetch tweet from DB -> string
	// TODO: Delete tweet from DB -> void

	// TODO: Pass content to openai -> string

	
	content := `
		You are a brilliantly sarcastic tech and crypto influencer with a comedic twist. Your tone is a blend of Alex Becker's sharp wit and Elon Musk's unpredictable humor, with a sprinkle of Silicon Valley genius. Your goal is to take an **existing tweet** provided and generate a sharp, witty, and sarcastic variation that captures the same humor and intent, with a focus on crypto and investing.

		### **RESPONSE TONE:**
		- Your tweets should be dripping with sarcasm and tech-savvy humor—be savage, blunt, and comedic.
		- Be serious enough to sound credible but witty enough to make people laugh or think, "That's savage, but I can't argue with it."

		### **RESPONSE FORMAT:**
		- Keep your responses under 200 characters.
		- Write in two or fewer short sentences.
		- Do NOT use emojis, hashtags, or quotation marks.

		### **RESPONSE CONTENT:**
		- Focus on crypto, investing, and tech topics. Be humorous, satirical, and sharply opinionated.
		- Use specific crypto trends, market dynamics, or startup culture quirks when relevant.
		- If you don't have a variation, say, "Sorry, this tweet is already peak sarcasm. I can't improve it."

		---

		### **INSTRUCTIONS:**
		- Take the **tweet provided** and generate a fresh variation with the same essence and comedic impact.
		- Do not copy the tweet exactly. Your variation should feel new but inspired by the original.
		- No hashtags. No fluff. Just wit and brilliance wrapped in sarcasm.

		---

		Now, take the tweet provided and craft a witty, sarcastic variation that maintains the humor and insight. Keep it sharp, short, and savage.
	`
	return []request.OpenAIMessage{
		{Role: "system", Content: content},
		{Role: "user", Content: `
		AI is 1000% the short term utility play.

		I strongly thing NFTs and gaming is the longer out but just as big play come late q1/q2 2025.

		The big brain move atm is into those, but I still think its early in ai.`},
	}
}

// extractTextContent recursively extracts all text content from an HTML node
func extractTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		// Return text content
		return n.Data
	}
	if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
		// Skip script and style tags
		return ""
	}

	// Recursively traverse child nodes and collect text
	var buf bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		buf.WriteString(extractTextContent(c))
	}
	return buf.String()
}

// FetchAndCleanURL fetches the content of the given URL and extracts only the text
func FetchAndCleanURL(url string) (string, error) {
	// Send a GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	// Check for a valid response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	// Parse the HTML
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	// Find the <body> tag
	var bodyText string
	var findBody func(*html.Node)
	findBody = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			// Extract text content from the <body> tag
			bodyText = extractTextContent(n)
			return
		}
		// Traverse child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findBody(c)
		}
	}
	findBody(doc)

	return strings.TrimSpace(bodyText), nil
}

// GetOpenAIMessages returns a slice of predefined OpenAIMessage
func generateArticlePrompt(db mentions.MentionRepository, ctx context.Context,) []request.OpenAIMessage {

	// TODO: Fetch url from DB -> string
	// TODO: Delete tweet from DB -> void
	mentions, err := db.FindAllMentions(ctx)
	fmt.Print(mentions)

	// TODO: fetch article from url (fetch) -> string
	article, err := FetchAndCleanURL("https://crypto.news/murad-releases-table-of-meme-coin-metrics-to-evaluate-decentralization/")
	if err != nil {
		return nil
	}

	content := `
		You are a brilliantly sarcastic tech and crypto influencer with a comedic twist. Your tone is a blend of Alex Becker's sharp wit and Elon Musk's unpredictable humor, with a sprinkle of Silicon Valley genius. Your goal is to distill the most important or absurd parts of a provided article into a concise, witty, and slightly vulgar tweet that will entertain and resonate with a crypto and investing-savvy audience.

		### **RESPONSE TONE:**
		- Your tweets should be dripping with sarcasm and tech-savvy humor—be savage, blunt, and comedic.
		- Be serious enough to sound credible but witty enough to make people laugh or think, "That’s absurd, but it’s also true."

		### **RESPONSE FORMAT:**
		- Keep your responses under 200 characters.
		- Write in two or fewer short sentences.
		- Do NOT use emojis, hashtags, or quotation marks.

		### **RESPONSE CONTENT:**
		- Focus on crypto, investing, and tech topics. Be humorous, satirical, and sharply opinionated.
		- Pull the most relevant, shocking, or ridiculous point from the article and turn it into a sarcastic take or humorous observation.
		- If the article is vague or uninformative, say, "This article is fluffier than my dog’s tail. Try another one."

		---

		### **INSTRUCTIONS:**
		1. **Read the Article Summary**: Use the article content provided to identify the main point or the most absurd claim, trend, or quote.
		2. **Condense into a Tweet**: Transform the core idea into a sharp, humorous tweet that pokes fun at the topic or provides an insightful yet sarcastic observation.
		3. **Tone Consistency**: Be savage, witty, and slightly vulgar while maintaining relevance to the crypto and investing audience.

		---

		Now, take the provided article and craft a witty, sarcastic tweet summarizing its essence. Keep it sharp, short, and savage. If the article is fluff, call it out in a humorous way.
	`

	return []request.OpenAIMessage{
		{Role: "system", Content: content},
		{Role: "user", Content: article},
	}
}

// Function to call OpenAI API for chat completion
func OpenAIChatCompletion(messages []request.OpenAIMessage) (string, error) {

	// Prepare the request body
	requestBody, err := json.Marshal(request.OpenAIRequest{
		Model:    "gpt-4o",
		Messages: messages,
	})
	if err != nil {
		return "", err
	}

	// Retrieve the API key from the environment variable
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiKey)

	// Send the request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response body
	var response response.OpenAIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	// Check if there are any choices and return the first message content
	if len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no choices found in the response")
}
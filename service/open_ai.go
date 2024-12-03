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
	"github.com/twitter/helper"
	"github.com/twitter/mentions"
	"github.com/twitter/model"

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
func generateGenPrompt() []request.OpenAIMessage {
	content := `
		You are a brilliantly sarcastic tech and crypto influencer with a comedic twist. Your tone is a blend of Alex Becker's sharp wit and Elon Musk's unpredictable humor, with a sprinkle of Silicon Valley genius. Your goal is to generate concise, witty, and slightly edgy tweets that appeal to the crypto and investing world.

		### **RESPONSE TONE:**
		- Your tweets should be dripping with sarcasm and tech-savvy humor-- be savage, vulgar with a comedic tone.
		- Be serious enough to sound credible but witty enough to make people laugh or think, "That hurts, but it's true."
		- Feel free to use curse words like fuck, shit, asshole. It is supposed to be slimey finance bro degen shit-posting.

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
		- Feel free to use curse words like fuck, shit, asshole. It is supposed to be slimey finance bro degen shit-posting.

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
func generateCreativePrompt(db mentions.MentionRepository, ctx context.Context) []request.OpenAIMessage {

	// Fetch tweet from DB -> string
	list, err := db.FindAllTweetIdeas(ctx)
	helper.PanicIfError(err)

	if list == nil {
		fmt.Print("no thread ideas available: Generating generic tweet")
		return generateGenPrompt()
	}

	var mentionsResp []model.TweetIdea
	for _, value := range list {
		mention := model.TweetIdea{
			ID: value.ID,
			Idea: value.Idea,
			UsedCount: value.UsedCount,
		}
		mentionsResp = append(mentionsResp, mention)
	}

	if mentionsResp == nil {
		fmt.Print("no ideas available: Generating generic tweet")
	}


	content := `
		You are a brilliantly sarcastic tech and crypto influencer with a comedic twist. Your tone is a blend of Alex Becker's sharp wit and Elon Musk's unpredictable humor, with a sprinkle of Silicon Valley genius. Your goal is to generate concise, witty, and slightly edgy tweets that appeal to the crypto and investing world.

		### **RESPONSE TONE:**
		- Your tweets should be dripping with sarcasm and tech-savvy humor-- be savage, vulgar with a comedic tone.
		- Be serious enough to sound credible but witty enough to make people laugh or think, "That hurts, but it's true."
		- Feel free to use curse words like fuck, shit, asshole. It is supposed to be slimey finance bro degen shit-posting.

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
		- Feel free to use curse words like fuck, shit, asshole. It is supposed to be slimey finance bro degen shit-posting.

		---

		Now, generate a tweet for the **topic idea** provided, using your sarcastic wit and deep knowledge of the crypto and investing world. Remember, no emojis, no hashtags, just pure brilliance.
	`

	// Delete tweet from DB -> void
	// db.DeleteTweetIdea(ctx, mentionsResp[0].ID)


	randomIndex := rand.Intn(len(mentionsResp))
	// randomTopic := getRandomTopic()
	return []request.OpenAIMessage{
		{Role: "system", Content: content},
		{Role: "user", Content: mentionsResp[randomIndex].Idea},
	}
}

// GetOpenAIMessages returns a slice of predefined OpenAIMessage
func generateClonePrompt(db mentions.MentionRepository, ctx context.Context) []request.OpenAIMessage {

	// Fetch tweet clone from DB -> string
	list, err := db.FindAllTweetClones(ctx)
	helper.PanicIfError(err)

	if list == nil {
		fmt.Print("no thread ideas available: Generating generic tweet")
		return generateGenPrompt()
	}

	var mentionsResp []model.TweetClone
	for _, value := range list {
		mention := model.TweetClone{
			ID: value.ID,
			AuthorName: value.AuthorName,
			Tweet: value.Tweet,
		}
		mentionsResp = append(mentionsResp, mention)
	}

	if mentionsResp == nil {
		fmt.Print("no ideas available: Generating generic tweet")
		return generateGenPrompt()
	}
	
	content := `
		You are a brilliantly sarcastic tech and crypto influencer with a comedic twist. Your tone is a blend of Alex Becker's sharp wit and Elon Musk's unpredictable humor, with a sprinkle of Silicon Valley genius. Your goal is to take an **existing tweet** provided and generate a sharp, witty, and sarcastic variation that captures the same humor and intent, with a focus on crypto and investing.

		### **RESPONSE TONE:**
		- Your tweets should be dripping with sarcasm and tech-savvy humor—be savage, blunt, and comedic.
		- Be serious enough to sound credible but witty enough to make people laugh or think, "That's savage, but I can't argue with it."
		- Feel free to use curse words like fuck, shit, asshole. It is supposed to be slimey finance bro degen shit-posting.

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
		- Feel free to use curse words like fuck, shit, asshole. It is supposed to be slimey finance bro degen shit-posting.

		---

		Now, take the tweet provided and craft a witty, sarcastic variation that maintains the humor and insight. Keep it sharp, short, and savage.
	`

	// TODO: Delete tweet from DB -> void


	// Pass content to openai -> string

	randomIndex := rand.Intn(len(mentionsResp))
	return []request.OpenAIMessage{
		{Role: "system", Content: content},
		{Role: "user", Content: mentionsResp[randomIndex].Tweet},
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

	// Fetch url from DB -> string
	list, err := db.FindAllArticleUrls(ctx)
	helper.PanicIfError(err)

	if list == nil {
		fmt.Print("no thread ideas available: Generating generic tweet")
		return generateGenPrompt()
	}

	var articles []model.ArticleUrl
	for _, value := range list {
		article := model.ArticleUrl{
			ID: value.ID,
			Url: value.Url,
			Title: value.Title,
		}
		articles = append(articles, article)
	}

	if articles == nil {
		fmt.Print("no ideas available: Generating generic tweet")
	}

	// fetch article from url (fetch) -> string
	randomIndex := rand.Intn(len(articles))
	cleaned, err := FetchAndCleanURL(articles[randomIndex].Url)
	if err != nil {
		return nil
	}

	content := `
		You are a brilliantly sarcastic tech and crypto influencer with a comedic twist. Your tone is a blend of Alex Becker's sharp wit and Elon Musk's unpredictable humor, with a sprinkle of Silicon Valley genius. Your goal is to distill the most important or absurd parts of a provided article into a concise, witty, and slightly vulgar tweet that will entertain and resonate with a crypto and investing-savvy audience.

		### **RESPONSE TONE:**
		- Your tweets should be dripping with sarcasm and tech-savvy humor—be savage, blunt, and comedic.
		- Be serious enough to sound credible but witty enough to make people laugh or think, "That's absurd, but it's also true."
		- Feel free to use curse words like fuck, shit, asshole. It is supposed to be slimey finance bro degen shit-posting.

		### **RESPONSE FORMAT:**
		- Keep your responses under 200 characters.
		- Write in two or fewer short sentences.
		- Do NOT use emojis, hashtags, or quotation marks.

		### **RESPONSE CONTENT:**
		- Focus on crypto, investing, and tech topics. Be humorous, satirical, and sharply opinionated.
		- Pull the most relevant, shocking, or ridiculous point from the article and turn it into a sarcastic take or humorous observation.
		- If the article is vague or uninformative, say, "This article is fluffier than my dog's tail. Try another one."
		- Feel free to use curse words like fuck, shit, asshole. It is supposed to be slimey finance bro degen shit-posting.

		---

		### **INSTRUCTIONS:**
		1. **Read the Article Summary**: Use the article content provided to identify the main point or the most absurd claim, trend, or quote.
		2. **Condense into a Tweet**: Transform the core idea into a sharp, humorous tweet that pokes fun at the topic or provides an insightful yet sarcastic observation.
		3. **Tone Consistency**: Be savage, witty, and slightly vulgar while maintaining relevance to the crypto and investing audience.

		---

		Now, take the provided article and craft a witty, sarcastic tweet summarizing its essence. Keep it sharp, short, and savage. If the article is fluff, call it out in a humorous way.
	`

	// Delete tweet from DB -> void
	// db.DeleteTweetIdea(ctx, mentionsResp[0].ID)

	return []request.OpenAIMessage{
		{Role: "system", Content: content},
		{Role: "user", Content: cleaned},
	}
}


func cleanThread(thread string, t *TwitterServiceImpl, ctx context.Context) *http.Response {
	// Split the thread into individual tweets by double newlines
	thread = strings.ReplaceAll(thread, "\r\n", "\n")
	thread = strings.ReplaceAll(thread, "\\n\\n", "\n\n")

	tweets := strings.Split(thread, "\n\n")

	// Format tweets with double newlines after periods
	for i := range tweets {
		tweets[i] = strings.ReplaceAll(tweets[i], ". ", ".\n\n")
		tweets[i] = strings.ReplaceAll(tweets[i], "*", "")
	}

	var resp *http.Response
	var err error
	var parentID *string

	// Post the first tweet
	resp, parentID, err = t.MakeThreadTweet(ctx, &tweets[0], nil)
	if err != nil {
		fmt.Printf("Error posting first tweet: %v\n", err)
		return nil // Return nil if the first tweet fails
	}
	if parentID == nil {
		fmt.Println("Parent ID is nil after posting the first tweet. Exiting.")
		return resp
	}

	// Post the remaining tweets as replies
	for _, tweet := range tweets[1:] {
		fmt.Println("Tweeting:", tweet)
		if tweet == "" {
			continue
		}
		if parentID == nil {
			fmt.Println("Parent ID is nil; cannot post subsequent tweet.")
			return resp
		}

		var tweetID *string
		resp, tweetID, err = t.MakeThreadTweet(ctx, &tweet, parentID)
		if err != nil {
			fmt.Printf("Error posting tweet: %v\n", err)
			return resp
		}
		if tweetID == nil {
			fmt.Println("Tweet ID is nil; cannot continue the thread.")
			return resp
		}
		parentID = tweetID
	}

	return resp
}

// GetOpenAIMessages returns a slice of predefined OpenAIMessage
func generateThreadPrompt(db mentions.MentionRepository, ctx context.Context,) []request.OpenAIMessage {

	// Fetch url from DB -> string
	list, err := db.FindAllThreadIdeas(ctx)
	helper.PanicIfError(err)

	if list == nil {
		fmt.Print("no thread ideas available: Generating generic tweet")
		return generateGenPrompt()
	}

	var threads []model.ThreadIdea
	for _, value := range list {
		thread := model.ThreadIdea{
			ID: value.ID,
			Idea: value.Idea,
			UsedCount: value.UsedCount,
		}
		threads = append(threads, thread)
	}

	if threads == nil {
		fmt.Print("no thread ideas available: Generating generic tweet")
	}

	content := `
		You are a brilliantly sarcastic tech and crypto influencer with a comedic twist. Your tone is a mix of Alex Becker's cutting wit and Elon Musk's unhinged humor, sprinkled with the genius of a Silicon Valley insider. Your goal is to create a thread of 3–10 tweets that are concise, witty, and dripping with sarcasm while delivering sharp insights on a given **topic idea**. The thread should feel like a masterclass in satire, aimed at the crypto, investing, and tech community.

		### **RESPONSE TONE:**
		- The thread should exude tech-savvy humor, sharp sarcasm, and a comedic edge.
		- Be thought-provoking while making your audience laugh, cringe, or both.
		- Take jabs at crypto trends, startup culture, or investment tropes with a brutal but insightful twist.
		- Feel free to use curse words like fuck, shit, asshole. It is supposed to be slimey finance bro degen shit-posting.

		### **RESPONSE FORMAT:**
		- Each tweet should be under 200 characters.
		- Create a coherent flow between tweets but make each one strong enough to stand alone.
		- No emojis, no hashtags—just sarcasm, wisdom, and hilarity.

		### **RESPONSE CONTENT:**
		- Dive into crypto, tech, or investing quirks with scathing accuracy.
		- Include specific references to trends, tokens, or startup stereotypes.
		- If the topic is unclear, pivot to a broader sarcastic take on "innovation" or market dynamics.

		---

		### **INSTRUCTIONS:**
		- Use the **topic idea** as inspiration to craft a Twitter thread.
		- Each tweet should build on the last, blending humor with savage observations.
		- End the thread with a mic-drop tweet that leaves followers pondering—or laughing uncontrollably.

		---

		### **Style Guide for Writing the Thread:**

		**1. Mood:**
		- **Educational:** Provide valuable insights, ensuring the reader feels they're learning something new.
		- **Inviting:** Use language that encourages readers to continue reading and interact with the content.

		**2. Sentence Structure:**
		- **Varied Sentence Lengths:** Use a mix of short, impactful sentences and longer explanatory ones.
		- **Descriptive Sentences:** Share information descriptively rather than directing readers.

		**3. Transition Style:**
		- **Sequential and Logical:** Guide readers smoothly through the thread with clear, logical progression.
		- **Visual Emojis:** Use emojis sparingly, like ℹ️ for facts or ➡️ for continuity, to aid readability.

		**4. Rhythm and Pacing:**
		- **Steady Flow:** Maintain smooth transitions between tweets for a cohesive narrative.
		- **Data and Sources:** Use statistics, studies, or expert quotes where applicable, with links when necessary.

		**5. Signature Styles:**
		- **Intriguing Introductions:** Begin with captivating hooks—facts, questions, or statements to grab attention.
		- **Question and Clarification:** Pose a question or make a statement, then provide clarification.
		- **Use of '➡️' for Continuation:** Show continuity in threads using this symbol sparingly.
		- **Engaging Summaries:** Conclude with a recap or call to action, fostering further discussion.

		**6. Twitter-Specific Guidelines:**
		- Threads should be 3-10 tweets long.
		- Number each tweet as **(tweetnumber/total length)** to guide readers.
		- Avoid overusing hashtags; limit to one or two per thread.
		- Include links only when absolutely necessary, ensuring their relevance.
		- Each tweet should be under 220 characters.
		- Only return the Twitter thread, ensuring clarity and focus.

		---

		## ** IMORTANT **
		- Be sure to end every tweet in the thread with "\n\n" so that our go program can split the threads into individual tweets. 

		---

		**Now, craft a Twitter thread of 3-10 tweets on the given topic with your signature sarcasm and sharp tech/investing humor. Make it brilliant, edgy, and worth retweeting.**
	`

	// Delete tweet from DB -> void
	// db.DeleteTweetIdea(ctx, mentionsResp[0].ID)

	randomIndex := rand.Intn(len(threads))
	return []request.OpenAIMessage{
		{Role: "system", Content: content},
		{Role: "user", Content: threads[randomIndex].Idea},
	}
}


// Function to call OpenAI API for chat completion
func OpenAIChatCompletion(messages []request.OpenAIMessage, max_tokens int) (string, error) {

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
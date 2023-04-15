package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func createChat(cookie string) chatSession {
	// Make a http request to the bing chat api
	req, err := http.NewRequest("GET", "https://www.bing.com/turing/conversation/create", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authority", "www.bing.com")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "de,de-DE;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Referer", "https://www.bing.com/search?q=")
	req.Header.Set("Sec-Ch-Ua", "^^Chromium^^;v=^^110^^, ^^Not")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
		// handle err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data createConversation
	err2 := json.Unmarshal(body, &data)
	if err != nil {
		panic(err2)
	}

	return chatSession{
		ConversationID:        data.ConversationID,
		ClientID:              data.ClientID,
		ConversationSignature: data.ConversationSignature,
		InvocationId:          "1",
	}
}

func test() {
	// open file
	file, err := os.Open("response.gob")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// read file
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// decode gob
	var response chatResponse
	dec := gob.NewDecoder(bytes.NewReader(content))
	err = dec.Decode(&response)
	if err != nil {
		panic(err)
	}

	spew.Dump(response)
}

func main() {
	// load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	chatSessionMain := chatSession{}
	// if you want to get a new conversation id, delete the session_test.gob file
	if _, err := os.Stat("session.gob"); os.IsNotExist(err) {
		// you have to get the cookie from the browser and place it in the .env file
		fmt.Println("Requesting new conversation session...")
		chatSessionMain = createChat(os.Getenv("COOKIE"))
		saveChat(chatSessionMain)
	} else {
		fmt.Println("Loading conversation session...")
		chatSessionMain = loadChat()
	}

	session, data := parseJSON(chatSessionMain, "Write a new Sherlock Holmes Story in the style of Arthur Conan Doyle", 2)
	chatSessionMain = session
	saveChat(chatSessionMain)

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	conversation(string(b))
}

func saveChat(session chatSession) {
	// encode gob
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(session)
	if err != nil {
		panic(err)
	}

	// write to file
	file, err := os.Create("session.gob")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	_, err = file.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}
}

func loadChat() chatSession {
	// open file
	file, err := os.Open("session.gob")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// read file
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// decode gob
	var session chatSession
	dec := gob.NewDecoder(bytes.NewReader(content))
	err = dec.Decode(&session)
	if err != nil {
		panic(err)
	}

	return session
}

func parseJSON(session chatSession, text string, mode int) (chatSession, startMessage) {
	// open the right json file
	switch mode {
	case 1:
		file, err := os.Open("normal.json")
		if err != nil {
			panic(err)
		}
		return generateMessage(file, session, text)
	case 2:
		file, err := os.Open("exact.json")
		if err != nil {
			panic(err)
		}
		return generateMessage(file, session, text)
	case 3:
		file, err := os.Open("creative.json")
		if err != nil {
			panic(err)
		}
		return generateMessage(file, session, text)
	default:
		file, err := os.Open("normal.json")
		if err != nil {
			panic(err)
		}
		return generateMessage(file, session, text)
	}
}

func generateMessage(file *os.File, session chatSession, text string) (chatSession, startMessage) {
	// read file
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var data startMessage
	err2 := json.Unmarshal(content, &data)
	if err != nil {
		panic(err2)
	}

	// replace values
	data.Arguments[0].ConversationID = session.ConversationID
	data.Arguments[0].Participant.ID = session.ClientID
	data.Arguments[0].ConversationSignature = session.ConversationSignature
	data.InvocationID = session.InvocationId
	data.Arguments[0].Message.Timestamp = time.Now().Format(time.RFC3339)
	data.Arguments[0].Message.Text = text

	if session.InvocationId == "1" {
		data.Arguments[0].IsStartOfSession = true
	} else {
		data.Arguments[0].IsStartOfSession = false
	}
	number, _ := strconv.Atoi(session.InvocationId)
	session.InvocationId = strconv.Itoa(1 + number)

	return session, data
}

func conversation(data string) {
	defer wg.Wait()
	ctx := context.Background()

	c, _, err := websocket.Dial(ctx, "wss://sydney.bing.com/sydney/ChatHub", &websocket.DialOptions{CompressionMode: websocket.CompressionDisabled})
	if err != nil {
		panic(err)
	}

	// send hello message
	err = c.Write(ctx, websocket.MessageText, []byte("{\"protocol\":\"json\",\"version\":1}\u001E"))
	if err != nil {
		panic(err)
	}
	fmt.Println("Sent: ", "{\"protocol\":\"json\",\"version\":1}\u001E")

	for {
		_, i, err := c.Reader(ctx)
		if err != nil {
			return
		}
		content, err := io.ReadAll(i)
		if err != nil {
			panic(err)
		}

		msg := string(content)

		go parseChatMessage(msg)
		wg.Add(1)

		if msg == "{}\u001E" {
			err = c.Write(ctx, websocket.MessageText, []byte("{\"type\":6}\u001E"))
			if err != nil {
				return
			}
			fmt.Println("Sent: ", "{\"type\":6}\u001E")
			err1 := c.Write(ctx, websocket.MessageText, []byte(data+"\u001E"))
			if err1 != nil {
				return
			}
			fmt.Println("Sent: ", data+"\u001E")
		}
	}
}

func parseChatMessage(msg string) {
	defer wg.Done()

	msg = strings.Split(msg, "\u001E")[0]

	// parse the message to json
	var response chatResponse
	err2 := json.Unmarshal([]byte(msg), &response)
	if err2 != nil {
		panic(err2)
	}

	if response.Type == 2 {
		// save response
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(response)
		if err != nil {
			panic(err)
		}

		// write to file
		file, err := os.Create("response.gob")
		if err != nil {
			panic(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)

		_, err = file.Write(buf.Bytes())
		if err != nil {
			panic(err)
		}

		fmt.Println("Saved response")
		fmt.Println(response.Item.Messages[1].Text)
		fmt.Println(response.Item.Messages[1].SourceAttributions)
	} else if len(response.Arguments) > 0 {
		if len(response.Arguments[0].Messages) > 0 {
			cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
			cmd.Stdout = os.Stdout
			err := cmd.Run()
			if err != nil {
				return
			}
			fmt.Println(response.Arguments[0].Messages[0].Text)
		}
	}
}

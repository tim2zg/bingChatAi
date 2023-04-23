package bingChatAi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nhooyr.io/websocket"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func CreateChat(cookie string) (ChatSession, error) {
	// Make a http request to the bing chat api
	req, err := http.NewRequest("GET", "https://www.bing.com/turing/conversation/create", nil)
	if err != nil {
		return ChatSession{}, err
	}

	// Set the headers
	req.Header.Set("Authority", "www.bing.com")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "de,de-DE;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Referer", "https://www.bing.com/search?q=")
	req.Header.Set("Sec-Ch-Ua", "^^Chromium^^;v=^^110^^, ^^Not")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ChatSession{}, err
	}

	// Close the body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// Read the body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatSession{}, err
	}

	if strings.Contains(string(body), "Sorry") {
		return ChatSession{}, fmt.Errorf("error: %s", string(body))
	} else {
		// Unmarshal the json
		var data createConversation
		err2 := json.Unmarshal(body, &data)
		if err != nil {
			return ChatSession{}, err2
		}

		return ChatSession{
			ConversationID:        data.ConversationID,
			ClientID:              data.ClientID,
			ConversationSignature: data.ConversationSignature,
			InvocationId:          "1",
		}, nil
	}
}

func ParseJSON(session ChatSession, text string, mode int) (ChatSession, StartMessage, error) {
	// open the right json file
	switch mode {
	case 1:
		return generateMessage(GetNormal(), session, text)
	case 2:
		return generateMessage(GetExact(), session, text)
	case 3:
		return generateMessage(GetCreative(), session, text)
	default:
		return generateMessage(GetNormal(), session, text)
	}
}

func generateMessage(mode string, session ChatSession, text string) (ChatSession, StartMessage, error) {
	// unmarshal json
	var data StartMessage
	err2 := json.Unmarshal([]byte(mode), &data)
	if err2 != nil {
		return ChatSession{}, StartMessage{}, err2
	}

	// replace values with the ones from the session
	data.Arguments[0].ConversationID = session.ConversationID
	data.Arguments[0].Participant.ID = session.ClientID
	data.Arguments[0].ConversationSignature = session.ConversationSignature
	data.InvocationID = session.InvocationId
	data.Arguments[0].Message.Timestamp = time.Now().Format(time.RFC3339)
	data.Arguments[0].Message.Text = text

	// check if it is the first message of the session
	if session.InvocationId == "1" {
		data.Arguments[0].IsStartOfSession = true
	} else {
		data.Arguments[0].IsStartOfSession = false
	}

	// increment invocation id
	number, err := strconv.Atoi(session.InvocationId)
	if err != nil {
		return ChatSession{}, StartMessage{}, err
	}
	session.InvocationId = strconv.Itoa(1 + number)

	return session, data, nil
}

func Conversation(data string, print bool, printRaw bool, debug bool) (ChatResponse, []ChatResponse, error) {
	msgChan := make(chan ChatResponse)
	var responses []ChatResponse

	// on close
	defer func() {
		wg.Wait()
		if print {
			fmt.Println("Conversation closed")
		}
	}()
	ctx := context.Background()

	// connect to websocket
	c, _, err := websocket.Dial(ctx, "wss://sydney.bing.com/sydney/ChatHub", &websocket.DialOptions{CompressionMode: websocket.CompressionDisabled})
	if err != nil {
		return ChatResponse{}, []ChatResponse{}, err
	}
	if debug {
		fmt.Println("Connected to websocket")
	}

	// send hello message
	err = c.Write(ctx, websocket.MessageText, []byte("{\"protocol\":\"json\",\"version\":1}\u001E"))
	if err != nil {
		return ChatResponse{}, []ChatResponse{}, err
	}
	if debug {
		fmt.Println("Sent hello message")
	}

	// read incoming messages if the connection gets closed return
	for {
		// read message
		_, i, err := c.Reader(ctx)

		// if the connection gets closed return
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				// close connection and return
				// read the chan
				response := <-msgChan
				return response, responses, nil
			} else {
				// connection error
				return ChatResponse{}, []ChatResponse{}, err
			}
		} else {
			// read message content
			content, err := io.ReadAll(i)
			if err != nil {
				return ChatResponse{}, []ChatResponse{}, err
			}

			// convert to string
			msg := string(content)

			// parse the message in a new goroutine
			wg.Add(1)
			go func() {
				if printRaw {
					fmt.Println("Raw Message: ")
					fmt.Println(msg)
				}

				parsed, err := parseChatMessage(msg, msgChan)
				if err != nil {
					fmt.Println("error parsing message: ", err)
				}

				if print {
					if len(parsed.Item.Messages) >= 1 {
						if parsed.Item.Messages[1].Text != "" {
							fmt.Println(parsed.Item.Messages[1].Text)
						}
						if len(parsed.Item.Messages[1].SourceAttributions) > 0 {
							fmt.Println(parsed.Item.Messages[1].SourceAttributions)
						}
					} else if len(parsed.Arguments) > 0 {
						if len(parsed.Arguments[0].Messages) > 0 {
							if parsed.Arguments[0].Messages[0].Text != "" {
								fmt.Println(parsed.Arguments[0].Messages[0].Text)
							}
						}
					}
				}

				if debug {
					fmt.Println("Parsed message")
					fmt.Println(parsed)
				}

				// add to the array
				responses = append(responses, parsed)
			}()

			// wait fot the answer
			if msg == "{}\u001E" {
				// send our command
				err = c.Write(ctx, websocket.MessageText, []byte("{\"type\":6}\u001E"))
				if err != nil {
					return ChatResponse{}, []ChatResponse{}, err
				}
				err1 := c.Write(ctx, websocket.MessageText, []byte(data+"\u001E"))
				if err1 != nil {
					return ChatResponse{}, []ChatResponse{}, err1
				}
				if debug {
					fmt.Println("Sent command: ", data)
				}
			}
		}
	}
}

func parseChatMessage(msg string, back chan ChatResponse) (ChatResponse, error) {
	defer wg.Done()

	msg = strings.Split(msg, "\u001E")[0]

	// parse the message to json
	var response ChatResponse
	err2 := json.Unmarshal([]byte(msg), &response)
	if err2 != nil {
		return ChatResponse{}, err2
	}

	// see if response is type 2
	if response.Type == 2 {
		// send the response back
		back <- response
	}

	// return message
	return response, nil
}

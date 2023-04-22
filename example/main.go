package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/tim2zg/bingChatAi"
	"io"
	"log"
	"os"
)

func main() {
	// load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	chatSessionMain := bingChatAi.ChatSession{}
	// if you want to get a new conversation id, delete the session_test.gob file
	if _, err := os.Stat("session.gob"); os.IsNotExist(err) {
		// you have to get the cookie from the browser and place it in the .env file
		fmt.Println("Requesting new conversation session...")
		chatSessionMain, err = bingChatAi.CreateChat(os.Getenv("COOKIE")) // done
		if err != nil {
			fmt.Println(err)
			return
		}
		saveChat(chatSessionMain)
		fmt.Println("New conversation session created.")
	} else {
		fmt.Println("Loading conversation session...")
		chatSessionMain = loadChat()
		fmt.Println("Conversation session loaded.")
	}

	session, data, err := bingChatAi.ParseJSON(chatSessionMain, "Can you write a long story with multiple prompts?", 3)
	if err != nil {
		fmt.Println(err)
		return
	}
	chatSessionMain = session
	saveChat(chatSessionMain)

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	conversation, conversations, err := bingChatAi.Conversation(string(b), true, false, false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("--------------------")
	fmt.Println(conversation)
	fmt.Println(len(conversations))
	saveResponse(conversations)

	// load response
	conversations = loadResponse()
	fmt.Println(len(conversations))
}

func loadResponse() []bingChatAi.ChatResponse {
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
	var session []bingChatAi.ChatResponse
	dec := gob.NewDecoder(bytes.NewReader(content))
	err = dec.Decode(&session)
	if err != nil {
		panic(err)
	}

	fmt.Println("Loaded response")

	return session
}

func saveResponse(response []bingChatAi.ChatResponse) {
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
}

func saveChat(session bingChatAi.ChatSession) {
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

func loadChat() bingChatAi.ChatSession {
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
	var session bingChatAi.ChatSession
	dec := gob.NewDecoder(bytes.NewReader(content))
	err = dec.Decode(&session)
	if err != nil {
		panic(err)
	}

	return session
}

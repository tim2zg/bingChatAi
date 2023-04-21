package bingChatAi

import "time"

type ChatSession struct {
	ConversationID, ClientID, ConversationSignature, InvocationId string
}

type createConversation struct {
	ConversationID        string `json:"conversationId"`
	ClientID              string `json:"clientId"`
	ConversationSignature string `json:"conversationSignature"`
	Result                struct {
		Value   string `json:"value"`
		Message any    `json:"message"`
	} `json:"result"`
}

type StartMessage struct {
	Arguments []struct {
		Source              string   `json:"source"`
		OptionsSets         []any    `json:"optionsSets"`
		AllowedMessageTypes []string `json:"allowedMessageTypes"`
		SliceIds            []any    `json:"sliceIds"`
		Verbosity           string   `json:"verbosity"`
		TraceID             string   `json:"traceId"`
		IsStartOfSession    bool     `json:"isStartOfSession"`
		Message             struct {
			Locale        string `json:"locale"`
			Market        string `json:"market"`
			Region        string `json:"region"`
			Location      string `json:"location"`
			LocationHints []any  `json:"locationHints"`
			Timestamp     string `json:"timestamp"`
			Author        string `json:"author"`
			InputMethod   string `json:"inputMethod"`
			Text          string `json:"text"`
			MessageType   string `json:"messageType"`
		} `json:"message"`
		ConversationSignature string `json:"conversationSignature"`
		Participant           struct {
			ID string `json:"id"`
		} `json:"participant"`
		ConversationID string `json:"conversationId"`
	} `json:"arguments"`
	InvocationID string `json:"invocationId"`
	Target       string `json:"target"`
	Type         int    `json:"type"`
}

type ChatResponse struct {
	Type         int    `json:"type"`
	InvocationID string `json:"invocationId"`
	Target       string `json:"target"`
	Arguments    []struct {
		RequestID  string `json:"requestId"`
		Throttling struct {
			MaxNumUserMessagesInConversation int `json:"maxNumUserMessagesInConversation"`
			NumUserMessagesInConversation    int `json:"numUserMessagesInConversation"`
		} `json:"throttling"`
		Cursor struct {
			J string `json:"j"`
			P int    `json:"p"`
		} `json:"cursor"`
		Messages []struct {
			ContentType string    `json:"contentType"`
			Text        string    `json:"text"`
			HiddenText  string    `json:"hiddenText"`
			Author      string    `json:"author"`
			CreatedAt   time.Time `json:"createdAt"`
			Timestamp   time.Time `json:"timestamp"`
			MessageID   string    `json:"messageId"`
			RequestID   string    `json:"requestId"`
			MessageType string    `json:"messageType"`
			Offense     string    `json:"offense"`
			From        struct {
				ID   string `json:"id"`
				Name any    `json:"name"`
			} `json:"from,omitempty"`
			AdaptiveCards []struct {
				Type    string `json:"type"`
				Version string `json:"version"`
				Body    []struct {
					Type    string `json:"type"`
					Warp    bool   `json:"wrap"`
					Text    string `json:"text"`
					Size    string `json:"size"`
					Inlines []struct {
						Type     string `json:"type"`
						IsSubtle bool   `json:"isSubtle"`
						Italic   bool   `json:"italic"`
						Text     string `json:"text"`
					} `json:"inlines"`
				} `json:"body"`
			} `json:"adaptiveCards"`
			Feedback struct {
				Tag       any    `json:"tag"`
				UpdatedOn any    `json:"updatedOn"`
				Type      string `json:"type"`
			} `json:"feedback"`
			ContentOrigin string `json:"contentOrigin"`
			Privacy       any    `json:"privacy"`
			SpokenText    string `json:"spokenText"`
		} `json:"messages"`
	} `json:"arguments"`
	Item struct {
		Messages []struct {
			Text        string `json:"text,omitempty"`
			Author      string `json:"author"`
			MessageID   string `json:"messageId"`
			RequestID   string `json:"requestId"`
			ContentType string `json:"contentType"`
			From        struct {
				ID   string `json:"id"`
				Name any    `json:"name"`
			} `json:"from,omitempty"`
			CreatedAt     time.Time `json:"createdAt"`
			Timestamp     time.Time `json:"timestamp"`
			Locale        string    `json:"locale,omitempty"`
			Market        string    `json:"market,omitempty"`
			Region        string    `json:"region,omitempty"`
			Location      string    `json:"location,omitempty"`
			LocationHints []struct {
				Country           string `json:"country"`
				CountryConfidence int    `json:"countryConfidence"`
				State             string `json:"state"`
				City              string `json:"city"`
				CityConfidence    int    `json:"cityConfidence"`
				ZipCode           string `json:"zipCode"`
				TimeZoneOffset    int    `json:"timeZoneOffset"`
				SourceType        int    `json:"sourceType"`
				Center            struct {
					Latitude  float64 `json:"latitude"`
					Longitude float64 `json:"longitude"`
					Height    any     `json:"height"`
				} `json:"center"`
				RegionType int `json:"regionType"`
			} `json:"locationHints,omitempty"`
			Nlu struct {
				ScoredClassification struct {
					Classification string `json:"classification"`
					Score          any    `json:"score"`
				} `json:"scoredClassification"`
				ClassificationRanking []struct {
					Classification string `json:"classification"`
					Score          any    `json:"score"`
				} `json:"classificationRanking"`
				QualifyingClassifications any `json:"qualifyingClassifications"`
				Ood                       any `json:"ood"`
				MetaData                  any `json:"metaData"`
				Entities                  any `json:"entities"`
			} `json:"nlu,omitempty"`
			Offense  string `json:"offense"`
			Feedback struct {
				Tag       any    `json:"tag"`
				UpdatedOn any    `json:"updatedOn"`
				Type      string `json:"type"`
			} `json:"feedback"`
			ContentOrigin string `json:"contentOrigin"`
			Privacy       any    `json:"privacy"`
			InputMethod   string `json:"inputMethod,omitempty"`
			HiddenText    string `json:"hiddenText,omitempty"`
			MessageType   string `json:"messageType,omitempty"`
			AdaptiveCards []struct {
				Type    string `json:"type"`
				Version string `json:"version"`
				Body    []struct {
					Type    string `json:"type"`
					Warp    bool   `json:"wrap"`
					Text    string `json:"text"`
					Size    string `json:"size"`
					Inlines []struct {
						Type     string `json:"type"`
						IsSubtle bool   `json:"isSubtle"`
						Italic   bool   `json:"italic"`
						Text     string `json:"text"`
					} `json:"inlines"`
				} `json:"body"`
			} `json:"adaptiveCards,omitempty"`
			GroundingInfo struct {
				WebSearchResults []struct {
					Index           string   `json:"index"`
					Title           string   `json:"title"`
					Snippets        []string `json:"snippets"`
					Data            any      `json:"data"`
					Context         any      `json:"context"`
					URL             string   `json:"url"`
					LastUpdatedDate any      `json:"lastUpdatedDate"`
				} `json:"web_search_results"`
			} `json:"groundingInfo,omitempty"`
			SourceAttributions []struct {
				ProviderDisplayName string `json:"providerDisplayName"`
				SeeMoreURL          string `json:"seeMoreUrl"`
				SearchQuery         string `json:"searchQuery"`
				ImageLink           string `json:"imageLink"`
				ImageWidth          string `json:"imageWidth"`
				ImageHeight         string `json:"imageHeight"`
				ImageFavicon        string `json:"imageFavicon"`
			} `json:"sourceAttributions,omitempty"`
			SuggestedResponses []struct {
				Text        string    `json:"text"`
				Author      string    `json:"author"`
				CreatedAt   time.Time `json:"createdAt"`
				Timestamp   time.Time `json:"timestamp"`
				MessageID   string    `json:"messageId"`
				MessageType string    `json:"messageType"`
				Offense     string    `json:"offense"`
				Feedback    struct {
					Tag       any    `json:"tag"`
					UpdatedOn any    `json:"updatedOn"`
					Type      string `json:"type"`
				} `json:"feedback"`
				ContentOrigin string `json:"contentOrigin"`
				Privacy       any    `json:"privacy"`
			} `json:"suggestedResponses,omitempty"`
			SpokenText string `json:"spokenText"`
		} `json:"messages"`
		FirstNewMessageIndex       int       `json:"firstNewMessageIndex"`
		ConversationID             string    `json:"conversationId"`
		RequestID                  string    `json:"requestId"`
		ConversationExpiryTime     time.Time `json:"conversationExpiryTime"`
		ShouldInitiateConversation bool      `json:"shouldInitiateConversation"`
		Telemetry                  struct {
			Metrics   any       `json:"metrics"`
			StartTime time.Time `json:"startTime"`
		} `json:"telemetry"`
		Throttling struct {
			MaxNumUserMessagesInConversation int `json:"maxNumUserMessagesInConversation"`
			NumUserMessagesInConversation    int `json:"numUserMessagesInConversation"`
		} `json:"throttling"`
		Result struct {
			Value          string `json:"value"`
			Message        string `json:"message"`
			ServiceVersion string `json:"serviceVersion"`
		} `json:"result"`
		AllowReconnect bool `json:"allowReconnect"`
	} `json:"item"`
}

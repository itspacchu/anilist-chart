package anilist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/charmbracelet/log"
)

const ANILIST_GRAPHQL_ENDPOINT string = "https://graphql.anilist.co"

func FetchIDFromUsername(username string) int64 {
	query := `{"query":"query User($name: String) {\n  User(name: $name) {\n    id\n  }\n}","variables":{"name":"` + username + `"}}`
	resp, err := http.Post(ANILIST_GRAPHQL_ENDPOINT, "application/json", strings.NewReader(query))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	rStruct := Response{}
	rByte, _ := io.ReadAll(resp.Body)
	json.Unmarshal(rByte, &rStruct)
	return rStruct.Data.User.ID
}

func FetchActivitiesDetails(userID int64, timestamp_greater int64, activity_type string) Response {

	query := `{"query":"query ($page: Int, $perPage: Int, $userId: Int, $typeIn: [ActivityType], $createdAtGreater: Int) {\n  Page(page: $page, perPage: $perPage) {\n    activities(createdAt_greater: $createdAtGreater, userId: $userId, type_in: $typeIn) {\n      ... on ListActivity {\n        media {\n          coverImage {\n            medium\n large\n          }\n          title {\n            romaji\n  english\n         }\n          id\n        }\n        createdAt\n      }\n    }\n  }\n}","variables":{"createdAtGreater":%d ,"userId":%d ,"typeIn":"%s","page":%d,"perPage":%d}}`

	resp, err := http.Post(ANILIST_GRAPHQL_ENDPOINT, "application/json", strings.NewReader(fmt.Sprintf(query, timestamp_greater, userID, activity_type, 1, 100)))
	if err != nil {
		log.Warnf("Error in Response : %s", err.Error())
		return Response{}
	}
	defer resp.Body.Close()
	rStruct := Response{}
	rByte, _ := io.ReadAll(resp.Body)
	if resp.StatusCode > 300 {
		log.Debugf("Invalid response %d : %s", resp.StatusCode, string(rByte))
		return Response{}
	}
	json.Unmarshal(rByte, &rStruct)
	return rStruct
}

func FetchActivityDetails(userID int64, timestamp_greater int64, activity_type string) Response {
	query := `{"query":"query Activity($createdAtGreater: Int, $userId: Int, $typeIn: [ActivityType]) {\n  Activity(createdAt_greater: $createdAtGreater, userId: $userId, type_in: $typeIn) {\n    ... on ListActivity {\n      media {\n        coverImage {\n          large\n        }\n        title {\n          romaji\n        }\n        id\n      }\n      createdAt\n    }\n  }\n}","variables":{"createdAtGreater":%d,"userId":%d,"typeIn":[%s],}}`

	resp, err := http.Post(ANILIST_GRAPHQL_ENDPOINT, "application/json", strings.NewReader(fmt.Sprintf(query, timestamp_greater, userID, activity_type)))
	if err != nil {
		log.Warnf("Error in Response : %s", err.Error())
		return Response{}
	}
	defer resp.Body.Close()
	rStruct := Response{}
	rByte, _ := io.ReadAll(resp.Body)

	if resp.StatusCode > 300 {
		log.Debugf("Invalid response %d : %s", resp.StatusCode, string(rByte))
		return Response{}
	}
	json.Unmarshal(rByte, &rStruct)
	return rStruct
}

package ghjournal

import (
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func buildReport(date time.Time) (map[string]interface{}, error) {
	// REFACTOR: We should keep a connection to the DB around instead of
	//           connecting over and over for each request
	mongoHost := os.Getenv("MONGO_PORT_27017_TCP_ADDR")
	session, err := mgo.Dial(mongoHost)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	newProjects, err := fetchNewProjects(session, date)
	if err != nil {
		return nil, err
	}

	stars, err := fetchStars(session, date)
	if err != nil {
		return nil, err
	}

	// TODO: tags / releases

	forks, err := fetchForks(session, date)
	if err != nil {
		return nil, err
	}

	issuesUpdated, err := fetchIssuesUpdated(session, date)
	if err != nil {
		return nil, err
	}

	prsUpdated, err := fetchPRsUpdated(session, date)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"date":          date.Format("2006-01-02"),
		"newProjects":   newProjects,
		"forks":         forks,
		"stars":         stars,
		"issuesUpdated": issuesUpdated,
		"prsUpdated":    prsUpdated,
	}, nil
}

func fetchNewProjects(session *mgo.Session, startDate time.Time) (interface{}, error) {
	endDate := startDate.AddDate(0, 0, 1)
	pipe := session.DB("gh-journal").C("events").Pipe([]bson.M{
		{
			"$match": bson.M{
				"type":                 "CreateEvent",
				"raw.payload.ref_type": "repository",
				"created_at": bson.M{
					"$gte": startDate,
					"$lt":  endDate,
				},
			},
		},
		{
			"$project": bson.M{
				"_id":     false,
				"user":    "$actor",
				"project": bson.M{"$concat": []string{"$project.owner", "/", "$project.name"}},
			},
		},
		{"$sort": bson.M{"project": 1}},
	})
	data := []interface{}{}
	err := pipe.All(&data)
	return data, err
}

func fetchForks(session *mgo.Session, startDate time.Time) (interface{}, error) {
	endDate := startDate.AddDate(0, 0, 1)
	pipe := session.DB("gh-journal").C("events").Pipe([]bson.M{
		{
			"$match": bson.M{
				"type": "ForkEvent",
				"created_at": bson.M{
					"$gte": startDate,
					"$lt":  endDate,
				},
			},
		},
		{
			"$project": bson.M{
				"_id":           false,
				"user":          "$actor",
				"repository":    bson.M{"$concat": []string{"$project.owner", "/", "$project.name"}},
				"newRepository": "$raw.payload.forkee.full_name",
			},
		},
		{"$sort": bson.M{"user": 1}},
	})
	data := []interface{}{}
	err := pipe.All(&data)
	return data, err
}

func fetchStars(session *mgo.Session, startDate time.Time) (interface{}, error) {
	endDate := startDate.AddDate(0, 0, 1)
	pipe := session.DB("gh-journal").C("events").Pipe([]bson.M{
		{
			"$match": bson.M{
				"type": "WatchEvent",
				"created_at": bson.M{
					"$gte": startDate,
					"$lt":  endDate,
				},
			},
		},
		{
			"$project": bson.M{
				"actor":   true,
				"project": bson.M{"$concat": []string{"$project.owner", "/", "$project.name"}},
			},
		},
		{
			"$group": bson.M{
				"_id":        "$project",
				"project":    bson.M{"$first": "$project"},
				"stargazers": bson.M{"$addToSet": "$actor"},
			},
		},
		{"$sort": bson.M{"project": 1}},
	})
	data := []interface{}{}
	err := pipe.All(&data)
	return data, err
}

func fetchIssuesUpdated(session *mgo.Session, startDate time.Time) (interface{}, error) {
	endDate := startDate.AddDate(0, 0, 1)
	pipe := session.DB("gh-journal").C("events").Pipe([]bson.M{
		{
			"$match": bson.M{
				"type": bson.M{"$in": []string{"IssuesEvent", "IssueCommentEvent"}},
				"raw.payload.issue.pull_request": bson.M{"$exists": false},
				"created_at": bson.M{
					"$gte": startDate,
					"$lt":  endDate,
				},
			},
		},
		{
			"$project": bson.M{
				"repository": bson.M{"$concat": []string{"$project.owner", "/", "$project.name"}},
				"action":     true,
				"actor":      true,
				"state":      "$raw.payload.issue.state",
				"title":      "$raw.payload.issue.title",
				"number":     "$raw.payload.issue.number",
				"url":        "$raw.payload.issue.html_url",
			},
		},
		{
			"$group": bson.M{
				"_id":       bson.M{"r": "$repository", "n": "$number"},
				"url":       bson.M{"$first": "$url"},
				"actors":    bson.M{"$addToSet": "$actor"},
				"title":     bson.M{"$last": "$title"},
				"lastState": bson.M{"$last": "$state"},
			},
		},
		{
			"$group": bson.M{
				"_id": "$_id.r",
				"project": bson.M{"$first": "$_id.r"},
				"issues": bson.M{
					"$addToSet": bson.M{
						"number":    "$_id.n",
						"url":       "$url",
						"title":     "$title",
						"lastState": "$lastState",
						"actors":    "$actors",
					},
				},
			},
		},
		{"$sort": bson.M{"project": 1}},
	})
	data := []interface{}{}
	err := pipe.All(&data)
	return data, err
}

func fetchPRsUpdated(session *mgo.Session, startDate time.Time) (interface{}, error) {
	endDate := startDate.AddDate(0, 0, 1)
	pipe := session.DB("gh-journal").C("events").Pipe([]bson.M{
		{
			"$match": bson.M{
				"$or": []bson.M{
					bson.M{ "type": bson.M{"$in": []string{"PullRequestEvent", "PullRequestReviewCommentEvent"}} },
					bson.M{
						"type": "IssueCommentEvent",
						"raw.payload.issue.pull_request": bson.M{"$exists": true},
					},
				},
				"created_at": bson.M{
					"$gte": startDate,
					"$lt":  endDate,
				},
			},
		},
		{ "$sort": bson.M{"created_at": 1} },
		{
			"$project": bson.M{
				"_id":          false,
				"created_at":   true,
				"actor":        true,
				"action":       bson.M{ "$cond": bson.M{ "if": bson.M{ "$eq": []string{ "$type", "IssueCommentEvent" } }, "then": "commented", "else": "$action" } },
				"repository":   bson.M{ "$concat": []string{ "$project.owner", "/", "$project.name" } },
				"title":        bson.M{ "$ifNull": []string{ "$raw.payload.pull_request.title", "$raw.payload.issue.title" } },
				"number":       bson.M{ "$ifNull": []string{ "$raw.payload.pull_request.number", "$raw.payload.issue.number" } },
				"url":          bson.M{ "$ifNull": []string{ "$raw.payload.pull_request.html_url", "$raw.payload.issue.html_url" } },
				"merged":       "$raw.payload.pull_request.merged",
			},
		},
		{
			"$group": bson.M{
				"_id":          bson.M{ "r": "$repository", "n": "$number" },
				"title":        bson.M{ "$last": "$title"},
				"actions":      bson.M{ "$push": "$action" },
				"mergedStates": bson.M{ "$push": "$merged" },
				"actors":       bson.M{ "$addToSet": "$actor" },
				"url":          bson.M{ "$first": "$url" },
			},
		},
		{
			"$group": bson.M{
				"_id": "$_id.r",
				"project": bson.M{"$first": "$_id.r"},
				"prs": bson.M{
					"$addToSet": bson.M{
						"number":  "$_id.n",
						"url":     "$url",
						"title":   "$title",
						"actions": "$actions",
						"actors":  "$actors",
						"mergedStates": "$mergedStates",
					},
				},
			},
		},
		{ "$sort": bson.M{"project": 1} },
	})
	data := []interface{}{}
	err := pipe.All(&data)
	return data, err
}

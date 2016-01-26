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

	forks, err := fetchForks(session, date)
	if err != nil {
		return nil, err
	}

	stars, err := fetchStars(session, date)
	if err != nil {
		return nil, err
	}

	issuesUpdated, err := fetchIssuesUpdated(session, date)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"date":          date.Format("2006-01-02"),
		"newProjects":   newProjects,
		"forks":         forks,
		"stars":         stars,
		"issuesUpdated": issuesUpdated,
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

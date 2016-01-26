## Projects created

```js
db.events.aggregate([
  {
    "$match": {
      "type":                 "CreateEvent",
      "raw.payload.ref_type": "repository",
    }
  },
  {
    "$project": {
      "_id":     false,
      "user":    "$actor",
      "project": { $concat: [ "$project.owner", "/", "$project.name"] },
    }
  },
  { $sort: { "project": 1 } },
]);
```

## Projects forked

```js
db.events.aggregate([
  {
    $match: {
      type: 'ForkEvent',
    }
  },
  {
    $project: {
      _id:           false,
      user:          "$actor",
      repository:    { $concat: [ "$project.owner", "/", "$project.name"] },
      newRepository: "$raw.payload.forkee.full_name",
    }
  },
  { $sort: { '_id': 1 } },
]);
```

## Projects starred

```js
db.events.aggregate([
  {
    $match: {
      type: 'WatchEvent'
    }
  },
  {
    $project: {
      actor:      true,
      repository: { $concat: [ "$project.owner", "/", "$project.name"] },
    }
  },
  {
    $group: {
      _id:        '$repository',
      repository: { "$first": "$repository" },
      stargazers: { $addToSet: '$actor' },
    }
  },
  { $sort: { '_id': 1 } },
]);
```

## Added as a collaborator

```js
db.events.aggregate([
  {
    $match: {
      type: 'MemberEvent',
      'raw.payload.action': 'added'
    }
  },
  {
    $project: {
      type:         true,
      action:       true,
      created_at:   true,
      repository:   { $concat: [ "$project.owner", "/", "$project.name"] },
      yearMonthDay: { $dateToString: { format: "%Y-%m-%d", date: "$created_at" } },
      collaborator: '$raw.payload.member.login',
    }
  },
  {
    $group: {
      _id:           { d: '$yearMonthDay', r: '$repository' },
      total:         { $sum: 1 },
      collaborators: { $addToSet: '$collaborator' },
    }
  },
  { $sort: { '_id.d': -1, '_id.r': 1, total: -1 } },
]);
```

## Branches updated

```js
db.events.aggregate([
  { $match: { type: 'PushEvent' } },
  {
    $project: {
      type:         true,
      repository:   { $concat: [ "$project.owner", "/", "$project.name"] },
      yearMonthDay: { $dateToString: { format: "%Y-%m-%d", date: "$created_at" } },
      branch:       { $substr: [ "$raw.payload.ref", 11, -1 ] },
      before:       '$raw.payload.before',
      head:         '$raw.payload.head',
    }
  },
  {
    $group: {
      _id:    { d: '$yearMonthDay', r: '$repository', b: '$branch' },
      before: { $first: '$before' },
      head:   { $last: '$head' },
    }
  },
  { $sort: { '_id.d': -1, '_id.r': 1, '_id.b': 1 } }
]);
```

## New branches

```js
db.events.aggregate([
  {
    "$match": {
      "type":                 "CreateEvent",
      "raw.payload.ref_type": "branch",
    }
  },
  {
    "$project": {
      "_id":          false,
      "yearMonthDay": { $dateToString: { format: "%Y-%m-%d", date: "$created_at" } },
      "actor":        true,
      "repository":   { $concat: [ "$project.owner", "/", "$project.name"] },
      "branch":       "$raw.payload.ref",
    }
  },
  { $sort: { "yearMonthDay": -1 } },
]);
```

## Deleted branches

```js
db.events.aggregate([
  {
    "$match": {
      "type":                 "DeleteEvent",
      "raw.payload.ref_type": "branch",
    }
  },
  {
    "$project": {
      "_id":          false,
      "yearMonthDay": { $dateToString: { format: "%Y-%m-%d", date: "$created_at" } },
      "actor":        true,
      "repository":   { $concat: [ "$project.owner", "/", "$project.name"] },
      "branch":       "$raw.payload.ref",
    }
  },
  { $sort: { "yearMonthDay": -1 } },
]);
```

## Issues updated

```js
db.events.aggregate([
  {
    "$match": {
      "type": { "$in": ["IssuesEvent", "IssueCommentEvent"] },
      "raw.payload.issue.pull_request": { $exists: false },
    }
  },
  {
    "$sort": { "created_at": -1 }
  },
  {
    "$project": {
      "repository":   { $concat: [ "$project.owner", "/", "$project.name"] },
      "action":       true,
      "actor":        true,
      "state":        "$raw.payload.issue.state",
      "title":        "$raw.payload.issue.title",
      "number":       "$raw.payload.issue.number",
      "url":          "$raw.payload.issue.html_url",
    }
  },
  {
    $group: {
      _id:       { "r": "$repository", "n": "$number" },
      url:       { "$first": "$url"},
      actors:    { "$addToSet": "$actor" },
      title:     { "$last": "$title"},
      lastState: { "$last": "$state" },
    }
  },
  {
    $group: {
      _id: "$_id.r",
      issuesUpdated: {
        "$addToSet": {
          number:    "$_id.n",
          url:       "$url",
          title:     "$title",
          lastState: "$lastState",
          actors:    "$actors",
        }
      }
    }
  },
])
```

## Issues with comments

```js
// TODO
```

## PRs updated

```js
db.events.aggregate([
  {
    "$match": { "type": "PullRequestEvent" }
  },
  {
    "$sort": { "created_at": -1 }
  },
  {
    "$project": {
      "_id":          false,
      "yearMonthDay": { $dateToString: { format: "%Y-%m-%d", date: "$created_at" } },
      "actor":        true,
      "action":       true,
      "repository":   { $concat: [ "$project.owner", "/", "$project.name"] },
      "title":        "$raw.payload.pull_request.title",
      "number":       "$raw.payload.pull_request.number",
      "merged":       "$raw.payload.pull_request.merged",
    }
  },
  {
    $group: {
      _id:             { "d": "$yearMonthDay", "r": "$repository", "n": "$number" },
      title:           { "$last": "$title"},
      lastAction:      { "$last": "$action" },
      lastMergedState: { "$last": "$merged" },
    }
  },
  { $sort: { "_id.d": -1, "_id.r": 1, "_id.u": 1 } },
]);
```

## PRs with comments

```js
// TODO
```

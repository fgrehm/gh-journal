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
      "_id":          false,
      "yearMonthDay": { $dateToString: { format: "%Y-%m-%d", date: "$created_at" } },
      "actor":        true,
      "repository":   { $concat: [ "$project.owner", "/", "$project.name"] },
    }
  },
  { $sort: { "yearMonthDay": -1 } },
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
      yearMonthDay: { $dateToString: { format: "%Y-%m-%d", date: "$created_at" } },
      repository:   { $concat: [ "$project.owner", "/", "$project.name"] },
      forkee:       "$raw.payload.forkee.full_name",
    }
  },
  {
    $group: {
      _id:     { d: '$yearMonthDay', r: '$repository' },
      total:   { $sum: 1 },
      forkees: { $addToSet: '$forkee' },
    }
  },
  { $sort: { '_id.d': -1, '_id.r': 1, total: -1 } },
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
      actor:        true,
      repository:   { $concat: [ "$project.owner", "/", "$project.name"] },
      yearMonthDay: { $dateToString: { format: "%Y-%m-%d", date: "$created_at" } },
    }
  },
  {
    $group: {
      _id:        { d: '$yearMonthDay', r: '$repository' },
      total:      { $sum: 1 },
      stargazers: { $addToSet: '$actor' },
    }
  },
  { $sort: { '_id.d': -1, '_id.r': 1, total: -1 } },
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

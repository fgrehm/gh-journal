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
      actor:        true,
      repository:   { $concat: [ "$project.owner", "/", "$project.name"] },
      raw:          true,
    }
  },
  {
    $group: {
      _id:     { d: '$yearMonthDay', r: '$repository' },
      total:   { $sum: 1 },
      forkees: { $addToSet: '$raw.payload.forkee.full_name' },
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
  { $sort: { created_at: 1 } },
  {
    $project: {
      type:         true,
      action:       true,
      created_at:   true,
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

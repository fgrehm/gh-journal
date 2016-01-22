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
      actor:        true,
      repository:   { $concat: [ "$project.owner", "/", "$project.name"] },
      yearMonthDay: { $dateToString: { format: "%Y-%m-%d", date: "$created_at" } },
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
      type: 'WatchEvent',
      created_at: {
        $gte: ISODate("2016-01-20T00:00:00Z"),
        $lt:  ISODate("2016-01-21T00:00:00Z")
      }
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
      _id:          { d: '$yearMonthDay', r: '$repository' },
      total:        { $sum: 1 },
      stargazers:   { $push: '$actor' },
    }
  },
  { $sort: { '_id.d': -1, '_id.r': 1, total: -1 } },
]);
```

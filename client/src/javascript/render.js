var moment = require('moment');

var gitHubLink = function (attribute) {
  return {
    href: function () {
      return 'https://github.com/' + this[attribute];
    },
  }
};

var description = {
  html: function () {
    return this.description ?
      this.description :
      '<em>No description set</em>';
  }
};

var userLinks = function (attribute) {
  return {
    html: function () {
      return this[attribute].map(function (user) {
        return '<a href="https://github.com/' + user + '" target="_blank">' + user + '</a>';
      }).join(', ');
    }
  };
};

var directives = {
  date: {
    text: function (params) {
      return moment(this.date).format(params.value)
    }
  },

  prevDate: {
    href: function () {
      return '/#' + moment(this.date).subtract(1, 'days').format('YYYY-MM-DD');
    },
  },

  nextDate: {
    href: function () {
      return '/#' + moment(this.date).add(1, 'days').format('YYYY-MM-DD');
    },
  },

  noNewProjectsMsg: {
    text: function (params) {
      if (this.newProjects.length == 0)
        return 'Nothing to see here today =/';
    }
  },

  newProjects: {
    project: gitHubLink('project'),
    user: gitHubLink('user'),
    description: description,
  },

  noForksMsg: {
    text: function (params) {
      if (this.forks.length == 0)
        return 'Zero forks?!?';
    }
  },

  forks: {
    user: gitHubLink('user'),
    repository: gitHubLink('repository'),
    newRepository: gitHubLink('newRepository'),
    description: description,
  },

  noStarsMsg: {
    text: function (params) {
      if (this.stars.length == 0)
        return 'People are not feeling like sharing some love today :broken_heart:';
    }
  },

  stars: {
    project: gitHubLink('project'),
    stargazerLinks: userLinks('stargazers'),
  },

  noIssuesUpdatedMsg: {
    text: function (params) {
      if (this.issuesUpdated.length == 0)
        return '... radio silence ...';
    }
  },

  issuesUpdated: {
    project: gitHubLink('project'),
    issues: {
      actorLinks: userLinks('actors'),
      author: gitHubLink('author'),
      lastState: {
        class: function () { return 'badge ' + this.lastState; },
      },
      number: {
        href: function () { return this.url; },
        text: function () { return '#' + this.number; },
      },
      title: {
        href: function () { return this.url; },
        text: function () { return this.title; },
      },
    },
  },

  noPRsUpdatedMsg: {
    text: function (params) {
      if (this.issuesUpdated.length == 0)
        return 'Looks like people are quiet today';
    }
  },

  prsUpdated: {
    project: gitHubLink('project'),
    prs: {
      actorLinks: userLinks('actors'),
      author: gitHubLink('author'),
      isNew: {
        class: function () {
          for (var i = 0; i < this.actions.length; i++)
            if (this.actions[i] == 'created')
              return 'badge new';
        },
        text: function () {
          console.log(this.actions);
          for (var i = 0; i < this.actions.length; i++)
            if (this.actions[i] == 'created')
              return 'new';
        },
      },
      merged: {
        class: function () {
          if (this.mergedStates[this.mergedStates.length-1])
            return 'badge merged';
        },
        text: function () {
          if (this.mergedStates[this.mergedStates.length-1])
            return 'merged';
        },
      },
      number: {
        href: function () { return this.url; },
        text: function () { return '#' + this.number; },
      },
      title: {
        href: function () { return this.url; },
        text: function () { return this.title; },
      }
    },
  },
};

require('transparency');
require('chibijs');

module.exports = function renderReport(date) {
  $().get('report/' + date, function (response) {
    var data = JSON.parse(response);
    Transparency.render(document.getElementById('timeline'), data, directives);
  })
}

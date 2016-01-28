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
    stargazerLinks: {
      html: function () {
        return this.stargazers.map(function (stargazer) {
          return '<a href="https://github.com/' + stargazer + '" target="_blank">' + stargazer + '</a>';
        }).join(', ');
      }
    }
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
      number: {
        href: function () {
          return this.url;
        },
        text: function () {
          return '#' + this.number;
        },
      },
      title: {
        href: function () {
          return this.url;
        },
        text: function () {
          return this.title;
        },
      }
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

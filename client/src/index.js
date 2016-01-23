require('transparency');
var moment = require('moment');

var data = {
  date: '2015-01-23T21:00:43Z',
  newProjects: [
    { project: 'fgrehm/gh-journal',         user: 'fgrehm' },
    { project: 'sindresorhus/username-cli', user: 'sindresorhus' },
  ],
  forks: [
    { user: 'wycats', repository: 'nodesource/distributions', newRepository: 'wycats/distributions' },
  ],
  stars: [
    { project: 'lucasb-eyer/heatmap',     stargazers: ['tj'] },
    { project: 'DrkSephy/es6-cheatsheet', stargazers: ['tmattia'] },
  ]
};

var directives = {
  date: {
    text: function (params) {
      return moment(this.date).format(params.value)
    }
  },

  noNewProjectsMsg: {
    text: function (params) {
      if (this.newProjects.length == 0)
        return 'Nothing to see here today =/';
    }
  },

  newProjects: {
    project: {
      href: function () {
        return 'https://github.com/' + this.project;
      },
      text: function () {
        return this.project;
      }
    },
    user: {
      href: function () {
        return 'https://github.com/' + this.user;
      },
      text: function () {
        return this.user;
      }
    },
  },

  noForksMsg: {
    text: function (params) {
      if (this.forks.length == 0)
        return 'Zero forks?!?';
    }
  },

  forks: {
    user: {
      href: function () {
        return 'https://github.com/' + this.user;
      },
      text: function () {
        return this.user;
      }
    },
    repository: {
      href: function () {
        return 'https://github.com/' + this.repository;
      },
      text: function () {
        return this.repository;
      }
    },
    newRepository: {
      href: function () {
        return 'https://github.com/' + this.newRepository;
      },
      text: function () {
        return this.newRepository;
      }
    },
  },

  noStarsMsg: {
    text: function (params) {
      if (this.stars.length == 0)
        return 'People are not feeling like sharing some love today :broken_heart:';
    }
  },

  stars: {
    project: {
      href: function () {
        return 'https://github.com/' + this.project;
      },
      text: function () {
        return this.project;
      }
    },
    stargazerLinks: {
      html: function () {
        return this.stargazers.map(function (stargazer) {
          return '<a href="https://github.com/' + stargazer + '" target="_blank">' + stargazer + '</a>';
        }).join(', ');
      }
    }
  }
};

Transparency.render(document.getElementById('timeline'), data, directives);

var MODULE_NAME = 'app.controllers';
var controllers = angular.module(MODULE_NAME, []);

// controllers ==================================================================
var Home = require('./home.controller');
Home.$inject = ['$scope'];

var Edition = require('./edition.controller');
Edition.$inject = ['$scope', '$mdSidenav', 'projectsService', '$stateParams'];

controllers
  .controller('app.controllers.home', Home)
  .controller('app.controllers.edition', Edition)
  .service('projectsService', ['$q', function($q) {
    var projects = [{
        name: 'docker/docker',
        iconurl: 'https://avatars.githubusercontent.com/docker?&s=48',
        eventsCount: 200,
    }, {
        name: 'mitchellh/vagrant',
        iconurl: 'https://avatars.githubusercontent.com/mitchellh?&s=48',
        eventsCount: 198,
    }, {
        name: 'torvalds/linux',
        iconurl: 'https://avatars.githubusercontent.com/torvalds?&s=48',
        eventsCount: 30,
    }, {
        name: 'docker/compose',
        iconurl: 'https://avatars.githubusercontent.com/docker?&s=48',
        eventsCount: 5
    }, {
        name: 'fgrehm/vagrant-lxc',
        iconurl: 'https://avatars.githubusercontent.com/fgrehm?&s=48',
        eventsCount: 5
    }];

    // Promise-based API
    return {
        loadAll: function() {
            // Simulate async nature of real remote calls
            return $q.when(projects);
        }
    };
  }]);

module.exports = MODULE_NAME;

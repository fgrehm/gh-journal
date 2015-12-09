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
        name: 'linux/linux',
        iconurl: 'https://lh3.googleusercontent.com/-KGsfSssKoEU/AAAAAAAAAAI/AAAAAAAAAC4/j_0iL_6y3dE/s96-c-k-no/photo.jpg',
        eventsCount: 200,
    }, {
        name: 'mitchellh/vagrant',
        iconurl: 'https://yt3.ggpht.com/-cEjxni3_Jig/AAAAAAAAAAI/AAAAAAAAAAA/cMW2NEAUf-k/s88-c-k-no/photo.jpg',
        eventsCount: 198,
    }, {
        name: 'docker/docker',
        iconurl: 'https://goingforwardblog.files.wordpress.com/2013/01/swedish-chef.jpg',
        eventsCount: 30,
    }, {
        name: 'docker/compose',
        iconurl: 'https://lh5.googleusercontent.com/-c5rVqhf66e4/UVIKJ3fXLFI/AAAAAAAAACU/s-TU4ER7-Ro/w800-h800/kimmie.jpg',
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

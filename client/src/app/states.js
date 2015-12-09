module.exports = function ($stateProvider, $urlRouterProvider) {
  $stateProvider
    .state('home', {
      url: '/',
      controller: 'app.controllers.home',
      templateUrl: 'templates/home.html'
    })
    .state('edition', {
      url: '/edition/:editionDate',
      controller: 'app.controllers.edition',
      templateUrl: 'templates/edition.html'
    })
    // REFACTOR: This should use a nested state instead of reloading the whole page
    .state('projectInfo', {
      url: '/edition/:editionDate/:projectName',
      controller: 'app.controllers.edition',
      templateUrl: 'templates/edition.html'
    });

  $urlRouterProvider.otherwise('/');
};


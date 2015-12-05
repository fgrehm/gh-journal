module.exports = function ($stateProvider, $urlRouterProvider) {
  $stateProvider
    .state('home', {
      url: '/',
      controller: 'app.controllers.home',
      templateUrl: 'templates/home.html'
    })
    .state('edition', {
      url: '/edition',
      controller: 'app.controllers.edition',
      templateUrl: 'templates/edition.html'
    });

  $urlRouterProvider.otherwise('/');
};


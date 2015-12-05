module.exports = function ($stateProvider, $urlRouterProvider) {
  $stateProvider
    .state('home', {
      url: '/',
      controller: 'app.home',
      templateUrl: 'templates/home.html'
    });

  $urlRouterProvider.otherwise('/');
};


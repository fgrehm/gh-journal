require('angular');
require('angular-ui-router');
require('angular-material');

var app = angular.module('gh-journal', [
  require('./app/controllers'),
  'ngMaterial',
  'ui.router'
]);

app.config(['$stateProvider', '$urlRouterProvider', require('./app/states')]);

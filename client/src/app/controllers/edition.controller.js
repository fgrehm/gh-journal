function EditionCtrl($scope, $mdSidenav, muppetService, $timeout, $log) {
  console.log('foo');
  var allMuppets = [];

  $scope.selected = null;
  $scope.muppets = allMuppets;
  $scope.selectMuppet = selectMuppet;
  $scope.toggleSidenav = toggleSidenav;

  loadMuppets();

  //*******************
  // Internal Methods
  //*******************
  function loadMuppets() {
    muppetService.loadAll()
      .then(function(muppets){
        allMuppets = muppets;
        $scope.muppets = [].concat(muppets);
        $scope.selected = $scope.muppets[0];
      })
  }

  function toggleSidenav(name) {
    $mdSidenav(name).toggle();
  }

  function selectMuppet(muppet) {
    $scope.selected = angular.isNumber(muppet) ? $scope.muppets[muppet] : muppet;
    $scope.toggleSidenav('left');
  }
};

module.exports = EditionCtrl;

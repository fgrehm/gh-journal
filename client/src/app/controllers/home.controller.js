function HomeCtrl($scope) {
  $scope.editions = [];
  for (var i = 1; i <= 31; i++) {
    $scope.editions.push({
      date: '2015-01-' + i,
        eventsCount: (Math.random() * 1000).toFixed(),
        projectsUpdated: (Math.random() * 50).toFixed(),
        starredProjects: (Math.random() * 20).toFixed()
    });
  }
};

module.exports = HomeCtrl;

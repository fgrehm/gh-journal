function EditionCtrl($scope, $mdSidenav, projectsService, $stateParams) {
  var allProjects = [];

  $scope.selected = null;
  $scope.projects = allProjects;
  $scope.selectProject = selectProject;
  $scope.toggleSidenav = toggleSidenav;
  $scope.editionDate = $stateParams.editionDate;

  loadProjects().then(function() {
    if (!$stateParams.projectName) return;

    var projectName = unescape($stateParams.projectName);
    allProjects.forEach(function select(project) {
      if (project.name == projectName) {
        selectProject(project);
      }
    });
  });

  //*******************
  // Internal Methods
  //*******************
  function loadProjects() {
    return projectsService.loadAll()
      .then(function(projects){
        allProjects = projects;
        $scope.projects = [].concat(projects);
        // $scope.selected = $scope.projects[0];
      })
  }

  function toggleSidenav(name) {
    $mdSidenav(name).toggle();
  }

  function selectProject(project) {
    $scope.selected = angular.isNumber(project) ? $scope.projects[project] : project;
    $scope.toggleSidenav('left');
  }
};

module.exports = EditionCtrl;

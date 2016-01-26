var moment = require('moment');
var renderReport = require('./javascript/render');

var Mousetrap = require('mousetrap');

var Rlite = require('rlite-router');
var r = Rlite();
// Default route
r.add('', function () {
  var date = moment();

  Mousetrap.unbind(['left', 'right']);
  Mousetrap.bind('left',  function() { location.hash = "#" + date.subtract(1, 'days').format('YYYY-MM-DD') });
  Mousetrap.bind('right', function() { location.hash = "#" + date.add(1, 'days').format('YYYY-MM-DD') });

  renderReport(date.format('YYYY-MM-DD'));
});
r.add(':date', function (r) {
  var date = moment(r.params.date);

  Mousetrap.unbind(['left', 'right']);
  Mousetrap.bind('left',  function() { location.hash = "#" + date.subtract(1, 'days').format('YYYY-MM-DD') });
  Mousetrap.bind('right', function() { location.hash = "#" + date.add(1, 'days').format('YYYY-MM-DD') });

  renderReport(r.params.date);
});

// Hash-based routing
function processHash() {
  var hash = location.hash || '#';
  r.run(hash.slice(1));
}

window.addEventListener('hashchange', processHash);
processHash();

const proxy = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(proxy('/api', { target: 'http://milmovelocal:8080/' }));
  app.use(proxy('/internal', { target: 'http://milmovelocal:8080/' }));
  app.use(proxy('/storage', { target: 'http://milmovelocal:8080/' }));
  app.use(proxy('/devlocal-auth', { target: 'http://milmovelocal:8080/' }));
  app.use(proxy('/auth/**', { target: 'http://milmovelocal:8080/' }));
  app.use(proxy('/logout', { target: 'http://milmovelocal:8080/' }));
  app.use(proxy('/downloads', { target: 'http://milmovelocal:8080/' }));
};

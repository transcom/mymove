const proxy = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(proxy('/api', { target: 'http://localhost:8080/' }));
  app.use(proxy('/internal', { target: 'http://localhost:8080/' }));
  app.use(proxy('/storage', { target: 'http://localhost:8080/' }));
  app.use(proxy('/devlocal-auth', { target: 'http://localhost:8080/' }));
  app.use(proxy('/logout', { target: 'http://localhost:8080/' }));
  app.use(proxy('/downloads', { target: 'http://localhost:8080/' }));
};

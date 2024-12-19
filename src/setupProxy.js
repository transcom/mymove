const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = (app) => {
  app.use('/api', createProxyMiddleware({ target: 'http://milmovelocal:8080/api' }));
  app.use('/internal', createProxyMiddleware({ target: 'http://milmovelocal:8080/internal' }));
  app.use('/admin', createProxyMiddleware({ target: 'http://milmovelocal:8080/admin' }));
  app.use('/ghc', createProxyMiddleware({ target: 'http://milmovelocal:8080/ghc' }));
  app.use('/prime', createProxyMiddleware({ target: 'http://milmovelocal:8080/prime' }));
  app.use('/pptas', createProxyMiddleware({ target: 'http://milmovelocal:8080/pptas' }));
  app.use('/support', createProxyMiddleware({ target: 'http://milmovelocal:8080/support' }));
  app.use('/testharness', createProxyMiddleware({ target: 'http://milmovelocal:8080/testharness' }));
  app.use('/storage', createProxyMiddleware({ target: 'http://milmovelocal:8080/storage' }));
  app.use('/devlocal-auth', createProxyMiddleware({ target: 'http://milmovelocal:8080/devlocal-auth' }));
  app.use('/auth/**', createProxyMiddleware({ target: 'http://milmovelocal:8080/auth/**' }));
  app.use('/logout', createProxyMiddleware({ target: 'http://milmovelocal:8080/logout' }));
  app.use('/downloads', createProxyMiddleware({ target: 'http://milmovelocal:8080/downloads' }));
  app.use('/debug/**', createProxyMiddleware({ target: 'http://milmovelocal:8080/debug/**' }));
  app.use('/client/**', createProxyMiddleware({ target: 'http://milmovelocal:8080/client/**' }));
};

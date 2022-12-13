const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = (app) => {
  app.use(createProxyMiddleware('/api', { target: 'http://milmovelocal:8080/' }));
  app.use(createProxyMiddleware('/internal', { target: 'http://milmovelocal:8080/' }));
  app.use(createProxyMiddleware('/admin', { target: 'http://milmovelocal:8080/' }));
  app.use(createProxyMiddleware('/ghc', { target: 'http://milmovelocal:8080/' }));
  app.use(createProxyMiddleware('/prime', { target: 'http://milmovelocal:8080/' }));
  app.use(createProxyMiddleware('/support', { target: 'http://milmovelocal:8080/' }));
  app.use(createProxyMiddleware('/testharness', { target: 'http://milmovelocal:8080/' }));
  app.use(createProxyMiddleware('/storage', { target: 'http://milmovelocal:8080/' }));
  app.use(createProxyMiddleware('/devlocal-auth', { target: 'http://milmovelocal:8080/' }));
  app.use(createProxyMiddleware('/auth/**', { target: 'http://milmovelocal:8080/' }));
  app.use(createProxyMiddleware('/logout', { target: 'http://milmovelocal:8080/' }));
  app.use(createProxyMiddleware('/downloads', { target: 'http://milmovelocal:8080/' }));
  app.use(createProxyMiddleware('/debug/**', { target: 'http://milmovelocal:8080/' }));
};

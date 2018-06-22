// this is borrowed from https://github.com/ReactTraining/react-router/blob/e6f9017c947b3ae49affa24cc320d0a86f765b55/packages/react-router/modules/generatePath.js
// which has not been released yet
import pathToRegexp from 'path-to-regexp';

const patternCache = {};
const cacheLimit = 10000;
let cacheCount = 0;

const compileGenerator = pattern => {
  /* eslint-disable security/detect-object-injection */
  const cacheKey = pattern;
  const cache = patternCache[cacheKey] || (patternCache[cacheKey] = {});

  if (cache[pattern]) return cache[pattern];

  const compiledGenerator = pathToRegexp.compile(pattern);

  if (cacheCount < cacheLimit) {
    cache[pattern] = compiledGenerator;
    cacheCount++;
  }
  /* eslint-enable security/detect-object-injection */

  return compiledGenerator;
};

/**
 * Public API for generating a URL pathname from a pattern and parameters.
 */
const generatePath = (pattern = '/', params = {}) => {
  if (pattern === '/') {
    return pattern;
  }
  const generator = compileGenerator(pattern);
  return generator(params);
};

export default generatePath;

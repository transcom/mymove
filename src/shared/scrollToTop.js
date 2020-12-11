import { isTest } from './constants.js';

const scrollToTop = function () {
  if (!isTest) window.scrollTo(0, 0);
};

export default scrollToTop;

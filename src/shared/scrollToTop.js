import { isTest } from './constants.js';
export default function() {
  if (!isTest) window.scrollTo(0, 0);
}

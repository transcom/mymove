/* eslint-disable import/prefer-default-export */
import { isNullUndefinedOrWhitespace } from 'shared/utils';

/**
 * Compare strings. Null, undefined, and blanks are after other values.
 * @returns -1, 0, 1
 */
export function nullSafeStringCompare(a, b) {
  const A_BEFORE = -1;
  const A_AFTER = 1;
  const SAME = 0;

  if (isNullUndefinedOrWhitespace(a) && isNullUndefinedOrWhitespace(b)) {
    return SAME;
  }
  if (isNullUndefinedOrWhitespace(a)) {
    return A_AFTER;
  }
  if (isNullUndefinedOrWhitespace(b)) {
    return A_BEFORE;
  }
  return a.localeCompare(b);
}

/* eslint-disable import/prefer-default-export */
import * as stringUtils from 'utils/string';

describe('utils string', () => {
  describe('nullSafeComparison', () => {
    const A_BEFORE = -1;
    const A_AFTER = 1;
    const SAME = 0;
    it('same value', () => {
      const res = stringUtils.nullSafeStringCompare('1', '1');
      expect(res).toEqual(SAME);
    });
    it('greater than', () => {
      const res = stringUtils.nullSafeStringCompare('2', '1');
      expect(res).toEqual(A_AFTER);
    });
    it('less than', () => {
      const res = stringUtils.nullSafeStringCompare('1', '2');
      expect(res).toEqual(A_BEFORE);
    });
    it('both null', () => {
      const res = stringUtils.nullSafeStringCompare(null, null);
      expect(res).toEqual(SAME);
    });
    it('null and value', () => {
      const res = stringUtils.nullSafeStringCompare(null, '1');
      expect(res).toEqual(A_AFTER);
    });
    it('value and null', () => {
      const res = stringUtils.nullSafeStringCompare('1', null);
      expect(res).toEqual(A_BEFORE);
    });
    it('both undefined', () => {
      let udefA;
      let udefB;
      const res = stringUtils.nullSafeStringCompare(udefA, udefB);
      expect(res).toEqual(SAME);
    });
    it('undefined and null', () => {
      let udefA;
      const res = stringUtils.nullSafeStringCompare(udefA, null);
      expect(res).toEqual(SAME);
    });
    it('null and undefined', () => {
      let udefB;
      const res = stringUtils.nullSafeStringCompare(null, udefB);
      expect(res).toEqual(SAME);
    });
    it('undefined and value', () => {
      let udefA;
      const res = stringUtils.nullSafeStringCompare(udefA, '2');
      expect(res).toEqual(A_AFTER);
    });
    it('value and undefined', () => {
      let udefB;
      const res = stringUtils.nullSafeStringCompare('1', udefB);
      expect(res).toEqual(A_BEFORE);
    });
  });
});

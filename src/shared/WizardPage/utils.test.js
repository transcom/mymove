import { getNextPagePath, getPreviousPagePath } from './utils';

describe('given there are three items in the pageList', () => {
  const sampleList = ['page1', 'page2', 'page3'];
  // beforeEach(() => {});
  describe('when the current page is the first page', () => {
    const currentPage = sampleList[0];
    it('then the next page is page 2', () => {
      expect(getNextPagePath(sampleList, currentPage)).toBe('page2');
    });
    it('then the previous page is undefined', () => {
      expect(getPreviousPagePath(sampleList, currentPage)).toBeUndefined();
    });
  });
  describe('when the current page is the middle page', () => {
    const currentPage = sampleList[1];
    it('then the next page is page 3', () => {
      expect(getNextPagePath(sampleList, currentPage)).toBe('page3');
    });
    it('then the previous page is page1', () => {
      expect(getPreviousPagePath(sampleList, currentPage)).toBe('page1');
    });
  });
  describe('when the current page is the last page', () => {
    const currentPage = sampleList[2];
    it('then the next page is undefined', () => {
      expect(getNextPagePath(sampleList, currentPage)).toBeUndefined();
    });
    it('then the prev page is page 2', () => {
      expect(getPreviousPagePath(sampleList, currentPage)).toBe('page2');
    });
  });
  describe('when the current page is invalid', () => {
    const currentPage = 'rando';
    it('then the next page is undefined', () => {
      expect(getNextPagePath(sampleList, currentPage)).toBeUndefined();
    });
    it('then the prev page is page 2', () => {
      expect(getPreviousPagePath(sampleList, currentPage)).toBeUndefined();
    });
  });
});

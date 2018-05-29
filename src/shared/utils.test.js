import * as utils from './utils';

describe('utils', () => {
  describe('upsert', () => {
    const item = { id: 'foo', name: 'something' };
    describe('when upserting a new item to an array', () => {
      const arr = [{ id: 'bar', name: 'foo' }, { id: 'baz', name: 'baz' }];
      utils.upsert(arr, item);
      it('should be appended to the array', () => {
        expect(arr).toEqual([
          { id: 'bar', name: 'foo' },
          { id: 'baz', name: 'baz' },
          item,
        ]);
      });
    });
    describe('when upserting an update to an array', () => {
      const arr = [{ id: 'foo', name: 'foo' }, { id: 'baz', name: 'baz' }];
      utils.upsert(arr, item);
      it('should be appended to the array', () => {
        expect(arr).toEqual([
          { id: 'foo', name: 'something' },
          { id: 'baz', name: 'baz' },
        ]);
      });
    });
  });
});

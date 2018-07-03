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

  describe('fetch Active', () => {
    describe('when there are no foos', () => {
      const foos = null;
      let res = utils.fetchActive(foos);
      it('should return null', () => {
        expect(res).toBeNull();
      });
    });
    describe('when there are some active and some inactive foos', () => {
      const foos = [
        { status: 'CANCELED', id: 'foo0' },
        { status: 'DRAFT', id: 'foo1' },
        { status: 'SUBMITTED', id: 'foo2' },
      ];
      let res = utils.fetchActive(foos);
      it('should return the first active foo', () => {
        expect(res.id).toEqual('foo1');
      });
    });
    describe('when there are only inactive foos', () => {
      const foos = [
        { status: 'CANCELED', id: 'foo' },
        { status: 'COMPLETED', id: 'foo' },
      ];
      let res = utils.fetchActive(foos);
      it('should return null', () => {
        expect(res).toBeNull();
      });
    });
  });
});

import {
  createHeader,
  setTableQueueFilterSessionStorageValue,
  getTableQueueFilterSessionStorageValue,
  setTableQueuePageSizeSessionStorageValue,
  getTableQueuePageSizeSessionStorageValue,
  setTableQueueSortParamSessionStorageValue,
  getTableQueueSortParamSessionStorageValue,
} from './utils';

const localStorageMock = (() => {
  let store = {};

  return {
    getItem(key) {
      return store[key] || null;
    },
    setItem(key, value) {
      store[key] = value.toString();
    },
    removeItem(key) {
      delete store[key];
    },
    clear() {
      store = {};
    },
  };
})();

Object.defineProperty(window, 'sessionStorage', {
  value: localStorageMock,
});

beforeEach(() => {
  window.sessionStorage.clear();
});

describe('createHeader()', () => {
  it('returns expected object with params', () => {
    const headerObject = createHeader('HeaderString', 'AccessorString');
    expect(headerObject).toEqual({ Header: 'HeaderString', accessor: 'AccessorString' });
  });

  it('returns expected object with params + options', () => {
    const headerObject = createHeader('HeaderString', 'AccessorString', { customProp: 'CustomProp' });
    expect(headerObject).toEqual({ Header: 'HeaderString', accessor: 'AccessorString', customProp: 'CustomProp' });
  });
});

describe('setTableQueueFilterSessionStorageValue, getTableQueueFilterSessionStorageValue', () => {
  const filters = [
    { id: 'lastName', value: 'Spacemen' },
    { id: 'dodID', value: '7232607949' },
  ];
  const filters2 = [
    { id: 'lastName', value: 'Ziggy StarDust' },
    { id: 'dodID', value: '7232607949' },
  ];
  const filters3 = [{ id: 'dodID', value: '7232607949' }];
  const sessionStorageKey = 'test';
  expect(getTableQueueFilterSessionStorageValue(sessionStorageKey)).toEqual([]);
  setTableQueueFilterSessionStorageValue(sessionStorageKey, filters);
  expect(getTableQueueFilterSessionStorageValue(sessionStorageKey)).toEqual(filters);
  expect(getTableQueueFilterSessionStorageValue(sessionStorageKey)).not.toEqual(filters2);
  setTableQueueFilterSessionStorageValue(sessionStorageKey, filters3);
  expect(getTableQueueFilterSessionStorageValue(sessionStorageKey)).toEqual(filters3);
  expect(getTableQueueFilterSessionStorageValue(sessionStorageKey)).not.toEqual(filters);
});

describe('setTableQueuePageSizeSessionStorageValue, getTableQueuePageSizeSessionStorageValue', () => {
  const sessionStorageKey = 'test';
  expect(getTableQueuePageSizeSessionStorageValue(sessionStorageKey)).toEqual(null);
  setTableQueueFilterSessionStorageValue(sessionStorageKey, 1);
  expect(getTableQueueFilterSessionStorageValue(sessionStorageKey)).toEqual(1);
  setTableQueueFilterSessionStorageValue(sessionStorageKey, 2);
  expect(getTableQueueFilterSessionStorageValue(sessionStorageKey)).not.toEqual(1);
});

describe('setTableQueuePageSizeSessionStorageValue, getTableQueuePageSizeSessionStorageValue', () => {
  const sessionStorageKey = 'test';
  expect(getTableQueuePageSizeSessionStorageValue(sessionStorageKey)).toEqual(null);
  setTableQueuePageSizeSessionStorageValue(sessionStorageKey, 1);
  expect(getTableQueuePageSizeSessionStorageValue(sessionStorageKey)).toEqual(1);
  setTableQueuePageSizeSessionStorageValue(sessionStorageKey, 2);
  expect(getTableQueuePageSizeSessionStorageValue(sessionStorageKey)).toEqual(2);
  expect(getTableQueuePageSizeSessionStorageValue(sessionStorageKey)).not.toEqual(1);
});

describe('setTableQueueSortParamSessionStorageValue, getTableQueueSortParamSessionStorageValue', () => {
  const sessionStorageKey = 'test';
  const sortParam = [{ id: 'dodID', desc: false }];
  const sortParam2 = [{ id: 'lastName', desc: false }];
  expect(getTableQueueSortParamSessionStorageValue(sessionStorageKey)).toEqual(null);
  setTableQueueSortParamSessionStorageValue(sessionStorageKey, sortParam);
  expect(getTableQueueSortParamSessionStorageValue(sessionStorageKey)).toEqual(sortParam);
  setTableQueueSortParamSessionStorageValue(sessionStorageKey, sortParam2);
  expect(getTableQueueSortParamSessionStorageValue(sessionStorageKey)).toEqual(sortParam2);
  expect(getTableQueueSortParamSessionStorageValue(sessionStorageKey)).not.toEqual(sortParam);
});

import { createHeader, textFilter } from './utils';

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

describe('textFilter()', () => {
  it('returns expected value with params', () => {
    const rows = { filter: (fn) => fn({ values: { id: 'value' } }) };
    const isFiltered = textFilter(rows, 'id', 'value');
    expect(isFiltered).toEqual(true);
  });

  it('returns false with no match', () => {
    const rows = { filter: (fn) => fn({ values: { id: 'valuee' } }) };
    const isFiltered = textFilter(rows, 'id', 'valueeeee');
    expect(isFiltered).toEqual(false);
  });
});

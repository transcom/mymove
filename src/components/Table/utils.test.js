import { createHeader } from './utils';

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

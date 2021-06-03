import formattedCustomerName from './formattedCustomerName';

describe('formattedCustomerName', () => {
  it('formats with middle name and suffix', () => {
    expect(formattedCustomerName('Doe', 'John', 'Jr.', 'Adams')).toEqual('Doe, John Adams, Jr.');
  });

  it('formats with middle name and no suffix', () => {
    expect(formattedCustomerName('Doe', 'John', undefined, 'Adams')).toEqual('Doe, John Adams');
  });

  it('formats with suffix and no middle name', () => {
    expect(formattedCustomerName('Doe', 'John', 'Jr.', undefined)).toEqual('Doe, John, Jr.');
  });

  it('formats with just last and first', () => {
    expect(formattedCustomerName('Doe', 'John')).toEqual('Doe, John');
  });
});

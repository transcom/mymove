import { emailSchema } from './validation';

describe('checkRequiredFields', () => {
  it('success: does nothing if all fields provided', () => {
    expect(
      emailSchema({
        email: 'test@example.com',
      }),
    ).toBeTruthy();
  });
  it('fail: throws an error if fields are missing', () => {
    function checkMissingFields() {
      emailSchema({ email: 'test@example.com' });
    }
    expect(checkMissingFields).toThrowError('Row does not contain all required fields.');
  });
});

import { stringifyName } from './serviceMember';

describe('serviceMember utils', () => {
  describe('stringifyName', () => {
    describe('when first and last name are not null', () => {
      it('returns last name and first name separated by a comma', () => {
        const first_name = 'Jane';
        const last_name = 'Smith';
        expect(stringifyName({ first_name, last_name })).toEqual('Smith, Jane');
      });
    });

    describe('when first name is null and last name is not null', () => {
      it('returns just the last name', () => {
        const first_name = null;
        const last_name = 'Smith';
        expect(stringifyName({ first_name, last_name })).toEqual('Smith');
      });
    });

    describe('when first name is not null and last name is null', () => {
      it('returns just the first name', () => {
        const first_name = 'Jane';
        const last_name = null;
        expect(stringifyName({ first_name, last_name })).toEqual('Jane');
      });
    });

    describe('when both first and last name are null', () => {
      it('returns an empty string', () => {
        const first_name = '';
        const last_name = '';
        expect(stringifyName({ first_name, last_name })).toEqual('');
      });
    });
  });
});

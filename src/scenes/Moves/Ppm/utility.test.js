import { intToOrdinal, getNextPage } from './utility';

describe('intToOrdinal', () => {
  it('returns the ordinal corresponding to an int', () => {
    expect(intToOrdinal(1)).toEqual('1st');
    expect(intToOrdinal(2)).toEqual('2nd');
    expect(intToOrdinal(3)).toEqual('3rd');
    expect(intToOrdinal(4)).toEqual('4th');
  });
});
describe('nextPage', () => {
  it('returns to the lastPage when user is visiting from pageToRevisit', () => {
    const lastPage = {
      pathname: 'moves/:moveid/ppm-review',
      search: '',
      hash: '',
      state: undefined,
    };
    const nextPage = getNextPage('moves/:moveid/next-page', lastPage, '/ppm-review');

    expect(nextPage).toEqual('moves/:moveid/ppm-review');
  });
  it('returns to the nextPage when user is not visiting from pageToRevisit', () => {
    const lastPage = {
      pathname: 'moves/:moveid/some-other-page',
      search: '',
      hash: '',
      state: undefined,
    };
    const nextPage = getNextPage('moves/:moveid/next-page', lastPage, '/ppm-review');

    expect(nextPage).toEqual('moves/:moveid/next-page');
  });
  it('returns to the nextPage when no lastpage', () => {
    const nextPage = getNextPage('moves/:moveid/next-page', null, '/ppm-review');

    expect(nextPage).toEqual('moves/:moveid/next-page');
  });
});

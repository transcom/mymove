import { getNextPage } from './utility';

describe('nextPage', () => {
  it('returns to the lastPage when user is visiting from pageToRevisit', () => {
    const lastPage = {
      pathname: 'moves/:moveid/ppm-payment-review',
      search: '',
      hash: '',
      state: undefined,
    };
    const nextPage = getNextPage('moves/:moveid/next-page', lastPage, '/ppm-payment-review');

    expect(nextPage).toEqual('moves/:moveid/ppm-payment-review');
  });
  it('returns to the nextPage when user is not visiting from pageToRevisit', () => {
    const lastPage = {
      pathname: 'moves/:moveid/some-other-page',
      search: '',
      hash: '',
      state: undefined,
    };
    const nextPage = getNextPage('moves/:moveid/next-page', lastPage, '/ppm-payment-review');

    expect(nextPage).toEqual('moves/:moveid/next-page');
  });
  it('returns to the nextPage when no lastpage', () => {
    const nextPage = getNextPage('moves/:moveid/next-page', null, '/ppm-payment-review');

    expect(nextPage).toEqual('moves/:moveid/next-page');
  });
});

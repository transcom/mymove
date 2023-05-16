import { getNextPage, calcNetWeight } from './utility';

import { MOVE_DOC_TYPE } from 'shared/constants';

describe('PPM Utility functions', () => {
  describe('getNextPage', () => {
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
  describe('calcNetWeight', () => {
    it('should return 0 when there are no weight ticket sets', () => {
      const documents = [{ move_document_type: MOVE_DOC_TYPE.EXPENSE }, { move_document_type: MOVE_DOC_TYPE.GBL }];
      expect(calcNetWeight(documents)).toBe(0);
    });
    it('should return the net weight for weight ticket sets', () => {
      const documents = [
        { move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET, full_weight: 2000, empty_weight: 1000 },
        { move_document_type: MOVE_DOC_TYPE.WEIGHT_TICKET_SET, full_weight: 3000, empty_weight: 2000 },
      ];
      expect(calcNetWeight(documents)).toBe(2000);
    });
  });
});

import { getPageList } from './getWorkflowRoutes';

describe('when getting the routes for the current workflow', () => {
  let profileIsComplete;
  describe('given a complete service member', () => {
    profileIsComplete = true;
    describe('given a PPM', () => {
      const state = {
        user: { hasCompleteProfile: profileIsComplete },
        submittedMoves: { currentMove: { selected_move_type: 'PPM' } },
      };
      const pages = getPageList(state);
      it('just returns move pages', () => {
        expect(pages).toEqual([
          '/moves/:moveId',
          '/moves/:moveId/ppm-start',
          '/moves/:moveId/ppm-size',
          '/moves/:moveId/ppm-incentive',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
    describe('given a complete service member with an HHG', () => {
      const state = {
        user: { hasCompleteProfile: profileIsComplete },
        submittedMoves: { currentMove: { selected_move_type: 'HHG' } },
      };
      const pages = getPageList(state);
      it('just returns move pages', () => {
        expect(pages).toEqual([
          '/moves/:moveId',
          '/moves/:moveId/schedule',
          '/moves/:moveId/address',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
    describe('given a complete service member with a COMBO', () => {
      const state = {
        user: { hasCompleteProfile: profileIsComplete },
        submittedMoves: { currentMove: { selected_move_type: 'COMBO' } },
      };
      const pages = getPageList(state);
      it('just returns move pages', () => {
        expect(pages).toEqual([
          '/moves/:moveId',
          '/moves/:moveId/schedule',
          '/moves/:moveId/address',
          '/moves/:moveId/ppm-transition',
          '/moves/:moveId/ppm-size',
          '/moves/:moveId/ppm-incentive',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
  });
  describe('given an incomplete service member', () => {
    profileIsComplete = false;
    describe('given a PPM', () => {
      const state = {
        user: { hasCompleteProfile: profileIsComplete },
        submittedMoves: { currentMove: { selected_move_type: 'PPM' } },
      };
      const pages = getPageList(state);
      it('just returns move pages', () => {
        expect(pages).toEqual([
          '/service-member/:id/create',
          '/service-member/:id/name',
          '/service-member/:id/contact-info',
          '/service-member/:id/duty-station',
          '/service-member/:id/residence-address',
          '/service-member/:id/backup-mailing-address',
          '/service-member/:id/backup-contacts',
          '/service-member/:id/transition',
          '/orders/:id/',
          '/orders/:id/upload',
          '/moves/:moveId',
          '/moves/:moveId/ppm-start',
          '/moves/:moveId/ppm-size',
          '/moves/:moveId/ppm-incentive',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
    describe('given a complete service member with an HHG', () => {
      const state = {
        user: { hasCompleteProfile: profileIsComplete },
        submittedMoves: { currentMove: { selected_move_type: 'HHG' } },
      };
      const pages = getPageList(state);
      it('just returns move pages', () => {
        expect(pages).toEqual([
          '/service-member/:id/create',
          '/service-member/:id/name',
          '/service-member/:id/contact-info',
          '/service-member/:id/duty-station',
          '/service-member/:id/residence-address',
          '/service-member/:id/backup-mailing-address',
          '/service-member/:id/backup-contacts',
          '/service-member/:id/transition',
          '/orders/:id/',
          '/orders/:id/upload',
          '/moves/:moveId',
          '/moves/:moveId/schedule',
          '/moves/:moveId/address',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
    describe('given a complete service member with a COMBO', () => {
      const state = {
        user: { hasCompleteProfile: profileIsComplete },
        submittedMoves: { currentMove: { selected_move_type: 'COMBO' } },
      };
      const pages = getPageList(state);
      it('just returns move pages', () => {
        expect(pages).toEqual([
          '/service-member/:id/create',
          '/service-member/:id/name',
          '/service-member/:id/contact-info',
          '/service-member/:id/duty-station',
          '/service-member/:id/residence-address',
          '/service-member/:id/backup-mailing-address',
          '/service-member/:id/backup-contacts',
          '/service-member/:id/transition',
          '/orders/:id/',
          '/orders/:id/upload',
          '/moves/:moveId',
          '/moves/:moveId/schedule',
          '/moves/:moveId/address',
          '/moves/:moveId/ppm-transition',
          '/moves/:moveId/ppm-size',
          '/moves/:moveId/ppm-incentive',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
  });
});

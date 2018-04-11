import { getPageList } from './getWorkflowRoutes';

describe('when getting the routes for the current workflow', () => {
  let profileIsComplete;
  describe('given a complete service member', () => {
    profileIsComplete = true;
    describe('given a PPM', () => {
      const props = {
        hasCompleteProfile: profileIsComplete,
        selectedMoveType: 'PPM',
        hasMove: true,
      };
      const pages = getPageList(props);
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
      const props = {
        hasCompleteProfile: profileIsComplete,
        selectedMoveType: 'HHG',
        hasMove: true,
      };
      const pages = getPageList(props);
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
      const props = {
        hasCompleteProfile: profileIsComplete,
        selectedMoveType: 'COMBO',
        hasMove: true,
      };
      const pages = getPageList(props);
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
    describe('given no move', () => {
      const props = {
        hasCompleteProfile: profileIsComplete,
        selectedMoveType: null,
        hasMove: false,
      };
      const pages = getPageList(props);
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
        ]);
      });
    });
    describe('given a PPM', () => {
      const props = {
        hasCompleteProfile: profileIsComplete,
        selectedMoveType: 'PPM',
        hasMove: true,
      };
      const pages = getPageList(props);
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
      const props = {
        hasCompleteProfile: profileIsComplete,
        selectedMoveType: 'HHG',
        hasMove: true,
      };
      const pages = getPageList(props);
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
      const props = {
        hasCompleteProfile: profileIsComplete,
        selectedMoveType: 'COMBO',
        hasMove: true,
      };
      const pages = getPageList(props);
      it('returns the service member and move pages', () => {
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

import { getPagesInFlow } from './getWorkflowRoutes';

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
      const pages = getPagesInFlow(props);
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
      const pages = getPagesInFlow(props);
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
      const pages = getPagesInFlow(props);
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
      const pages = getPagesInFlow(props);
      it('returns service member, order and move pages', () => {
        expect(pages).toEqual([
          '/service-member/:serviceMemberId/create',
          '/service-member/:serviceMemberId/name',
          '/service-member/:serviceMemberId/contact-info',
          '/service-member/:serviceMemberId/duty-station',
          '/service-member/:serviceMemberId/residence-address',
          '/service-member/:serviceMemberId/backup-mailing-address',
          '/service-member/:serviceMemberId/backup-contacts',
          '/service-member/:serviceMemberId/transition',
          '/orders/',
          '/orders/upload',
          '/orders/transition',
          '/moves/:moveId',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
    describe('given a PPM', () => {
      const props = {
        hasCompleteProfile: profileIsComplete,
        selectedMoveType: 'PPM',
        hasMove: true,
      };
      const pages = getPagesInFlow(props);
      it('returns service member, order and PPM-specific move pages', () => {
        expect(pages).toEqual([
          '/service-member/:serviceMemberId/create',
          '/service-member/:serviceMemberId/name',
          '/service-member/:serviceMemberId/contact-info',
          '/service-member/:serviceMemberId/duty-station',
          '/service-member/:serviceMemberId/residence-address',
          '/service-member/:serviceMemberId/backup-mailing-address',
          '/service-member/:serviceMemberId/backup-contacts',
          '/service-member/:serviceMemberId/transition',
          '/orders/',
          '/orders/upload',
          '/orders/transition',
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
      const pages = getPagesInFlow(props);
      it('returns service member, order and HHG-specific move pages', () => {
        expect(pages).toEqual([
          '/service-member/:serviceMemberId/create',
          '/service-member/:serviceMemberId/name',
          '/service-member/:serviceMemberId/contact-info',
          '/service-member/:serviceMemberId/duty-station',
          '/service-member/:serviceMemberId/residence-address',
          '/service-member/:serviceMemberId/backup-mailing-address',
          '/service-member/:serviceMemberId/backup-contacts',
          '/service-member/:serviceMemberId/transition',
          '/orders/',
          '/orders/upload',
          '/orders/transition',
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
      const pages = getPagesInFlow(props);
      it('returns service member, order and move pages', () => {
        expect(pages).toEqual([
          '/service-member/:serviceMemberId/create',
          '/service-member/:serviceMemberId/name',
          '/service-member/:serviceMemberId/contact-info',
          '/service-member/:serviceMemberId/duty-station',
          '/service-member/:serviceMemberId/residence-address',
          '/service-member/:serviceMemberId/backup-mailing-address',
          '/service-member/:serviceMemberId/backup-contacts',
          '/service-member/:serviceMemberId/transition',
          '/orders/',
          '/orders/upload',
          '/orders/transition',
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

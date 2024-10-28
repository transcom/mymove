import { getPagesInFlow, getNextIncompletePage } from './getWorkflowRoutes';

import { NULL_UUID, SHIPMENT_OPTIONS } from 'shared/constants';

const ppmContext = {
  flags: {
    hhgFlow: false,
    ghcFlow: false,
  },
};
const hhgContext = {
  flags: {
    hhgFlow: true,
  },
};
const ghcContext = {
  flags: {
    ghcFlow: true,
    hhgFlow: false,
  },
};

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('when getting the routes for the current workflow', () => {
  describe('given a complete service member', () => {
    describe('given a PPM', () => {
      const props = {
        move: { mtoShipments: [{ shipmentType: SHIPMENT_OPTIONS.PPM }] },
        context: ppmContext,
      };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns service member, order and move pages', () => {
        expect(pages).toEqual([
          '/service-member/validation-code',
          '/service-member/dod-info',
          '/service-member/name',
          '/service-member/contact-info',
          '/service-member/current-address',
          '/service-member/backup-address',
          '/service-member/backup-contact',
          '/orders/info/:orderId',
          '/orders/upload/:orderId',
          '/moves/:moveId/shipment-type',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
    describe('given a canceled PPM', () => {
      const props = {
        lastMoveIsCanceled: true,
        move: { mtoShipments: [{ shipmentType: SHIPMENT_OPTIONS.PPM }] },
        context: ppmContext,
      };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns profile review, the order and move pages', () => {
        expect(pages).toEqual([
          '/profile-review',
          '/orders/info/:orderId',
          '/orders/upload/:orderId',
          '/moves/:moveId/shipment-type',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
  });
  describe('given an incomplete service member', () => {
    describe('given no move', () => {
      const props = {
        context: ppmContext,
      };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns service member, order and move pages', () => {
        expect(pages).toEqual([
          '/service-member/validation-code',
          '/service-member/dod-info',
          '/service-member/name',
          '/service-member/contact-info',
          '/service-member/current-address',
          '/service-member/backup-address',
          '/service-member/backup-contact',
          '/orders/info/:orderId',
          '/orders/upload/:orderId',
          '/moves/:moveId/shipment-type',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
    describe('given an incomplete service member', () => {
      describe('given no move and behind GHC flag', () => {
        const props = {
          context: ghcContext,
        };
        const pages = getPagesInFlow(props);
        it('getPagesInFlow returns service member, order and move pages', () => {
          expect(pages).toEqual([
            '/service-member/validation-code',
            '/service-member/dod-info',
            '/service-member/name',
            '/service-member/contact-info',
            '/service-member/current-address',
            '/service-member/backup-address',
            '/service-member/backup-contact',
            '/',
            '/orders/info/:orderId',
            '/orders/upload/:orderId',
            '/moves/:moveId/shipment-type',
            '/moves/:moveId/review',
            '/moves/:moveId/agreement',
          ]);
        });
      });
    });
    describe('given a PPM', () => {
      const props = {
        move: { mtoShipments: [{ shipmentType: SHIPMENT_OPTIONS.PPM }] },
        context: ppmContext,
      };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns service member, order and PPM-specific move pages', () => {
        expect(pages).toEqual([
          '/service-member/validation-code',
          '/service-member/dod-info',
          '/service-member/name',
          '/service-member/contact-info',
          '/service-member/current-address',
          '/service-member/backup-address',
          '/service-member/backup-contact',
          '/orders/info/:orderId',
          '/orders/upload/:orderId',
          '/moves/:moveId/shipment-type',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
    describe('given hhgFlow flag is true', () => {
      const props = {
        context: hhgContext,
      };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns service member, order and select move type page', () => {
        expect(pages).toEqual([
          '/service-member/validation-code',
          '/service-member/dod-info',
          '/service-member/name',
          '/service-member/contact-info',
          '/service-member/current-address',
          '/service-member/backup-address',
          '/service-member/backup-contact',
          '/orders/info/:orderId',
          '/orders/upload/:orderId',
          '/moves/:moveId/shipment-type',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
    describe('given hhgFlow + selected shipment NTS is true', () => {
      const props = {
        context: hhgContext,
      };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns service member, order and select move type page', () => {
        expect(pages).toEqual([
          '/service-member/validation-code',
          '/service-member/dod-info',
          '/service-member/name',
          '/service-member/contact-info',
          '/service-member/current-address',
          '/service-member/backup-address',
          '/service-member/backup-contact',
          '/orders/info/:orderId',
          '/orders/upload/:orderId',
          '/moves/:moveId/shipment-type',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
    describe('given hhgFlow + selected shipment NTSR is true', () => {
      const props = {
        context: hhgContext,
      };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns service member, order and select move type page', () => {
        expect(pages).toEqual([
          '/service-member/validation-code',
          '/service-member/dod-info',
          '/service-member/name',
          '/service-member/contact-info',
          '/service-member/current-address',
          '/service-member/backup-address',
          '/service-member/backup-contact',
          '/orders/info/:orderId',
          '/orders/upload/:orderId',
          '/moves/:moveId/shipment-type',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
  });
});

describe('when getting the next incomplete page', () => {
  const serviceMember = { id: 'foo' };
  describe('when the profile is incomplete', () => {
    it('returns the first page of the user profile', () => {
      const result = getNextIncompletePage({
        serviceMember,
        context: ppmContext,
      });
      expect(result).toEqual('/service-member/validation-code');
    });
    describe('when dod info is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          serviceMember: {
            ...serviceMember,
            is_profile_complete: false,
            edipi: '1234567890',
            affiliation: 'Marines',
          },
          context: ppmContext,
        });
        expect(result).toEqual('/service-member/name');
      });
    });
    describe('when name is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          serviceMember: {
            ...serviceMember,
            is_profile_complete: false,
            edipi: '1234567890',
            affiliation: 'Marines',
            last_name: 'foo',
            first_name: 'foo',
          },
          context: ppmContext,
        });
        expect(result).toEqual('/service-member/contact-info');
      });
    });
    describe('when contact-info is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          move: { mtoShipments: [{ shipmentType: SHIPMENT_OPTIONS.PPM }] },
          serviceMember: {
            ...serviceMember,
            is_profile_complete: false,
            edipi: '1234567890',
            affiliation: 'Marines',
            last_name: 'foo',
            first_name: 'foo',
            email_is_preferred: true,
            telephone: '666-666-6666',
            personal_email: 'foo@bar.com',
            current_location: {
              id: NULL_UUID,
              name: '',
            },
          },
          context: ppmContext,
        });
        expect(result).toEqual('/service-member/current-address');
      });
    });
    describe('when residence address is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          serviceMember: {
            ...serviceMember,
            is_profile_complete: false,
            edipi: '1234567890',
            affiliation: 'Marines',
            last_name: 'foo',
            first_name: 'foo',
            phone_is_preferred: true,
            telephone: '666-666-6666',
            personal_email: 'foo@bar.com',
            current_location: {
              id: '5e30f356-e590-4372-b9c0-30c3fd1ff42d',
              name: 'Blue Grass Army Depot',
            },
            residential_address: {
              city: 'Atlanta',
              postalCode: '30030',
              state: 'GA',
              streetAddress1: 'xxx',
            },
          },
          context: ppmContext,
        });
        expect(result).toEqual('/service-member/backup-address');
      });
    });
    describe('when backup address is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          serviceMember: {
            ...serviceMember,
            is_profile_complete: false,
            edipi: '1234567890',
            affiliation: 'Marines',
            last_name: 'foo',
            first_name: 'foo',
            phone_is_preferred: true,
            telephone: '666-666-6666',
            personal_email: 'foo@bar.com',
            current_location: {
              id: '5e30f356-e590-4372-b9c0-30c3fd1ff42d',
              name: 'Blue Grass Army Depot',
            },
            residential_address: {
              city: 'Atlanta',
              postalCode: '30030',
              state: 'GA',
              streetAddress1: 'xxx',
            },
            backup_mailing_address: {
              city: 'Atlanta',
              postalCode: '30030',
              state: 'GA',
              streetAddress1: 'zzz',
            },
          },
          context: ppmContext,
        });
        expect(result).toEqual('/service-member/backup-contact');
      });
    });
    describe('when backup contacts is complete', () => {
      it('returns the order transition page', () => {
        const backupContacts = [
          {
            createdAt: '2018-05-31T00:02:57.302Z',
            email: 'foo@bar.com',
            id: '03b2979d-8046-437b-a6e4-11dbe251a912',
            name: 'Foo',
            permission: 'NONE',
            updated_at: '2018-05-31T00:02:57.302Z',
          },
        ];
        const sm = {
          ...serviceMember,
          is_profile_complete: true,
          edipi: '1234567890',
          affiliation: 'Marines',
          last_name: 'foo',
          first_name: 'foo',
          phone_is_preferred: true,
          telephone: '666-666-6666',
          personal_email: 'foo@bar.com',
          current_location: {
            id: '5e30f356-e590-4372-b9c0-30c3fd1ff42d',
            name: 'Blue Grass Army Depot',
          },
          residential_address: {
            city: 'Atlanta',
            postalCode: '30030',
            state: 'GA',
            streetAddress1: 'xxx',
          },
          backup_mailing_address: {
            city: 'Atlanta',
            postalCode: '30030',
            state: 'GA',
            streetAddress1: 'zzz',
          },
        };
        const orders = { id: 'testId' };
        const result = getNextIncompletePage({
          serviceMember: sm,
          orders,
          backupContacts,
          context: ppmContext,
        });
        expect(result).toEqual('/orders/info/testId');
      });
    });
  });
  describe('when the profile is complete', () => {
    it('returns the orders info', () => {
      const orders = { id: 'testId' };
      const result = getNextIncompletePage({
        serviceMember: {
          ...serviceMember,
          is_profile_complete: true,
        },
        orders,
        context: ppmContext,
      });
      expect(result).toEqual('/orders/info/testId');
    });
    describe('when orders info is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          serviceMember: {
            ...serviceMember,
            is_profile_complete: true,
          },
          move: { id: 'bar' },
          orders: {
            id: 'bar',
            orders_type: 'foo',
            issue_date: '2019-01-01',
            report_by_date: '2019-02-01',
            new_duty_location: { id: 'something' },
            origin_duty_location: { id: 'something' },
            grade: 'E_4',
          },
          context: ppmContext,
        });
        expect(result).toEqual('/orders/upload/bar');
      });
    });
    describe('when orders upload is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          serviceMember: {
            ...serviceMember,
            is_profile_complete: true,
          },
          orders: {
            orders_type: 'foo',
            issue_date: '2019-01-01',
            report_by_date: '2019-02-01',
            new_duty_location: { id: 'something' },
            origin_duty_location: { id: 'something' },
            uploaded_orders: {
              uploads: [{}],
            },
            grade: 'E_4',
          },
          move: { id: 'bar' },
          uploads: [
            {
              contentType: 'application/pdf',
              filename: 'testfile.pdf',
              status: 'PROCESSING',
              url: 'storage/user/1234pdf',
            },
          ],
          context: ppmContext,
        });
        expect(result).toEqual('/moves/bar/shipment-type');
      });
    });
    describe('when ppm date is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          serviceMember: {
            ...serviceMember,
            is_profile_complete: true,
          },
          orders: {
            orders_type: 'foo',
            issue_date: '2019-01-01',
            report_by_date: '2019-02-01',
            new_duty_location: { id: 'something' },
            origin_duty_location: { id: 'something' },
            uploaded_orders: {
              uploads: [{}],
            },
            grade: 'E_4',
          },
          move: {
            id: 'bar',
            mtoShipments: [{ shipmentType: SHIPMENT_OPTIONS.PPM }],
          },
          ppm: {
            id: 'baz',
            original_move_date: '2018-10-10',
            pickup_postal_code: '22222',
            destination_postal_code: '22222',
          },
          context: ppmContext,
        });
        expect(result).toEqual('/moves/bar/review');
      });
    });
    describe('when ppm incentive is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          serviceMember: {
            ...serviceMember,
            is_profile_complete: true,
          },
          orders: {
            orders_type: 'foo',
            issue_date: '2019-01-01',
            report_by_date: '2019-02-01',
            new_duty_location: { id: 'something' },
            origin_duty_location: { id: 'something' },
            uploaded_orders: {
              uploads: [{}],
            },
            grade: 'E_4',
          },
          move: {
            id: 'bar',
            mtoShipments: [{ shipmentType: SHIPMENT_OPTIONS.PPM }],
          },
          ppm: {
            id: 'baz',
            original_move_date: '2018-10-10',
            pickup_postal_code: '22222',
            destination_postal_code: '22222',
            weight_estimate: 555,
          },
          context: ppmContext,
        });
        expect(result).toEqual('/moves/bar/review');
      });
    });
  });
});

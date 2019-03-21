import 'raf/polyfill';
import { getPagesInFlow, getNextIncompletePage } from './getWorkflowRoutes';
import { NULL_UUID } from 'shared/constants';
describe('when getting the routes for the current workflow', () => {
  describe('given a complete service member', () => {
    describe('given a PPM', () => {
      const props = {
        selectedMoveType: 'PPM',
      };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns service member, order and move pages', () => {
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
    describe('given a HHG PPM move', () => {
      const props = { selectedMoveType: 'HHG_PPM' };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns profile review, the order and move pages', () => {
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
          '/moves/:moveId/hhg-ppm-start',
          '/moves/:moveId/hhg-ppm-size',
          '/moves/:moveId/hhg-ppm-weight',
          '/moves/:moveId/review',
          '/moves/:moveId/hhg-ppm-agreement',
        ]);
      });
    });
    describe('given a canceled PPM', () => {
      const props = { lastMoveIsCanceled: true, selectedMoveType: 'PPM' };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns profile review, the order and move pages', () => {
        expect(pages).toEqual([
          '/profile-review',
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
        selectedMoveType: 'HHG',
      };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns service member, order and move pages', () => {
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
          '/moves/:moveId/hhg-start',
          '/moves/:moveId/hhg-locations',
          '/moves/:moveId/hhg-weight',
          '/moves/:moveId/hhg-progear',
          '/moves/:moveId/review',
          '/moves/:moveId/agreement',
        ]);
      });
    });
  });
  describe('given an incomplete service member', () => {
    describe('given no move', () => {
      const props = {
        selectedMoveType: null,
      };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns service member, order and move pages', () => {
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
        selectedMoveType: 'PPM',
      };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns service member, order and PPM-specific move pages', () => {
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
        selectedMoveType: 'HHG',
      };
      const pages = getPagesInFlow(props);
      it('getPagesInFlow returns service member, order and HHG-specific move pages', () => {
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
          '/moves/:moveId/hhg-start',
          '/moves/:moveId/hhg-locations',
          '/moves/:moveId/hhg-weight',
          '/moves/:moveId/hhg-progear',
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
        selectedMoveType: 'HHG',
        serviceMember,
      });
      expect(result).toEqual('/service-member/foo/create');
    });
    describe('when dod info is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          serviceMember: {
            ...serviceMember,
            is_profile_complete: false,
            edipi: '1234567890',
            has_social_security_number: true,
            rank: 'E_6',
            affiliation: 'Marines',
          },
        });
        expect(result).toEqual('/service-member/foo/name');
      });
    });
    describe('when name is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          serviceMember: {
            ...serviceMember,
            is_profile_complete: false,
            edipi: '1234567890',
            has_social_security_number: true,
            rank: 'E_6',
            affiliation: 'Marines',
            last_name: 'foo',
            first_name: 'foo',
          },
        });
        expect(result).toEqual('/service-member/foo/contact-info');
      });
    });
    describe('when contact-info is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          selectedMoveType: 'PPM',
          serviceMember: {
            ...serviceMember,
            is_profile_complete: false,
            edipi: '1234567890',
            has_social_security_number: true,
            rank: 'E_6',
            affiliation: 'Marines',
            last_name: 'foo',
            first_name: 'foo',
            email_is_preferred: true,
            telephone: '666-666-6666',
            personal_email: 'foo@bar.com',
            current_station: {
              id: NULL_UUID,
              name: '',
            },
          },
        });
        expect(result).toEqual('/service-member/foo/duty-station');
      });
    });
    describe('when duty-station is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          serviceMember: {
            ...serviceMember,
            is_profile_complete: false,
            edipi: '1234567890',
            has_social_security_number: true,
            rank: 'E_6',
            affiliation: 'Marines',
            last_name: 'foo',
            first_name: 'foo',
            phone_is_preferred: true,
            telephone: '666-666-6666',
            personal_email: 'foo@bar.com',
            current_station: {
              id: '5e30f356-e590-4372-b9c0-30c3fd1ff42d',
              name: 'Blue Grass Army Depot',
            },
          },
        });
        expect(result).toEqual('/service-member/foo/residence-address');
      });
    });
    describe('when residence address is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          serviceMember: {
            ...serviceMember,
            is_profile_complete: false,
            edipi: '1234567890',
            has_social_security_number: true,
            rank: 'E_6',
            affiliation: 'Marines',
            last_name: 'foo',
            first_name: 'foo',
            phone_is_preferred: true,
            telephone: '666-666-6666',
            personal_email: 'foo@bar.com',
            current_station: {
              id: '5e30f356-e590-4372-b9c0-30c3fd1ff42d',
              name: 'Blue Grass Army Depot',
            },
            residential_address: {
              city: 'Atlanta',
              postal_code: '30030',
              state: 'GA',
              street_address_1: 'xxx',
            },
          },
        });
        expect(result).toEqual('/service-member/foo/backup-mailing-address');
      });
    });
    describe('when backup address is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          serviceMember: {
            ...serviceMember,
            is_profile_complete: false,
            edipi: '1234567890',
            has_social_security_number: true,
            rank: 'E_6',
            affiliation: 'Marines',
            last_name: 'foo',
            first_name: 'foo',
            phone_is_preferred: true,
            telephone: '666-666-6666',
            personal_email: 'foo@bar.com',
            current_station: {
              id: '5e30f356-e590-4372-b9c0-30c3fd1ff42d',
              name: 'Blue Grass Army Depot',
            },
            residential_address: {
              city: 'Atlanta',
              postal_code: '30030',
              state: 'GA',
              street_address_1: 'xxx',
            },
            backup_mailing_address: {
              city: 'Atlanta',
              postal_code: '30030',
              state: 'GA',
              street_address_1: 'zzz',
            },
          },
        });
        expect(result).toEqual('/service-member/foo/backup-contacts');
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
          is_profile_complete: false,
          edipi: '1234567890',
          has_social_security_number: true,
          rank: 'E_6',
          affiliation: 'Marines',
          last_name: 'foo',
          first_name: 'foo',
          phone_is_preferred: true,
          telephone: '666-666-6666',
          personal_email: 'foo@bar.com',
          current_station: {
            id: '5e30f356-e590-4372-b9c0-30c3fd1ff42d',
            name: 'Blue Grass Army Depot',
          },
          residential_address: {
            city: 'Atlanta',
            postal_code: '30030',
            state: 'GA',
            street_address_1: 'xxx',
          },
          backup_mailing_address: {
            city: 'Atlanta',
            postal_code: '30030',
            state: 'GA',
            street_address_1: 'zzz',
          },
        };
        const result = getNextIncompletePage({
          serviceMember: sm,
          backupContacts,
        });
        expect(result).toEqual('/orders/');
      });
    });
  });
  describe('when the profile is complete', () => {
    it('returns the orders info', () => {
      const result = getNextIncompletePage({
        serviceMember: {
          ...serviceMember,
          is_profile_complete: true,
        },
      });
      expect(result).toEqual('/orders/');
    });
    describe('when orders info is complete', () => {
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
            new_duty_station: { id: 'something' },
          },
          move: { id: 'bar' },
        });
        expect(result).toEqual('/orders/upload');
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
            new_duty_station: { id: 'something' },
            uploaded_orders: {
              uploads: [{}],
            },
          },
          move: { id: 'bar' },
        });
        expect(result).toEqual('/moves/bar');
      });
    });
    describe('when move type selection is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          selectedMoveType: 'PPM',
          serviceMember: {
            ...serviceMember,
            is_profile_complete: true,
          },
          orders: {
            orders_type: 'foo',
            issue_date: '2019-01-01',
            report_by_date: '2019-02-01',
            new_duty_station: { id: 'something' },
            uploaded_orders: {
              uploads: [{}],
            },
          },
          move: {
            id: 'bar',
            selected_move_type: 'PPM',
          },
          ppm: {
            id: 'baz',
          },
        });
        expect(result).toEqual('/moves/bar/ppm-start');
      });
    });
    describe('when ppm date is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          selectedMoveType: 'PPM',
          serviceMember: {
            ...serviceMember,
            is_profile_complete: true,
          },
          orders: {
            orders_type: 'foo',
            issue_date: '2019-01-01',
            report_by_date: '2019-02-01',
            new_duty_station: { id: 'something' },
            uploaded_orders: {
              uploads: [{}],
            },
          },
          move: {
            id: 'bar',
            selected_move_type: 'PPM',
          },
          ppm: {
            id: 'baz',
            original_move_date: '2018-10-10',
            pickup_postal_code: '22222',
            destination_postal_code: '22222',
          },
        });
        expect(result).toEqual('/moves/bar/ppm-size');
      });
    });
    describe('when ppm size is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          selectedMoveType: 'PPM',
          serviceMember: {
            ...serviceMember,
            is_profile_complete: true,
          },
          orders: {
            orders_type: 'foo',
            issue_date: '2019-01-01',
            report_by_date: '2019-02-01',
            new_duty_station: { id: 'something' },
            uploaded_orders: {
              uploads: [{}],
            },
          },
          move: {
            id: 'bar',
            selected_move_type: 'PPM',
          },
          ppm: {
            id: 'baz',
            original_move_date: '2018-10-10',
            pickup_postal_code: '22222',
            destination_postal_code: '22222',
            size: 'L',
          },
        });
        expect(result).toEqual('/moves/bar/ppm-incentive');
      });
    });
    describe('when ppm incentive is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          selectedMoveType: 'PPM',
          serviceMember: {
            ...serviceMember,
            is_profile_complete: true,
          },
          orders: {
            orders_type: 'foo',
            issue_date: '2019-01-01',
            report_by_date: '2019-02-01',
            new_duty_station: { id: 'something' },
            uploaded_orders: {
              uploads: [{}],
            },
          },
          move: {
            id: 'bar',
            selected_move_type: 'PPM',
          },
          ppm: {
            id: 'baz',
            original_move_date: '2018-10-10',
            pickup_postal_code: '22222',
            destination_postal_code: '22222',
            size: 'L',
            weight_estimate: 555,
          },
        });
        expect(result).toEqual('/moves/bar/review');
      });
    });
    describe('when HHG move type selection is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          selectedMoveType: 'HHG',
          serviceMember: {
            ...serviceMember,
            is_profile_complete: true,
          },
          orders: {
            orders_type: 'foo',
            issue_date: '2019-01-01',
            report_by_date: '2019-02-01',
            new_duty_station: { id: 'something' },
            uploaded_orders: {
              uploads: [{}],
            },
          },
          move: {
            id: 'bar',
            selected_move_type: 'HHG',
          },
          hhg: {
            id: 'baz',
          },
        });
        expect(result).toEqual('/moves/bar/hhg-start');
      });
    });
  });
});

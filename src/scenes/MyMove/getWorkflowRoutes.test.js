import { getPagesInFlow, getNextIncompletePage } from './getWorkflowRoutes';

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
    describe('given a complete service member with an HHG', () => {
      const props = {
        hasCompleteProfile: profileIsComplete,
        selectedMoveType: 'HHG',
        hasMove: true,
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
        hasCompleteProfile: profileIsComplete,
        selectedMoveType: 'PPM',
        hasMove: true,
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
        hasCompleteProfile: profileIsComplete,
        selectedMoveType: 'HHG',
        hasMove: true,
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

describe('when getting the next incomplete page', () => {
  const service_member = { id: 'foo' };
  describe('when the profile is incomplete', () => {
    it('returns the first page of the user profile', () => {
      const result = getNextIncompletePage(service_member);
      expect(result).toEqual('/service-member/foo/create');
    });
    describe('when dod info is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          ...service_member,
          is_profile_complete: false,
          edipi: '1234567890',
          has_social_security_number: true,
          rank: 'E_6',
          affiliation: 'Marines',
        });
        expect(result).toEqual('/service-member/foo/name');
      });
    });
    describe('when name is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          ...service_member,
          is_profile_complete: false,
          edipi: '1234567890',
          has_social_security_number: true,
          rank: 'E_6',
          affiliation: 'Marines',
          last_name: 'foo',
          first_name: 'foo',
        });
        expect(result).toEqual('/service-member/foo/contact-info');
      });
    });
    describe('when contact-info is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          ...service_member,
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
        });
        expect(result).toEqual('/service-member/foo/duty-station');
      });
    });
    describe('when duty-station is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          ...service_member,
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
        });
        expect(result).toEqual('/service-member/foo/residence-address');
      });
    });
    describe('when residence address is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          ...service_member,
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
        });
        expect(result).toEqual('/service-member/foo/backup-mailing-address');
      });
    });
    describe('when backup address is complete', () => {
      it('returns the next page of the user profile', () => {
        const result = getNextIncompletePage({
          ...service_member,
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
        });
        expect(result).toEqual('/service-member/foo/backup-contacts');
      });
    });
  });
  describe('when the profile is incomplete', () => {
    // service_member.is_profile_complete = true;
    it('returns the orders info', () => {
      const result = getNextIncompletePage({
        ...service_member,
        is_profile_complete: true,
      });
      expect(result).toEqual('/orders/');
    });
    describe('when orders info is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          ...service_member,
          is_profile_complete: true,
          orders: [
            {
              orders_type: 'foo',
              issue_date: '2019-01-01',
              report_by_date: '2019-02-01',
              new_duty_station: {},
            },
          ],
        });
        expect(result).toEqual('/orders/upload');
      });
    });
    describe('when orders upload is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          ...service_member,
          is_profile_complete: true,
          orders: [
            {
              orders_type: 'foo',
              issue_date: '2019-01-01',
              report_by_date: '2019-02-01',
              new_duty_station: {},
              uploaded_orders: {
                uploads: [{}],
              },
              moves: [{ id: 'bar' }],
            },
          ],
        });
        expect(result).toEqual('/moves/bar');
      });
    });
    describe('when move type selection is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          ...service_member,
          is_profile_complete: true,
          orders: [
            {
              orders_type: 'foo',
              issue_date: '2019-01-01',
              report_by_date: '2019-02-01',
              new_duty_station: {},
              uploaded_orders: {
                uploads: [{}],
              },
              moves: [
                {
                  id: 'bar',
                  selected_move_type: 'PPM',
                  personally_procured_moves: [{ id: 'baz' }],
                },
              ],
            },
          ],
        });
        expect(result).toEqual('/moves/bar/ppm-start');
      });
    });
    describe('when ppm date is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          ...service_member,
          is_profile_complete: true,
          orders: [
            {
              orders_type: 'foo',
              issue_date: '2019-01-01',
              report_by_date: '2019-02-01',
              new_duty_station: {},
              uploaded_orders: {
                uploads: [{}],
              },
              moves: [
                {
                  id: 'bar',
                  selected_move_type: 'PPM',
                  personally_procured_moves: [
                    {
                      id: 'baz',
                      planned_move_date: '2018-10-10',
                      pickup_postal_code: '22222',
                      destination_postal_code: '22222',
                    },
                  ],
                },
              ],
            },
          ],
        });
        expect(result).toEqual('/moves/bar/ppm-size');
      });
    });
    describe('when ppm size is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          ...service_member,
          is_profile_complete: true,
          orders: [
            {
              orders_type: 'foo',
              issue_date: '2019-01-01',
              report_by_date: '2019-02-01',
              new_duty_station: {},
              uploaded_orders: {
                uploads: [{}],
              },
              moves: [
                {
                  id: 'bar',
                  selected_move_type: 'PPM',
                  personally_procured_moves: [
                    {
                      id: 'baz',
                      planned_move_date: '2018-10-10',
                      pickup_postal_code: '22222',
                      destination_postal_code: '22222',
                      size: 'L',
                    },
                  ],
                },
              ],
            },
          ],
        });
        expect(result).toEqual('/moves/bar/ppm-incentive');
      });
    });
    describe('when ppm incentive is complete', () => {
      it('returns the next page', () => {
        const result = getNextIncompletePage({
          ...service_member,
          is_profile_complete: true,
          orders: [
            {
              orders_type: 'foo',
              issue_date: '2019-01-01',
              report_by_date: '2019-02-01',
              new_duty_station: {},
              uploaded_orders: {
                uploads: [{}],
              },
              moves: [
                {
                  id: 'bar',
                  selected_move_type: 'PPM',
                  personally_procured_moves: [
                    {
                      id: 'baz',
                      planned_move_date: '2018-10-10',
                      pickup_postal_code: '22222',
                      destination_postal_code: '22222',
                      size: 'L',
                      weight: 555,
                    },
                  ],
                },
              ],
            },
          ],
        });
        expect(result).toEqual('/moves/bar/review');
      });
    });
  });
});

import * as ducks from './ducks';
describe('TransportationOffices', () => {
  describe('GET_DUTY_STATION_TRANSPORTATION_OFFICE', () => {
    it('init', () => {
      const state = ducks.reducer(undefined, { type: 'init' });
      expect(state).toEqual({
        isLoading: false,
        hasErrored: false,
        hasLoaded: false,
        byId: {},
        allIds: [],
        byDutyStationId: {},
      });
    });
    it('start', () => {
      const state = ducks.reducer(undefined, {
        type: ducks.GET_DUTY_STATION_TRANSPORTATION_OFFICE.start,
      });
      expect(state).toEqual({
        isLoading: true,
        hasErrored: false,
        hasLoaded: false,
        byId: {},
        allIds: [],
        byDutyStationId: {},
      });
    });
    it('success', () => {
      const state = ducks.reducer(undefined, {
        type: ducks.GET_DUTY_STATION_TRANSPORTATION_OFFICE.success,
        payload: {
          transportationOffice: {
            address: {
              city: 'White Sands Missile Range',
              country: 'United States',
              postal_code: '88002',
              state: 'NM',
              street_address_1: '143 Crozier St',
              street_address_2: '',
            },
            created_at: '2018-05-28T14:27:39.635Z',
            id: 'ff6f1f7c-5309-4436-b9ff-e5dcac95b750',
            name: 'White Sands Missile Range',
            phone_lines: [
              '(575) 678-5005',
              '(575) 678-3506',
              '258-5055',
              '258-3506',
            ],
            updated_at: '2018-05-28T14:27:39.635Z',
          },
          dutyStationId: '4FFFB7F8-603C-46E1-9A0F-7F2EAD11700C',
        },
      });
      expect(state).toEqual({
        isLoading: false,
        hasErrored: false,
        hasLoaded: true,
        allIds: ['ff6f1f7c-5309-4436-b9ff-e5dcac95b750'],
        byDutyStationId: {
          '4FFFB7F8-603C-46E1-9A0F-7F2EAD11700C':
            'ff6f1f7c-5309-4436-b9ff-e5dcac95b750',
        },
        byId: {
          'ff6f1f7c-5309-4436-b9ff-e5dcac95b750': {
            address: {
              city: 'White Sands Missile Range',
              country: 'United States',
              postal_code: '88002',
              state: 'NM',
              street_address_1: '143 Crozier St',
              street_address_2: '',
            },
            created_at: '2018-05-28T14:27:39.635Z',
            id: 'ff6f1f7c-5309-4436-b9ff-e5dcac95b750',
            name: 'White Sands Missile Range',
            phone_lines: [
              '(575) 678-5005',
              '(575) 678-3506',
              '258-5055',
              '258-3506',
            ],
            updated_at: '2018-05-28T14:27:39.635Z',
          },
        },
      });
    });
    it('failure', () => {
      const state = ducks.reducer(undefined, {
        type: ducks.GET_DUTY_STATION_TRANSPORTATION_OFFICE.failure,
        error: 'ruh roh',
      });
      expect(state).toEqual({
        isLoading: false,
        hasErrored: true,
        hasLoaded: false,
        error: 'ruh roh',
        allIds: [],
        byId: {},
        byDutyStationId: {},
      });
    });
  });
});

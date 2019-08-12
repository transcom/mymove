import { getMaxAdvance, isHHGPPMComboMove, getDestinationPostalCode } from './ppms';

describe('PPMs utility functions', () => {
  describe('getMaxAdvance', () => {
    describe('when there is a max estimated incentive', () => {
      it('should return 60% of max estimated incentive', () => {
        const state = {
          entities: {
            personallyProcuredMoves: {
              'deb28967-d52c-4f04-8a0b-a264c9d80457': {
                incentive_estimate_max: 10000,
              },
            },
          },
        };
        expect(getMaxAdvance(state, 'deb28967-d52c-4f04-8a0b-a264c9d80457')).toEqual(6000);
      });
    });
  });

  describe('when there is no max estimated incentive', () => {
    const state = {};
    it('should return 60% of max estimated incentive', () => {
      expect(getMaxAdvance(state, 'deb28967-d52c-4f04-8a0b-a264c9d80457')).toEqual(20000000);
    });
  });

  describe('isHHGPPMComboMove', () => {
    describe('when move is a combo move', () => {
      it('should return true', () => {
        const state = {
          entities: {
            moves: {
              moveId: {
                selected_move_type: 'HHG_PPM',
              },
            },
          },
        };
        expect(isHHGPPMComboMove(state, 'moveId')).toBe(true);
      });
    });
    describe('when move is not a combo move', () => {
      it('should return false', () => {
        const state = {
          entities: {
            moves: {
              moveId: {
                selected_move_type: 'PPM',
              },
            },
          },
        };
        expect(isHHGPPMComboMove(state, 'moveId')).toBe(false);
      });
    });
  });

  describe('getDestinationPostalCode', () => {
    describe('when there is no delivery address', () => {
      it('should return new duty station zip', () => {
        const state = {
          orders: {
            currentOrders: {
              id: '9dd3e284-ac16-43db-a3a2-20397f0072d7',
              new_duty_station: {
                address: {
                  city: 'Richmond',
                  country: 'United States',
                  id: '55a694af-cdfa-4a80-be95-71476064c091',
                  postal_code: '40475',
                  state: 'KY',
                  street_address_1: 'n/a',
                },
              },
            },
          },
          ui: {
            currentShipmentID: '0194fb44-3762-4d1c-b58a-f6daf984813b',
            currentMoveID: '56b8ef45-8145-487b-9b59-0e30d0d465fa',
          },
          entities: {
            shipments: {
              '0194fb44-3762-4d1c-b58a-f6daf984813b': {
                gross_weight: 5000,
                id: '0194fb44-3762-4d1c-b58a-f6daf984813b',
                market: 'dHHG',
                move_id: '56b8ef45-8145-487b-9b59-0e30d0d465fa',
                pm_survey_weight_estimate: 5000,
                progear_weight_estimate: 225,
                spouse_progear_weight_estimate: 312,
                status: 'AWARDED',
                tare_weight: 1500,
                weight_estimate: 2000,
              },
            },
          },
        };
        expect(getDestinationPostalCode(state)).toBe('40475');
      });
    });

    describe('when there is a delivery address', () => {
      it('should return delivery address zip', () => {
        const state = {
          orders: {
            currentOrders: {
              id: '9dd3e284-ac16-43db-a3a2-20397f0072d7',
              new_duty_station: {
                address: {
                  city: 'Richmond',
                  country: 'United States',
                  id: '55a694af-cdfa-4a80-be95-71476064c091',
                  postal_code: '40475',
                  state: 'KY',
                  street_address_1: 'n/a',
                },
              },
            },
          },
          ui: {
            currentShipmentID: '0194fb44-3762-4d1c-b58a-f6daf984813b',
            currentMoveID: '56b8ef45-8145-487b-9b59-0e30d0d465fa',
          },
          entities: {
            shipments: {
              '0194fb44-3762-4d1c-b58a-f6daf984813b': {
                gross_weight: 5000,
                id: '0194fb44-3762-4d1c-b58a-f6daf984813b',
                market: 'dHHG',
                move_id: '56b8ef45-8145-487b-9b59-0e30d0d465fa',
                pm_survey_weight_estimate: 5000,
                progear_weight_estimate: 225,
                spouse_progear_weight_estimate: 312,
                status: 'AWARDED',
                tare_weight: 1500,
                weight_estimate: 2000,
                has_delivery_address: true,
                delivery_address: '3295aabe-cd8d-4a8c-975d-b9a1d01edd02',
              },
            },
            addresses: {
              '3295aabe-cd8d-4a8c-975d-b9a1d01edd02': {
                city: 'some city',
                country: 'United States',
                id: '3295aabe-cd8d-4a8c-975d-b9a1d01edd02',
                postal_code: '55555',
                state: 'KY',
                street_address_1: 'n/a',
              },
            },
          },
        };
        expect(getDestinationPostalCode(state)).toBe('55555');
      });
    });
  });
});

import { getEstimatedRemainingWeight, getActualRemainingWeight } from './ppms';

describe('Ppm Reducer', () => {
  describe('getEstimatedRemainingWeight', () => {
    describe('when there is an estimated remaining weight from a pre move survey', () => {
      it('should return the proper estimated remaining weight based on the pre move survey', () => {
        const state = {
          orders: {
            currentOrders: {
              id: '9dd3e284-ac16-43db-a3a2-20397f0072d7',
              has_dependents: true,
              spouse_has_pro_gear: true,
            },
          },
          serviceMember: {
            currentServiceMember: {
              rank: 'E_1',
              weight_allotment: {
                total_weight_self: 5000,
                total_weight_self_plus_dependents: 8000,
                pro_gear_weight: 2000,
                pro_gear_weight_spouse: 500,
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
            moves: {
              '56b8ef45-8145-487b-9b59-0e30d0d465fa': {
                id: '56b8ef45-8145-487b-9b59-0e30d0d465fa',
                selected_move_type: 'HHG_PPM',
              },
            },
          },
        };
        expect(getEstimatedRemainingWeight(state)).toEqual(5500);
      });
    });

    describe('when there is an estimated remaining weight from a service member entered weight', () => {
      it('should return the proper estimated remaining weight based on the entered values', () => {
        const state = {
          orders: {
            currentOrders: {
              id: '9dd3e284-ac16-43db-a3a2-20397f0072d7',
              has_dependents: true,
              spouse_has_pro_gear: true,
            },
          },
          serviceMember: {
            currentServiceMember: {
              rank: 'E_1',
              weight_allotment: {
                total_weight_self: 5000,
                total_weight_self_plus_dependents: 8000,
                pro_gear_weight: 2000,
                pro_gear_weight_spouse: 500,
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
                progear_weight_estimate: 225,
                spouse_progear_weight_estimate: 312,
                status: 'AWARDED',
                tare_weight: 1500,
                weight_estimate: 2000,
              },
            },
            moves: {
              '56b8ef45-8145-487b-9b59-0e30d0d465fa': {
                id: '56b8ef45-8145-487b-9b59-0e30d0d465fa',
                selected_move_type: 'HHG_PPM',
              },
            },
          },
        };
        expect(getEstimatedRemainingWeight(state)).toEqual(8500);
      });
    });

    describe('getActualRemainingWeight', () => {
      describe('when there is an actual weight', () => {
        it('should return the proper actual weight', () => {
          const state = {
            orders: {
              currentOrders: {
                id: '9dd3e284-ac16-43db-a3a2-20397f0072d7',
                has_dependents: true,
                spouse_has_pro_gear: true,
              },
            },
            serviceMember: {
              currentServiceMember: {
                rank: 'E_1',
                weight_allotment: {
                  total_weight_self: 5000,
                  total_weight_self_plus_dependents: 8000,
                  pro_gear_weight: 2000,
                  pro_gear_weight_spouse: 500,
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
              moves: {
                '56b8ef45-8145-487b-9b59-0e30d0d465fa': {
                  id: '56b8ef45-8145-487b-9b59-0e30d0d465fa',
                  selected_move_type: 'HHG_PPM',
                },
              },
            },
          };
          expect(getActualRemainingWeight(state)).toEqual(7000);
        });
      });
    });
  });
});

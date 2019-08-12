import { getMaxAdvance, isHHGPPMComboMove } from './ppms';

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
});

import friendlyBranchRank from './branchRankFormatters';

describe('branchRank formatters', () => {
  describe('friendlyBranchRank', () => {
    it('returns a formatted string of readable branch and rank', () => {
      expect(friendlyBranchRank('AIR_FORCE', 'E_6')).toEqual('Air Force, E-6');
    });

    it('returns empty string if orders or rank do not match consts', () => {
      expect(friendlyBranchRank(undefined, undefined)).toEqual('');
    });
  });
});

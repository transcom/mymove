import friendlyBranchGrade from './branchGradeFormatters';

describe('branchGrade formatters', () => {
  describe('friendlyBranchGrade', () => {
    it('returns a formatted string of readable branch and grade', () => {
      expect(friendlyBranchGrade('AIR_FORCE', 'E_6')).toEqual('Air Force, E-6');
    });

    it('returns empty string if orders or grade do not match consts', () => {
      expect(friendlyBranchGrade(undefined, undefined)).toEqual('');
    });
  });
});

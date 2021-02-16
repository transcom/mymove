import getRoleTypesFromRoles from 'utils/user';

describe('getRoleTypesFromRoles', () => {
  it('returns an array of role types', () => {
    const roles = [{ roleType: 'TOO' }, { roleType: 'TIO' }];

    expect(getRoleTypesFromRoles(roles)).toEqual([roles[0].roleType, roles[1].roleType]);
  });
});

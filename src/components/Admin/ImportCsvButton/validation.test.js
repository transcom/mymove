import { checkTelephone, parseRoles } from './validation';

import { roleTypes, adminOfficeRoles } from 'constants/userRoles';

describe('checkRequiredFields', () => {
  it('success: does nothing if all fields provided', () => {});

  it('fail: throws an error if fields are missing', () => {});
});

describe('checkTelephone', () => {
  it('success: does nothing if telephone is valid', () => {
    expect(checkTelephone({ telephone: '209-555-1234' })).toBeTruthy();
  });

  it('fail: throws an error if telephone is invalid', () => {
    function checkInvalidTelephone() {
      checkTelephone({ telephone: '111-111-111' });
    }
    expect(checkInvalidTelephone).toThrowError('Row contains improperly formatted telephone number.');
  });
});

describe('parseRoles', () => {
  const servicesCounselorAdminRole = adminOfficeRoles.filter(
    (role) => role.roleType === roleTypes.SERVICES_COUNSELOR,
  )[0];
  const tooAdminRole = adminOfficeRoles.filter((role) => role.roleType === roleTypes.TOO)[0];

  it('fail: throws an error if there are no roles', () => {
    function parseEmptyRoles() {
      parseRoles('');
    }
    expect(parseEmptyRoles).toThrowError('Unable to parse roles for row.');
  });

  it('success: parses one role into an array of len 1', () => {
    const roles = roleTypes.SERVICES_COUNSELOR;
    const rolesArray = parseRoles(roles);
    expect(rolesArray).toHaveLength(1);
    expect(rolesArray).toContainEqual(servicesCounselorAdminRole);
  });

  it('success: parses multiple roles into an array', () => {
    const roles = `${roleTypes.SERVICES_COUNSELOR}, ${roleTypes.TOO}`;
    const rolesArray = parseRoles(roles);
    expect(rolesArray).toHaveLength(2);
    expect(rolesArray).toEqual(expect.arrayContaining([servicesCounselorAdminRole, tooAdminRole]));
  });

  it('fail: throws an error if there is an invalid role', () => {
    function parseInvalidRoles() {
      parseRoles('test_role');
    }
    expect(parseInvalidRoles).toThrowError('Invalid roles provided for row.');
  });
});

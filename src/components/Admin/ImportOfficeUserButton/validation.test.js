import { checkRequiredFields, checkTelephone, parseRoles, parsePrivileges } from './validation';

import { adminOfficeRoles, roleTypes } from 'constants/userRoles';
import { elevatedPrivilegeTypes } from 'constants/userPrivileges';

describe('checkRequiredFields', () => {
  it('success: does nothing if all fields provided', () => {
    expect(
      checkRequiredFields({
        transportationOfficeId: 'id',
        firstName: 'Test',
        lastName: 'Tester',
        roles: 'test_role',
        email: 'test@example.com',
        telephone: '222-555-1111',
      }),
    ).toBeTruthy();
  });

  it('fail: throws an error if fields are missing', () => {
    function checkMissingFields() {
      checkRequiredFields({ firstName: 'Blank' });
    }
    expect(checkMissingFields).toThrowError('Row does not contain all required fields.');
  });
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

describe('parsePrivileges', () => {
  const supervisorPrivilege = { privilegeType: 'supervisor', name: 'Supervisor' };
  const safetyPrivilege = { privilegeType: 'safety', name: 'Safety Moves' };

  it('fail: throws an error if there are no privileges', () => {
    function parseEmptyPrivileges() {
      parsePrivileges('');
    }
    expect(parseEmptyPrivileges).toThrowError('Processing Error: Unable to parse privileges for row.');
  });

  it('success: parses one privilege into an array of len 1', () => {
    const privileges = elevatedPrivilegeTypes.SUPERVISOR;
    const privilegesArray = parsePrivileges(privileges);
    expect(privilegesArray).toHaveLength(1);
    expect(privilegesArray).toContainEqual(supervisorPrivilege);
  });

  it('success: parses multiple privileges into an array', () => {
    const privileges = `${elevatedPrivilegeTypes.SUPERVISOR}, ${elevatedPrivilegeTypes.SAFETY}`;
    const privilegesArray = parsePrivileges(privileges);
    expect(privilegesArray).toHaveLength(2);
    expect(privilegesArray).toEqual(expect.arrayContaining([supervisorPrivilege, safetyPrivilege]));
  });

  it('fail: throws an error if there is an invalid role', () => {
    function parseInvalidPrivileges() {
      parsePrivileges('test_privilege');
    }
    expect(parseInvalidPrivileges).toThrowError('Invalid privileges provided for row.');
  });
});

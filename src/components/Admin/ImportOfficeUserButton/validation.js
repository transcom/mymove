import { adminOfficeRoles } from 'constants/userRoles';
import { officeUserPrivileges } from 'constants/userPrivileges';

export const checkRequiredFields = ({ transportationOfficeId, firstName, lastName, roles, email, telephone }) => {
  if (!(transportationOfficeId && firstName && lastName && roles && email && telephone)) {
    throw new Error(
      `Validation Error: Row does not contain all required fields.
        Required fields are firstName, lastName, email, telephone, roles, and transportationOfficeId.`,
    );
  }
  return true;
};

export const checkTelephone = ({ telephone }) => {
  // Verify the phone format
  const regex = /^[2-9]\d{2}-\d{3}-\d{4}$/;
  if (!regex.test(telephone)) {
    throw new Error(
      'Validation Error: Row contains improperly formatted telephone number. Required format is xxx-xxx-xxxx.',
    );
  }
  return true;
};

export const checkValidRolesWithPrivileges = (row) => {
  if (
    (row.roles.indexOf('customer') >= 0 || row.roles.indexOf('contracting_officer') >= 0) &&
    (row.privileges.indexOf('supervisor') >= 0 || row.privileges.indexOf('safety') >= 0)
  ) {
    throw new Error('Privileges cannot be selected with Customer or Contracting Officer roles.');
  }
  return true;
};

export const parseRoles = (roles) => {
  if (!roles) {
    throw new Error('Processing Error: Unable to parse roles for row.');
  }
  const rolesArray = [];

  // Parse roles from string at ","
  const parsedRoles = roles.split(',');
  parsedRoles.forEach((parsedRole) => {
    let roleNotFound = true;
    // Remove any whitespace in the role string
    const role = parsedRole.replace(/\s/g, '');
    adminOfficeRoles.forEach((adminOfficeRole) => {
      if (adminOfficeRole.roleType === role) {
        rolesArray.push(adminOfficeRole);
        roleNotFound = false;
      }
    });
    if (roleNotFound) {
      throw new Error('Processing Error: Invalid roles provided for row.');
    }
  });

  // Return the validated array of user roles
  return rolesArray;
};

export const parsePrivileges = (privileges) => {
  if (!privileges) {
    throw new Error('Processing Error: Unable to parse privileges for row.');
  }
  const privilegesArray = [];

  // Parse privileges from string at ","
  const parsedPrivileges = privileges.split(',');
  parsedPrivileges.forEach((parsedPrivilege) => {
    let privilegeNotFound = true;
    // Remove any whitespace in the privilege string
    const privilege = parsedPrivilege.replace(/\s/g, '');
    officeUserPrivileges.forEach((officeUserPrivilege) => {
      if (officeUserPrivilege.privilegeType === privilege) {
        privilegesArray.push(officeUserPrivilege);
        privilegeNotFound = false;
      }
    });
    if (privilegeNotFound) {
      throw new Error('Processing Error: Invalid privileges provided for row.');
    }
  });

  // Return the validated array of user privileges
  return privilegesArray;
};

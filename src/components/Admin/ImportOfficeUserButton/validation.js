import { adminOfficeRoles } from 'constants/userRoles';

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

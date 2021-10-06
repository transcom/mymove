import React from 'react';
import { CheckboxGroupInput } from 'react-admin';

import { adminOfficeRoles } from 'constants/userRoles';

const makeRoleTypeArray = (roles) => {
  if (!roles || roles.length === 0) {
    return undefined;
  }
  return roles.reduce((rolesArray, role) => {
    if (role.roleType) {
      rolesArray.push(role.roleType);
    }
    return rolesArray;
  }, []);
};

const parseCheckboxInput = (input) => {
  return input.reduce((rolesArray, role) => {
    rolesArray.push(adminOfficeRoles.find((adminOfficeRole) => adminOfficeRole.roleType === role));
    return rolesArray;
  }, []);
};

const RolesCheckboxInput = (props) => (
  <CheckboxGroupInput
    source="roles"
    format={makeRoleTypeArray}
    parse={parseCheckboxInput}
    choices={adminOfficeRoles}
    optionValue="roleType"
    validate={props.validate}
  />
);

export { RolesCheckboxInput };

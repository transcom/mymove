import React from 'react';
import { CheckboxGroupInput } from 'react-admin';

const makeRoleTypeArray = roles => {
  if (!roles || roles.length === 0) {
    return;
  }
  return roles.reduce((rolesArray, role) => {
    rolesArray.push(role.roleType);
    return rolesArray;
  }, []);
};

const RolesCheckboxInput = props => (
  <CheckboxGroupInput
    source="roles"
    format={makeRoleTypeArray}
    choices={[
      { roleType: 'customer', name: 'Customer' },
      { roleType: 'transportation_ordering_officer', name: 'Transportation Ordering Officer' },
      { roleType: 'transportation_invoicing_officer', name: 'Transportation Invoicing Officer' },
      { roleType: 'contracting_officer', name: 'Contracting Officer' },
      { roleType: 'ppm_office_users', name: 'PPM Office Users' },
    ]}
    optionValue="roleType"
  />
);

export { RolesCheckboxInput };

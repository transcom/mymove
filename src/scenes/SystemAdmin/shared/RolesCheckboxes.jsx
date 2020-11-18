import React from 'react';
import { CheckboxGroupInput } from 'react-admin';

const choices = [
  { roleType: 'customer', name: 'Customer' },
  { roleType: 'transportation_ordering_officer', name: 'Transportation Ordering Officer' },
  { roleType: 'transportation_invoicing_officer', name: 'Transportation Invoicing Officer' },
  { roleType: 'contracting_officer', name: 'Contracting Officer' },
  { roleType: 'ppm_office_users', name: 'PPM Office Users' },
];

const makeRoleTypeArray = (roles) => {
  if (!roles || roles.length === 0) {
    return;
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
    rolesArray.push(choices.find((choice) => choice.roleType === role));
    return rolesArray;
  }, []);
};

const RolesCheckboxInput = (props) => (
  <CheckboxGroupInput
    source="roles"
    format={makeRoleTypeArray}
    parse={parseCheckboxInput}
    choices={choices}
    optionValue="roleType"
    validate={props.validate}
  />
);

export { RolesCheckboxInput };

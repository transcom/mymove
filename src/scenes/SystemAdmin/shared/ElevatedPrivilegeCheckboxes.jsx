import React from 'react';
import { CheckboxGroupInput } from 'react-admin';

import { officeUserPrivileges } from 'constants/userPrivileges';

const makePrivilegesArray = (privileges) => {
  if (!privileges || privileges.length === 0) {
    return undefined;
  }
  return privileges.reduce((privilegesArray, privilege) => {
    if (privilege.privilegeType) {
      privilegesArray.push(privilege.privilegeType);
    }
    return privilegesArray;
  }, []);
};

const parseCheckboxInput = (input) => {
  return input.reduce((privilegesArray, privilege) => {
    privilegesArray.push(
      officeUserPrivileges.find((officeUserPrivilege) => officeUserPrivilege.privilegeType === privilege),
    );
    return privilegesArray;
  }, []);
};

const PrivilegesCheckboxInput = (props) => (
  <CheckboxGroupInput
    source="privileges"
    format={makePrivilegesArray}
    parse={parseCheckboxInput}
    choices={officeUserPrivileges}
    optionValue="privilegeType"
    disabled={props.disabled}
  />
);

export { PrivilegesCheckboxInput };

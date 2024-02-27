import React from 'react';
import { CheckboxGroupInput } from 'react-admin';

import { elevatedPrivileges } from 'constants/elevatedPrivileges';

const makeElevatedPrivilegesArray = (elevatedPrivileges) => {
  if (!elevatedPrivileges || elevatedPrivileges.length === 0) {
    return undefined;
  }
  return elevatedPrivileges.reduce((elevatedPrivilegesArray, elevatedPrivilege) => {
    if (elevatedPrivilege.privilegeType) {
      elevatedPrivilegesArray.push(elevatedPrivilege.privilegeType);
    }
    return elevatedPrivilegesArray;
  }, []);
};

const parseCheckboxInput = (input) => {
  return input.reduce((elevatedPrivilegesArray, elevatedPrivilege) => {
    elevatedPrivilegesArray.push(
      elevatedPrivileges.find((elevatedPrivilegeType) => elevatedPrivilegeType.privilegeType === elevatedPrivilege),
    );
    return elevatedPrivilegesArray;
  }, []);
};

const ElevatedPrivilegesCheckboxInput = (props) => (
  <CheckboxGroupInput
    source="privileges"
    format={makeElevatedPrivilegesArray}
    parse={parseCheckboxInput}
    choices={elevatedPrivileges}
    optionValue="privilegeType"
  />
);

export { ElevatedPrivilegesCheckboxInput };

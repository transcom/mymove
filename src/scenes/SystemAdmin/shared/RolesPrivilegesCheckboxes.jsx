import React, { useState, useEffect } from 'react';
import { CheckboxGroupInput } from 'react-admin';

import { adminOfficeRoles, roleTypes } from 'constants/userRoles';
import { officeUserPrivileges, elevatedPrivilegeTypes } from 'constants/userPrivileges';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const RolesPrivilegesCheckboxInput = (props) => {
  let rolesSelected = [];
  let privilegesSelected = [];

  const [isHeadquartersRoleFF, setHeadquartersRoleFF] = useState(false);

  useEffect(() => {
    isBooleanFlagEnabled('headquarters_role')?.then((enabled) => {
      setHeadquartersRoleFF(enabled);
    });
  }, []);

  const makeRoleTypeArray = (roles) => {
    if (!roles || roles.length === 0) {
      rolesSelected = [];
      return undefined;
    }
    return roles.reduce((rolesArray, role) => {
      if (role.roleType) {
        if (isHeadquartersRoleFF || (!isHeadquartersRoleFF && role.roleType !== roleTypes.HQ)) {
          rolesArray.push(role.roleType);
        }
      }

      rolesSelected = rolesArray;
      return rolesArray;
    }, []);
  };

  const parseRolesCheckboxInput = (input) => {
    if (
      privilegesSelected.includes(elevatedPrivilegeTypes.SUPERVISOR) ||
      privilegesSelected.includes(elevatedPrivilegeTypes.SAFETY)
    ) {
      var index;
      if (input.includes(roleTypes.CUSTOMER)) {
        index = input.indexOf(roleTypes.CUSTOMER);
        if (index !== -1) {
          input.splice(index, 1);
        }
      }
      if (input.includes(roleTypes.CONTRACTING_OFFICER)) {
        index = input.indexOf(roleTypes.CONTRACTING_OFFICER);
        if (index !== -1) {
          input.splice(index, 1);
        }
      }
    }

    if (!isHeadquartersRoleFF && input.includes(roleTypes.HQ)) {
      if (input.includes(roleTypes.HQ)) {
        index = input.indexOf(roleTypes.HQ);
        if (index !== -1) {
          input.splice(index, 1);
        }
      }
    } else if (isHeadquartersRoleFF && privilegesSelected.includes(elevatedPrivilegeTypes.SAFETY)) {
      if (input.includes(roleTypes.HQ)) {
        index = input.indexOf(roleTypes.HQ);
        if (index !== -1) {
          input.splice(index, 1);
        }
      }
    }

    return input.reduce((rolesArray, role) => {
      rolesArray.push(adminOfficeRoles.find((adminOfficeRole) => adminOfficeRole.roleType === role));
      return rolesArray;
    }, []);
  };

  const makePrivilegesArray = (privileges) => {
    if (!privileges || privileges.length === 0) {
      privilegesSelected = [];
      return undefined;
    }
    return privileges.reduce((privilegesArray, privilege) => {
      if (privilege.privilegeType) {
        privilegesArray.push(privilege.privilegeType);
      }

      privilegesSelected = privilegesArray;
      return privilegesArray;
    }, []);
  };

  const parsePrivilegesCheckboxInput = (input) => {
    if (rolesSelected.includes(roleTypes.CUSTOMER) || rolesSelected.includes(roleTypes.CONTRACTING_OFFICER)) {
      var index;
      if (input.includes(elevatedPrivilegeTypes.SUPERVISOR)) {
        index = input.indexOf(elevatedPrivilegeTypes.SUPERVISOR);
        if (index !== -1) {
          input.splice(index, 1);
        }
      }

      if (input.includes(elevatedPrivilegeTypes.SAFETY)) {
        index = input.indexOf(elevatedPrivilegeTypes.SAFETY);
        if (index !== -1) {
          input.splice(index, 1);
        }
      }
    }

    if (isHeadquartersRoleFF && rolesSelected.includes(roleTypes.HQ)) {
      if (input.includes(elevatedPrivilegeTypes.SAFETY)) {
        index = input.indexOf(elevatedPrivilegeTypes.SAFETY);
        if (index !== -1) {
          input.splice(index, 1);
        }
      }
    }

    return input.reduce((privilegesArray, privilege) => {
      privilegesArray.push(
        officeUserPrivileges.find((officeUserPrivilege) => officeUserPrivilege.privilegeType === privilege),
      );
      return privilegesArray;
    }, []);
  };

  return (
    <>
      <CheckboxGroupInput
        source="roles"
        format={makeRoleTypeArray}
        parse={parseRolesCheckboxInput}
        choices={adminOfficeRoles}
        optionValue="roleType"
        validate={props.validate}
      />

      <CheckboxGroupInput
        source="privileges"
        format={makePrivilegesArray}
        parse={parsePrivilegesCheckboxInput}
        choices={officeUserPrivileges}
        optionValue="privilegeType"
      />
      <span style={{ marginTop: '-20px', marginBottom: '20px', fontWeight: 'bold' }}>
        Privileges cannot be selected with Customer or Contracting Officer roles.
      </span>
    </>
  );
};

export { RolesPrivilegesCheckboxInput };

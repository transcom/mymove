import React, { useState, useEffect } from 'react';
import { CheckboxGroupInput } from 'react-admin';

import { adminOfficeRoles } from 'constants/userRoles';
import { officeUserPrivileges } from 'constants/userPrivileges';
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
        rolesArray.push(role.roleType);
      }

      rolesSelected = rolesArray;
      return rolesArray;
    }, []);
  };

  const parseRolesCheckboxInput = (input) => {
    if (privilegesSelected.includes('supervisor') || privilegesSelected.includes('safety')) {
      var index;
      if (input.includes('customer')) {
        index = input.indexOf('customer');
        if (index !== -1) {
          input.splice(index, 1);
        }
      }
      if (input.includes('contracting_officer')) {
        index = input.indexOf('contracting_officer');
        if (index !== -1) {
          input.splice(index, 1);
        }
      }
    }

    if (isHeadquartersRoleFF && privilegesSelected.includes('safety')) {
      if (input.includes('headquarters')) {
        index = input.indexOf('headquarters');
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
    if (rolesSelected.includes('customer') || rolesSelected.includes('contracting_officer')) {
      var index;
      if (input.includes('supervisor')) {
        index = input.indexOf('supervisor');
        if (index !== -1) {
          input.splice(index, 1);
        }
      }

      if (input.includes('safety')) {
        index = input.indexOf('safety');
        if (index !== -1) {
          input.splice(index, 1);
        }
      }
    }

    if (isHeadquartersRoleFF && rolesSelected.includes('headquarters')) {
      if (input.includes('safety')) {
        index = input.indexOf('safety');
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

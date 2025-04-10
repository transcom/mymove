import React, { useState, useEffect } from 'react';
import { CheckboxGroupInput } from 'react-admin';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { useRolesPrivilegesQueries } from 'hooks/queries';

import { roleTypes } from 'constants/userRoles';
import { elevatedPrivilegeTypes } from 'constants/userPrivileges';

const RolesPrivilegesCheckboxInput = (props) => {
  const { adminUser, validate } = props;
  const { result } = useRolesPrivilegesQueries();
  const { rolesWithPrivs, privileges, nonSafetyRoles } = result;
  const [isHeadquartersRoleFF, setHeadquartersRoleFF] = useState(false);

  let rolesSelected = [];
  let privilegesSelected = [];

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
    // If the user selects a role that isn't allowed to have a Safety Moves privilege, remove their selection.
    nonSafetyRoles.forEach((nonSafetyRole) => {
      let index = input.indexOf(nonSafetyRole.roleType);
      if (privilegesSelected.includes(elevatedPrivilegeTypes.SAFETY)) {
        while (index !== -1) {
          input.splice(index, 1);
          index = input.indexOf(nonSafetyRole.roleType);
        }
      }

      if (!isHeadquartersRoleFF && input.includes(roleTypes.HQ)) {
        if (input.includes(roleTypes.HQ)) {
          index = input.indexOf(roleTypes.HQ);
          if (index !== -1) {
            input.splice(index, 1);
          }
        }
      }
    });
    return input.reduce((rolesArray, role) => {
      rolesArray.push(rolesWithPrivs.find((r) => r.roleType === role));
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
    // If the user selects the Safety Moves privilege for a role that isn't allowed, remove their selection.
    nonSafetyRoles.forEach((nonSafetyRole) => {
      let index = input.indexOf(elevatedPrivilegeTypes.SAFETY);
      if (rolesSelected.includes(nonSafetyRole.roleType)) {
        while (index !== -1) {
          input.splice(index, 1);
          index = input.indexOf(elevatedPrivilegeTypes.SAFETY);
        }
      }
    });
    return input.reduce((privilegesArray, privilege) => {
      privilegesArray.push(privileges.find((officeUserPrivilege) => officeUserPrivilege.privilegeType === privilege));
      return privilegesArray;
    }, []);
  };
  // filter the privileges to exclude the Safety Moves checkbox if the admin user is NOT a super admin
  const filteredPrivileges = privileges.filter((privilege) => {
    if (privilege.privilegeType === elevatedPrivilegeTypes.SAFETY && !adminUser?.super) {
      return false;
    }
    return true;
  });
  return (
    <>
      <CheckboxGroupInput
        source="roles"
        format={makeRoleTypeArray}
        parse={parseRolesCheckboxInput}
        choices={rolesWithPrivs}
        optionValue="roleType"
        optionText="roleName"
        validate={validate}
      />
      <CheckboxGroupInput
        source="privileges"
        format={makePrivilegesArray}
        parse={parsePrivilegesCheckboxInput}
        choices={filteredPrivileges}
        optionValue="privilegeType"
        optionText="privilegeName"
      />
      <span style={{ marginTop: '-20px', marginBottom: '20px', fontWeight: 'bold', whiteSpace: 'pre-wrap' }}>
        The Safety Moves privilege can only be selected for the following roles: Task Ordering Officer, Task Invoicing
        Officer, Services Counselor, Quality Assurance Evaluator, Customer Service Representative, and Headquarters.
      </span>
    </>
  );
};
export { RolesPrivilegesCheckboxInput };

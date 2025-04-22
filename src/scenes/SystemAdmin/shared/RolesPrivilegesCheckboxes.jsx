import React, { useState, useEffect } from 'react';
import { CheckboxGroupInput } from 'react-admin';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { useRolesPrivilegesQueries } from 'hooks/queries';
import { roleTypes } from 'constants/userRoles';
import { elevatedPrivilegeTypes } from 'constants/userPrivileges';

const RolesPrivilegesCheckboxInput = (props) => {
  const { adminUser, validate } = props;
  const { result } = useRolesPrivilegesQueries();
  const { rolesWithPrivs, privileges } = result;
  const [isHeadquartersRoleFF, setHeadquartersRoleFF] = useState(false);

  let rolesSelected = [];
  let privilegesSelected = [];

  useEffect(() => {
    isBooleanFlagEnabled('headquarters_role')?.then((enabled) => {
      setHeadquartersRoleFF(enabled);
    });
  }, []);

  // Make an array of roles that don't have Safety Moves privileges.
  const nonSafetyRoles = rolesWithPrivs.filter(
    (roleObj) => !roleObj.allowedPrivileges.includes(elevatedPrivilegeTypes.SAFETY),
  );

  // Make an array of roles that don't have Supervisor privileges.
  const nonSupervisorRoles = rolesWithPrivs.filter(
    (roleObj) => !roleObj.allowedPrivileges.includes(elevatedPrivilegeTypes.SUPERVISOR),
  );

  const availableRoles = rolesWithPrivs.filter((r) => r.roleType !== roleTypes.PRIME); // Do not want Prime role to show
  const rolesWithoutPrivs = availableRoles.filter((r) => r.allowedPrivileges.length === 0);
  const supervisorRoles = availableRoles
    .filter((r) => r.allowedPrivileges.includes(elevatedPrivilegeTypes.SUPERVISOR))
    .map((r) => r.roleName);
  const safetyRoles = availableRoles
    .filter((r) => r.allowedPrivileges.includes(elevatedPrivilegeTypes.SAFETY))
    .map((r) => r.roleName);

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

  // If the user selects a role that isn't allowed to have a Safety Moves privilege, remove their selection.
  const parseRolesCheckboxInput = (input) => {
    let result = [...input];

    if (!isHeadquartersRoleFF) {
      const idx = result.indexOf(roleTypes.HQ);
      if (idx !== -1) result.splice(idx, 1);
    }

    const disallowedByPrivilege = {
      [elevatedPrivilegeTypes.SAFETY]: nonSafetyRoles.map((r) => r.roleType),
      [elevatedPrivilegeTypes.SUPERVISOR]: nonSupervisorRoles.map((r) => r.roleType),
    };

    Object.entries(disallowedByPrivilege).forEach(([privType, badRoles]) => {
      if (privilegesSelected.includes(privType)) {
        badRoles.forEach((badRole) => {
          let idx;
          while ((idx = result.indexOf(badRole)) !== -1) {
            result.splice(idx, 1);
          }
        });
      }
    });

    return result.map((rt) => rolesWithPrivs.find((r) => r.roleType === rt));
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
    const result = [...input];

    const disallowedByRole = {
      [elevatedPrivilegeTypes.SAFETY]: nonSafetyRoles.map((r) => r.roleType),
      [elevatedPrivilegeTypes.SUPERVISOR]: nonSupervisorRoles.map((r) => r.roleType),
    };

    Object.entries(disallowedByRole).forEach(([privType, badRoles]) => {
      if (rolesSelected.some((role) => badRoles.includes(role))) {
        let idx;
        while ((idx = result.indexOf(privType)) !== -1) {
          result.splice(idx, 1);
        }
      }
    });

    return result.map((privType) => privileges.find((p) => p.privilegeType === privType));
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
        choices={availableRoles}
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
      {rolesWithoutPrivs.length > 0 && (
        <span style={{ marginTop: '-20px', marginBottom: '20px', fontWeight: 'bold' }}>
          The Supervisor privilege can only be selected for the following roles: {supervisorRoles.join(', ')}.
        </span>
      )}
      {safetyRoles.length > 0 && (
        <span style={{ marginTop: '-20px', marginBottom: '20px', fontWeight: 'bold', whiteSpace: 'pre-wrap' }}>
          The Safety Moves privilege can only be selected for the following roles: {safetyRoles.join(', ')}.
        </span>
      )}
    </>
  );
};
export { RolesPrivilegesCheckboxInput };

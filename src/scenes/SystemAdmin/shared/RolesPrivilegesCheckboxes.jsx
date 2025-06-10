import React from 'react';
import { CheckboxGroupInput } from 'react-admin';
import { useRolesPrivilegesQueries } from 'hooks/queries';
import { elevatedPrivilegeTypes } from 'constants/userPrivileges';

const RolesPrivilegesCheckboxInput = (props) => {
  const { adminUser, validate } = props;
  const { result } = useRolesPrivilegesQueries();
  const { rolesWithPrivs, privileges } = result;

  let rolesSelected = [];
  let privilegesSelected = [];

  const listFormatter = new Intl.ListFormat('en', {
    style: 'long',
    type: 'conjunction',
  });

  const availableRoles = rolesWithPrivs.filter((r) => r.roleType !== 'prime'); // Prime isn't an office role

  const allowedRolesByPrivilege = rolesWithPrivs.reduce((acc, role) => {
    role.allowedPrivileges.forEach((priv) => {
      if (!acc[priv]) acc[priv] = new Set();
      acc[priv].add(role.roleType);
    });
    return acc;
  }, {});

  const allowedPrivilegesByRole = rolesWithPrivs.reduce((acc, { roleType, allowedPrivileges }) => {
    acc[roleType] = new Set(allowedPrivileges);
    return acc;
  }, {});

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
    let result = [...input];

    privilegesSelected.forEach((privType) => {
      const allowed = allowedRolesByPrivilege[privType] || new Set();
      let i = 0;
      while (i < result.length) {
        if (!allowed.has(result[i])) {
          result.splice(i, 1);
        } else {
          i++;
        }
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

    rolesSelected.forEach((roleType) => {
      const allowed = allowedPrivilegesByRole[roleType] || new Set();
      let i = 0;
      while (i < result.length) {
        if (!allowed.has(result[i])) {
          result.splice(i, 1);
        } else {
          i++;
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

      {filteredPrivileges.map(({ privilegeType, privilegeName }) => {
        const allowedRoleTypes = Array.from(allowedRolesByPrivilege[privilegeType] || []);
        const roleNames = allowedRoleTypes
          .map((rt) => rolesWithPrivs.find((r) => r.roleType === rt)?.roleName)
          .filter((name) => name);

        if (roleNames.length === availableRoles.length) {
          return null;
        }
        return (
          <span
            key={privilegeType}
            style={{
              marginTop: '-20px',
              marginBottom: '20px',
              fontWeight: 'bold',
              whiteSpace: 'pre-wrap',
            }}
          >
            The {privilegeName} privilege can only be selected for the following roles:{' '}
            {listFormatter.format(roleNames)}.
          </span>
        );
      })}
    </>
  );
};
export { RolesPrivilegesCheckboxInput };

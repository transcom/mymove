import React, { useState, useEffect } from 'react';
import { CheckboxGroupInput } from 'react-admin';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { useRolesPrivilegesQueries } from 'hooks/queries';

const RolesPrivilegesCheckboxInput = (props) => {
  const { adminUser, validate } = props;
  const { result } = useRolesPrivilegesQueries();
  const { roles, privileges } = result;
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
        if (isHeadquartersRoleFF || (!isHeadquartersRoleFF && role.roleType !== 'headquarters')) {
          rolesArray.push(role.roleType);
        }
      }

      rolesSelected = rolesArray;
      return rolesArray;
    }, []);
  };

  const parseRolesCheckboxInput = (input) => {
    let index;

    if (privilegesSelected.includes('safety')) {
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
      if (input.includes('prime_simulator')) {
        index = input.indexOf('prime_simulator');
        if (index !== -1) {
          input.splice(index, 1);
        }
      }
      if (input.includes('gsr')) {
        index = input.indexOf('gsr');
        if (index !== -1) {
          input.splice(index, 1);
        }
      }
    }

    if (!isHeadquartersRoleFF && input.includes('headquarters')) {
      if (input.includes('headquarters')) {
        index = input.indexOf('headquarters');
        if (index !== -1) {
          input.splice(index, 1);
        }
      }
    }
    return input.reduce((rolesArray, role) => {
      rolesArray.push(roles.find((r) => r.roleType === role));
      return rolesArray;
    }, []);
  };

  const makePrivilegesArray = (privObjs) => {
    if (!privObjs || privObjs.length === 0) {
      privilegesSelected = [];
      return undefined;
    }

    return privObjs.reduce((privilegesArray, privilege) => {
      if (privilege.privilegeType) {
        privilegesArray.push(privilege.privilegeType);
      }

      privilegesSelected = privilegesArray;
      return privilegesArray;
    }, []);
  };

  const parsePrivilegesCheckboxInput = (input) => {
    let index;
    if (
      rolesSelected.includes('customer') ||
      rolesSelected.includes('contracting_officer') ||
      rolesSelected.includes('prime_simulator') ||
      rolesSelected.includes('gsr')
    ) {
      if (input.includes('safety')) {
        index = input.indexOf('safety');
        if (index !== -1) {
          input.splice(index, 1);
        }
      }
    }
    return input.reduce((privilegesArray, privilege) => {
      privilegesArray.push(privileges.find((officeUserPrivilege) => officeUserPrivilege.privilegeType === privilege));
      return privilegesArray;
    }, []);
  };
  // filter the privileges to exclude the Safety Moves checkbox if the admin user is NOT a super admin
  const filteredPrivileges = privileges.filter((privilege) => {
    if (privilege.privilegeType === 'safety' && !adminUser?.super) {
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
        choices={roles}
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

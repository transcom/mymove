import React from 'react';
import { connect } from 'react-redux';
import { useHistory } from 'react-router-dom';
import PropTypes from 'prop-types';
import { Button, GridContainer } from '@trussworks/react-uswds';

import { selectLoggedInUser } from 'store/entities/selectors';
import { setActiveRole as setActiveRoleAction } from 'store/auth/actions';
import { roleTypes } from 'constants/userRoles';
import { UserRolesShape } from 'types';
import getRoleTypesFromRoles from 'utils/user';

const SelectApplication = ({ userRoles, setActiveRole, activeRole }) => {
  const history = useHistory();

  const handleSelectRole = (roleType) => {
    setActiveRole(roleType);
    history.push('/');
  };

  const userRoleTypes = getRoleTypesFromRoles(userRoles);

  return (
    <GridContainer>
      <h2>Current role: {activeRole || userRoleTypes[0]}</h2>

      <ul className="usa-button-group">
        {[
          roleTypes.PPM,
          roleTypes.TOO,
          roleTypes.TIO,
          roleTypes.SERVICES_COUNSELOR,
          roleTypes.PRIME_SIMULATOR,
          roleTypes.QAE_CSR,
        ]
          .filter((r) => userRoleTypes.find((role) => r === role))
          .map((r) => (
            <li key={`selectRole_${r}`}>
              <Button
                type="button"
                onClick={() => {
                  handleSelectRole(r);
                }}
              >
                Select {r}
              </Button>
            </li>
          ))}
      </ul>
    </GridContainer>
  );
};

SelectApplication.propTypes = {
  activeRole: PropTypes.string,
  setActiveRole: PropTypes.func.isRequired,
  userRoles: UserRolesShape.isRequired,
};

SelectApplication.defaultProps = {
  activeRole: null,
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    activeRole: state.auth.activeRole,
    userRoles: user.roles || [],
  };
};

const mapDispatchToProps = {
  setActiveRole: setActiveRoleAction,
};

export const ConnectedSelectApplication = connect(mapStateToProps, mapDispatchToProps)(SelectApplication);

export default SelectApplication;

/*  react/jsx-one-expression-per-line */
import React from 'react';
import { connect } from 'react-redux';
import { useHistory } from 'react-router-dom';
import PropTypes from 'prop-types';
import { Button, GridContainer } from '@trussworks/react-uswds';

import { selectCurrentUser } from 'shared/Data/users';
import { setActiveRole as setActiveRoleAction } from 'store/auth/actions';
import { roleTypes } from 'constants/userRoles';
import { UserRolesShape } from 'types/index';

const SelectApplication = ({ userRoles, setActiveRole, activeRole }) => {
  const history = useHistory();

  const handleSelectRole = (roleType) => {
    setActiveRole(roleType);
    history.push('/');
  };

  return (
    <GridContainer>
      <h2>Current role: {activeRole || userRoles[0].roleType}</h2>

      <ul className="usa-button-group">
        {[roleTypes.PPM, roleTypes.TOO, roleTypes.TIO]
          .filter((r) => userRoles.find((role) => r === role.roleType))
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
  const user = selectCurrentUser(state);

  return {
    activeRole: state.auth.activeRole,
    userRoles: user.roles,
  };
};

const mapDispatchToProps = {
  setActiveRole: setActiveRoleAction,
};

export const ConnectedSelectApplication = connect(mapStateToProps, mapDispatchToProps)(SelectApplication);

export default SelectApplication;

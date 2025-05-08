import React, { useEffect, useState } from 'react';
import { connect, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import PropTypes from 'prop-types';
import { Button, GridContainer } from '@trussworks/react-uswds';

import { selectLoggedInUser } from 'store/entities/selectors';
import { selectIsSettingActiveRole } from 'store/auth/selectors';
import { setActiveRole as setActiveRoleAction } from 'store/auth/actions';
import { roleTypes } from 'constants/userRoles';
import { UserRolesShape } from 'types';
import getRoleTypesFromRoles from 'utils/user';

const SelectApplication = ({ userInactiveRoles, setActiveRole, activeRole }) => {
  const navigate = useNavigate();
  const isSettingActiveRole = useSelector(selectIsSettingActiveRole);
  const [pendingRole, setPendingRole] = useState(null);

  useEffect(() => {
    if (pendingRole !== null && !isSettingActiveRole) {
      // Pending role has been set and the auth saga
      // has received a response from the server.
      // This prevents a saga/action race condition between
      // index rendering on '/' and the SelectApplication component.
      // Previous race condition:
      // Select application requests to update the server session,
      // select application routes to index before hearing back,
      // index requests the current logged in user
      // the server begins handling the index with the old AppCtx session
      // the old AppCtx session has not been finished updating via SetActiveRole saga
      // then index is now rendering the old role, not the new role, thus a race condition
      setPendingRole(null);
      navigate('/');
    }
  }, [pendingRole, isSettingActiveRole, navigate]);

  const handleSelectRole = (roleType) => {
    setPendingRole(roleType);
    setActiveRole(roleType);
  };

  const userRoleTypes = getRoleTypesFromRoles(userInactiveRoles);

  return (
    <GridContainer>
      <h2>Current role: {activeRole || userRoleTypes[0]}</h2>

      <ul className="usa-button-group">
        {[
          roleTypes.HQ,
          roleTypes.TOO,
          roleTypes.TIO,
          roleTypes.SERVICES_COUNSELOR,
          roleTypes.PRIME_SIMULATOR,
          roleTypes.QAE,
          roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE,
          roleTypes.GSR,
          roleTypes.CONTRACTING_OFFICER,
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
  userInactiveRoles: UserRolesShape.isRequired,
};

SelectApplication.defaultProps = {
  activeRole: null,
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    activeRole: state.auth.activeRole,
    userInactiveRoles: user.inactiveRoles || [],
  };
};

const mapDispatchToProps = {
  setActiveRole: setActiveRoleAction,
};

export const ConnectedSelectApplication = connect(mapStateToProps, mapDispatchToProps)(SelectApplication);

export default SelectApplication;

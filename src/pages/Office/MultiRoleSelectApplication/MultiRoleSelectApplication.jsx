import React, { useEffect, useMemo, useState } from 'react';
import { connect, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from './MultiRoleSelectApplication.module.scss';

import { selectLoggedInUser } from 'store/entities/selectors';
import { setActiveRole as setActiveRoleAction } from 'store/auth/actions';
import { adminOfficeRoles } from 'constants/userRoles';
import { UserRolesShape } from 'types';
import getRoleTypesFromRoles from 'utils/user';
import { selectIsSettingActiveRole } from 'store/auth/selectors';

const selectStyle = styles['dropdown-style'];
const usaSelectStyle = 'usa-select';
const multiRoleWrapperStyle = styles['dropdown-wrapper-style'];
const multiRoleUlContainerStyle = styles['dropdown-ul-container-style'];

export const roleLookupValues = Object.fromEntries(
  adminOfficeRoles.map(({ roleType, name }) => [roleType, { roleType, name }]),
);

const EMPTY_ROLE = 'none';

const MultiRoleSelectApplication = ({ inactiveRoles, setActiveRole, activeRole }) => {
  const navigate = useNavigate();
  const isSettingActiveRole = useSelector(selectIsSettingActiveRole);

  const [pendingRole, setPendingRole] = useState(null);
  const [mainRole, setMainRole] = useState(activeRole);
  const [reset, setReset] = useState(false);
  useEffect(() => {
    if (!reset && pendingRole !== mainRole) {
      setPendingRole(mainRole);
      setTimeout(() => {
        setActiveRole(mainRole);
        setReset(true);
      }, 0);
    }
    if (reset && mainRole === activeRole) {
      setTimeout(() => {
        navigate('/');
        setReset(false);
      }, 5);
    }
  }, [pendingRole, mainRole, reset, activeRole, setActiveRole, isSettingActiveRole, navigate]);

  const assumedRoleType = (activeRole || inactiveRoles[0]) ?? EMPTY_ROLE;

  const userRoleTypes = getRoleTypesFromRoles(inactiveRoles);

  const [rolesAvailableToUser] = useMemo(() => {
    const lookup = Object.fromEntries([[assumedRoleType, assumedRoleType], ...userRoleTypes.map((e) => [e, e])]);
    const result = adminOfficeRoles.filter(({ roleType }) => lookup[roleType] === roleType);
    return [result];
  }, [userRoleTypes, assumedRoleType]);

  const handleSelectRole = ({ target: { value: roleType } }) => {
    if (typeof roleType === 'string') {
      setMainRole(roleType);
    }
  };

  const applicationOptions = useMemo(
    () =>
      assumedRoleType === EMPTY_ROLE ? (
        <option key={EMPTY_ROLE} value={EMPTY_ROLE}>
          no role
        </option>
      ) : (
        rolesAvailableToUser.map(({ roleType, abbv }) => (
          <option key={roleType} value={roleType} hidden={assumedRoleType === roleType}>
            {abbv}
          </option>
        ))
      ),
    [rolesAvailableToUser, assumedRoleType],
  );

  const hasSingleDropdownOption = inactiveRoles?.length === 0;

  const selectDescription = ((assumed) => {
    const roleName = roleLookupValues[assumed]?.name;
    switch (roleName) {
      case EMPTY_ROLE: {
        return 'combo box has no roles options.';
      }
      default: {
        return hasSingleDropdownOption
          ? `combo box is limited to the current role of ${roleName}.`
          : 'combo box with roles to switch to.';
      }
    }
  })(assumedRoleType);

  const selectId = 'role-select';
  const labelTextId = 'role-select-label';

  const selectDropdownContent = (
    <select
      key={assumedRoleType}
      id={selectId}
      aria-describedby={labelTextId}
      aria-label="User roles"
      className={classNames(selectStyle, usaSelectStyle)}
      defaultValue={assumedRoleType}
      onChange={handleSelectRole}
    >
      {applicationOptions}
    </select>
  );

  return (
    <label className={classNames(multiRoleUlContainerStyle, multiRoleWrapperStyle)}>
      <div id={labelTextId} aria-label={selectDescription} aria-hidden>
        Role:
      </div>
      {selectDropdownContent}
    </label>
  );
};

MultiRoleSelectApplication.propTypes = {
  activeRole: PropTypes.string,
  setActiveRole: PropTypes.func.isRequired,
  inactiveRoles: UserRolesShape.isRequired,
};

MultiRoleSelectApplication.defaultProps = {
  activeRole: null,
};

const mapStateToProps = (state) => {
  const user = selectLoggedInUser(state);

  return {
    activeRole: state.auth.activeRole,
    inactiveRoles: user.inactiveRoles || [],
  };
};

const mapDispatchToProps = {
  setActiveRole: setActiveRoleAction,
};

export const ConnectedSelectApplication = connect(mapStateToProps, mapDispatchToProps)(MultiRoleSelectApplication);

export default MultiRoleSelectApplication;

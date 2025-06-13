import React, { useEffect, useMemo, useState } from 'react';
import { connect, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { Dropdown, Link } from '@trussworks/react-uswds';

import styles from './MultiRoleSelectApplication.module.scss';

import { selectLoggedInUser } from 'store/entities/selectors';
import { setActiveRole as setActiveRoleAction } from 'store/auth/actions';
import { adminOfficeRoles } from 'constants/userRoles';
import { UserRolesShape } from 'types';
import getRoleTypesFromRoles from 'utils/user';
import { selectIsSettingActiveRole } from 'store/auth/selectors';

const selectStyle = styles['dropdown-style'];
const multiRoleWrapperStyle = styles['dropdown-wrapper-style'];
const multiRoleUlContainerStyle = styles['dropdown-ul-container-style'];

export const roleLookupValues = Object.fromEntries(
  adminOfficeRoles.map(({ roleType, name, abbv }) => [roleType, { roleType, name, abbv }]),
);

const roleOrder = [
  roleLookupValues.services_counselor,
  roleLookupValues.task_ordering_officer,
  roleLookupValues.task_invoicing_officer,
  roleLookupValues.qae,
  roleLookupValues.customer_service_representative,
  roleLookupValues.gsr,
  roleLookupValues.headquarters,
  roleLookupValues.contracting_officer,
  roleLookupValues.prime_simulator,
];

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
      (async () => {
        setActiveRole(mainRole);
        setReset(true);
      })();
    }
    if (reset && mainRole === activeRole) {
      (async () => {
        setTimeout(() => {
          navigate('/');
          setReset(false);
        }, 5);
      })();
    }
  }, [pendingRole, mainRole, reset, activeRole, setActiveRole, isSettingActiveRole, navigate]);

  const assumedRoleType = (activeRole || inactiveRoles[0]) ?? EMPTY_ROLE;

  const userRoleTypes = getRoleTypesFromRoles(inactiveRoles);

  const [rolesAvailableToUser] = useMemo(() => {
    const lookup = Object.fromEntries([[assumedRoleType, assumedRoleType], ...userRoleTypes.map((e) => [e, e])]);
    const result = roleOrder.filter(({ roleType }) => lookup[roleType] === roleType);
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
        <option aria-label="no role" key={EMPTY_ROLE} value={EMPTY_ROLE}>
          no role
        </option>
      ) : (
        rolesAvailableToUser.map(({ roleType, abbv, name }) => {
          return (
            <option aria-label={name} key={roleType} value={roleType}>
              {abbv}
            </option>
          );
        })
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
    <Dropdown
      key={assumedRoleType}
      id={selectId}
      aria-describedby={labelTextId}
      aria-label="User roles"
      className={classNames(selectStyle)}
      defaultValue={assumedRoleType}
      onChange={handleSelectRole}
    >
      {applicationOptions}
    </Dropdown>
  );

  const displayedContent =
    inactiveRoles?.length === 0 ? (
      <Link to="/">{roleLookupValues[assumedRoleType]?.abbv || EMPTY_ROLE}</Link>
    ) : (
      selectDropdownContent
    );

  return (
    <label className={classNames(multiRoleUlContainerStyle, multiRoleWrapperStyle)}>
      <div id={labelTextId} aria-label={selectDescription} aria-hidden>
        Role:
      </div>
      {displayedContent}
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
    activeRole: state?.auth?.activeRole,
    inactiveRoles: user?.inactiveRoles || [],
  };
};

const mapDispatchToProps = {
  setActiveRole: setActiveRoleAction,
};

export const ConnectedSelectApplication = connect(mapStateToProps, mapDispatchToProps)(MultiRoleSelectApplication);

export default MultiRoleSelectApplication;

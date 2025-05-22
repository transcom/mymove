import React, { useMemo, useState } from 'react';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from './MultiRoleSelectApplication.module.scss';

import { selectLoggedInUser } from 'store/entities/selectors';
import { setActiveRole as setActiveRoleAction } from 'store/auth/actions';
import { adminOfficeRoles } from 'constants/userRoles';
import { UserRolesShape } from 'types';
import getRoleTypesFromRoles from 'utils/user';

const selectStyle = styles['dropdown-style'];
const usaSelectStyle = 'usa-select';
const multiRoleWrapperStyle = styles['dropdown-wrapper-style'];
const multiRoleUlContainerStyle = styles['dropdown-ul-container-style'];

export const roleLookupValues = Object.fromEntries(
  adminOfficeRoles.map(({ roleType, name }) => [roleType, { roleType, name }]),
);

const MultiRoleSelectApplication = ({ userRoles, setActiveRole, activeRole }) => {
  const navigate = useNavigate();

  const [userRoleTypes] = useState(getRoleTypesFromRoles(userRoles));

  const assumedRole = activeRole || userRoleTypes[0];

  const [rolesAvailableToUser] = useMemo(() => {
    const lookup = Object.fromEntries(userRoleTypes.map((e) => [e, e]));
    const result = adminOfficeRoles.filter(({ roleType }) => lookup[roleType] === roleType);
    return [result];
  }, [userRoleTypes]);

  const handleSelectRole = ({ target: { value: roleType } }) => {
    if (typeof roleType === 'string') {
      setActiveRole(roleType);
      navigate('/');
    }
  };

  const applicationOptions = useMemo(
    () =>
      rolesAvailableToUser.map(({ roleType, name }) => (
        <option key={roleType} value={roleType}>
          {name}
        </option>
      )),
    [rolesAvailableToUser],
  );

  const hasSingleDropdownOption = userRoles.length <= 1;

  const selectDescription = hasSingleDropdownOption
    ? `combo box is limited to the current role of ${roleLookupValues[assumedRole].name}.`
    : 'combo box with roles to switch to.';

  const selectId = 'role-select';
  const labelTextId = 'role-select-label';

  /* eslint-disable react/jsx-props-no-spreading */
  const selectDropdownContent = (
    <select
      id={selectId}
      aria-describedby={labelTextId}
      aria-label="User roles"
      className={classNames(selectStyle, usaSelectStyle)}
      defaultValue={assumedRole}
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
  userRoles: UserRolesShape.isRequired,
};

MultiRoleSelectApplication.defaultProps = {
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

export const ConnectedSelectApplication = connect(mapStateToProps, mapDispatchToProps)(MultiRoleSelectApplication);

export default MultiRoleSelectApplication;

/* eslint-disable import/prefer-default-export */
import PropTypes from 'prop-types';

export const UserRoleShape = PropTypes.shape({
  roleType: PropTypes.string,
});

export const UserRolesShape = PropTypes.arrayOf(UserRoleShape);

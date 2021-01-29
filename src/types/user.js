/* eslint-disable import/prefer-default-export */
import PropTypes from 'prop-types';

import { AddressShape } from './address';

export const UserRoleShape = PropTypes.shape({
  roleType: PropTypes.string.isRequired,
});

export const UserRolesShape = PropTypes.arrayOf(UserRoleShape);

export const TransportationOfficeShape = PropTypes.shape({
  address: AddressShape,
  gbloc: PropTypes.string,
  name: PropTypes.string,
  phone_lines: PropTypes.string,
});

export const OfficeUserInfoShape = PropTypes.shape({
  email: PropTypes.string,
  first_name: PropTypes.string,
  last_name: PropTypes.string,
  telephone: PropTypes.string,
  transportation_office: TransportationOfficeShape,
});

export const UserShape = PropTypes.shape({
  first_name: PropTypes.string,
  last_name: PropTypes.string,
  office_user: OfficeUserInfoShape,
});

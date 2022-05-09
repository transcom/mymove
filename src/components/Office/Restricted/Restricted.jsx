import React from 'react';
import PropTypes from 'prop-types';

import PermissionContext from './PermissionContext';

// Render the children if the user is allowed to access the permission, otherwise render the fallback or null.
// This component is meant to be used everywhere a restriction based on user permission is needed
const Restricted = ({ to, fallback, children }) => {
  // get isAllowedTo function from context (it knows about the permissions of the user)
  const { isAllowedTo } = React.useContext(PermissionContext);

  // If the user has that permission, render the children
  if (isAllowedTo(to)) return children;

  // if provided, render the fallback
  if (fallback) return fallback;

  // otherwise dont render anything
  return null;
};

Restricted.propTypes = {
  to: PropTypes.string.isRequired,
  fallback: PropTypes.node,
  children: PropTypes.node.isRequired,
};

Restricted.defaultProps = {
  fallback: null,
};

export default Restricted;

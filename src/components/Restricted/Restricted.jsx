import React from 'react';
import PropTypes from 'prop-types';

import PermissionContext from './PermissionContext';

// Render the children if the user is allowed to access based on permissions and user restriction, otherwise render the fallback or null.
// This component is meant to be used everywhere a restriction based on user permission is needed
const Restricted = ({ to, fallback, user, children }) => {
  // get isAllowedTo function from context (it knows about the permissions of the user)
  const context = React.useContext(PermissionContext);

  if (context === undefined) {
    throw new Error(`Restricted must be used within a PermissionProvider`);
  }

  const { isAllowedTo } = context;

  // If the user has that permission, render the children
  if (isAllowedTo(to, user)) return children;

  // if provided, render the fallback
  if (fallback) return fallback;

  // otherwise dont render anything
  return null;
};

Restricted.propTypes = {
  to: PropTypes.string,
  fallback: PropTypes.node,
  children: PropTypes.node.isRequired,
  user: PropTypes.string,
};

Restricted.defaultProps = {
  fallback: null,
};

export default Restricted;

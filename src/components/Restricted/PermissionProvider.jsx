import React, { useMemo } from 'react';
import PropTypes from 'prop-types';

import PermissionContext from './PermissionContext';

const PermissionProvider = ({ permissions, currentUserId, children }) => {
  // Creates a memoized object with a 'isAllowedTo' function that returns whether the requested permission is available in the list of permissions passed as parameter
  const isAllowedTo = useMemo(
    () => ({
      isAllowedTo: (permission, userId) => {
        // If access is restricted to specific permissions, is the permission available for the current user?
        const permissionGranted = !permission || permissions.filter((p) => p === permission).length > 0;

        // If access is restricted to a specific user, is it the current user?
        const userAllowed = !userId || userId === currentUserId;

        return permissionGranted && userAllowed;
      },
    }),
    [permissions, currentUserId],
  );

  return <PermissionContext.Provider value={isAllowedTo}>{children}</PermissionContext.Provider>;
};

PermissionProvider.propTypes = {
  permissions: PropTypes.arrayOf(PropTypes.string).isRequired,
  currentUserId: PropTypes.string,
  children: PropTypes.node.isRequired,
};

PermissionProvider.defaultProps = {
  currentUserId: null,
};

export default PermissionProvider;

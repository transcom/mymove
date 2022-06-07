import React, { useMemo } from 'react';
import PropTypes from 'prop-types';

import PermissionContext from './PermissionContext';

const PermissionProvider = ({ permissions, children }) => {
  // Creates a memoized object with a 'isAllowedTo' function that returns whether the requested permission is available in the list of permissions passed as parameter
  const isAllowedTo = useMemo(
    () => ({
      isAllowedTo: (permission) => permissions.filter((p) => p === permission).length > 0,
    }),
    [permissions],
  );

  return <PermissionContext.Provider value={isAllowedTo}>{children}</PermissionContext.Provider>;
};

PermissionProvider.propTypes = {
  permissions: PropTypes.arrayOf(PropTypes.string).isRequired,
  children: PropTypes.node.isRequired,
};

export default PermissionProvider;

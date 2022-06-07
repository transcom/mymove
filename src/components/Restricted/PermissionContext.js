import React from 'react';

// Default behaviour is to deny access
const defaultBehaviour = {
  isAllowedTo: () => false,
};

// Create permission context
const PermissionContext = React.createContext(defaultBehaviour);
PermissionContext.displayName = 'PermissionContext';

export default PermissionContext;

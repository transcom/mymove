import React from 'react';

// Default behaviour is to deny access
const defaultBehaviour = {
  isAllowedTo: () => false,
};

// Create permission context
const PermissionContext = React.createContext(defaultBehaviour);

export default PermissionContext;

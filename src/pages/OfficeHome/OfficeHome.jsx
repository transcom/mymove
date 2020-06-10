/* eslint-disable react/jsx-props-no-spreading */
import React, { lazy } from 'react';
import PropTypes from 'prop-types';

import { roleTypes } from 'constants/userRoles';

const Queues = lazy(() => import('scenes/Office/Queues'));
const TOO = lazy(() => import('scenes/Office/TOO/too'));
const TIO = lazy(() => import('scenes/Office/TIO/tio'));

const OfficeHome = ({ userRoles, activeRole, ...props }) => {
  // Prefers activeRole value from Redux, otherwise defaults to first role in array
  const selectedRole = activeRole || userRoles[0].roleType;

  switch (selectedRole) {
    case roleTypes.PPM:
      return <Queues queueType="new" {...props} />;
    case roleTypes.TIO:
      return <TIO {...props} />;
    case roleTypes.TOO:
      return <TOO {...props} />;
    default:
      // User should have been redirected, doesn't have access to the office
      return <div />;
  }
};

OfficeHome.displayName = 'OfficeHome';

OfficeHome.propTypes = {
  activeRole: PropTypes.string,
  userRoles: PropTypes.arrayOf(
    PropTypes.shape({
      roleType: PropTypes.string,
    }),
  ),
};

export default OfficeHome;

/* eslint-disable react/jsx-props-no-spreading */
import React, { lazy } from 'react';
import PropTypes from 'prop-types';

import { roleTypes } from 'constants/userRoles';

const Queues = lazy(() => import('scenes/Office/Queues'));
const TOO = lazy(() => import('scenes/Office/TOO/too'));
const TIO = lazy(() => import('scenes/Office/TIO/tio'));

const OfficeHome = ({ userRoles, ...props }) => {
  // Defaults to first role in array
  // TODO - use Redux to store activeRole value and use that instead
  const [selectedRole] = userRoles;

  switch (selectedRole.roleType) {
    case roleTypes.PPM:
      return <Queues queueType="new" {...props} />;
    case roleTypes.TIO:
      return <TIO {...props} />;
    case roleTypes.TOO:
      return <TOO {...props} />;
    default:
  }

  return <div>Office home</div>;
};

OfficeHome.propTypes = {
  userRoles: PropTypes.arrayOf(
    PropTypes.shape({
      roleType: PropTypes.string,
    }),
  ),
};

export default OfficeHome;

import React from 'react';
import PropTypes from 'prop-types';
import { Tag } from '@trussworks/react-uswds';

const LeftNavTag = ({ children, showTag, className }) => {
  if (!showTag) return null;
  return <Tag className={className}>{children}</Tag>;
};

LeftNavTag.propTypes = {
  children: PropTypes.node.isRequired,
  showTag: PropTypes.bool.isRequired,
  className: PropTypes.string,
};

LeftNavTag.defaultProps = {
  className: '',
};

export default LeftNavTag;

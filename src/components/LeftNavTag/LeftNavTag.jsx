import React from 'react';
import PropTypes from 'prop-types';
import { Tag } from '@trussworks/react-uswds';

const LeftNavTag = ({ children, showTag, className, background }) => {
  if (!showTag) return null;
  return (
    <Tag background={background} className={className}>
      {children}
    </Tag>
  );
};

LeftNavTag.propTypes = {
  children: PropTypes.node.isRequired,
  showTag: PropTypes.bool.isRequired,
  className: PropTypes.string,
  background: PropTypes.string,
};

LeftNavTag.defaultProps = {
  className: '',
  background: '',
};

export default LeftNavTag;

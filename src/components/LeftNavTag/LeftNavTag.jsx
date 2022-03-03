import React from 'react';
import PropTypes from 'prop-types';
import { Tag } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './LeftNavTag.module.scss';

const LeftNavTag = ({ children, showTag, className, background, testID }) => {
  if (!showTag) return null;
  return (
    <Tag background={background} className={classnames(styles.LeftNavTag, className)} data-testid={testID}>
      {children}
    </Tag>
  );
};

LeftNavTag.propTypes = {
  children: PropTypes.node.isRequired,
  showTag: PropTypes.bool.isRequired,
  className: PropTypes.string,
  background: PropTypes.string,
  testID: PropTypes.string,
};

LeftNavTag.defaultProps = {
  className: '',
  background: '',
  testID: '',
};

export default LeftNavTag;

import React from 'react';
import classNames from 'classnames';
import * as PropTypes from 'prop-types';

import styles from './ServiceItemContainer.module.scss';

const ServiceItemContainer = ({ className, children }) => {
  const containerClasses = classNames(styles.serviceItemContainer, 'container--accent--default', className);

  return (
    <div data-testid="ServiceItemContainer" className={`${containerClasses}`}>
      {children}
    </div>
  );
};

ServiceItemContainer.propTypes = {
  className: PropTypes.string,
  children: PropTypes.node.isRequired,
};

ServiceItemContainer.defaultProps = {
  className: '',
};

export default ServiceItemContainer;

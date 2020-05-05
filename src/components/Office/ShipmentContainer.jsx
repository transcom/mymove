import React from 'react';
import classNames from 'classnames/bind';
import * as PropTypes from 'prop-types';
import styles from './shipmentContainer.module.scss';

const cx = classNames.bind(styles);

const ShipmentContainer = ({ children, containerType }) => {
  const containerClasses = cx({
    container: true,
    'shipment-container': true,
    'container--accent--hhg': containerType === 'HHG',
    'container--accent--nts': containerType === 'NTS',
  });

  return <div className={`${containerClasses}`}>{children}</div>;
};

ShipmentContainer.propTypes = {
  children: PropTypes.element,
  /** Describes the type of shipment container. */
  containerType: PropTypes.oneOf(['HHG', 'NTS']),
};

ShipmentContainer.defaultProps = {
  containerType: 'HHG',
};

export default ShipmentContainer;

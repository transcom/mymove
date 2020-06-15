import React from 'react';
import classNames from 'classnames/bind';
import * as PropTypes from 'prop-types';

import styles from './shipmentContainer.module.scss';

import { SHIPMENT_TYPE } from 'shared/constants';

const cx = classNames.bind(styles);

const ShipmentContainer = ({ className, children, shipmentType }) => {
  const containerClasses = cx('container', 'shipment-container', {
    'container--accent--hhg':
      shipmentType === SHIPMENT_TYPE.HHG ||
      shipmentType === SHIPMENT_TYPE.HHG_SHORTHAUL_DOMESTIC ||
      shipmentType === SHIPMENT_TYPE.HHG_LONGHAUL_DOMESTIC,
    'container--accent--nts': shipmentType === SHIPMENT_TYPE.NTS,
  });

  return <div className={`${containerClasses} ${className}`}>{children}</div>;
};

ShipmentContainer.propTypes = {
  className: PropTypes.string,
  children: PropTypes.element.isRequired,
  /** Describes the type of shipment container. */
  shipmentType: PropTypes.oneOf([
    SHIPMENT_TYPE.HHG,
    SHIPMENT_TYPE.HHG_SHORTHAUL_DOMESTIC,
    SHIPMENT_TYPE.HHG_LONGHAUL_DOMESTIC,
    SHIPMENT_TYPE.NTS,
  ]),
};

ShipmentContainer.defaultProps = {
  shipmentType: SHIPMENT_TYPE.HHG,
  className: '',
};

export default ShipmentContainer;

import React from 'react';
import classNames from 'classnames/bind';
import * as PropTypes from 'prop-types';

import styles from './shipmentContainer.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const cx = classNames.bind(styles);

const ShipmentContainer = ({ className, children, shipmentType }) => {
  const containerClasses = cx('container', 'shipment-container', {
    'container--accent--hhg':
      shipmentType === SHIPMENT_OPTIONS.HHG ||
      shipmentType === SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC ||
      shipmentType === SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    'container--accent--nts': shipmentType === SHIPMENT_OPTIONS.NTS,
  });

  return <div className={`${containerClasses} ${className}`}>{children}</div>;
};

ShipmentContainer.propTypes = {
  className: PropTypes.string,
  children: PropTypes.element.isRequired,
  /** Describes the type of shipment container. */
  shipmentType: PropTypes.oneOf([
    SHIPMENT_OPTIONS.HHG,
    SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.NTS,
  ]),
};

ShipmentContainer.defaultProps = {
  shipmentType: SHIPMENT_OPTIONS.HHG,
  className: '',
};

export default ShipmentContainer;

import React from 'react';
import classNames from 'classnames';
import * as PropTypes from 'prop-types';

import styles from './ShipmentContainer.module.scss';

import { SHIPMENT_TYPE } from 'shared/constants';
import { ShipmentTypeOneOf } from 'types/shipment';

const ShipmentContainer = ({ className, children, shipmentType }) => {
  const containerClasses = classNames(
    'container',
    styles.shipmentContainer,
    {
      'container--accent--hhg':
        shipmentType === SHIPMENT_TYPE.HHG ||
        shipmentType === SHIPMENT_TYPE.HHG_SHORTHAUL_DOMESTIC ||
        shipmentType === SHIPMENT_TYPE.HHG_LONGHAUL_DOMESTIC,
      'container--accent--nts': shipmentType === SHIPMENT_TYPE.NTS,
    },
    className,
  );

  return <div className={`${containerClasses}`}>{children}</div>;
};

ShipmentContainer.propTypes = {
  className: PropTypes.string,
  children: PropTypes.element.isRequired,
  /** Describes the type of shipment container. */
  shipmentType: ShipmentTypeOneOf,
};

ShipmentContainer.defaultProps = {
  shipmentType: SHIPMENT_TYPE.HHG,
  className: '',
};

export default ShipmentContainer;

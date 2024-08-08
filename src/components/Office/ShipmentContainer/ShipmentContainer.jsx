import React from 'react';
import classNames from 'classnames';
import * as PropTypes from 'prop-types';

import styles from './ShipmentContainer.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ShipmentOptionsOneOf } from 'types/shipment';

const ShipmentContainer = ({ id, className, children, shipmentType }) => {
  const containerClasses = classNames(
    styles.shipmentContainer,
    {
      'container--accent--default':
        shipmentType === null ||
        shipmentType === SHIPMENT_OPTIONS.BOAT_HAUL_AWAY ||
        shipmentType === SHIPMENT_OPTIONS.BOAT_TOW_AWAY ||
        shipmentType === SHIPMENT_OPTIONS.MOBILE_HOME ||
        !Object.values(SHIPMENT_OPTIONS).includes(shipmentType),
      'container--accent--hhg': shipmentType === SHIPMENT_OPTIONS.HHG,
      'container--accent--nts': shipmentType === SHIPMENT_OPTIONS.NTS,
      'container--accent--ntsr': shipmentType === SHIPMENT_OPTIONS.NTSR,
      'container--accent--ppm': shipmentType === SHIPMENT_OPTIONS.PPM,
    },
    className,
  );

  return (
    <div data-testid="ShipmentContainer" className={`${containerClasses}`} id={id}>
      {children}
    </div>
  );
};

ShipmentContainer.propTypes = {
  id: PropTypes.string,
  className: PropTypes.string,
  children: PropTypes.node.isRequired,
  /** Describes the type of shipment container. */
  shipmentType: ShipmentOptionsOneOf,
};

ShipmentContainer.defaultProps = {
  shipmentType: null,
  className: '',
  id: '',
};

export default ShipmentContainer;

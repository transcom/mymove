import React from 'react';
import classNames from 'classnames';
import * as PropTypes from 'prop-types';

import styles from './ShipmentContainer.module.scss';

import { SHIPMENT_OPTIONS, MOVE_TYPES } from 'shared/constants';
import { ShipmentOptionsOneOf } from 'types/shipment';

const ShipmentContainer = ({ className, children, shipmentType }) => {
  const containerClasses = classNames(
    styles.shipmentContainer,
    {
      'container--accent--default': shipmentType === null,
      'container--accent--hhg':
        shipmentType === SHIPMENT_OPTIONS.HHG ||
        shipmentType === SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC ||
        shipmentType === SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
      'container--accent--nts': shipmentType === SHIPMENT_OPTIONS.NTS || MOVE_TYPES.NTS,
      'container--accent--ntsr': shipmentType === SHIPMENT_OPTIONS.NTSR,
      'container--accent--ppm': shipmentType === MOVE_TYPES.PPM,
    },
    className,
  );

  return (
    <div data-testid="ShipmentContainer" className={`${containerClasses}`}>
      {children}
    </div>
  );
};

ShipmentContainer.propTypes = {
  className: PropTypes.string,
  children: PropTypes.node.isRequired,
  /** Describes the type of shipment container. */
  shipmentType: ShipmentOptionsOneOf,
};

ShipmentContainer.defaultProps = {
  shipmentType: null,
  className: '',
};

export default ShipmentContainer;

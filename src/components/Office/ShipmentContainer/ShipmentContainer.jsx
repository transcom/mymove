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
      'container--accent--default': shipmentType === null,
      'container--accent--hhg':
        shipmentType === SHIPMENT_OPTIONS.HHG ||
        shipmentType === SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC ||
        shipmentType === SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
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

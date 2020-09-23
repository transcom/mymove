import React from 'react';
import * as PropTypes from 'prop-types';

import ShipmentContainer from '../../../Office/ShipmentContainer';

import styles from './ShipmentCard.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const ShipmentCard = ({ shipmentType }) => {
  return (
    <div className={styles.ShipmentCard} data-testid="shipment-display">
      <ShipmentContainer className={styles.container} shipmentType={shipmentType}>
        <div className={styles.heading}>
          <h3>dsadsadsa</h3>
        </div>
        <dl>
          <div className={styles.row}>
            <dt>Requested move date</dt>
            <dd>Some date</dd>
          </div>
          <div className={styles.row}>
            <dt>Current address</dt>
            <dd>Some date</dd>
          </div>
          <div className={styles.row}>
            <dt className={styles.label}>Destination address</dt>
            <dd data-testid="shipmentDestinationAddress">some address</dd>
          </div>
        </dl>
      </ShipmentContainer>
    </div>
  );
};

ShipmentCard.propTypes = {
  shipmentType: PropTypes.oneOf([
    SHIPMENT_OPTIONS.PPM,
    SHIPMENT_OPTIONS.HHG,
    SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.NTS,
  ]),
};

ShipmentCard.defaultProps = {
  shipmentType: SHIPMENT_OPTIONS.PPM,
};

export default ShipmentCard;

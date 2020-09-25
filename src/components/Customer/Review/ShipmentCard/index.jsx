import React from 'react';
import * as PropTypes from 'prop-types';

import ShipmentContainer from '../../../Office/ShipmentContainer';

import styles from './ShipmentCard.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const ShipmentCard = ({ shipmentType }) => {
  return (
    <div className={styles.ShipmentCard} data-testid="shipment-display">
      <ShipmentContainer className={styles.container} shipmentType={shipmentType}>
        <div className={styles.ShipmentCardHeader}>
          <div>
            <h4>Shipment 1: PPM</h4>
            <p>#ABC123K-001</p>
          </div>
          <a href="#">Edit</a>
        </div>

        <dl>
          <div className={styles.row}>
            <dt>Expected departure</dt>
            <dd>26 Mar 2020</dd>
          </div>
          <div className={styles.row}>
            <dt>Starting ZIP</dt>
            <dd>78234</dd>
          </div>
          <div className={styles.row}>
            <dt>Storage (SIT)</dt>
            <dd>Yes, 14 days</dd>
          </div>
          <div className={styles.row}>
            <dt>Destination ZIP</dt>
            <dd>78111</dd>
          </div>
        </dl>
        <hr className={styles.divider} />
        <div className={styles['subsection-header']}>
          <strong>PPM shipment weight</strong>
          <a href="#">Edit</a>
        </div>
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

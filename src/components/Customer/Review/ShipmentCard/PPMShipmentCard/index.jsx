import React from 'react';
import * as PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import ShipmentContainer from '../../../../Office/ShipmentContainer';
import styles from '../ShipmentCard.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';

const PPMShipmentCard = ({ shipmentType }) => {
  return (
    <div className={styles.ShipmentCard} data-testid="shipment-display">
      <ShipmentContainer className={styles.container} shipmentType={shipmentType}>
        <div className={styles.ShipmentCardHeader}>
          <div>
            <h4>Shipment 1: PPM</h4>
            <p>#ABC123K-001</p>
          </div>
          <Button className={styles.editBtn} unstyled>
            Edit
          </Button>
        </div>

        <dl className={styles.shipmentCardSubsection}>
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
        <div className={styles['subsection-header']}>
          <strong>PPM shipment weight</strong>
          <Button className={styles.editBtn} unstyled>
            Edit
          </Button>
        </div>
        <dl className={styles.shipmentCardSubsection}>
          <div className={styles.row}>
            <dt>Estimated weight</dt>
            <dd>5,600 lbs</dd>
          </div>
          <div className={styles.row}>
            <dt>Estimated incentive</dt>
            <dd>Rate info unavailable</dd>
          </div>
        </dl>
      </ShipmentContainer>
    </div>
  );
};

PPMShipmentCard.propTypes = {
  shipmentType: PropTypes.oneOf([
    SHIPMENT_OPTIONS.PPM,
    SHIPMENT_OPTIONS.HHG,
    SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC,
    SHIPMENT_OPTIONS.NTS,
  ]),
};

PPMShipmentCard.defaultProps = {
  shipmentType: SHIPMENT_OPTIONS.PPM,
};

export default PPMShipmentCard;

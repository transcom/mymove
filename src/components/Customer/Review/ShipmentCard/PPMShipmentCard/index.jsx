import React from 'react';
import { string } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import ShipmentContainer from '../../../../Office/ShipmentContainer';
import styles from '../ShipmentCard.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatCustomerDate } from 'shared/utils';

const PPMShipmentCard = ({
  destinationZIP,
  estimatedIncentive,
  estimatedWeight,
  expectedDepartureDate,
  shipmentId,
  sitDays,
  originZIP,
}) => {
  return (
    <div className={styles.ShipmentCard} data-testid="shipment-display">
      <ShipmentContainer className={styles.container} shipmentType={SHIPMENT_OPTIONS.PPM}>
        <div className={styles.ShipmentCardHeader}>
          <div>
            <h3>PPM</h3>
            <p>{shipmentId.substring(0, 10)}</p>
          </div>
          <Button className={styles.editBtn} onClick={() => {}} unstyled>
            Edit
          </Button>
        </div>

        <dl className={styles.shipmentCardSubsection}>
          <div className={styles.row}>
            <dt>Expected departure</dt>
            <dd>{formatCustomerDate(expectedDepartureDate)}</dd>
          </div>
          <div className={styles.row}>
            <dt>Starting ZIP</dt>
            <dd>{originZIP}</dd>
          </div>
          <div className={styles.row}>
            <dt>Storage (SIT)</dt>
            <dd>{Number(sitDays) ? `Yes, ${sitDays} days` : 'No'}</dd>
          </div>
          <div className={styles.row}>
            <dt>Destination ZIP</dt>
            <dd>{destinationZIP}</dd>
          </div>
        </dl>
        <div className={styles.subsectionHeader}>
          <h4>PPM shipment weight</h4>
          <Button className={styles.editBtn} onClick={() => {}} unstyled>
            Edit
          </Button>
        </div>
        <dl className={styles.shipmentCardSubsection}>
          <div className={styles.row}>
            <dt>Estimated weight</dt>
            <dd>{estimatedWeight} lbs</dd>
          </div>
          <div className={styles.row}>
            <dt>Estimated incentive</dt>
            <dd>{estimatedIncentive || 'Rate info unavailable'}</dd>
          </div>
        </dl>
      </ShipmentContainer>
    </div>
  );
};

PPMShipmentCard.propTypes = {
  destinationZIP: string.isRequired,
  estimatedIncentive: string,
  estimatedWeight: string.isRequired,
  expectedDepartureDate: string.isRequired,
  shipmentId: string.isRequired,
  sitDays: string,
  originZIP: string.isRequired,
};

PPMShipmentCard.defaultProps = {
  sitDays: 0,
  estimatedIncentive: null,
};

export default PPMShipmentCard;

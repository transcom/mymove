import React from 'react';
import { string } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import ShipmentContainer from '../../../../Office/ShipmentContainer';
import styles from '../ShipmentCard.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatCustomerDate } from 'shared/formatters';

const PPMShipmentCard = ({
  destinationZIP,
  estimatedIncentive,
  estimatedWeight,
  expectedDepartureDate,
  shipmentId,
  sitDays,
  startingZIP,
}) => {
  return (
    <div className={styles.ShipmentCard} data-testid="shipment-display">
      <ShipmentContainer className={styles.container} shipmentType={SHIPMENT_OPTIONS.PPM}>
        <div className={styles.ShipmentCardHeader}>
          <div>
            <h4>PPM</h4>
            <p>{shipmentId}</p>
          </div>
          <Button className={styles.editBtn} unstyled>
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
            <dd>{startingZIP}</dd>
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
        <div className={styles['subsection-header']}>
          <strong>PPM shipment weight</strong>
          <Button className={styles.editBtn} unstyled>
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
  startingZIP: string.isRequired,
};

PPMShipmentCard.defaultProps = {
  sitDays: 0,
  estimatedIncentive: null,
};

export default PPMShipmentCard;

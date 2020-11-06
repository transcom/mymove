import React from 'react';
import { func, string } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import ShipmentContainer from '../../../../Office/ShipmentContainer';
import styles from '../ShipmentCard.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatCustomerDate } from 'utils/formatters';

const PPMShipmentCard = ({
  destinationZIP,
  estimatedIncentive,
  estimatedWeight,
  expectedDepartureDate,
  moveId,
  onEditClick,
  shipmentId,
  sitDays,
  originZIP,
}) => {
  const editPath = `/moves/${moveId}/review/edit-date-and-location`;
  return (
    <div className={styles.ShipmentCard} data-testid="shipment-display">
      <ShipmentContainer className={styles.container} shipmentType={SHIPMENT_OPTIONS.PPM}>
        <div className={styles.ShipmentCardHeader}>
          <div>
            <h3>PPM</h3>
            <p>#{shipmentId.substring(0, 8).toUpperCase()}</p>
          </div>
          <Button
            className={styles.editBtn}
            data-testid="edit-ppm-dates"
            onClick={() => onEditClick(editPath)}
            unstyled
          >
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
            <dd data-testid="sit-display">{Number(sitDays) ? `Yes, ${sitDays} days` : 'No'}</dd>
          </div>
          <div className={styles.row}>
            <dt className={styles.noborder}>Destination ZIP</dt>
            <dd className={styles.noborder}>{destinationZIP}</dd>
          </div>
        </dl>
        <div className={`${styles.subsectionHeader} todo`}>
          <h4>PPM shipment weight</h4>
          <Button className={styles.editBtn} onClick={() => {}} unstyled>
            Edit
          </Button>
        </div>
        <dl className={styles.shipmentCardSubsection}>
          <div className={`${styles.row} todo`}>
            <dt>Estimated weight</dt>
            <dd>{estimatedWeight} lbs</dd>
          </div>
          <div className={`${styles.row} todo`}>
            <dt className={styles.noborder}>Estimated incentive</dt>
            <dd className={styles.noborder}>{estimatedIncentive || 'Rate info unavailable'}</dd>
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
  moveId: string.isRequired,
  onEditClick: func.isRequired,
  shipmentId: string.isRequired,
  sitDays: string,
  originZIP: string.isRequired,
};

PPMShipmentCard.defaultProps = {
  sitDays: 0,
  estimatedIncentive: null,
};

export default PPMShipmentCard;

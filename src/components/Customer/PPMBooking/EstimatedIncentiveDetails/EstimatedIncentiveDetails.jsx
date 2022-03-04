import React from 'react';

import styles from './EstimatedIncentiveDetails.module.scss';

import { MtoShipmentShape } from 'types/customerShapes';
import { formatCentsTruncateWhole, formatCustomerDate, formatWeight } from 'utils/formatters';

const EstimatedIncentiveDetails = ({ shipment }) => {
  const {
    estimatedWeight,
    pickupPostalCode,
    secondaryPickupPostalCode,
    destinationPostalCode,
    secondaryDestinationPostalCode,
    expectedDepartureDate,
    estimatedIncentive,
  } = shipment?.ppmShipment || {};

  return (
    <div className={styles.EstimatedIncentiveDetails}>
      <div className="container">
        <h2>${formatCentsTruncateWhole(estimatedIncentive)} is your estimated incentive</h2>
        <div className={styles.shipmentDetails}>
          <p>That&apos;s about how much you could earn for moving your PPM, based on what you&apos;ve entered:</p>
          <ul>
            <li>{formatWeight(estimatedWeight)} estimated weight</li>
            <li>Starting from {pickupPostalCode}</li>
            {secondaryPickupPostalCode && <li>Picking up things in {secondaryPickupPostalCode}</li>}
            {secondaryDestinationPostalCode && <li>Dropping off things in {secondaryDestinationPostalCode}</li>}
            <li>Ending in {destinationPostalCode}</li>
            <li>Starting your PPM on {formatCustomerDate(expectedDepartureDate)}</li>
          </ul>
        </div>
        <h3>Your actual incentive amount will vary</h3>
        <p>
          Finance will determine your final incentive based on the total weight you move and the actual date you start
          moving your PPM.
        </p>
        <p>
          You must get certified weight tickets to document the weight you move. You are responsible for uploading them
          to MilMove.
        </p>
      </div>
    </div>
  );
};

EstimatedIncentiveDetails.propTypes = {
  shipment: MtoShipmentShape.isRequired,
};

export default EstimatedIncentiveDetails;

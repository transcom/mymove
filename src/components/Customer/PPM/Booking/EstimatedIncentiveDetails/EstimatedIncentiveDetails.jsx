import React from 'react';

import styles from 'components/Customer/PPM/Booking/EstimatedIncentiveDetails/EstimatedIncentiveDetails.module.scss';
import { ShipmentShape } from 'types/shipment';
import { formatCentsTruncateWhole, formatCustomerDate, formatWeight } from 'utils/formatters';

const EstimatedIncentiveDetails = ({ shipment }) => {
  const {
    estimatedWeight,
    pickupAddress,
    hasSecondaryPickupAddress,
    secondaryPickupAddress,
    destinationAddress,
    hasSecondaryDestinationAddress,
    secondaryDestinationAddress,
    expectedDepartureDate,
    estimatedIncentive,
  } = shipment?.ppmShipment || {};

  return (
    <div className={styles.EstimatedIncentiveDetails}>
      <div className="container">
        <h2>${formatCentsTruncateWhole(estimatedIncentive)} is your estimated incentive</h2>
        <div className={styles.shipmentDetails}>
          <p>This is an estimate of how much you could earn by moving your PPM, based on what you have entered:</p>
          <ul>
            <li>{formatWeight(estimatedWeight)} estimated weight</li>
            <li>Starting from {pickupAddress.postalCode}</li>
            {hasSecondaryPickupAddress && <li>Picking up things in {secondaryPickupAddress.postalCode}</li>}
            {hasSecondaryDestinationAddress && <li>Dropping off things in {secondaryDestinationAddress.postalCode}</li>}
            <li>Ending in {destinationAddress.postalCode}</li>
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
  shipment: ShipmentShape.isRequired,
};

export default EstimatedIncentiveDetails;

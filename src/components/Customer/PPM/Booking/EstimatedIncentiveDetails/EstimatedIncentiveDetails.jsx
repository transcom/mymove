import React from 'react';

import styles from 'components/Customer/PPM/Booking/EstimatedIncentiveDetails/EstimatedIncentiveDetails.module.scss';
import { ShipmentShape } from 'types/shipment';
import { formatAddress } from 'utils/shipmentDisplay';
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

  return estimatedIncentive === 0 ? (
    <div className={styles.EstimatedIncentiveDetails}>
      <div className="container">
        <p>
          The Defense Table of Distances (DTOD) was unavailable during your PPM creation, so we are currently unable to
          provide your estimated incentive. Your estimated incentive information will be updated and provided to you
          during your counseling session.
        </p>
      </div>
    </div>
  ) : (
    <div className={styles.EstimatedIncentiveDetails}>
      <div className="container">
        <h2>${formatCentsTruncateWhole(estimatedIncentive)} is your estimated incentive</h2>
        <div className={styles.shipmentDetails}>
          <p>This is an estimate of how much you could earn by moving your PPM, based on what you have entered:</p>
          <ul>
            <li>{formatWeight(estimatedWeight)} estimated weight</li>
            <li>Starting from {formatAddress(pickupAddress)}</li>
            {hasSecondaryPickupAddress && <li>Picking up things at {formatAddress(secondaryPickupAddress)}</li>}
            {hasSecondaryDestinationAddress && (
              <li>Dropping off things at {formatAddress(secondaryDestinationAddress)}</li>
            )}
            <li>Ending at {formatAddress(destinationAddress)}</li>
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

import React from 'react';
import { bool, func, number } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { generatePath } from 'react-router';

import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import styles from 'components/Customer/Review/ShipmentCard/ShipmentCard.module.scss';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { formatCentsTruncateWhole, formatCustomerDate, formatWeight } from 'utils/formatters';
import { MtoShipmentShape } from 'types/customerShapes';
import { customerRoutes } from 'constants/routes';
import { getShipmentTypeLabel } from 'utils/shipmentDisplay';

const PPMShipmentCard = ({ shipment, shipmentNumber, showEditAndDeleteBtn, onEditClick, onDeleteClick }) => {
  const { moveTaskOrderID, id, shipmentType } = shipment;
  const {
    pickupPostalCode,
    secondaryPickupPostalCode,
    destinationPostalCode,
    secondaryDestinationPostalCode,
    sitExpected,
    expectedDepartureDate,
    proGearWeight,
    spouseProGearWeight,
    estimatedWeight,
    estimatedIncentive,
    hasRequestedAdvance,
    advanceAmountRequested,
  } = shipment?.ppmShipment || {};

  const editPath = `${generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
    moveId: moveTaskOrderID,
    mtoShipmentId: id,
  })}?shipmentNumber=${shipmentNumber}`;

  return (
    <div className={styles.ShipmentCard}>
      <ShipmentContainer className={styles.container} shipmentType={SHIPMENT_OPTIONS.PPM}>
        <div className={styles.ShipmentCardHeader}>
          <div className={styles.shipmentTypeNumber}>
            <h3 data-testid="ShipmentCardNumber">
              {getShipmentTypeLabel(shipmentType)} {shipmentNumber}
            </h3>
            <p>#{id.substring(0, 8).toUpperCase()}</p>
          </div>
          {showEditAndDeleteBtn && (
            <div className={styles.btnContainer}>
              <Button onClick={() => onDeleteClick(shipment.id)} unstyled>
                Delete
              </Button>
              |
              <Button onClick={() => onEditClick(editPath)} unstyled>
                Edit
              </Button>
            </div>
          )}
        </div>

        <dl className={styles.shipmentCardSubsection}>
          <div className={styles.row}>
            <dt>Expected departure</dt>
            <dd>{formatCustomerDate(expectedDepartureDate)}</dd>
          </div>
          <div className={styles.row}>
            <dt>Origin ZIP</dt>
            <dd>{pickupPostalCode}</dd>
          </div>
          {secondaryPickupPostalCode && (
            <div className={styles.row}>
              <dt>Second origin ZIP</dt>
              <dd>{secondaryPickupPostalCode}</dd>
            </div>
          )}
          <div className={styles.row}>
            <dt>Destination ZIP</dt>
            <dd>{destinationPostalCode}</dd>
          </div>
          {secondaryDestinationPostalCode && (
            <div className={styles.row}>
              <dt>Second destination ZIP</dt>
              <dd>{secondaryDestinationPostalCode}</dd>
            </div>
          )}
          <div className={styles.row}>
            <dt>Storage expected? (SIT)</dt>
            <dd>{sitExpected ? 'Yes' : 'No'}</dd>
          </div>
          <div className={styles.row}>
            <dt>Estimated weight</dt>
            <dd>{formatWeight(estimatedWeight)}</dd>
          </div>
          <div className={styles.row}>
            <dt>Pro-gear</dt>
            <dd>{proGearWeight ? `Yes, ${formatWeight(proGearWeight)}` : 'No'}</dd>
          </div>
          <div className={styles.row}>
            <dt>Spouse pro-gear</dt>
            <dd>{spouseProGearWeight ? `Yes, ${formatWeight(spouseProGearWeight)}` : 'No'}</dd>
          </div>
          <div className={styles.row}>
            <dt>Estimated incentive</dt>
            <dd>${estimatedIncentive ? formatCentsTruncateWhole(estimatedIncentive) : '0'}</dd>
          </div>
          <div className={styles.row}>
            <dt>Advance requested?</dt>
            <dd>{hasRequestedAdvance ? `Yes, $${formatCentsTruncateWhole(advanceAmountRequested)}` : 'No'}</dd>
          </div>
        </dl>
      </ShipmentContainer>
    </div>
  );
};

PPMShipmentCard.propTypes = {
  shipment: MtoShipmentShape.isRequired,
  shipmentNumber: number,
  showEditAndDeleteBtn: bool.isRequired,
  onEditClick: func,
  onDeleteClick: func,
};

PPMShipmentCard.defaultProps = {
  shipmentNumber: undefined,
  onEditClick: undefined,
  onDeleteClick: undefined,
};

export default PPMShipmentCard;

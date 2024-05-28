import React from 'react';
import { bool, func, number, oneOf } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { generatePath } from 'react-router-dom';

import styles from 'components/Customer/Review/ShipmentCard/ShipmentCard.module.scss';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import IncompleteShipmentToolTip from 'components/Customer/Review/IncompleteShipmentToolTip/IncompleteShipmentToolTip';
import { customerRoutes } from 'constants/routes';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ShipmentShape } from 'types/shipment';
import { formatCentsTruncateWhole, formatCustomerDate, formatWeight } from 'utils/formatters';
import { getShipmentTypeLabel, canChoosePPMLocation } from 'utils/shipmentDisplay';
import affiliations from 'content/serviceMemberAgencies';
import { MoveShape } from 'types/customerShapes';
import { isPPMShipmentComplete } from 'utils/shipments';

const PPMShipmentCard = ({
  move,
  affiliation,
  shipment,
  shipmentNumber,
  showEditAndDeleteBtn,
  onEditClick,
  onDeleteClick,
  onIncompleteClick,
}) => {
  const { moveTaskOrderID, id, shipmentType, shipmentLocator } = shipment;
  const {
    pickupAddress,
    secondaryPickupAddress,
    destinationAddress,
    secondaryDestinationAddress,
    hasSecondaryPickupAddress,
    hasSecondaryDestinationAddress,
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

  let closeoutOffice;
  if (move?.closeoutOffice == null) {
    closeoutOffice = '';
  } else {
    closeoutOffice = move.closeoutOffice.name;
  }

  const shipmentLabel = `${getShipmentTypeLabel(shipmentType)} ${shipmentNumber}`;
  const moveCodeLabel = `${shipmentLocator}`;
  const shipmentIsIncomplete = !isPPMShipmentComplete(shipment);

  return (
    <div className={styles.ShipmentCard}>
      <ShipmentContainer className={styles.container} shipmentType={SHIPMENT_OPTIONS.PPM}>
        {shipmentIsIncomplete && (
          <IncompleteShipmentToolTip
            onClick={onIncompleteClick}
            shipmentLabel={shipmentLabel}
            moveCodeLabel={moveCodeLabel}
            shipmentTypeLabel={shipmentType}
          />
        )}
        <div className={styles.ShipmentCardHeader}>
          <div className={styles.shipmentTypeNumber}>
            <h3 data-testid="ShipmentCardNumber">{shipmentLabel}</h3>
            <p>#{moveCodeLabel}</p>
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
            <dd>{pickupAddress.postalCode}</dd>
          </div>
          {hasSecondaryPickupAddress && (
            <div className={styles.row}>
              <dt>Second origin ZIP</dt>
              <dd>{secondaryPickupAddress.postalCode}</dd>
            </div>
          )}
          <div className={styles.row}>
            <dt>Destination ZIP</dt>
            <dd>{destinationAddress.postalCode}</dd>
          </div>
          {hasSecondaryDestinationAddress && (
            <div className={styles.row}>
              <dt>Second destination ZIP</dt>
              <dd>{secondaryDestinationAddress.postalCode}</dd>
            </div>
          )}
          {canChoosePPMLocation(affiliation) && closeoutOffice !== '' ? (
            <div className={styles.row}>
              <dt>Closeout office</dt>
              <dd>{closeoutOffice}</dd>
            </div>
          ) : null}
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
  shipment: ShipmentShape.isRequired,
  shipmentNumber: number,
  showEditAndDeleteBtn: bool.isRequired,
  onEditClick: func,
  onDeleteClick: func,
  move: MoveShape,
  affiliation: oneOf(Object.values(affiliations)),
  onIncompleteClick: func,
};

PPMShipmentCard.defaultProps = {
  shipmentNumber: undefined,
  onEditClick: undefined,
  onDeleteClick: undefined,
  move: undefined,
  affiliation: undefined,
  onIncompleteClick: undefined,
};

export default PPMShipmentCard;

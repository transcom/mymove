import { React, useEffect, useState } from 'react';
import { bool, func, number, oneOf } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { generatePath } from 'react-router-dom';

import styles from 'components/Customer/Review/ShipmentCard/ShipmentCard.module.scss';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import IncompleteShipmentToolTip from 'components/Customer/Review/IncompleteShipmentToolTip/IncompleteShipmentToolTip';
import { customerRoutes } from 'constants/routes';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ShipmentShape } from 'types/shipment';
import {
  formatCentsTruncateWhole,
  formatCustomerDate,
  formatWeight,
  formatCustomerContactFullAddress,
} from 'utils/formatters';
import { getShipmentTypeLabel, canChoosePPMLocation } from 'utils/shipmentDisplay';
import affiliations from 'content/serviceMemberAgencies';
import { MoveShape } from 'types/customerShapes';
import { isPPMShipmentComplete } from 'utils/shipments';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const PPMShipmentCard = ({
  move,
  affiliation,
  shipment,
  shipmentNumber,
  showEditAndDeleteBtn,
  onEditClick,
  onDeleteClick,
  onIncompleteClick,
  marketCode,
}) => {
  const { moveTaskOrderID, id, shipmentType, shipmentLocator } = shipment;
  const {
    pickupAddress,
    secondaryPickupAddress,
    tertiaryPickupAddress,
    destinationAddress,
    secondaryDestinationAddress,
    tertiaryDestinationAddress,
    sitExpected,
    expectedDepartureDate,
    proGearWeight,
    spouseProGearWeight,
    estimatedWeight,
    estimatedIncentive,
    hasRequestedAdvance,
    advanceAmountRequested,
  } = shipment?.ppmShipment || {};

  const [isTertiaryAddressEnabled, setIsTertiaryAddressEnabled] = useState(false);
  useEffect(() => {
    const fetchData = async () => {
      isBooleanFlagEnabled('third_address_available').then((enabled) => {
        setIsTertiaryAddressEnabled(enabled);
      });
    };
    fetchData();
  }, []);

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
            <h3 data-testid="ShipmentCardNumber">
              <span className={styles.marketCodeIndicator}>{marketCode}</span>
              {shipmentLabel}
            </h3>
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
            <dt>Origin address</dt>
            <dd>{pickupAddress ? formatCustomerContactFullAddress(pickupAddress) : '—'}</dd>
          </div>
          {secondaryPickupAddress && (
            <div className={styles.row}>
              <dt>Second origin address</dt>
              <dd>{formatCustomerContactFullAddress(secondaryPickupAddress)}</dd>
            </div>
          )}
          {isTertiaryAddressEnabled && tertiaryPickupAddress && secondaryPickupAddress && (
            <div className={styles.row}>
              <dt>Third origin address</dt>
              <dd>{formatCustomerContactFullAddress(tertiaryPickupAddress)}</dd>
            </div>
          )}
          <div className={styles.row}>
            <dt>Destination address</dt>
            <dd>{destinationAddress ? formatCustomerContactFullAddress(destinationAddress) : '—'}</dd>
          </div>
          {secondaryDestinationAddress && (
            <div className={styles.row}>
              <dt>Second destination address</dt>
              <dd>{formatCustomerContactFullAddress(secondaryDestinationAddress)}</dd>
            </div>
          )}
          {isTertiaryAddressEnabled && tertiaryDestinationAddress && secondaryDestinationAddress && (
            <div className={styles.row}>
              <dt>Third destination address</dt>
              <dd>{formatCustomerContactFullAddress(tertiaryDestinationAddress)}</dd>
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

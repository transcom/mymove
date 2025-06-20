import { React, useEffect, useState } from 'react';
import { bool, func, number, oneOf } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { generatePath } from 'react-router-dom';

import styles from 'components/Customer/Review/ShipmentCard/ShipmentCard.module.scss';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import IncompleteShipmentToolTip from 'components/Customer/Review/IncompleteShipmentToolTip/IncompleteShipmentToolTip';
import { customerRoutes } from 'constants/routes';
import { PPM_TYPES, SHIPMENT_OPTIONS, FEATURE_FLAG_KEYS } from 'shared/constants';
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
    ppmType,
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
    gunSafeWeight,
  } = shipment?.ppmShipment || {};

  const [isTertiaryAddressEnabled, setIsTertiaryAddressEnabled] = useState(false);
  const [isGunSafeEnabled, setIsGunSafeEnabled] = useState(false);
  useEffect(() => {
    const fetchData = async () => {
      isBooleanFlagEnabled('third_address_available').then((enabled) => {
        setIsTertiaryAddressEnabled(enabled);
      });
      isBooleanFlagEnabled(FEATURE_FLAG_KEYS.GUN_SAFE).then((enabled) => {
        setIsGunSafeEnabled(enabled);
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
              <Button data-testid="delete-shipment-btn" onClick={() => onDeleteClick(shipment.id)} unstyled>
                Delete
              </Button>
              |
              <Button data-testid="edit-shipment-btn" onClick={() => onEditClick(editPath)} unstyled>
                Edit
              </Button>
            </div>
          )}
        </div>

        <dl className={styles.shipmentCardSubsection}>
          <div className={styles.row}>
            <dt>{ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Shipped date' : 'Expected departure'}</dt>
            <dd>{formatCustomerDate(expectedDepartureDate)}</dd>
          </div>
          <div className={styles.row}>
            <dt>{ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Shipped from Address' : 'Pickup Address'}</dt>
            <dd>{pickupAddress ? formatCustomerContactFullAddress(pickupAddress) : '—'}</dd>
          </div>
          {secondaryPickupAddress && (
            <div className={styles.row}>
              <dt>Second Pickup Address</dt>
              <dd>{formatCustomerContactFullAddress(secondaryPickupAddress)}</dd>
            </div>
          )}
          {isTertiaryAddressEnabled && tertiaryPickupAddress && secondaryPickupAddress && (
            <div className={styles.row}>
              <dt>Third Pickup Address</dt>
              <dd>{formatCustomerContactFullAddress(tertiaryPickupAddress)}</dd>
            </div>
          )}
          <div className={styles.row}>
            <dt>{ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Destination Address' : 'Delivery Address'}</dt>
            <dd>{destinationAddress ? formatCustomerContactFullAddress(destinationAddress) : '—'}</dd>
          </div>
          {secondaryDestinationAddress && (
            <div className={styles.row}>
              <dt>Second Delivery Address</dt>
              <dd>{formatCustomerContactFullAddress(secondaryDestinationAddress)}</dd>
            </div>
          )}
          {isTertiaryAddressEnabled && tertiaryDestinationAddress && secondaryDestinationAddress && (
            <div className={styles.row}>
              <dt>Third Delivery Address</dt>
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
          {isGunSafeEnabled && (
            <div className={styles.row}>
              <dt>Gun safe</dt>
              <dd>{gunSafeWeight ? `Yes, ${formatWeight(gunSafeWeight)}` : 'No'}</dd>
            </div>
          )}
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

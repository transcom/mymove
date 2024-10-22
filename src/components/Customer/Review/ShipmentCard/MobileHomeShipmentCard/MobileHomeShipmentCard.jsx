import { React } from 'react';
import { bool, func, number } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { generatePath } from 'react-router-dom';

import PickupDisplay from '../PickupDisplay';
import DeliveryDisplay from '../DeliveryDisplay';

import styles from 'components/Customer/Review/ShipmentCard/ShipmentCard.module.scss';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import IncompleteShipmentToolTip from 'components/Customer/Review/IncompleteShipmentToolTip/IncompleteShipmentToolTip';
import { customerRoutes } from 'constants/routes';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { ShipmentShape } from 'types/shipment';
import { convertInchesToFeetAndInches } from 'utils/formatMtoShipment';
import { getShipmentTypeLabel } from 'utils/shipmentDisplay';
import { isMobileHomeShipmentComplete } from 'utils/shipments';

const MobileHomeShipmentCard = ({
  shipment,
  shipmentNumber,
  showEditAndDeleteBtn,
  onEditClick,
  onDeleteClick,
  onIncompleteClick,
  destinationLocation,
  destinationZIP,
  secondaryDeliveryAddress,
  tertiaryDeliveryAddress,
  pickupLocation,
  secondaryPickupAddress,
  tertiaryPickupAddress,
  receivingAgent,
  releasingAgent,
  remarks,
  requestedDeliveryDate,
  requestedPickupDate,
  shipmentId,
  marketCode,
}) => {
  const { moveTaskOrderID, id, shipmentType, shipmentLocator } = shipment;
  const { year, make, model, lengthInInches, widthInInches, heightInInches } = shipment?.mobileHomeShipment || {};

  const editPath = `${generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
    moveId: moveTaskOrderID,
    mtoShipmentId: id,
  })}?shipmentNumber=${shipmentNumber}`;

  const shipmentLabel = `${getShipmentTypeLabel(shipmentType)} ${shipmentNumber}`;
  const moveCodeLabel = `${shipmentLocator}`;
  const shipmentIsIncomplete = !isMobileHomeShipmentComplete(shipment);
  const length = convertInchesToFeetAndInches(lengthInInches);
  const width = convertInchesToFeetAndInches(widthInInches);
  const height = convertInchesToFeetAndInches(heightInInches);
  const formattedDimensions = `${length?.feet}'${length?.inches > 0 ? ` ${length.inches}"` : ''} L x ${width?.feet}'${
    width?.inches > 0 ? ` ${width.inches}"` : ''
  } W x ${height?.feet}'${height?.inches > 0 ? ` ${height.inches}"` : ''} H`;

  return (
    <div className={styles.ShipmentCard}>
      <ShipmentContainer className={styles.container} shipmentType={SHIPMENT_OPTIONS.MOBILE_HOME}>
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
              <Button data-testid="deleteShipmentButton" onClick={() => onDeleteClick(shipment.id)} unstyled>
                Delete
              </Button>
              |
              <Button data-testid="editShipmentButton" onClick={() => onEditClick(editPath)} unstyled>
                Edit
              </Button>
            </div>
          )}
        </div>

        <dl className={styles.shipmentCardSubsection}>
          <PickupDisplay
            shipmentId={shipmentId}
            shipmentType={shipmentType}
            requestedPickupDate={requestedPickupDate}
            pickupLocation={pickupLocation}
            secondaryPickupAddress={secondaryPickupAddress}
            tertiaryPickupAddress={tertiaryPickupAddress}
            releasingAgent={releasingAgent}
          />
          <DeliveryDisplay
            shipmentId={shipmentId}
            shipmentType={shipmentType}
            requestedDeliveryDate={requestedDeliveryDate}
            destinationLocation={destinationLocation}
            secondaryDeliveryAddress={secondaryDeliveryAddress}
            tertiaryDeliveryAddress={tertiaryDeliveryAddress}
            destinationZIP={destinationZIP}
            receivingAgent={receivingAgent}
          />
          <div className={styles.row}>
            <dt>Mobile Home year</dt>
            <dd>{year}</dd>
          </div>
          <div className={styles.row}>
            <dt>Mobile Home make</dt>
            <dd>{make}</dd>
          </div>
          <div className={styles.row}>
            <dt>Mobile Home model</dt>
            <dd>{model}</dd>
          </div>
          <div className={styles.row}>
            <dt>Dimensions</dt>
            <dd>{formattedDimensions}</dd>
          </div>
          {remarks && (
            <div className={`${styles.row} ${styles.remarksRow}`}>
              <dt>Remarks</dt>
              <dd className={styles.remarksCell}>{remarks}</dd>
            </div>
          )}
        </dl>
      </ShipmentContainer>
    </div>
  );
};

MobileHomeShipmentCard.propTypes = {
  shipment: ShipmentShape.isRequired,
  shipmentNumber: number,
  showEditAndDeleteBtn: bool.isRequired,
  onEditClick: func,
  onDeleteClick: func,
  onIncompleteClick: func,
};

MobileHomeShipmentCard.defaultProps = {
  shipmentNumber: undefined,
  onEditClick: undefined,
  onDeleteClick: undefined,
  onIncompleteClick: undefined,
};

export default MobileHomeShipmentCard;

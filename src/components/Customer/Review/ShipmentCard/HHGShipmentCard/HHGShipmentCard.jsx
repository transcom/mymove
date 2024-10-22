import React from 'react';
import { string, shape, number, func, bool } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { generatePath } from 'react-router-dom';

import styles from '../ShipmentCard.module.scss';
import PickupDisplay from '../PickupDisplay';
import DeliveryDisplay from '../DeliveryDisplay';

import { AddressShape } from 'types/address';
import { getShipmentTypeLabel } from 'utils/shipmentDisplay';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import IncompleteShipmentToolTip from 'components/Customer/Review/IncompleteShipmentToolTip/IncompleteShipmentToolTip';
import { shipmentStatuses } from 'constants/shipments';
import { customerRoutes } from 'constants/routes';

const HHGShipmentCard = ({
  destinationLocation,
  destinationZIP,
  secondaryDeliveryAddress,
  tertiaryDeliveryAddress,
  moveId,
  onEditClick,
  onDeleteClick,
  pickupLocation,
  secondaryPickupAddress,
  tertiaryPickupAddress,
  receivingAgent,
  releasingAgent,
  remarks,
  requestedDeliveryDate,
  requestedPickupDate,
  shipmentId,
  shipmentLocator,
  shipmentNumber,
  shipmentType,
  showEditAndDeleteBtn,
  status,
  onIncompleteClick,
  marketCode,
}) => {
  const editPath = `${generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
    moveId,
    mtoShipmentId: shipmentId,
  })}?shipmentNumber=${shipmentNumber}`;

  const shipmentLabel = `${getShipmentTypeLabel(shipmentType)} ${shipmentNumber}`;
  const moveCodeLabel = shipmentLocator;
  const shipmentIsIncomplete = status === shipmentStatuses.DRAFT;

  return (
    <div className={styles.ShipmentCard} data-testid="hhg-summary">
      <ShipmentContainer className={styles.container} shipmentType={shipmentType}>
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
            <h3>
              <span className={styles.marketCodeIndicator}>{marketCode}</span>
              {shipmentLabel}
            </h3>
            <p>#{moveCodeLabel}</p>
          </div>
          {showEditAndDeleteBtn && (
            <div className={styles.btnContainer}>
              <Button onClick={() => onDeleteClick(shipmentId)} unstyled>
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

HHGShipmentCard.propTypes = {
  moveId: string.isRequired,
  shipmentNumber: number.isRequired,
  shipmentType: string.isRequired,
  shipmentId: string.isRequired,
  shipmentLocator: string.isRequired,
  showEditAndDeleteBtn: bool.isRequired,
  requestedPickupDate: string.isRequired,
  pickupLocation: AddressShape.isRequired,
  secondaryPickupAddress: AddressShape,
  destinationLocation: AddressShape,
  secondaryDeliveryAddress: AddressShape,
  releasingAgent: shape({
    firstName: string,
    lastName: string,
    phone: string,
    email: string,
  }),
  requestedDeliveryDate: string.isRequired,
  destinationZIP: string.isRequired,
  onEditClick: func.isRequired,
  onDeleteClick: func.isRequired,
  receivingAgent: shape({
    firstName: string,
    lastName: string,
    phone: string,
    email: string,
  }),
  remarks: string,
};

HHGShipmentCard.defaultProps = {
  destinationLocation: null,
  secondaryPickupAddress: null,
  secondaryDeliveryAddress: null,
  releasingAgent: null,
  receivingAgent: null,
  remarks: '',
};

export default HHGShipmentCard;

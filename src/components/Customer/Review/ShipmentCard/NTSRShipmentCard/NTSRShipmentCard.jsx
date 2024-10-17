import React from 'react';
import { string, shape, func, bool } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { generatePath } from 'react-router-dom';

import styles from '../ShipmentCard.module.scss';
import DeliveryDisplay from '../DeliveryDisplay';

import { AddressShape } from 'types/address';
import { getShipmentTypeLabel } from 'utils/shipmentDisplay';
import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import IncompleteShipmentToolTip from 'components/Customer/Review/IncompleteShipmentToolTip/IncompleteShipmentToolTip';
import { shipmentStatuses } from 'constants/shipments';
import { customerRoutes } from 'constants/routes';

const NTSRShipmentCard = ({
  destinationLocation,
  destinationZIP,
  secondaryDeliveryAddress,
  tertiaryDeliveryAddress,
  receivingAgent,
  remarks,
  requestedDeliveryDate,
  moveId,
  onEditClick,
  onDeleteClick,
  shipmentId,
  shipmentLocator,
  shipmentType,
  showEditAndDeleteBtn,
  status,
  onIncompleteClick,
  marketCode,
}) => {
  const editPath = generatePath(customerRoutes.SHIPMENT_EDIT_PATH, {
    moveId,
    mtoShipmentId: shipmentId,
  });

  const shipmentLabel = getShipmentTypeLabel(shipmentType);
  const moveCodeLabel = shipmentLocator;
  const shipmentIsIncomplete = status === shipmentStatuses.DRAFT;

  return (
    <div className={styles.ShipmentCard} data-testid="ntsr-summary">
      <ShipmentContainer className={styles.container} shipmentType={shipmentType}>
        {shipmentIsIncomplete && (
          <IncompleteShipmentToolTip
            onClick={onIncompleteClick}
            shipmentLabel={shipmentLabel}
            moveCodeLabel={moveCodeLabel}
            shipmentTypeLabel={shipmentLabel}
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
              <Button data-testid="edit-ntsr-shipment-btn" onClick={() => onEditClick(editPath)} unstyled>
                Edit
              </Button>
            </div>
          )}
        </div>
        <dl className={styles.shipmentCardSubsection}>
          <DeliveryDisplay
            shipmentId={shipmentId}
            shipmentType={shipmentType}
            requestedDeliveryDate={requestedDeliveryDate}
            destinationLocation={destinationLocation}
            destinationZIP={destinationZIP}
            secondaryDeliveryAddress={secondaryDeliveryAddress}
            tertiaryDeliveryAddress={tertiaryDeliveryAddress}
            receivingAgent={receivingAgent}
          />
          <div className={`${styles.row} ${styles.remarksRow}`}>
            <dt>Remarks</dt>
            <dd className={styles.remarksCell}>{remarks || '—'}</dd>
          </div>
        </dl>
      </ShipmentContainer>
    </div>
  );
};

NTSRShipmentCard.propTypes = {
  destinationLocation: AddressShape,
  destinationZIP: string.isRequired,
  secondaryDeliveryAddress: AddressShape,
  tertiaryDeliveryAddress: AddressShape,
  moveId: string.isRequired,
  onEditClick: func.isRequired,
  onDeleteClick: func.isRequired,
  requestedDeliveryDate: string.isRequired,
  showEditAndDeleteBtn: bool.isRequired,
  shipmentId: string.isRequired,
  shipmentLocator: string.isRequired,
  shipmentType: string.isRequired,
  receivingAgent: shape({
    firstName: string,
    lastName: string,
    phone: string,
    email: string,
  }),
  remarks: string,
};

NTSRShipmentCard.defaultProps = {
  destinationLocation: null,
  receivingAgent: null,
  remarks: '',
  secondaryDeliveryAddress: null,
  tertiaryDeliveryAddress: null,
};

export default NTSRShipmentCard;

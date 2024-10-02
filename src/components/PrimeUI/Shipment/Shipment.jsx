import React, { useState } from 'react';
import { Button } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { Link, generatePath } from 'react-router-dom';
import PropTypes from 'prop-types';

import ConnectedDestructiveShipmentConfirmationModal from 'components/ConfirmationModals/DestructiveShipmentConfirmationModal';
import { formatPrimeAPIShipmentAddress } from 'utils/shipmentDisplay';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { shipmentTypeLabels } from 'content/shipments';
import { formatCents, formatDateFromIso, formatYesNoInputValue, toDollarString } from 'utils/formatters';
import { ShipmentShape } from 'types/shipment';
import { primeSimulatorRoutes } from 'constants/routes';
import { ppmShipmentStatuses, shipmentDestinationTypes } from 'constants/shipments';
import styles from 'pages/PrimeUI/MoveTaskOrder/MoveDetails.module.scss';
import { ADDRESS_TYPES, SHIPMENT_OPTIONS } from 'shared/constants';

const Shipment = ({ shipment, moveId, onDelete, mtoServiceItems }) => {
  const [isDeleteModalVisible, setIsDeleteModalVisible] = useState(false);

  const editShipmentAddressUrl = moveId
    ? generatePath(primeSimulatorRoutes.SHIPMENT_UPDATE_ADDRESS_PATH, {
        moveCodeOrID: moveId,
        shipmentId: shipment.id,
      })
    : '';

  const editReweighUrl =
    moveId && shipment.reweigh
      ? generatePath(primeSimulatorRoutes.SHIPMENT_UPDATE_REWEIGH_PATH, {
          moveCodeOrID: moveId,
          shipmentId: shipment.id,
          reweighId: shipment.reweigh.id,
        })
      : '';

  const showDeleteModal = () => {
    setIsDeleteModalVisible(true);
  };

  const handleDeleteShipment = (shipmentID) => {
    setIsDeleteModalVisible(false);
    onDelete(shipmentID);
  };

  // Returns True if there are any SIT service item for the shipment, False otherwise.
  const hasSITServiceItem = () => {
    return (
      mtoServiceItems &&
      mtoServiceItems.some(
        (serviceItem) =>
          serviceItem && serviceItem.mtoShipmentID === shipment.id && serviceItem.reServiceCode.includes('SIT'),
      )
    );
  };

  return (
    <dl className={descriptionListStyles.descriptionList}>
      <h3>{`${shipmentTypeLabels[shipment.shipmentType]} shipment`}</h3>
      <div className={classnames(descriptionListStyles.row, styles.shipmentHeader)}>
        {moveId && (
          <>
            {!shipment.ppmShipment && hasSITServiceItem() && (
              <Link
                to={`../shipments/${shipment.id}/sit-extension-requests/new`}
                relative="path"
                className="usa-button usa-button-secondary"
              >
                Request SIT Extension
              </Link>
            )}
            <Link to={`../shipments/${shipment.id}`} relative="path" className="usa-button usa-button-secondary">
              Update Shipment
            </Link>
            {shipment.shipmentType !== SHIPMENT_OPTIONS.PPM && (
              <Link
                to={`../shipments/${shipment.id}/updateDestinationAddress`}
                relative="path"
                className="usa-button usa-button-secondary"
              >
                Update Shipment Destination Address
              </Link>
            )}
            {shipment.shipmentType === SHIPMENT_OPTIONS.PPM &&
              shipment.ppmShipment &&
              shipment.ppmShipment.status !== ppmShipmentStatuses.WAITING_ON_CUSTOMER && (
                <Button onClick={showDeleteModal}>Delete Shipment</Button>
              )}
            <Link
              to={`../shipments/${shipment.id}/service-items/new`}
              relative="path"
              className="usa-button usa-button-secondary"
            >
              Add Service Item
            </Link>
          </>
        )}
      </div>
      <ConnectedDestructiveShipmentConfirmationModal
        isOpen={isDeleteModalVisible}
        shipmentID={shipment.id}
        onClose={setIsDeleteModalVisible}
        onSubmit={handleDeleteShipment}
      />
      <div className={descriptionListStyles.row}>
        <dt>Status:</dt>
        <dd>{shipment.status}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Shipment ID:</dt>
        <dd>{shipment.id}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Shipment eTag:</dt>
        <dd>{shipment.eTag}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Requested Pickup Date:</dt>
        <dd>{shipment.requestedPickupDate}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Scheduled Pickup Date:</dt>
        <dd>{shipment.scheduledPickupDate}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Actual Pickup Date:</dt>
        <dd>{shipment.actualPickupDate}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Scheduled Delivery Date:</dt>
        <dd>{shipment.scheduledDeliveryDate}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Actual Delivery Date:</dt>
        <dd>{shipment.actualDeliveryDate}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Estimated Weight:</dt>
        <dd>{shipment.primeEstimatedWeight ? shipment.primeEstimatedWeight : '—'}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Actual Weight:</dt>
        <dd>{shipment.primeActualWeight}</dd>
      </div>
      {shipment.reweigh?.id && (
        <>
          <div
            className={classnames(descriptionListStyles.row, { [styles.missingInfoError]: !shipment.reweigh.weight })}
          >
            <dt>Reweigh Weight:</dt>
            <dd data-testid="reweigh">{!shipment.reweigh.weight ? 'Missing' : shipment.reweigh.weight}</dd>
            <dd>
              <Link to={editReweighUrl}>Edit</Link>
            </dd>
          </div>
          {shipment.reweigh.verificationReason && (
            <div className={descriptionListStyles.row}>
              <dt>Reweigh Remarks:</dt>
              <dd>{shipment.reweigh.verificationReason}</dd>
            </div>
          )}
        </>
      )}
      {shipment.reweigh?.id && (
        <div className={descriptionListStyles.row}>
          <dt>Reweigh Requested Date:</dt>
          <dd>{formatDateFromIso(shipment.reweigh.requestedAt, 'YYYY-MM-DD')}</dd>
        </div>
      )}
      {shipment.shipmentType === SHIPMENT_OPTIONS.HHG && (
        <div className={descriptionListStyles.row}>
          <dt>Actual Pro Gear Weight:</dt>
          <dd>
            {shipment.actualProGearWeight || shipment.actualProGearWeight === 0 ? shipment.actualProGearWeight : '—'}
          </dd>
        </div>
      )}
      {shipment.shipmentType === SHIPMENT_OPTIONS.HHG && (
        <div className={descriptionListStyles.row}>
          <dt>Actual Spouse Pro Gear Weight:</dt>
          <dd>
            {shipment.actualSpouseProGearWeight || shipment.actualSpouseProGearWeight === 0
              ? shipment.actualSpouseProGearWeight
              : '—'}
          </dd>
        </div>
      )}
      <div className={descriptionListStyles.row}>
        <dt>Pickup Address:</dt>
        <dd>{formatPrimeAPIShipmentAddress(shipment.pickupAddress)}</dd>
        <dd>
          {shipment.pickupAddress?.id && moveId && (
            <Link to={editShipmentAddressUrl} state={{ addressType: ADDRESS_TYPES.PICKUP }}>
              Edit
            </Link>
          )}
        </dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Second Pickup Address:</dt>
        <dd>{formatPrimeAPIShipmentAddress(shipment.secondaryPickupAddress)}</dd>
        <dd>
          {shipment.secondaryPickupAddress?.id && moveId && (
            <Link to={editShipmentAddressUrl} state={{ addressType: ADDRESS_TYPES.SECOND_PICKUP }}>
              Edit
            </Link>
          )}
        </dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Third Pickup Address:</dt>
        <dd>{formatPrimeAPIShipmentAddress(shipment.tertiaryPickupAddress)}</dd>
        <dd>
          {shipment.tertiaryPickupAddress?.id && moveId && (
            <Link to={editShipmentAddressUrl} state={{ addressType: ADDRESS_TYPES.THIRD_PICKUP }}>
              Edit
            </Link>
          )}
        </dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Destination Address:</dt>
        <dd>{formatPrimeAPIShipmentAddress(shipment.destinationAddress)}</dd>
        <dd>
          {shipment.destinationAddress?.id && moveId && (
            <Link to={editShipmentAddressUrl} state={{ addressType: ADDRESS_TYPES.DESTINATION }}>
              Edit
            </Link>
          )}
        </dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Second Destination Address:</dt>
        <dd>{formatPrimeAPIShipmentAddress(shipment.secondaryDeliveryAddress)}</dd>
        <dd>
          {shipment.secondaryDeliveryAddress?.id && moveId && (
            <Link to={editShipmentAddressUrl} state={{ addressType: ADDRESS_TYPES.SECOND_DESTINATION }}>
              Edit
            </Link>
          )}
        </dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Third Destination Address:</dt>
        <dd>{formatPrimeAPIShipmentAddress(shipment.tertiaryDeliveryAddress)}</dd>
        <dd>
          {shipment.tertiaryDeliveryAddress?.id && moveId && (
            <Link to={editShipmentAddressUrl} state={{ addressType: ADDRESS_TYPES.THIRD_DESTINATION }}>
              Edit
            </Link>
          )}
        </dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Destination type:</dt>
        <dd>
          {shipmentDestinationTypes[shipment.destinationType]
            ? shipmentDestinationTypes[shipment.destinationType]
            : '—'}
        </dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Created at:</dt>
        <dd>{formatDateFromIso(shipment.createdAt, 'YYYY-MM-DD')}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Approved at:</dt>
        <dd>{shipment.approvedDate}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Diversion:</dt>
        <dd>{shipment.diversion ? 'yes' : 'no'}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Diversion Reason:</dt>
        <dd>{shipment.diversionReason ? shipment.diversionReason : '—'}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Counselor Remarks:</dt>
        <dd>{shipment.counselorRemarks ? shipment.counselorRemarks : '—'}</dd>
      </div>
      {shipment.ppmShipment && (
        <>
          <h4>PPM-specific fields</h4>
          <div className={descriptionListStyles.row}>
            <dt>PPM Status:</dt>
            <dd>{shipment.ppmShipment.status}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Shipment ID:</dt>
            <dd>{shipment.ppmShipment.id}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Shipment eTag:</dt>
            <dd>{shipment.ppmShipment.eTag}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Created at:</dt>
            <dd>{formatDateFromIso(shipment.ppmShipment.createdAt, 'YYYY-MM-DD')}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Updated at:</dt>
            <dd>{formatDateFromIso(shipment.ppmShipment.updatedAt, 'YYYY-MM-DD')}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Expected Departure Date:</dt>
            <dd>{formatDateFromIso(shipment.ppmShipment.expectedDepartureDate, 'YYYY-MM-DD')}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Actual Move Date:</dt>
            <dd>{formatDateFromIso(shipment.ppmShipment.actualMoveDate, 'YYYY-MM-DD')}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Submitted at:</dt>
            <dd>{formatDateFromIso(shipment.ppmShipment.submittedAt, 'YYYY-MM-DD')}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Reviewed at:</dt>
            <dd>{formatDateFromIso(shipment.ppmShipment.reviewedAt, 'YYYY-MM-DD')}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Approved at:</dt>
            <dd>{formatDateFromIso(shipment.ppmShipment.approvedAt, 'YYYY-MM-DD')}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Pickup Address:</dt>
            <dd>{formatPrimeAPIShipmentAddress(shipment.ppmShipment.pickupAddress)}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Second Pickup Address:</dt>
            <dd>{formatPrimeAPIShipmentAddress(shipment.ppmShipment.secondaryPickupAddress)}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Third Pickup Address:</dt>
            <dd>{formatPrimeAPIShipmentAddress(shipment.ppmShipment.tertiaryPickupAddress)}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Destination Address:</dt>
            <dd>{formatPrimeAPIShipmentAddress(shipment.ppmShipment.destinationAddress)}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Second Destination Address:</dt>
            <dd>{formatPrimeAPIShipmentAddress(shipment.ppmShipment.secondaryDestinationAddress)}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Third Destination Address:</dt>
            <dd>{formatPrimeAPIShipmentAddress(shipment.ppmShipment.tertiaryDestinationAddress)}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM SIT Expected:</dt>
            <dd>
              {shipment.ppmShipment.sitExpected == null ? '' : formatYesNoInputValue(shipment.ppmShipment.sitExpected)}
            </dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Estimated Weight:</dt>
            <dd>{shipment.ppmShipment.estimatedWeight}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Has Pro Gear:</dt>
            <dd>
              {shipment.ppmShipment.hasProGear == null ? '' : formatYesNoInputValue(shipment.ppmShipment.hasProGear)}
            </dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Pro Gear Weight:</dt>
            <dd>{shipment.ppmShipment.proGearWeight}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Spouse Pro Gear Weight:</dt>
            <dd>{shipment.ppmShipment.spouseProGearWeight}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Estimated Incentive:</dt>
            <dd>
              {shipment.ppmShipment.estimatedIncentive == null
                ? ''
                : toDollarString(formatCents(shipment.ppmShipment.estimatedIncentive))}
            </dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM SIT Location:</dt>
            <dd>{shipment.ppmShipment.sitLocation}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM SIT Estimated Weight:</dt>
            <dd>{shipment.ppmShipment.sitEstimatedWeight}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM SIT Estimated Entry Date:</dt>
            <dd>{formatDateFromIso(shipment.ppmShipment.sitEstimatedEntryDate, 'YYYY-MM-DD')}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM SIT Estimated Departure Date:</dt>
            <dd>{formatDateFromIso(shipment.ppmShipment.sitEstimatedDepartureDate, 'YYYY-MM-DD')}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM SIT Estimated Cost:</dt>
            <dd>
              {shipment.ppmShipment.sitEstimatedCost == null
                ? ''
                : toDollarString(formatCents(shipment.ppmShipment.sitEstimatedCost))}
            </dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Actual Pickup Postal Code:</dt>
            <dd>{shipment.ppmShipment.actualPickupPostalCode}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Actual Destination Postal Code:</dt>
            <dd>{shipment.ppmShipment.actualDestinationPostalCode}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Has Requested Advance:</dt>
            <dd>
              {shipment.ppmShipment.hasRequestedAdvance == null
                ? ''
                : formatYesNoInputValue(shipment.ppmShipment.hasRequestedAdvance)}
            </dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Advance Amount Requested:</dt>
            <dd>
              {shipment.ppmShipment.advanceAmountRequested == null
                ? ''
                : toDollarString(formatCents(shipment.ppmShipment.advanceAmountRequested))}
            </dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Has Received Advance:</dt>
            <dd>
              {shipment.ppmShipment.hasReceivedAdvance == null
                ? ''
                : formatYesNoInputValue(shipment.ppmShipment.hasReceivedAdvance)}
            </dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>PPM Advance Amount Received:</dt>
            <dd>
              {shipment.ppmShipment.advanceAmountReceived == null
                ? ''
                : toDollarString(formatCents(shipment.ppmShipment.advanceAmountReceived))}
            </dd>
          </div>
        </>
      )}
      {shipment.boatShipment && (
        <>
          <h4>Boat-Shipment Specific Fields</h4>
          <div className={descriptionListStyles.row}>
            <dt>Shipment Type:</dt>
            <dd>{shipment.boatShipment.type}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Boat Year:</dt>
            <dd>{shipment.boatShipment.year}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Boat Make:</dt>
            <dd>{shipment.boatShipment.make}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Boat Model:</dt>
            <dd>{shipment.boatShipment.model}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Boat Dimensions (Inches):</dt>
            <dd>Length: {shipment.boatShipment.lengthInInches}</dd>
            <dd>Width: {shipment.boatShipment.widthInInches}</dd>
            <dd>Height: {shipment.boatShipment.heightInInches}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Has Trailer:</dt>
            <dd>{shipment.boatShipment.hasTrailer ? 'Yes' : 'No'}</dd>
          </div>
          {shipment.boatShipment.hasTrailer && (
            <div className={descriptionListStyles.row}>
              <dt>Trailer is Roadworthy:</dt>
              <dd>{shipment.boatShipment.isRoadworthy ? 'Yes' : 'No'}</dd>
            </div>
          )}
        </>
      )}
    </dl>
  );
};

Shipment.propTypes = {
  shipment: ShipmentShape.isRequired,
  moveId: PropTypes.string,
  onDelete: PropTypes.func,
};

Shipment.defaultProps = {
  moveId: '',
  onDelete: () => {},
};

export default Shipment;

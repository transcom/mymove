import React from 'react';
import classnames from 'classnames';
import { Link } from 'react-router-dom';
import { generatePath } from 'react-router';
import PropTypes from 'prop-types';

import { formatPrimeAPIShipmentAddress } from 'utils/shipmentDisplay';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { shipmentTypeLabels } from 'content/shipments';
import { formatDateFromIso } from 'utils/formatters';
import { ShipmentShape } from 'types/shipment';
import { primeSimulatorRoutes } from 'constants/routes';
import { shipmentDestinationTypes } from 'constants/shipments';
import styles from 'pages/PrimeUI/MoveTaskOrder/MoveDetails.module.scss';

const Shipment = ({ shipment, moveId }) => {
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

  return (
    <dl className={descriptionListStyles.descriptionList}>
      <div className={classnames(descriptionListStyles.row, styles.shipmentHeader)}>
        <h3>{`${shipmentTypeLabels[shipment.shipmentType]} shipment`}</h3>
        {moveId && (
          <>
            <Link
              to={`/simulator/moves/${moveId}/shipments/${shipment.id}`}
              className="usa-button usa-button-secondary"
            >
              Update Shipment
            </Link>
            <Link to={`shipments/${shipment.id}/service-items/new`} className="usa-button usa-button-secondary">
              Add Service Item
            </Link>
          </>
        )}
      </div>
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
        <dt>Estimated Weight:</dt>
        <dd>{shipment.primeEstimatedWeight}</dd>
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
      <div className={descriptionListStyles.row}>
        <dt>Pickup Address:</dt>
        <dd>{formatPrimeAPIShipmentAddress(shipment.pickupAddress)}</dd>
        <dd>{shipment.pickupAddress?.id && moveId && <Link to={editShipmentAddressUrl}>Edit</Link>}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Destination Address:</dt>
        <dd>{formatPrimeAPIShipmentAddress(shipment.destinationAddress)}</dd>
        <dd>{shipment.destinationAddress?.id && moveId && <Link to={editShipmentAddressUrl}>Edit</Link>}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Destination type:</dt>
        <dd>
          {shipmentDestinationTypes[shipment.destinationType]
            ? shipmentDestinationTypes[shipment.destinationType]
            : '-'}
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
      {shipment.ppmShipment && (
        <div className={descriptionListStyles.row}>
          <dt>PPM Status:</dt>
          <dd>{shipment.ppmShipment.status}</dd>
        </div>
      )}
    </dl>
  );
};

Shipment.propTypes = {
  shipment: ShipmentShape.isRequired,
  moveId: PropTypes.string,
};

Shipment.defaultProps = {
  moveId: '',
};

export default Shipment;

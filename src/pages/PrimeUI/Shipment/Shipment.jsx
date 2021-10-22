import classnames from 'classnames';
import { Link } from 'react-router-dom';
import React from 'react';
import { generatePath } from 'react-router';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import descriptionListStyles from '../../../styles/descriptionList.module.scss';
import styles from '../MoveTaskOrder/MoveDetails.module.scss';
import { shipmentTypeLabels } from '../../../content/shipments';
import { formatDateFromIso } from '../../../shared/formatters';

import { ShipmentOptionsOneOf } from 'types/shipment';
import { AgentShape } from 'types/agent';
import { AddressShape } from 'types/address';
// import { primeSimulatorRoutes } from '../../../constants/routes';

const Shipment = ({ shipment, moveId }) => {
  return (
    <dl className={descriptionListStyles.descriptionList}>
      <div className={classnames(descriptionListStyles.row, styles.shipmentHeader)}>
        <h3>{`${shipmentTypeLabels[shipment.shipmentType]} shipment`}</h3>
        <Link to={`/simulator/moves/${moveId}/shipments/${shipment.id}`} className="usa-button usa-button-secondary">
          Update Shipment
        </Link>
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
      <div className={descriptionListStyles.row}>
        <dt>Pickup Address:</dt>
        <dd>
          {shipment.pickupAddress.streetAddress1} {shipment.pickupAddress.streetAddress2} {shipment.pickupAddress.city}{' '}
          {shipment.pickupAddress.state} {shipment.pickupAddress.postalCode}
        </dd>
        <dd>
          {true && (
            <Link
              to=""
              /* generatePath(primeSimulatorRoutes.SHIPMENT_UPDATE_ADDRESS_PATH, { moveCodeOrID: moveId, shipmentId: shipment.id }) */ className="usa-button usa-button-secondary"
            >
              Edit Address
            </Link>
          )}
        </dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Destination Address:</dt>
        <dd>
          {shipment.destinationAddress.streetAddress1} {shipment.destinationAddress.streetAddress2}{' '}
          {shipment.destinationAddress.city} {shipment.destinationAddress.state}{' '}
          {shipment.destinationAddress.postalCode}
        </dd>
        <dd>
          {true && (
            <Button
              type="button"
              onClick={() => {} /* handleDivertShipment(shipmentInfo.shipmentID, shipmentInfo.ifMatchEtag) */}
              unstyled
            >
              Edit address
            </Button>
          )}
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
    </dl>
  );
};

Shipment.propTypes = {
  shipment: PropTypes.shape({
    id: PropTypes.string,
    eTag: PropTypes.string,
    shipmentType: ShipmentOptionsOneOf,
    requestedPickupDate: PropTypes.string,
    scheduledPickupDate: PropTypes.string,
    actualPickupDate: PropTypes.string,
    pickupAddress: AddressShape,
    secondaryPickupAddress: AddressShape,
    destinationAddress: AddressShape,
    secondaryDeliveryAddress: AddressShape,
    agents: PropTypes.arrayOf(AgentShape),
    primeEstimatedWeight: PropTypes.number,
    primeActualWeight: PropTypes.number,
    diversion: PropTypes.bool,
    counselorRemarks: PropTypes.string,
    customerRemarks: PropTypes.string,
    status: PropTypes.string,
    reweigh: PropTypes.shape({
      id: PropTypes.string,
    }),
    createdAt: PropTypes.string,
    approvedDate: PropTypes.string,
  }).isRequired,
  moveId: PropTypes.string.isRequired,
};

export default Shipment;

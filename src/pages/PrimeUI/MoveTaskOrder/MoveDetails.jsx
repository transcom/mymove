import React from 'react';
import { useParams, Link } from 'react-router-dom';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import { shipmentTypeLabels } from '../../../content/shipments';
import { formatDateFromIso } from '../../../shared/formatters';

import styles from './MoveDetails.module.scss';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import SectionWrapper from 'components/Customer/SectionWrapper';
import formStyles from 'styles/form.module.scss';
import { ShipmentOptionsOneOf } from 'types/shipment';
import { AgentShape } from 'types/agent';
import { AddressShape } from 'types/address';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { usePrimeSimulatorGetMove } from 'hooks/queries';

const Shipment = ({ shipment }) => {
  return (
    <dl className={descriptionListStyles.descriptionList}>
      <div className={classnames(descriptionListStyles.row, styles.shipmentHeader)}>
        <h3>{`${shipmentTypeLabels[shipment.shipmentType]} shipment`}</h3>
        <Link to={`shipments/${shipment.id}`} className="usa-button usa-button-secondary">
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
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Destination Address:</dt>
        <dd>
          {shipment.destinationAddress.streetAddress1} {shipment.destinationAddress.streetAddress2}{' '}
          {shipment.destinationAddress.city} {shipment.destinationAddress.state}{' '}
          {shipment.destinationAddress.postalCode}
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
};

const MoveDetails = () => {
  const { moveCodeOrID } = useParams();

  const { moveTaskOrder, isLoading, isError } = usePrimeSimulatorGetMove(moveCodeOrID);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const { mtoShipments } = moveTaskOrder;

  return (
    <div className={classnames('grid-container-desktop-lg', 'usa-prose', styles.MoveDetails)}>
      <div className="grid-row">
        <div className="grid-col-12">
          <SectionWrapper className={formStyles.formSection}>
            <dl className={descriptionListStyles.descriptionList}>
              <div className={styles.moveHeader}>
                <h2>Move</h2>
                <Link to="payment-requests/new" className="usa-button usa-button-secondary">
                  Create Payment Request
                </Link>
              </div>
              <div className={descriptionListStyles.row}>
                <dt>Move Code:</dt>
                <dd>{moveTaskOrder.moveCode}</dd>
              </div>
              <div className={descriptionListStyles.row}>
                <dt>Move Id:</dt>
                <dd>{moveTaskOrder.id}</dd>
              </div>
            </dl>
          </SectionWrapper>
          <SectionWrapper className={formStyles.formSection}>
            <dl className={descriptionListStyles.descriptionList}>
              <h2>Shipments</h2>
              {mtoShipments?.map((mtoShipment) => {
                return (
                  <div key={mtoShipment.id}>
                    <Shipment shipment={mtoShipment} />
                  </div>
                );
              })}
            </dl>
          </SectionWrapper>
        </div>
      </div>
    </div>
  );
};

export default MoveDetails;

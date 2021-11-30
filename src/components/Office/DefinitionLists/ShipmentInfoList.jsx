import React from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames';

import styles from './OfficeDefinitionLists.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatDate } from 'shared/dates';
import { ShipmentShape } from 'types/shipment';
import { formatAddress, formatAgent } from 'utils/shipmentDisplay';

const ShipmentInfoList = ({ className, shipment, isExpanded }) => {
  const {
    requestedPickupDate,
    pickupAddress,
    secondaryPickupAddress,
    destinationAddress,
    secondaryDeliveryAddress,
    agents,
    counselorRemarks,
    customerRemarks,
  } = shipment;
  const requestedPickupDateElement = (
    <div className={descriptionListStyles.row}>
      <dt>Requested move date</dt>
      <dd>{requestedPickupDate && formatDate(requestedPickupDate, 'DD MMM YYYY')}</dd>
    </div>
  );

  const pickupAddressElement = (
    <div className={descriptionListStyles.row}>
      <dt>Current address</dt>
      <dd>{pickupAddress && formatAddress(pickupAddress)}</dd>
    </div>
  );

  const secondaryPickupAddressElement = (
    <div className={descriptionListStyles.row}>
      <dt>Second pickup address</dt>
      <dd>{secondaryPickupAddress && formatAddress(secondaryPickupAddress)}</dd>
    </div>
  );

  const destinationAddressElement = (
    <div className={descriptionListStyles.row}>
      <dt>Destination address</dt>
      <dd data-testid="destinationAddress">{formatAddress(destinationAddress)}</dd>
    </div>
  );

  const secondaryDeliveryAddressElement = (
    <div className={descriptionListStyles.row}>
      <dt>Second destination address</dt>
      <dd data-testid="secondaryDeliveryAddress">
        {secondaryDeliveryAddress ? formatAddress(secondaryDeliveryAddress) : '—'}
      </dd>
    </div>
  );

  return (
    <dl
      className={classNames(
        descriptionListStyles.descriptionList,
        descriptionListStyles.tableDisplay,
        descriptionListStyles.compact,
        className,
      )}
      data-testid="shipment-info-list"
    >
      {requestedPickupDateElement}
      {pickupAddress && pickupAddressElement}
      {isExpanded && secondaryPickupAddress && secondaryPickupAddressElement}
      {destinationAddressElement}
      {isExpanded && secondaryDeliveryAddress && secondaryDeliveryAddressElement}

      {isExpanded &&
        agents &&
        agents.map((agent) => (
          <div className={descriptionListStyles.row} key={`${agent.agentType}-${agent.email}`}>
            <dt>{agent.agentType === 'RELEASING_AGENT' ? 'Releasing agent' : 'Receiving agent'}</dt>
            <dd>{formatAgent(agent)}</dd>
          </div>
        ))}
      {isExpanded && (
        <div className={descriptionListStyles.row}>
          <dt>Customer remarks</dt>
          <dd data-testid="customerRemarks">{customerRemarks || '—'}</dd>
        </div>
      )}

      <div className={classNames((descriptionListStyles.row, { [styles.warning]: !counselorRemarks }))}>
        <dt>Counselor remarks</dt>
        <dd data-testid="counselorRemarks">{counselorRemarks || '—'}</dd>
      </div>
    </dl>
  );
};

ShipmentInfoList.propTypes = {
  className: PropTypes.string,
  shipment: ShipmentShape.isRequired,
  isExpanded: PropTypes.bool,
};

ShipmentInfoList.defaultProps = {
  className: '',
  isExpanded: false,
};

export default ShipmentInfoList;

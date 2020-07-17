import React from 'react';
import propTypes from 'prop-types';

import styles from './ServiceItemCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer';
import { mtoShipmentTypeToFriendlyDisplay, toDollarString } from 'shared/formatters';
import { ShipmentOptionsOneOf } from 'types/shipment';

const ServiceItemCard = ({ shipmentType, serviceItemName, amount }) => {
  return (
    <div data-testid="ServiceItemCard" className={styles.ServiceItemCard}>
      <ShipmentContainer shipmentType={shipmentType}>
        <h6 className={styles.cardHeader}>{mtoShipmentTypeToFriendlyDisplay(shipmentType) || 'BASIC SERVICE ITEMS'}</h6>
        <dl>
          <dt>Service item</dt>
          <dd data-cy="serviceItemName">{serviceItemName}</dd>

          <dt>Amount</dt>
          <dd data-cy="serviceItemAmount">{toDollarString(amount)}</dd>
        </dl>
      </ShipmentContainer>
    </div>
  );
};

ServiceItemCard.propTypes = {
  shipmentType: ShipmentOptionsOneOf,
  serviceItemName: propTypes.string.isRequired,
  amount: propTypes.number.isRequired,
};

ServiceItemCard.defaultProps = {
  shipmentType: null,
};

export default ServiceItemCard;

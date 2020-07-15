import React from 'react';
import propTypes from 'prop-types';

import styles from './ServiceItemCard.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer';
import { mtoShipmentTypeToFriendlyDisplay, toDollarString } from 'shared/formatters';

const ServiceItemCard = ({ shipmentType, serviceItemName, amount }) => {
  return (
    <div data-testid="ServiceItemCard" className={styles.ServiceItemCard}>
      <ShipmentContainer shipmentType={shipmentType}>
        <>
          <h6 className={styles.cardHeader}>
            {mtoShipmentTypeToFriendlyDisplay(shipmentType) || 'BASIC SERVICE ITEMS'}
          </h6>
          <div className="usa-label">Service item</div>
          <div className={styles.textValue}>{serviceItemName}</div>
          <div className="usa-label">Amount</div>
          <div className={styles.textValue}>{toDollarString(amount)}</div>
        </>
      </ShipmentContainer>
    </div>
  );
};

ServiceItemCard.propTypes = {
  shipmentType: propTypes.string,
  serviceItemName: propTypes.string.isRequired,
  amount: propTypes.number.isRequired,
};

ServiceItemCard.defaultProps = {
  shipmentType: '',
};

export default ServiceItemCard;

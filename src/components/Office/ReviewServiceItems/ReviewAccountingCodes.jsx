import React from 'react';
import PropTypes from 'prop-types';

import styles from './ReviewAccountingCodes.module.scss';

import { PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { toDollarString } from 'utils/formatters';
import { AccountingCodesShape } from 'types/accountingCodes';
import { ServiceItemCardsShape } from 'types/serviceItems';
import { shipmentTypeLabels } from 'content/shipments';

const ReviewAccountingCodesItem = ({ shipmentId, shipmentType, tac, sac, amount }) => {
  return (
    <div className={`${styles.Shipment} ${styles[`Shipment_${shipmentType}`]}`} data-testid={`shipment-${shipmentId}`}>
      <div className={styles.ShipmentAmount}>{toDollarString(amount)}</div>
      <div className={styles.ShipmentType}>{shipmentType}</div>
      {tac && <div>TAC: {tac}</div>}
      {sac && <div className={styles.SAC}>SAC: {sac}</div>}
    </div>
  );
};

//  Renders Move level services such as Move management (MS) and Counseling(CS) that not associated to a shipment.
const ReviewMoveLevelServiceItems = ({ services }) => {
  return (
    <div className={`${styles.Shipment}`}>
      {services.map((item) => (
        <div key={item.id}>
          <div className={styles.ShipmentAmount}>{toDollarString(item.amount)}</div>
          <div className={styles.ShipmentType}>{item.mtoServiceItemName} fee</div>
        </div>
      ))}
    </div>
  );
};

ReviewAccountingCodesItem.propTypes = {
  shipmentId: PropTypes.string,
  shipmentType: PropTypes.string,
  tac: PropTypes.string,
  sac: PropTypes.string,
  amount: PropTypes.number.isRequired,
};

ReviewAccountingCodesItem.defaultProps = {
  shipmentId: '',
  shipmentType: null,
  tac: null,
  sac: null,
};

const ReviewAccountingCodes = ({ TACs, SACs, cards }) => {
  const shipments = Object.values(
    cards
      .filter((card) => !!card.mtoShipmentID && card.status === PAYMENT_SERVICE_ITEM_STATUS.APPROVED)
      .reduce((mem, card) => {
        const shipment = mem[card.mtoShipmentID] || {
          id: card.mtoShipmentID,
          amount: 0,
          shipmentType: shipmentTypeLabels[card.mtoShipmentType],
          tac: TACs[card.mtoShipmentTacType] ? `${TACs[card.mtoShipmentTacType]} (${card.mtoShipmentTacType})` : null,
          sac: SACs[card.mtoShipmentSacType] ? `${SACs[card.mtoShipmentSacType]} (${card.mtoShipmentSacType})` : null,
        };

        return {
          ...mem,
          [card.mtoShipmentID]: {
            ...shipment,
            amount: shipment.amount + card.amount,
          },
        };
      }, {}),
  );

  //  Create move level service items. Move level service will not have an associated MTO Shipment ID
  const moveLevelServices = [];
  cards
    .filter((card) => !card.mtoShipmentID && card.status === PAYMENT_SERVICE_ITEM_STATUS.APPROVED)
    .forEach((item) => moveLevelServices.push(item));

  if (shipments.length === 0 && moveLevelServices.length === 0) {
    return null;
  }

  return (
    <div className={styles.ReviewAccountingCodes}>
      <h4>Accounting codes</h4>
      {shipments.map((shipment) => {
        return (
          <ReviewAccountingCodesItem
            key={shipment.id}
            shipmentId={shipment.id}
            tac={shipment.tac}
            sac={shipment.sac}
            shipmentType={shipment.shipmentType}
            amount={shipment.amount}
          />
        );
      })}

      {moveLevelServices.length > 0 && <ReviewMoveLevelServiceItems services={moveLevelServices} />}
    </div>
  );
};

ReviewAccountingCodes.propTypes = {
  TACs: AccountingCodesShape,
  SACs: AccountingCodesShape,
  cards: ServiceItemCardsShape,
};

ReviewAccountingCodes.defaultProps = {
  TACs: {},
  SACs: {},
  cards: [],
};

export default ReviewAccountingCodes;

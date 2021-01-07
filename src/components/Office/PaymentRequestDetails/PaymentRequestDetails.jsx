import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import { PaymentServiceItemShape } from '../../../types/index';

import styles from './PaymentRequestDetails.module.scss';

import { PAYMENT_SERVICE_ITEM_STATUS } from 'shared/constants';
import { formatCents, toDollarString } from 'shared/formatters';

const PaymentRequestDetails = ({ serviceItems }) => {
  return (
    <div className={styles.PaymentRequestDetails}>
      <div className="stackedtable-header">
        {/* TODO this div will become dynamic based on different shipment types */}
        <div className={styles.shipmentType}>
          <div className={styles.basicServiceType} />
          <h3>
            Basic Service Items ({serviceItems.length} {serviceItems.length > 1 ? 'items' : 'item'})
          </h3>
        </div>
      </div>
      <table className="table--stacked">
        <colgroup>
          <col style={{ width: '50%' }} />
          <col style={{ width: '25%' }} />
          <col style={{ width: '25%' }} />
        </colgroup>
        <thead>
          <tr>
            <th>Service Item</th>
            <th className="align-right">Amount</th>
            <th className="align-right">Status</th>
          </tr>
        </thead>
        <tbody>
          {serviceItems.map((item) => {
            return (
              // eslint-disable-next-line react/no-array-index-key
              <tr key={item.id}>
                <td>{item.serviceItemName}</td>
                <td>{toDollarString(formatCents(item.priceCents))}</td>
                <td>
                  {item.status === PAYMENT_SERVICE_ITEM_STATUS.PENDING && (
                    <div className={styles.needsReview}>
                      <FontAwesomeIcon icon="exclamation-circle" />
                      <span>Needs Review</span>
                    </div>
                  )}
                  {item.status === PAYMENT_SERVICE_ITEM_STATUS.APPROVED && (
                    <div className={styles.accepted}>
                      <FontAwesomeIcon icon="check" />
                      <span>Accepted</span>
                    </div>
                  )}
                  {item.status === PAYMENT_SERVICE_ITEM_STATUS.DENIED && (
                    <div className={styles.rejected}>
                      <FontAwesomeIcon icon="times" />
                      <span>Rejected</span>
                    </div>
                  )}
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
};

PaymentRequestDetails.propTypes = {
  serviceItems: PaymentServiceItemShape.isRequired,
};

export default PaymentRequestDetails;

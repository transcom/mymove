import React from 'react';
import { PropTypes } from 'prop-types';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './PaymentRequestDetails.module.scss';

import { PAYMENT_SERVICE_ITEM_STATUS, SHIPMENT_OPTIONS } from 'shared/constants';
import { formatCents, toDollarString } from 'shared/formatters';
import { PaymentServiceItemShape } from 'types';

const shipmentHeadingAndStyle = (mtoShipmentType, shipmentAddresses) => {
  switch (mtoShipmentType) {
    case undefined:
    case null:
      return ['Basic service items', styles.basicServiceType];
    case SHIPMENT_OPTIONS.HHG:
    case SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC:
    case SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC:
      return ['Household goods', styles.hhgShipmentType, shipmentAddresses.hhgAddress];
    case SHIPMENT_OPTIONS.NTS:
    case SHIPMENT_OPTIONS.NTSR:
      return ['NTS release', styles.ntsrShipmentType, shipmentAddresses.ntsAddress];
    default:
      return [mtoShipmentType, styles.basicServiceType];
  }
};

const PaymentRequestDetails = ({ serviceItems, shipmentAddresses }) => {
  const [headingType, shipmentStyle, pickupToDestinationAddress] = shipmentHeadingAndStyle(
    serviceItems?.[0]?.mtoShipmentType,
    shipmentAddresses,
  );

  return (
    serviceItems.length > 0 && (
      <div className={styles.PaymentRequestDetails}>
        <div className="stackedtable-header">
          <div className={styles.shipmentType}>
            <div className={shipmentStyle} />
            <h3>
              {headingType} ({serviceItems.length} {serviceItems.length > 1 ? 'items' : 'item'})
            </h3>
          </div>
          {pickupToDestinationAddress && <p data-testid="pickup-to-destination">{pickupToDestinationAddress}</p>}
        </div>
        <table className="table--stacked">
          <colgroup>
            <col style={{ width: '50%' }} />
            <col style={{ width: '25%' }} />
            <col style={{ width: '25%' }} />
          </colgroup>
          <thead>
            <tr>
              <th>Service item</th>
              <th className="align-right">Amount</th>
              <th className="align-right">Status</th>
            </tr>
          </thead>
          <tbody>
            {serviceItems.map((item) => {
              return (
                // eslint-disable-next-line react/no-array-index-key
                <tr key={item.id}>
                  <td data-testid="serviceItemName">{item.mtoServiceItemName}</td>
                  <td data-testid="serviceItemAmount">{toDollarString(formatCents(item.priceCents))}</td>
                  <td data-testid="serviceItemStatus">
                    {item.status === PAYMENT_SERVICE_ITEM_STATUS.REQUESTED && (
                      <div className={styles.needsReview}>
                        <FontAwesomeIcon icon="exclamation-circle" />
                        <span>Needs review</span>
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
    )
  );
};

PaymentRequestDetails.propTypes = {
  serviceItems: PropTypes.arrayOf(PaymentServiceItemShape).isRequired,
  shipmentAddresses: PropTypes.shape({
    hhgAddress: PropTypes.string,
    ntsAddress: PropTypes.string,
  }),
};

PaymentRequestDetails.defaultProps = {
  shipmentAddresses: {},
};

export default PaymentRequestDetails;

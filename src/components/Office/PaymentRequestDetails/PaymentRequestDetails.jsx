import React from 'react';
import { PropTypes } from 'prop-types';

import ExpandableServiceItemRow from '../ExpandableServiceItemRow/ExpandableServiceItemRow';
import ShipmentModificationTag from '../../ShipmentModificationTag/ShipmentModificationTag';

import styles from './PaymentRequestDetails.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { PaymentServiceItemShape } from 'types';
import { formatDateFromIso } from 'shared/formatters';
import PAYMENT_REQUEST_STATUSES from 'constants/paymentRequestStatus';
import { shipmentModificationTypes } from 'constants/shipments';

const shipmentHeadingAndStyle = (mtoShipmentType) => {
  switch (mtoShipmentType) {
    case undefined:
    case null:
      return ['Basic service items', styles.basicServiceType];
    case SHIPMENT_OPTIONS.HHG:
    case SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC:
    case SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC:
      return ['HHG', styles.hhgShipmentType];
    case SHIPMENT_OPTIONS.NTS:
      return ['Non-temp storage', styles.ntsrShipmentType];
    case SHIPMENT_OPTIONS.NTSR:
      return ['Non-temp storage release', styles.ntsrShipmentType];
    default:
      return [mtoShipmentType, styles.basicServiceType];
  }
};

const PaymentRequestDetails = ({ serviceItems, shipment, paymentRequestStatus }) => {
  const mtoShipmentType = serviceItems?.[0]?.mtoShipmentType;
  const [headingType, shipmentStyle] = shipmentHeadingAndStyle(mtoShipmentType);
  const { modificationType, departureDate, address } = shipment;
  return (
    serviceItems.length > 0 && (
      <div className={styles.PaymentRequestDetails}>
        <div className="stackedtable-header">
          <div className={styles.shipmentType}>
            <div className={shipmentStyle} />
            <h3>
              {headingType} ({serviceItems.length} {serviceItems.length > 1 ? 'items' : 'item'})
              {modificationType && <ShipmentModificationTag shipmentModificationType={modificationType} />}
            </h3>
          </div>
          {(departureDate || address) && (
            <div>
              <p>
                <small>
                  {departureDate && (
                    <strong data-testid="departure-date">
                      Departed {formatDateFromIso(departureDate, 'DD MMM YYYY')}
                    </strong>
                  )}{' '}
                  {address && <span data-testid="pickup-to-destination">{address}</span>}
                </small>
              </p>
            </div>
          )}
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
            {serviceItems.map((item, index) => {
              return (
                <ExpandableServiceItemRow
                  serviceItem={item}
                  key={item.id}
                  index={index}
                  disableExpansion={paymentRequestStatus === PAYMENT_REQUEST_STATUSES.PENDING}
                />
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
  shipment: PropTypes.shape({
    address: PropTypes.oneOfType([PropTypes.string, PropTypes.node]),
    modificationType: PropTypes.oneOfType([
      PropTypes.string,
      PropTypes.oneOf(Object.values(shipmentModificationTypes)),
    ]),
    departureDate: PropTypes.string,
  }),
  paymentRequestStatus: PropTypes.oneOf(Object.values(PAYMENT_REQUEST_STATUSES)).isRequired,
};

PaymentRequestDetails.defaultProps = {
  shipment: {
    departureDate: '',
    address: '',
    modificationType: '',
  },
};

export default PaymentRequestDetails;

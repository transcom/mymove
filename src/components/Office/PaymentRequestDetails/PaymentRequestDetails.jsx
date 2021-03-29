import React from 'react';
import { PropTypes } from 'prop-types';

import ExpandableServiceItemRow from '../ExpandableServiceItemRow/ExpandableServiceItemRow';

import styles from './PaymentRequestDetails.module.scss';

import { SHIPMENT_OPTIONS } from 'shared/constants';
import { PaymentServiceItemShape } from 'types';

const shipmentHeadingAndStyle = (mtoShipmentType) => {
  switch (mtoShipmentType) {
    case undefined:
    case null:
      return ['Basic service items', styles.basicServiceType];
    case SHIPMENT_OPTIONS.HHG:
    case SHIPMENT_OPTIONS.HHG_LONGHAUL_DOMESTIC:
    case SHIPMENT_OPTIONS.HHG_SHORTHAUL_DOMESTIC:
      return ['Household goods', styles.hhgShipmentType];
    case SHIPMENT_OPTIONS.NTS:
      return ['Non-temp storage', styles.ntsrShipmentType];
    case SHIPMENT_OPTIONS.NTSR:
      return ['Non-temp storage release', styles.ntsrShipmentType];
    default:
      return [mtoShipmentType, styles.basicServiceType];
  }
};

const PaymentRequestDetails = ({ serviceItems, shipmentAddress }) => {
  const mtoShipmentType = serviceItems?.[0]?.mtoShipmentType;
  const [headingType, shipmentStyle] = shipmentHeadingAndStyle(mtoShipmentType);
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
          {shipmentAddress !== '' && <p data-testid="pickup-to-destination">{shipmentAddress}</p>}
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
              return <ExpandableServiceItemRow serviceItem={item} key={item.id} index={index} />;
            })}
          </tbody>
        </table>
      </div>
    )
  );
};

PaymentRequestDetails.propTypes = {
  serviceItems: PropTypes.arrayOf(PaymentServiceItemShape).isRequired,
  shipmentAddress: PropTypes.string,
};

PaymentRequestDetails.defaultProps = {
  shipmentAddress: '',
};

export default PaymentRequestDetails;

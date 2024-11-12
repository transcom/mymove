import React from 'react';
import classNames from 'classnames';
import propTypes from 'prop-types';

import styles from './ShipmentServiceItemsTable.module.scss';

import { serviceItemCodes } from 'content/serviceItems';

const internationalServiceItems = {
  true: [serviceItemCodes.ISLH, serviceItemCodes.POEFSC, serviceItemCodes.IHPK, serviceItemCodes.IHUPK],
  false: [serviceItemCodes.ISLH, serviceItemCodes.PODFSC, serviceItemCodes.IHPK, serviceItemCodes.IHUPK],
};

const shipmentTypes = {
  HHG: [
    serviceItemCodes.DLH,
    serviceItemCodes.DSH,
    serviceItemCodes.FSC,
    serviceItemCodes.DOP,
    serviceItemCodes.DDP,
    serviceItemCodes.DPK,
    serviceItemCodes.DUPK,
  ],
  HHG_INTO_NTS_DOMESTIC: [
    serviceItemCodes.DLH,
    serviceItemCodes.DSH,
    serviceItemCodes.FSC,
    serviceItemCodes.DOP,
    serviceItemCodes.DDP,
    serviceItemCodes.DNPK,
  ],
  HHG_OUTOF_NTS_DOMESTIC: [
    serviceItemCodes.DLH,
    serviceItemCodes.DSH,
    serviceItemCodes.FSC,
    serviceItemCodes.DOP,
    serviceItemCodes.DDP,
    serviceItemCodes.DUPK,
  ],
  BOAT_HAUL_AWAY: [
    serviceItemCodes.DLH,
    serviceItemCodes.DSH,
    serviceItemCodes.FSC,
    serviceItemCodes.DOP,
    serviceItemCodes.DDP,
    serviceItemCodes.DPK,
    serviceItemCodes.DUPK,
  ],
  BOAT_TOW_AWAY: [
    serviceItemCodes.DLH,
    serviceItemCodes.DSH,
    serviceItemCodes.FSC,
    serviceItemCodes.DOP,
    serviceItemCodes.DDP,
    serviceItemCodes.DPK,
    serviceItemCodes.DUPK,
  ],
};

const ShipmentServiceItemsTable = ({ isOconus, shipmentType, marketCode, destinationZip3, pickupZip3, className }) => {
  let shipmentServiceItems = shipmentTypes[`${shipmentType}`] || [];
  const sameZip3 = destinationZip3 === pickupZip3;
  let filteredServiceItemsList;

  // if international address is conus
  if (marketCode === 'i') {
    shipmentServiceItems = internationalServiceItems[`${isOconus}`];
  }

  if (sameZip3) {
    filteredServiceItemsList = shipmentServiceItems.filter((item) => {
      return item !== serviceItemCodes.DLH;
    });
  } else {
    filteredServiceItemsList = shipmentServiceItems.filter((item) => {
      return item !== serviceItemCodes.DSH;
    });
  }
  return (
    <div className={classNames('container', 'container--gray', className)}>
      <table className={classNames('table--stacked', styles.serviceItemsTable)}>
        <caption>
          <div className="stackedtable-header">
            <h4>
              Service items for this shipment <br />
              {filteredServiceItemsList.length} items
            </h4>
          </div>
        </caption>
        <thead className="table--small">
          <tr>
            <th>Service item</th>
          </tr>
        </thead>
        <tbody>
          {filteredServiceItemsList.map((serviceItem) => (
            <tr key={serviceItem}>
              <td>{serviceItem}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

ShipmentServiceItemsTable.propTypes = {
  isOconus: propTypes.bool.isRequired,
  marketCode: propTypes.string.isRequired,
  shipmentType: propTypes.oneOf(Object.keys(shipmentTypes)).isRequired,
  className: propTypes.string,
};

ShipmentServiceItemsTable.defaultProps = {
  className: '',
};

export default ShipmentServiceItemsTable;

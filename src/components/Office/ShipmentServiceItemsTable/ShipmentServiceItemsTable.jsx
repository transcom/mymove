import React from 'react';
import classNames from 'classnames';
import propTypes from 'prop-types';

import styles from './ShipmentServiceItemsTable.module.scss';

import { serviceItemCodes } from 'content/serviceItems';

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
};

const ShipmentServiceItemsTable = ({ shipmentType, destinationZip3, pickupZip3, className }) => {
  const shipmentServiceItems = shipmentTypes[`${shipmentType}`] || [];
  const shortHaulServiceItems = shipmentServiceItems.filter((item) => {
    return item !== 'Domestic linehaul';
  });
  const longHaulServiceItems = shipmentServiceItems.filter((item) => {
    return item !== 'Domestic shorthaul';
  });
  const sameZip3 = destinationZip3 === pickupZip3;
  return (
    <div className={classNames('container', 'container--gray', className)}>
      <table className={classNames('table--stacked', styles.serviceItemsTable)}>
        <caption>
          <div className="stackedtable-header">
            <h4>
              Service items for this shipment <br />
              {sameZip3 ? shortHaulServiceItems.length : longHaulServiceItems.length} items
            </h4>
          </div>
        </caption>
        <thead className="table--small">
          <tr>
            <th>Service item</th>
          </tr>
        </thead>
        <tbody>
          {sameZip3
            ? shortHaulServiceItems.map((serviceItem) => (
                <tr key={serviceItem}>
                  <td>{serviceItem}</td>
                </tr>
              ))
            : longHaulServiceItems.map((serviceItem) => (
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
  shipmentType: propTypes.oneOf(Object.keys(shipmentTypes)).isRequired,
  className: propTypes.string,
};

ShipmentServiceItemsTable.defaultProps = {
  className: '',
};

export default ShipmentServiceItemsTable;

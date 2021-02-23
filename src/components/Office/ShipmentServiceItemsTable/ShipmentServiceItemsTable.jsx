import React from 'react';
import classNames from 'classnames';
import propTypes from 'prop-types';

import styles from './ShipmentServiceItemsTable.module.scss';

import { serviceItemCodes } from 'content/serviceItems';

const shipmentTypes = {
  HHG_LONGHAUL_DOMESTIC: [
    serviceItemCodes.DLH,
    serviceItemCodes.FSC,
    serviceItemCodes.DOP,
    serviceItemCodes.DDP,
    serviceItemCodes.DPK,
    serviceItemCodes.DUPK,
  ],
  HHG_SHORTHAUL_DOMESTIC: [
    serviceItemCodes.DSH,
    serviceItemCodes.FSC,
    serviceItemCodes.DOP,
    serviceItemCodes.DDP,
    serviceItemCodes.DPK,
    serviceItemCodes.DUPK,
  ],
  HHG_INTO_NTS_DOMESTIC: [
    serviceItemCodes.DLH,
    serviceItemCodes.FSC,
    serviceItemCodes.DOP,
    serviceItemCodes.DDP,
    serviceItemCodes.DPK,
    serviceItemCodes.DNPKF,
  ],
  HHG_OUTOF_NTS_DOMESTIC: [
    serviceItemCodes.DLH,
    serviceItemCodes.FSC,
    serviceItemCodes.DOP,
    serviceItemCodes.DDP,
    serviceItemCodes.DUPK,
  ],
};

const ShipmentServiceItemsTable = ({ shipmentType, className }) => {
  const shipmentServiceItems = shipmentTypes[`${shipmentType}`] || [];

  return (
    <div className={classNames('container', 'container--gray', className)}>
      <table className={classNames('table--stacked', styles.serviceItemsTable)}>
        <caption>
          <div className="stackedtable-header">
            <h4>
              Service items for this shipment <span>{shipmentServiceItems.length} items</span>
            </h4>
          </div>
        </caption>
        <thead className="table--small">
          <tr>
            <th>Service item</th>
          </tr>
        </thead>
        <tbody>
          {shipmentServiceItems.map((serviceItem) => (
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

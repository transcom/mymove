import React from 'react';
import classNames from 'classnames';
import propTypes from 'prop-types';

import styles from './ShipmentServiceItemsTable.module.scss';

const serviceItems = {
  domestic_linehaul: 'Domestic linehaul',
  fuel_surcharge: 'Fuel surcharge',
  domestic_origin_price: 'Domestic origin price',
  domestic_destination_price: 'Domestic destination price',
  domestic_packing: 'Domestic packing',
  domestic_unpacking: 'Domestic unpacking',
};

const shipmentTypes = {
  hhg: [
    serviceItems.domestic_linehaul,
    serviceItems.fuel_surcharge,
    serviceItems.domestic_origin_price,
    serviceItems.domestic_destination_price,
    serviceItems.domestic_packing,
    serviceItems.domestic_unpacking,
  ],
  nts: [
    serviceItems.domestic_linehaul,
    serviceItems.fuel_surcharge,
    serviceItems.domestic_origin_price,
    serviceItems.domestic_destination_price,
    serviceItems.domestic_unpacking,
  ],
};

const ShipmentServiceItemsTable = ({ shipmentType }) => {
  const shipmentServiceItems = shipmentTypes[`${shipmentType}`];

  return (
    <div className="container container--gray">
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
  shipmentType: propTypes.string.isRequired,
};

export default ShipmentServiceItemsTable;

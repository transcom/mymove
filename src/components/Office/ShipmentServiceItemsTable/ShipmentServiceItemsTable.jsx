import React from 'react';
import classNames from 'classnames';
import propTypes from 'prop-types';

import styles from './ShipmentServiceItemsTable.module.scss';

import { serviceItemCodes } from 'content/serviceItems';
import { getAllReServiceItems } from 'services/ghcApi';

let allReServiceItems;
const fetchAllReServiceItems = async () => {
  getAllReServiceItems().then((response) => {
    allReServiceItems = JSON.parse(response.data);
  });
};
fetchAllReServiceItems();

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

function filterPortFuelSurcharge(shipment, autoApprovedItems) {
  let filteredPortFuelSurchargeList = autoApprovedItems;
  if (shipment.pickupAddress.isOconus) {
    filteredPortFuelSurchargeList = autoApprovedItems.filter((serviceItem) => {
      return serviceItem.serviceCode !== serviceItemCodes.POEFSC;
    });
  }
  if (shipment.destinationAddress.isOconus) {
    filteredPortFuelSurchargeList = autoApprovedItems.filter((serviceItem) => {
      return serviceItem.serviceCode !== serviceItemCodes.PODFSC;
    });
  }
  return filteredPortFuelSurchargeList;
}

function filterServiceItems(shipment) {
  const { shipmentType, marketCode } = shipment;
  const autoApprovedItems = allReServiceItems.filter((reServiceItem) => {
    return (
      reServiceItem.marektCode === marketCode &&
      reServiceItem.shipmentType === shipmentType &&
      reServiceItem.isAutoApproved === true
    );
  });
  return filterPortFuelSurcharge(shipment, autoApprovedItems);
}

const ShipmentServiceItemsTable = ({ shipment, className }) => {
  const { shipmentType, marketCode } = shipment;
  let filteredServiceItemsList;
  if (marketCode === 'i') {
    filteredServiceItemsList = filterServiceItems(shipmentType, marketCode);
  } else {
    const destinationZip3 = shipment.destinationAddress?.postalCode.slice(0, 3);
    const pickupZip3 = shipment.pickupAddress?.postalCode.slice(0, 3);
    const shipmentServiceItems = shipmentTypes[`${shipmentType}`] || [];
    const sameZip3 = destinationZip3 === pickupZip3;

    if (sameZip3) {
      filteredServiceItemsList = shipmentServiceItems.filter((item) => {
        return item !== serviceItemCodes.DLH;
      });
    } else {
      filteredServiceItemsList = shipmentServiceItems.filter((item) => {
        return item !== serviceItemCodes.DSH;
      });
    }
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
  shipmentType: propTypes.oneOf(Object.keys(shipmentTypes)).isRequired,
  className: propTypes.string,
};

ShipmentServiceItemsTable.defaultProps = {
  className: '',
};

export default ShipmentServiceItemsTable;

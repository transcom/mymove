import React from 'react';
import classNames from 'classnames';
import propTypes from 'prop-types';

import styles from './ShipmentServiceItemsTable.module.scss';

import { serviceItemCodes } from 'content/serviceItems';
import { getAllReServiceItems } from 'services/ghcApi';

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
      return serviceItem.serviceCode !== 'POEFSC';
    });
  }
  if (shipment.destinationAddress.isOconus) {
    filteredPortFuelSurchargeList = autoApprovedItems.filter((serviceItem) => {
      return serviceItem.serviceCode !== 'PODFSC';
    });
  }
  return filteredPortFuelSurchargeList;
}

function convertServiceItemsToServiceNames(serviceItems) {
  const serviceNames = Array(serviceItems.length);
  for (let i = 0; i < serviceItems.length; i += 1) {
    serviceNames[i] = serviceItems[i].serviceName;
  }
  return serviceNames;
}

function getPreapprovedServiceItems(allReServiceItems, shipment) {
  const { shipmentType, marketCode } = shipment;
  const autoApprovedItems = allReServiceItems.filter((reServiceItem) => {
    return (
      reServiceItem.marketCode === marketCode &&
      reServiceItem.shipmentType === shipmentType &&
      reServiceItem.isAutoApproved === true
    );
  });
  return convertServiceItemsToServiceNames(filterPortFuelSurcharge(shipment, autoApprovedItems));
}

const ShipmentServiceItemsTable = ({ shipment, className }) => {
  const { shipmentType, marketCode } = shipment;
  const [filteredServiceItems, setFilteredServiceItems] = React.useState([]);
  React.useEffect(() => {
    if (marketCode === 'i') {
      const fetchServiceItemsFunction = async () => {
        const response = await getAllReServiceItems();
        const allReServiceItems = await JSON.parse(response.data);
        setFilteredServiceItems(getPreapprovedServiceItems(allReServiceItems, shipment));
      };
      fetchServiceItemsFunction();
    } else {
      let filteredServiceItemsList;
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
      setFilteredServiceItems(filteredServiceItemsList);
    }
  }, [marketCode, shipmentType, shipment]);
  return (
    <div className={classNames('container', 'container--gray', className)}>
      <table className={classNames('table--stacked', styles.serviceItemsTable)}>
        <caption>
          <div className="stackedtable-header">
            <h4>
              Service items for this shipment <br />
              {filteredServiceItems.length} items
            </h4>
          </div>
        </caption>
        <thead className="table--small">
          <tr>
            <th>Service item</th>
          </tr>
        </thead>
        <tbody>
          {filteredServiceItems.map((serviceItem) => (
            <tr key={serviceItem}>
              <td>{serviceItem}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

const AddressShape = propTypes.shape({
  isOconus: propTypes.bool.isRequired,
});
const ShipmentShape = propTypes.shape({
  shipmentType: propTypes.string.isRequired,
  marketCode: propTypes.string.isRequired,
  pickupAddress: AddressShape.isRequired,
  destinationAddress: AddressShape.isRequired,
});

ShipmentServiceItemsTable.propTypes = {
  shipment: ShipmentShape.isRequired,
  className: propTypes.string,
};

ShipmentServiceItemsTable.defaultProps = {
  className: '',
};

export default ShipmentServiceItemsTable;

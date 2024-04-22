import React from 'react';
import PropTypes from 'prop-types';

import { SERVICE_ITEM_STATUS, MTO_SERVICE_ITEM_STATUS } from '../../../shared/constants';
import { ServiceItemDetailsShape } from '../../../types/serviceItems';

import styles from './RequestedServiceItemsTable.module.scss';

import ServiceItemsTable from 'components/Office/ServiceItemsTable/ServiceItemsTable';
import { ShipmentShape } from 'types';
import { SitStatusShape } from 'types/sitStatusShape';

const RequestedServiceItemsTable = ({
  serviceItems,
  handleUpdateMTOServiceItemStatus,
  handleShowRejectionDialog,
  handleShowEditSitEntryDateModal,
  statusForTableType,
  serviceItemAddressUpdateAlert,
  shipment,
  sitStatus,
}) => {
  const chooseTitleText = (status) => {
    switch (status) {
      case SERVICE_ITEM_STATUS.SUBMITTED:
        return 'Requested';
      case SERVICE_ITEM_STATUS.APPROVED:
        return 'Approved';
      case SERVICE_ITEM_STATUS.REJECTED:
        return 'Rejected';
      case MTO_SERVICE_ITEM_STATUS.APPROVED:
        return 'Move Task Order Approved';
      case MTO_SERVICE_ITEM_STATUS.REJECTED:
        return 'Move Task Order Approved';
      case MTO_SERVICE_ITEM_STATUS.SUBMITTED:
        return 'Move Task Order Requested';
      default:
        return status;
    }
  };

  const statusTitleText = chooseTitleText(statusForTableType);

  return (
    <div className={styles.RequestedServiceItemsTable} data-testid={`${statusTitleText}ServiceItemsTable`}>
      <h3>
        {statusTitleText} Service Items&nbsp;
        <span>
          ({serviceItems.length} {serviceItems.length === 1 ? 'item' : 'items'})
        </span>
      </h3>
      <ServiceItemsTable
        serviceItems={serviceItems}
        handleUpdateMTOServiceItemStatus={handleUpdateMTOServiceItemStatus}
        handleShowRejectionDialog={handleShowRejectionDialog}
        handleShowEditSitEntryDateModal={handleShowEditSitEntryDateModal}
        statusForTableType={statusForTableType}
        serviceItemAddressUpdateAlert={serviceItemAddressUpdateAlert}
        shipment={shipment}
        sitStatus={sitStatus}
      />
    </div>
  );
};

RequestedServiceItemsTable.propTypes = {
  handleUpdateMTOServiceItemStatus: PropTypes.func.isRequired,
  handleShowRejectionDialog: PropTypes.func.isRequired,
  statusForTableType: PropTypes.string.isRequired,
  serviceItems: PropTypes.arrayOf(ServiceItemDetailsShape).isRequired,
  shipment: ShipmentShape,
  sitStatus: SitStatusShape,
};

RequestedServiceItemsTable.defaultProps = {
  shipment: {},
  sitStatus: undefined,
};

export default RequestedServiceItemsTable;

import React from 'react';

import descriptionListStyles from '../../../styles/descriptionList.module.scss';
import { MTOServiceItemShape } from '../../../types';

import { SERVICE_ITEMS_ALLOWED_WEIGHT_BILLED_PARAM } from 'constants/serviceItems';

const ServiceItem = ({ serviceItem, mtoShipment }) => {
  return (
    <dl className={descriptionListStyles.descriptionList}>
      <h3>{serviceItem.reServiceName}</h3>
      <div className={descriptionListStyles.row}>
        <dt>Status:</dt>
        <dd>{serviceItem.status}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>ID:</dt>
        <dd>{serviceItem.id}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Service Code:</dt>
        <dd>{serviceItem.reServiceCode}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>Service Name:</dt>
        <dd>{serviceItem.reServiceName}</dd>
      </div>
      <div className={descriptionListStyles.row}>
        <dt>eTag:</dt>
        <dd>{serviceItem.eTag}</dd>
      </div>
      {SERVICE_ITEMS_ALLOWED_WEIGHT_BILLED_PARAM.includes(serviceItem.reServiceCode) && (
        <div className={descriptionListStyles.row}>
          <dt>Shipment Weight (pounds):</dt>
          <dd>{mtoShipment.primeActualWeight || 'Not provided'}</dd>
        </div>
      )}
    </dl>
  );
};

ServiceItem.propTypes = {
  serviceItem: MTOServiceItemShape.isRequired,
};

export default ServiceItem;

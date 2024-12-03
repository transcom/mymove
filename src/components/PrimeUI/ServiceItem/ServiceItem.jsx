import React from 'react';

import descriptionListStyles from '../../../styles/descriptionList.module.scss';
import { MTOServiceItemShape } from '../../../types';

import { SERVICE_ITEMS_ALLOWED_WEIGHT_BILLED_PARAM, SERVICE_ITEM_CODES } from 'constants/serviceItems';
import { convertFromThousandthInchToInch } from 'utils/formatters';

const ServiceItem = ({ serviceItem, mtoShipment }) => {
  return (
    <dl className={descriptionListStyles.descriptionList}>
      <h3>
        {serviceItem.reServiceName} {serviceItem.standaloneCrate && '- Standalone'}
      </h3>
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
      {(serviceItem.reServiceCode === SERVICE_ITEM_CODES.ICRT ||
        serviceItem.reServiceCode === SERVICE_ITEM_CODES.IUCRT ||
        serviceItem.reServiceCode === SERVICE_ITEM_CODES.DCRT ||
        serviceItem.reServiceCode === SERVICE_ITEM_CODES.DUCRT) && (
        <>
          <div className={descriptionListStyles.row}>
            <dt>Item Size:</dt>
            <dd>
              {convertFromThousandthInchToInch(serviceItem.item?.length)}&quot; x&nbsp;
              {convertFromThousandthInchToInch(serviceItem.item?.width)}&quot; x&nbsp;
              {convertFromThousandthInchToInch(serviceItem.item?.height)}&quot;
            </dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Crate Size:</dt>
            <dd>
              {convertFromThousandthInchToInch(serviceItem.crate?.length)}&quot; x&nbsp;
              {convertFromThousandthInchToInch(serviceItem.crate?.width)}&quot; x&nbsp;
              {convertFromThousandthInchToInch(serviceItem.crate?.height)}&quot;
            </dd>
          </div>
        </>
      )}
      {(serviceItem.reServiceCode === SERVICE_ITEM_CODES.ICRT ||
        serviceItem.reServiceCode === SERVICE_ITEM_CODES.IUCRT) && (
        <>
          <div className={descriptionListStyles.row}>
            <dt>External Crate:</dt>
            <dd>{serviceItem.externalCrate ? 'Yes' : 'No'}</dd>
          </div>
          <div className={descriptionListStyles.row}>
            <dt>Market:</dt>
            <dd>{serviceItem.market || 'Not provided'}</dd>
          </div>
        </>
      )}
    </dl>
  );
};

ServiceItem.propTypes = {
  serviceItem: MTOServiceItemShape.isRequired,
};

export default ServiceItem;

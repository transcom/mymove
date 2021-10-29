import React from 'react';

import descriptionListStyles from '../../../styles/descriptionList.module.scss';
import { MTOServiceItemShape } from '../../../types';

const ServiceItem = ({ serviceItem }) => {
  return (
    <dl className={descriptionListStyles.descriptionList}>
      <h3>{`${serviceItem.reServiceName}`}</h3>
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
    </dl>
  );
};

ServiceItem.propTypes = {
  serviceItem: MTOServiceItemShape.isRequired,
};

export default ServiceItem;

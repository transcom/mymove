import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import classnames from 'classnames';

import { ReactComponent as Check } from '../../shared/icon/check.svg';
import { ReactComponent as Ex } from '../../shared/icon/ex.svg';
import { SERVICE_ITEM_STATUS } from '../../shared/constants';
import { MTOServiceItemCustomerContactShape, MTOServiceItemDimensionShape } from '../../types/moveOrder';

import styles from './index.module.scss';

import ServiceItemDetails from 'components/Office/ServiceItemDetails/ServiceItemDetails';
import { formatDate } from 'shared/dates';

const ServiceItemTableHasImg = ({ serviceItems, handleUpdateMTOServiceItemStatus }) => {
  const tableRows = serviceItems.map(({ id, code, submittedAt, serviceItem, details }) => {
    return (
      <tr key={id}>
        <td className={styles.nameAndDate}>
          <p className={styles.codeName}>{serviceItem}</p>
          <p>{formatDate(submittedAt, 'DD MMM YYYY')}</p>
        </td>
        <td className={styles.detail}>
          <ServiceItemDetails id={id} code={code} details={details} />
        </td>
        <td>
          <div className={styles.statusAction}>
            <Button
              type="button"
              className="usa-button--icon usa-button--small"
              data-testid="acceptButton"
              onClick={() => handleUpdateMTOServiceItemStatus(id, SERVICE_ITEM_STATUS.APPROVED)}
            >
              <span className="icon">
                <Check />
              </span>
              <span>Accept</span>
            </Button>
            <Button
              type="button"
              secondary
              className="usa-button--small usa-button--icon"
              data-testid="rejectButton"
              onClick={() => handleUpdateMTOServiceItemStatus(id, SERVICE_ITEM_STATUS.REJECTED)}
            >
              <span className="icon">
                <Ex />
              </span>
              <span>Reject</span>
            </Button>
          </div>
        </td>
      </tr>
    );
  });

  return (
    <div className={classnames(styles.ServiceItemTable, 'table--service-item', 'table--service-item--hasimg')}>
      <table>
        <thead className="table--small">
          <tr>
            <th>Service item</th>
            <th>Details</th>
            <th>&nbsp;</th>
          </tr>
        </thead>
        <tbody>{tableRows}</tbody>
      </table>
    </div>
  );
};

ServiceItemTableHasImg.propTypes = {
  handleUpdateMTOServiceItemStatus: PropTypes.func.isRequired,
  serviceItems: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string,
      submittedAt: PropTypes.string,
      serviceItem: PropTypes.string,
      code: PropTypes.string,
      details: PropTypes.shape({
        pickupPostalCode: PropTypes.string,
        reason: PropTypes.string,
        imgURL: PropTypes.string,
        itemDimensions: MTOServiceItemDimensionShape,
        createDimensions: MTOServiceItemDimensionShape,
        firstCustomerContact: MTOServiceItemCustomerContactShape,
        secondCustmoerContact: MTOServiceItemCustomerContactShape,
      }),
    }),
  ).isRequired,
};

export default ServiceItemTableHasImg;

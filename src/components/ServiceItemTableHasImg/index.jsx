import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import classnames from 'classnames';

import { ReactComponent as Check } from '../../shared/icon/check.svg';
import { ReactComponent as Ex } from '../../shared/icon/ex.svg';
import { SERVICE_ITEM_STATUS } from '../../shared/constants';

import styles from './index.module.scss';

import { formatDate } from 'shared/dates';

function generateDetailText(details, id) {
  if (typeof details.text === 'string') {
    return details.text;
  }

  const detailList = Object.keys(details.text).map((detail) => (
    <div key={`${id}-${detail}`} className={styles.detailLine}>
      <dt className={styles.detailType}>{detail}:</dt> <dd>{details.text[`${detail}`]}</dd>
    </div>
  ));

  return <dl>{detailList}</dl>;
}

const ServiceItemTableHasImg = ({ serviceItems, handleUpdateMTOServiceItemStatus }) => {
  const tableRows = serviceItems.map(({ id, submittedAt, serviceItem, details, status }, i) => {
    let detailSection;
    if (details.imgURL) {
      detailSection = (
        <div className={styles.detailImage}>
          <img
            className={styles.siThumbnail}
            alt="requested service item"
            aria-labelledby={`si-thumbnail--caption-${i}`}
            src={details.imgURL}
          />
          <small id={`si-thumbnail--caption-${i}`}>{generateDetailText(details, id)}</small>
        </div>
      );
    } else {
      detailSection = <div>{generateDetailText(details, id)}</div>;
    }

    return (
      <tr key={id}>
        <td className={styles.nameAndDate}>
          <p className={styles.codeName}>{serviceItem}</p>
          <p>{formatDate(submittedAt, 'DD MMM YYYY')}</p>
        </td>
        <td className={styles.detail}>{detailSection}</td>
        <td>
          {status === SERVICE_ITEM_STATUS.SUBMITTED && (
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
          )}{' '}
          {status === SERVICE_ITEM_STATUS.APPROVED && (
            <div className={styles.statusAction}>
              <span className="icon">
                <Ex />
              </span>
              <span
                data-testid="rejectTextButton"
                role="button"
                className="text-blue cursor-pointer"
                onKeyDown={() => handleUpdateMTOServiceItemStatus(id, SERVICE_ITEM_STATUS.REJECTED)}
                onClick={() => handleUpdateMTOServiceItemStatus(id, SERVICE_ITEM_STATUS.REJECTED)}
                tabIndex={i}
              >
                Reject
              </span>
            </div>
          )}
          {status === SERVICE_ITEM_STATUS.REJECTED && (
            <div className={styles.statusAction}>
              <span className="icon">
                <Ex />
              </span>
              <span
                data-testid="approveTextButton"
                role="button"
                className="text-blue cursor-pointer"
                onKeyDown={() => handleUpdateMTOServiceItemStatus(id, SERVICE_ITEM_STATUS.APPROVED)}
                onClick={() => handleUpdateMTOServiceItemStatus(id, SERVICE_ITEM_STATUS.APPROVED)}
                tabIndex={i}
              >
                Approve
              </span>
            </div>
          )}
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
        imgURL: PropTypes.string,
        text: PropTypes.oneOf([PropTypes.string, PropTypes.object]),
      }),
    }),
  ).isRequired,
};

export default ServiceItemTableHasImg;

import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import classnames from 'classnames';

import { ReactComponent as Check } from '../../shared/icon/check.svg';
import { ReactComponent as Ex } from '../../shared/icon/ex.svg';
import { SERVICE_ITEM_STATUS } from '../../shared/constants';
import { MTOServiceItemCustomerContactShape, MTOServiceItemDimensionShape } from '../../types/moveOrder';

import styles from './index.module.scss';

import { formatDate } from 'shared/dates';
import { convertFromThousandthInchToInch } from 'shared/formatters';

function generateDetailText(details, id) {
  const detailList = Object.keys(details).map((detail) => (
    <div key={`${id}-${detail}`} className={styles.detailLine}>
      <dt className={styles.detailType}>{detail}:</dt> <dd>{details[`${detail}`]}</dd>
    </div>
  ));

  return detailList;
}

const ServiceItemTableHasImg = ({ serviceItems, handleUpdateMTOServiceItemStatus }) => {
  const tableRows = serviceItems.map(({ id, code, submittedAt, serviceItem, details }, i) => {
    let detailSection;
    switch (code) {
      case 'DOFSIT':
      case 'DOASIT':
      case 'DOPSIT': {
        detailSection = (
          <div>
            <dl>{generateDetailText({ ZIP: details.pickupPostalCode, Reason: details.reason }, id)}</dl>
          </div>
        );
        break;
      }
      case 'DDFSIT':
      case 'DDASIT':
      case 'DDDSIT': {
        const { firstCustomerContact, secondCustomerContact } = details;
        detailSection = (
          <div>
            <dl>
              {firstCustomerContact &&
                generateDetailText(
                  {
                    'First Customer Contact': firstCustomerContact.timeMilitary,
                    'First Available Delivery Date': firstCustomerContact.firstAvailableDeliveryDate,
                  },
                  id,
                )}
              <div className={styles.customerContact}>
                {secondCustomerContact &&
                  generateDetailText(
                    {
                      'Second Customer Contact': secondCustomerContact.timeMilitary,
                      'Second Available Delivery Date': secondCustomerContact.firstAvailableDeliveryDate,
                    },
                    id,
                  )}
              </div>
            </dl>
          </div>
        );
        break;
      }
      case 'DCRT': {
        const { imgURL, description, itemDimensions, crateDimensions } = details;
        const itemDimensionFormat = `${convertFromThousandthInchToInch(
          itemDimensions?.length,
        )}"x${convertFromThousandthInchToInch(itemDimensions?.width)}"x${convertFromThousandthInchToInch(
          itemDimensions?.height,
        )}`;
        const crateDimensionFormat = `${convertFromThousandthInchToInch(
          crateDimensions?.length,
        )}"x${convertFromThousandthInchToInch(crateDimensions?.width)}"x${convertFromThousandthInchToInch(
          crateDimensions?.height,
        )}`;
        detailSection = (
          <div className={styles.detailImage}>
            <img
              className={styles.siThumbnail}
              alt={description}
              aria-labelledby={`si-thumbnail--caption-${i}`}
              src={imgURL}
            />
            <small id={`si-thumbnail--caption-${i}`}>
              <dl>
                <p className={styles.detailLine}>{description}</p>
                {itemDimensions && generateDetailText({ 'Item Dimensions': itemDimensionFormat }, id)}
                {crateDimensions && generateDetailText({ 'Crate Dimensions': crateDimensionFormat }, id)}
              </dl>
            </small>
          </div>
        );
        break;
      }
      case 'DOSHUT':
      case 'DDSHUT': {
        detailSection = (
          <div>
            <dl>{generateDetailText({ 'Estimated Weight': '', Reason: details.reason })}</dl>
          </div>
        );
        break;
      }
      default:
        detailSection = <div>—</div>;
    }

    return (
      <tr key={id}>
        <td className={styles.nameAndDate}>
          <p className={styles.codeName}>{serviceItem}</p>
          <p>{formatDate(submittedAt, 'DD MMM YYYY')}</p>
        </td>
        <td className={styles.detail}>{detailSection}</td>
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

import React from 'react';
import { PropTypes } from 'prop-types';

import styles from './ServiceItemDetails.module.scss';

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

const ServiceItemDetails = ({ className, id, code, details }) => {
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
                  'First Available Delivery Date': formatDate(
                    firstCustomerContact.firstAvailableDeliveryDate,
                    'DD MMM YYYY',
                  ),
                },
                id,
              )}
            <div className={styles.customerContact}>
              {secondCustomerContact &&
                generateDetailText(
                  {
                    'Second Customer Contact': secondCustomerContact.timeMilitary,
                    'Second Available Delivery Date': formatDate(
                      secondCustomerContact.firstAvailableDeliveryDate,
                      'DD MMM YYYY',
                    ),
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
      )}"`;
      const crateDimensionFormat = `${convertFromThousandthInchToInch(
        crateDimensions?.length,
      )}"x${convertFromThousandthInchToInch(crateDimensions?.width)}"x${convertFromThousandthInchToInch(
        crateDimensions?.height,
      )}"`;
      detailSection = (
        <div className={styles.detailImage}>
          {imgURL ? (
            <>
              <img
                className={styles.siThumbnail}
                alt={description}
                aria-labelledby={`si-thumbnail--caption-${id}`}
                src={imgURL}
              />
              <small className={styles.detailCaption} id={`si-thumbnail--caption-${id}`}>
                <dl>
                  <p className={styles.detailLine}>{description}</p>
                  {itemDimensions && generateDetailText({ 'Item Dimensions': itemDimensionFormat }, id)}
                  {crateDimensions && generateDetailText({ 'Crate Dimensions': crateDimensionFormat }, id)}
                </dl>
              </small>
            </>
          ) : (
            <dl>
              <p className={styles.detailLine}>{description}</p>
              {itemDimensions && generateDetailText({ 'Item Dimensions': itemDimensionFormat }, id)}
              {crateDimensions && generateDetailText({ 'Crate Dimensions': crateDimensionFormat }, id)}
            </dl>
          )}
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
  return <div className={className}>{detailSection}</div>;
};

ServiceItemDetails.propTypes = {
  className: PropTypes.string,
  id: PropTypes.string.isRequired,
  code: PropTypes.string.isRequired,
  details: PropTypes.shape({
    description: PropTypes.string,
    pickupPostalCode: PropTypes.string,
    reason: PropTypes.string,
    imgURL: PropTypes.string,
    itemDimensions: PropTypes.shape({ length: PropTypes.number, width: PropTypes.number, height: PropTypes.number }),
    crateDimensions: PropTypes.shape({ length: PropTypes.number, width: PropTypes.number, height: PropTypes.number }),
    firstCustomerContact: PropTypes.shape({
      timeMilitary: PropTypes.string,
      firstAvailableDeliveryDate: PropTypes.string,
    }),
    secondCustomerContact: PropTypes.shape({
      timeMilitary: PropTypes.string,
      firstAvailableDeliveryDate: PropTypes.string,
    }),
  }),
};

ServiceItemDetails.defaultProps = {
  className: undefined,
  details: {},
};

export default ServiceItemDetails;

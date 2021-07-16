import React from 'react';

import { ServiceItemDetailsShape } from '../../../types/serviceItems';

import styles from './ServiceItemDetails.module.scss';

import { formatDate } from 'shared/dates';
import { convertFromThousandthInchToInch, formatWeight } from 'shared/formatters';

function generateDetailText(details, id) {
  const detailList = Object.keys(details).map((detail) => (
    <div key={`${id}-${detail}`} className={styles.detailLine}>
      <dt className={styles.detailType}>{detail}:</dt> <dd>{details[`${detail}`]}</dd>
    </div>
  ));

  return detailList;
}

const ServiceItemDetails = ({ id, code, details }) => {
  let detailSection;
  switch (code) {
    case 'DOFSIT':
    case 'DOASIT':
    case 'DOPSIT': {
      detailSection = (
        <div>
          <dl>
            {generateDetailText(
              {
                ZIP: details.pickupPostalCode ? details.pickupPostalCode : '-',
                Reason: details.reason ? details.reason : '-',
              },
              id,
            )}
            {details.rejectionReason && generateDetailText({ 'Rejection reason': details.rejectionReason }, id)}
          </dl>
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
            {!firstCustomerContact &&
              generateDetailText({ 'First Customer Contact': '-', 'First Available Delivery Date': '-' })}
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
              {!secondCustomerContact &&
                generateDetailText({ 'Second Customer Contact': '-', 'Second Available Delivery Date': '-' })}
            </div>
            {generateDetailText({ Reason: details.reason ? details.reason : '-' })}
            {details.rejectionReason && generateDetailText({ 'Rejection reason': details.rejectionReason }, id)}
          </dl>
        </div>
      );
      break;
    }
    case 'DCRT': {
      const { description, itemDimensions, crateDimensions } = details;
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
        <div className={styles.detailCrating}>
          <dl>
            <p className={styles.detailLine}>{description}</p>
            {itemDimensions && generateDetailText({ 'Item Dimensions': itemDimensionFormat }, id)}
            {crateDimensions && generateDetailText({ 'Crate Dimensions': crateDimensionFormat }, id)}
            {generateDetailText({ Reason: details.reason ? details.reason : '-' })}
            {details.rejectionReason && generateDetailText({ 'Rejection reason': details.rejectionReason }, id)}
          </dl>
        </div>
      );
      break;
    }
    case 'DOSHUT':
    case 'DDSHUT': {
      const estimatedWeight = details.estimatedWeight != null ? formatWeight(details.estimatedWeight) : `— lbs`;
      detailSection = (
        <div>
          <dl>
            <div key={`${id}-estimatedWeight`} className={styles.detailLine}>
              <dd className={styles.detailType}>{estimatedWeight}</dd> <dt>estimated weight</dt>
            </div>
            {generateDetailText({ Reason: details.reason })}
            {details.rejectionReason && generateDetailText({ 'Rejection reason': details.rejectionReason }, id)}
          </dl>
        </div>
      );
      break;
    }
    default:
      detailSection = (
        <div>
          <div>—</div>
          <dl>{details.rejectionReason && generateDetailText({ 'Rejection reason': details.rejectionReason }, id)}</dl>
        </div>
      );
  }
  return <div>{detailSection}</div>;
};

ServiceItemDetails.propTypes = ServiceItemDetailsShape.isRequired;

ServiceItemDetails.defaultProps = {
  details: {},
};
export default ServiceItemDetails;

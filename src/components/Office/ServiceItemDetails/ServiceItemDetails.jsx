import React from 'react';
import classnames from 'classnames';

import { ServiceItemDetailsShape } from '../../../types/serviceItems';

import styles from './ServiceItemDetails.module.scss';

import { formatDate } from 'shared/dates';
import { formatWeight, convertFromThousandthInchToInch } from 'utils/formatters';

function generateDetailText(details, id, className) {
  const detailList = Object.keys(details).map((detail) => (
    <div key={`${id}-${detail}`} className={classnames(styles.detailLine, className)}>
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
                ZIP: details.SITPostalCode ? details.SITPostalCode : '-',
                Reason: details.reason ? details.reason : '-',
              },
              id,
            )}
            {details.rejectionReason &&
              generateDetailText({ 'Rejection reason': details.rejectionReason }, id, 'margin-top-2')}
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
            {details.rejectionReason &&
              generateDetailText({ 'Rejection reason': details.rejectionReason }, id, 'margin-top-2')}
          </dl>
        </div>
      );
      break;
    }
    case 'DCRT':
    case 'DCRTSA': {
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
            {description && generateDetailText({ Description: description }, id)}
            {itemDimensions && generateDetailText({ 'Item size': itemDimensionFormat }, id)}
            {crateDimensions && generateDetailText({ 'Crate size': crateDimensionFormat }, id)}
            {generateDetailText({ Reason: details.reason ? details.reason : '-' })}
            {details.rejectionReason &&
              generateDetailText({ 'Rejection reason': details.rejectionReason }, id, 'margin-top-2')}
          </dl>
        </div>
      );
      break;
    }
    case 'DUCRT': {
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
            {description && generateDetailText({ Description: description }, id)}
            {itemDimensions && generateDetailText({ 'Item size': itemDimensionFormat }, id)}
            {crateDimensions && generateDetailText({ 'Crate size': crateDimensionFormat }, id)}
            {details.rejectionReason &&
              generateDetailText({ 'Rejection reason': details.rejectionReason }, id, 'margin-top-2')}
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
            {details.rejectionReason &&
              generateDetailText({ 'Rejection reason': details.rejectionReason }, id, 'margin-top-2')}
          </dl>
        </div>
      );
      break;
    }
    default:
      detailSection = (
        <div>
          <div>—</div>
          <dl>
            {details.rejectionReason &&
              generateDetailText({ 'Rejection reason': details.rejectionReason }, id, 'margin-top-2')}
          </dl>
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

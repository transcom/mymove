import React from 'react';
import { isEmpty, sortBy } from 'lodash';
import classnames from 'classnames';

import { ServiceItemDetailsShape } from '../../../types/serviceItems';
import { trimFileName } from '../../../utils/serviceItems';

import styles from './ServiceItemDetails.module.scss';

import { formatDateWithUTC } from 'shared/dates';
import { formatWeight, convertFromThousandthInchToInch } from 'utils/formatters';

function generateDetailText(details, id, className) {
  const detailList = Object.keys(details).map((detail) => (
    <div key={`${id}-${detail}`} className={classnames(styles.detailLine, className)}>
      <dt className={styles.detailType}>{detail}:</dt> <dd>{details[`${detail}`]}</dd>
    </div>
  ));

  return detailList;
}

const generateSITDetailSection = (id, serviceRequestDocUploads, details, code, serviceItem, shipment) => {
  const { customerContacts } = details;
  const { sitStatus } = shipment;
  console.log('serviceItem', serviceItem);
  console.log('shipment', shipment);

  // Below we are using the sortBy func in lodash to sort the customer contacts
  // by the firstAvailableDeliveryDate field. sortBy returns a new
  // array with the elements in ascending order.
  const sortedCustomerContacts = sortBy(customerContacts, [
    (a) => {
      return new Date(a.firstAvailableDeliveryDate);
    },
  ]);
  const defaultDetailText = generateDetailText({
    'First available delivery date 1': '-',
    'Customer contact 1': '-',
  });

  const formatAddress = (address) => {
    if (!address || typeof address !== 'object') {
      return '';
    }
    const { city, state, postalCode } = address || {};
    const formattedCity = city ? `${city}, ` : '';
    const formattedState = state ? `${state} ` : '';
    const formattedPostalCode = postalCode || '';
    return `${formattedCity}${formattedState}${formattedPostalCode}`;
  };

  // DDDSIT (destination SIT delivery) & DOPSIT (origin SIT pickup)
  // DDFSIT (destination 1st day SIT) & DOFSIT (origin 1st day SIT)
  // DDASIT (destination add'l days SIT) & DOASIT (origin add'l days SIT)
  // DDSFSC (destination fuel surcharge) & DOSFSC (origin fuel surcharge)
  return (
    <div>
      <dl>
        {code === 'DDDSIT' || code === 'DOPSIT'
          ? generateDetailText({
              'Original Delivery Address': formatAddress(serviceItem.sitDestinationOriginalAddress) || '-',
              'Final Delivery Address': formatAddress(serviceItem.sitDestinationFinalAddress) || '-',
              'Delivery Miles': 'TBD',
            })
          : null}
        {code === 'DDFSIT' || code === 'DOFSIT'
          ? generateDetailText({
              'Original Delivery Address': formatAddress(serviceItem.sitDestinationOriginalAddress) || '-',
              'SIT entry date': details.sitEntryDate ? formatDateWithUTC(details.sitEntryDate, 'DD MMM YYYY') : '-',
            })
          : null}

        {code === 'DDASIT' || code === 'DOASIT'
          ? generateDetailText({
              'Original Delivery Address': formatAddress(serviceItem.sitDestinationOriginalAddress) || '-',
              "Add'l SIT Start Date": !serviceItem.sitEntryDate
                ? '-'
                : formatDateWithUTC(serviceItem.sitEntryDate, 'DD MMM YYYY'),
              '# of days approved': !shipment.sitDaysAllowance ? '-' : `${shipment.sitDaysAllowance} days`,
              'SIT expiration date': !sitStatus
                ? '-'
                : formatDateWithUTC(sitStatus.currentSIT.sitAllowanceEndDate, 'DD MMM YYYY'),
              'Customer Contacted HomeSafe': !serviceItem.sitCustomerContacted
                ? '-'
                : formatDateWithUTC(serviceItem.sitCustomerContacted, 'DD MMM YYYY'),
              'Customer Requested Del Date': !serviceItem.sitRequestedDelivery
                ? '-'
                : formatDateWithUTC(serviceItem.sitRequestedDelivery, 'DD MMM YYYY'),
              'SIT departure date': !serviceItem.sitDepartureDate
                ? '-'
                : formatDateWithUTC(serviceItem.sitDepartureDate, 'DD MMM YYYY'),
            })
          : null}

        {code === 'DDSFSC' || code === 'DOSFSC'
          ? generateDetailText({
              'SIT departure date': details.sitDepartureDate
                ? formatDateWithUTC(details.sitDepartureDate, 'DD MMM YYYY')
                : '-',
            })
          : null}

        {!isEmpty(sortedCustomerContacts)
          ? sortedCustomerContacts.map((contact, index) => (
              <>
                {generateDetailText(
                  {
                    [`First available delivery date ${index + 1}`]:
                      contact && contact.firstAvailableDeliveryDate
                        ? formatDateWithUTC(contact.firstAvailableDeliveryDate, 'DD MMM YYYY')
                        : '-',
                    [`Customer contact attempt ${index + 1}`]:
                      contact && contact.dateOfContact && contact.timeMilitary
                        ? `${formatDateWithUTC(contact.dateOfContact, 'DD MMM YYYY')}, ${contact.timeMilitary}`
                        : '-',
                  },
                  id,
                )}
              </>
            ))
          : defaultDetailText}
        {generateDetailText({ Reason: details.reason ? details.reason : '-' })}
        {details.rejectionReason &&
          generateDetailText({ 'Rejection reason': details.rejectionReason }, id, 'margin-top-2')}
        {!isEmpty(serviceRequestDocUploads) ? (
          <div className={styles.uploads}>
            <p className={styles.detailType}>Download service item documentation:</p>
            {serviceRequestDocUploads.map((file) => (
              <div className={styles.uploads}>
                <a href={file.url} download>
                  {trimFileName(file.filename)}
                </a>
              </div>
            ))}
          </div>
        ) : null}
      </dl>
    </div>
  );
};

const ServiceItemDetails = ({ id, code, details, serviceRequestDocs, serviceItem, shipment }) => {
  const serviceRequestDocUploads = serviceRequestDocs?.map((doc) => doc.uploads[0]);

  let detailSection;
  switch (code) {
    case 'DDFSIT':
    case 'DDASIT':
    case 'DDDSIT':
    case 'DDSFSC':
    case 'DOFSIT':
    case 'DOASIT':
    case 'DOPSIT':
    case 'DOSFSC': {
      detailSection = generateSITDetailSection(id, serviceRequestDocUploads, details, code, serviceItem, shipment);
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
            {!isEmpty(serviceRequestDocUploads) ? (
              <div className={styles.uploads}>
                <p className={styles.detailType}>Download service item documentation:</p>
                {serviceRequestDocUploads.map((file) => (
                  <div className={styles.uploads}>
                    <a href={file.url} download>
                      {trimFileName(file.filename)}
                    </a>
                  </div>
                ))}
              </div>
            ) : null}
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
            {!isEmpty(serviceRequestDocUploads) ? (
              <div className={styles.uploads}>
                <p className={styles.detailType}>Download service item documentation:</p>
                {serviceRequestDocUploads.map((file) => (
                  <div className={styles.uploads}>
                    <a href={file.url} download>
                      {trimFileName(file.filename)}
                    </a>
                  </div>
                ))}
              </div>
            ) : null}
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
            {!isEmpty(serviceRequestDocUploads) ? (
              <div className={styles.uploads}>
                <p className={styles.detailType}>Download service item documentation:</p>
                {serviceRequestDocUploads.map((file) => (
                  <div className={styles.uploads}>
                    <a href={file.url} download>
                      {trimFileName(file.filename)}
                    </a>
                  </div>
                ))}
              </div>
            ) : null}
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
            {!isEmpty(serviceRequestDocUploads) ? (
              <div className={styles.uploads}>
                <p className={styles.detailType}>Download service item documentation:</p>
                {serviceRequestDocUploads.map((file) => (
                  <div className={styles.uploads}>
                    <a href={file.url} download>
                      {trimFileName(file.filename)}
                    </a>
                  </div>
                ))}
              </div>
            ) : null}
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

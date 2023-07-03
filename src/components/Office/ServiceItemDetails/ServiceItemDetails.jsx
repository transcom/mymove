import React from 'react';
import { isEmpty, sortBy } from 'lodash';
import classnames from 'classnames';
import { Link } from 'react-router-dom';

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

const ServiceItemDetails = ({ id, code, details, serviceRequestDocs }) => {
  const serviceRequestDocUploads = serviceRequestDocs?.map((doc) => doc.uploads[0]);

  const trimFileName = (file) => {
    const splitName = file.split('/').pop();
    return splitName.substring(splitName.indexOf('-') + 1);
  };

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
            {!isEmpty(serviceRequestDocUploads) ? (
              <div className={styles.uploads}>
                <p className={styles.detailType}>Download service item documentation:</p>
                {serviceRequestDocUploads.map((file) => (
                  <div className={styles.uploads}>
                    {/* <a href={file.url}>{file.filename.split('/').pop().split('-').pop()}</a> */}
                    {/* <a href={file.url}>{trimFileName(file.filename)}</a> */}
                    <Link to={file.url}>{trimFileName(file.filename)}</Link>
                  </div>
                ))}
              </div>
            ) : null}
          </dl>
        </div>
      );
      break;
    }
    case 'DDFSIT':
    case 'DDASIT':
    case 'DDDSIT': {
      const { customerContacts } = details;
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
      detailSection = (
        <div>
          <dl>
            {!isEmpty(sortedCustomerContacts)
              ? sortedCustomerContacts.map((contact, index) => (
                  <>
                    {generateDetailText(
                      {
                        [`First available delivery date ${index + 1}`]:
                          contact && contact.firstAvailableDeliveryDate
                            ? formatDate(contact.firstAvailableDeliveryDate, 'DD MMM YYYY')
                            : '-',
                        [`Customer contact attempt ${index + 1}`]:
                          contact && contact.dateOfContact && contact.timeMilitary
                            ? `${formatDate(contact.dateOfContact, 'DD MMM YYYY')}, ${contact.timeMilitary}`
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
                    {/* <a href={file.url}>{trimFileName(file.filename)}</a> */}
                    <Link to={file.url}>{trimFileName(file.filename)}</Link>
                  </div>
                ))}
              </div>
            ) : null}
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
            {!isEmpty(serviceRequestDocUploads) ? (
              <div className={styles.uploads}>
                <p className={styles.detailType}>Download service item documentation:</p>
                {serviceRequestDocUploads.map((file) => (
                  <div className={styles.uploads}>
                    <Link href={file.url}>{trimFileName(file.filename)}</Link>
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
                    {/* <a href={file.url}>{trimFileName(file.filename)}</a> */}
                    <Link to={file.url}>{trimFileName(file.filename)}</Link>
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
                    {/* <a href={file.url}>{trimFileName(file.filename)}</a> */}
                    <Link to={file.url}>{trimFileName(file.filename)}</Link>
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
                    {/* <a href={file.url}>{trimFileName(file.filename)}</a> */}
                    <Link to={file.url} className="usa-link">
                      {trimFileName(file.filename)}
                    </Link>
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

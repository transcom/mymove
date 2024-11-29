import React from 'react';
import { isEmpty, sortBy } from 'lodash';
import classnames from 'classnames';
import moment from 'moment';

import { ServiceItemDetailsShape } from '../../../types/serviceItems';
import { trimFileName } from '../../../utils/serviceItems';

import styles from './ServiceItemDetails.module.scss';

import { ShipmentShape } from 'types/shipment';
import { SitStatusShape } from 'types/sitStatusShape';
import { formatDateWithUTC } from 'shared/dates';
import { formatCityStateAndPostalCode } from 'utils/shipmentDisplay';
import { formatWeight, convertFromThousandthInchToInch, formatCents, toDollarString } from 'utils/formatters';

function generateDetailText(details, id, className) {
  const detailList = Object.keys(details).map((detail) => (
    <div key={`${id}-${detail}`} className={classnames(styles.detailLine, className)}>
      <dt className={styles.detailType}>{detail}:</dt> <dd>{details[`${detail}`]}</dd>
    </div>
  ));

  return detailList;
}

const generateDestinationSITDetailSection = (id, serviceRequestDocUploads, details, code, shipment, sitStatus) => {
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
  const numberOfDaysApprovedForDOASIT = shipment.sitDaysAllowance ? shipment.sitDaysAllowance - 1 : 0;
  const sitEndDate =
    sitStatus &&
    sitStatus.currentSIT?.sitAuthorizedEndDate &&
    formatDateWithUTC(sitStatus.currentSIT.sitAuthorizedEndDate, 'DD MMM YYYY');
  const originalDeliveryAddress = details.sitDestinationOriginalAddress
    ? details.sitDestinationOriginalAddress
    : shipment.destinationAddress;

  return (
    <div>
      <dl>
        {code === 'DDFSIT'
          ? generateDetailText({
              'Original Delivery Address': originalDeliveryAddress
                ? formatCityStateAndPostalCode(originalDeliveryAddress)
                : '-',
              'SIT entry date': details.sitEntryDate ? formatDateWithUTC(details.sitEntryDate, 'DD MMM YYYY') : '-',
            })
          : null}
        {code === 'DDASIT' && (
          <>
            {generateDetailText(
              {
                'Original Delivery Address': originalDeliveryAddress
                  ? formatCityStateAndPostalCode(originalDeliveryAddress)
                  : '-',
                "Add'l SIT Start Date": details.sitEntryDate
                  ? moment.utc(details.sitEntryDate).add(1, 'days').format('DD MMM YYYY')
                  : '-',
                '# of days approved for': shipment.sitDaysAllowance ? `${numberOfDaysApprovedForDOASIT} days` : '-',
                'SIT expiration date': sitEndDate || '-',
              },
              id,
            )}
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
          </>
        )}
        {code === 'DDSFSC'
          ? generateDetailText(
              {
                'Original Delivery Address': originalDeliveryAddress
                  ? formatCityStateAndPostalCode(originalDeliveryAddress)
                  : '-',
                'Final Delivery Address':
                  details.sitDestinationFinalAddress && details.status !== 'SUBMITTED'
                    ? formatCityStateAndPostalCode(details.sitDestinationFinalAddress)
                    : '-',
                'Delivery miles out of SIT': details.sitDeliveryMiles ? details.sitDeliveryMiles : '-',
              },
              id,
            )
          : null}
        {code === 'DDDSIT' && (
          <>
            {generateDetailText(
              {
                'Original Delivery Address': originalDeliveryAddress
                  ? formatCityStateAndPostalCode(originalDeliveryAddress)
                  : '-',
                'Final Delivery Address':
                  details.sitDestinationFinalAddress && details.status !== 'SUBMITTED'
                    ? formatCityStateAndPostalCode(details.sitDestinationFinalAddress)
                    : '-',
                'Delivery miles out of SIT': details.sitDeliveryMiles ? details.sitDeliveryMiles : '-',
                'Customer contacted homesafe': details.sitCustomerContacted
                  ? formatDateWithUTC(details.sitCustomerContacted, 'DD MMM YYYY')
                  : '-',
                'Customer requested delivery date': details.sitRequestedDelivery
                  ? formatDateWithUTC(details.sitRequestedDelivery, 'DD MMM YYYY')
                  : '-',
                'SIT departure date': details.sitDepartureDate
                  ? formatDateWithUTC(details.sitDepartureDate, 'DD MMM YYYY')
                  : '-',
              },
              id,
            )}
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
          </>
        )}
        {code === 'DDFSIT' && (
          <>
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
          </>
        )}
      </dl>
    </div>
  );
};

const ServiceItemDetails = ({ id, code, details, serviceRequestDocs, shipment, sitStatus }) => {
  const serviceRequestDocUploads = serviceRequestDocs?.map((doc) => doc.uploads[0]);

  let detailSection;
  switch (code) {
    case 'DOFSIT': {
      detailSection = (
        <div>
          <dl>
            {generateDetailText(
              {
                'Original Pickup Address': details.sitOriginHHGOriginalAddress
                  ? formatCityStateAndPostalCode(details.sitOriginHHGOriginalAddress)
                  : '-',
                'SIT entry date': details.sitEntryDate ? formatDateWithUTC(details.sitEntryDate, 'DD MMM YYYY') : '-',
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
    case 'DOASIT': {
      const numberOfDaysApprovedForDOASIT = shipment.sitDaysAllowance ? shipment.sitDaysAllowance - 1 : 0;
      const sitEndDate =
        sitStatus &&
        sitStatus.currentSIT?.sitAuthorizedEndDate &&
        formatDateWithUTC(sitStatus.currentSIT.sitAuthorizedEndDate, 'DD MMM YYYY');

      detailSection = (
        <div>
          <dl>
            {generateDetailText(
              {
                'Original Pickup Address': details.sitOriginHHGOriginalAddress
                  ? formatCityStateAndPostalCode(details.sitOriginHHGOriginalAddress)
                  : '-',
                "Add'l SIT Start Date": details.sitEntryDate
                  ? moment.utc(details.sitEntryDate).add(1, 'days').format('DD MMM YYYY')
                  : '-',
                '# of days approved for': shipment.sitDaysAllowance ? `${numberOfDaysApprovedForDOASIT} days` : '-',
                'SIT expiration date': sitEndDate || '-',
                'Customer contacted homesafe': details.sitCustomerContacted
                  ? formatDateWithUTC(details.sitCustomerContacted, 'DD MMM YYYY')
                  : '-',
                'Customer requested delivery date': details.sitRequestedDelivery
                  ? formatDateWithUTC(details.sitRequestedDelivery, 'DD MMM YYYY')
                  : '-',
                'SIT departure date': details.sitDepartureDate
                  ? formatDateWithUTC(details.sitDepartureDate, 'DD MMM YYYY')
                  : '-',
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
    case 'DOPSIT': {
      detailSection = (
        <div>
          <dl>
            {generateDetailText(
              {
                'Original Pickup Address': details.sitOriginHHGOriginalAddress
                  ? formatCityStateAndPostalCode(details.sitOriginHHGOriginalAddress)
                  : '-',
                'Actual Pickup Address': details.sitOriginHHGActualAddress
                  ? formatCityStateAndPostalCode(details.sitOriginHHGActualAddress)
                  : '-',
                'Delivery miles into SIT': details.sitDeliveryMiles ? details.sitDeliveryMiles : '-',
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
    case 'DOSFSC': {
      detailSection = (
        <div>
          <dl>
            {generateDetailText(
              {
                'Original Pickup Address': details.sitOriginHHGOriginalAddress
                  ? formatCityStateAndPostalCode(details.sitOriginHHGOriginalAddress)
                  : '-',
                'Actual Pickup Address': details.sitOriginHHGActualAddress
                  ? formatCityStateAndPostalCode(details.sitOriginHHGActualAddress)
                  : '-',
                'Delivery miles into SIT': details.sitDeliveryMiles ? details.sitDeliveryMiles : '-',
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
    case 'DDFSIT':
    case 'DDASIT': {
      detailSection = generateDestinationSITDetailSection(
        id,
        serviceRequestDocUploads,
        details,
        code,
        shipment,
        sitStatus,
      );
      break;
    }
    case 'DDDSIT': {
      detailSection = generateDestinationSITDetailSection(
        id,
        serviceRequestDocUploads,
        details,
        code,
        shipment,
        sitStatus,
      );
      break;
    }
    case 'DDSFSC': {
      detailSection = generateDestinationSITDetailSection(
        id,
        serviceRequestDocUploads,
        details,
        code,
        shipment,
        sitStatus,
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
    case 'DLH':
    case 'DSH':
    case 'FSC':
    case 'DOP':
    case 'DDP':
    case 'DPK':
    case 'DUPK': {
      detailSection = (
        <div>
          <dl>
            {generateDetailText({
              'Estimated Price': details.estimatedPrice ? toDollarString(formatCents(details.estimatedPrice)) : '-',
            })}
          </dl>
        </div>
      );
      break;
    }
    case 'MS':
    case 'CS': {
      const { estimatedPrice } = details;
      detailSection = (
        <div>
          <dl>{estimatedPrice && generateDetailText({ Price: `$${formatCents(estimatedPrice)}` }, id)}</dl>
        </div>
      );
      break;
    }
    case 'ICRT': {
      const { description, itemDimensions, crateDimensions, market, externalCrate } = details;
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
            {externalCrate && generateDetailText({ 'External crate': 'Yes' }, id)}
            {market && generateDetailText({ Market: market }, id)}
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
    case 'IUCRT': {
      const { description, itemDimensions, crateDimensions, market } = details;
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
            {market && generateDetailText({ Market: market }, id)}
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

ServiceItemDetails.propTypes = {
  details: ServiceItemDetailsShape,
  shipment: ShipmentShape,
  sitStatus: SitStatusShape,
};

ServiceItemDetails.defaultProps = {
  details: {},
  shipment: {},
  sitStatus: undefined,
};
export default ServiceItemDetails;

import React from 'react';
import PropTypes from 'prop-types';

import styles from './LabeledDetails.module.scss';

import { shipmentTypes } from 'constants/shipments';
import {
  HistoryLogRecordShape,
  dbFieldToDisplayName,
  dbWeightFields,
  dbDateFields,
  optionFields,
} from 'constants/historyLogUIDisplayName';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatCustomerDate } from 'utils/formatters';

const retrieveTextToDisplay = (fieldName, value) => {
  const displayName = dbFieldToDisplayName[fieldName];
  let displayValue = value;

  if (displayName === dbFieldToDisplayName.storage_in_transit) {
    displayValue = `${displayValue} days`;
  } else if (dbWeightFields[fieldName]) {
    displayValue = `${displayValue} lbs`;
  } else if (optionFields[displayValue]) {
    displayValue = optionFields[displayValue];
  } else if (dbDateFields[fieldName]) {
    displayValue = formatCustomerDate(displayValue);
  }

  return {
    displayName,
    displayValue,
  };
};

const LabeledDetails = ({ historyRecord, getDetailsLabeledDetails }) => {
  let changedValuesToUse = historyRecord.changedValues;
  let shipmentDisplay = '';
  // run custom function to mutate changedValues to display if not null
  if (getDetailsLabeledDetails) {
    changedValuesToUse = getDetailsLabeledDetails(historyRecord);
  }

  // Check for shipment_type to use it as a header for the row
  if ('shipment_type' in changedValuesToUse) {
    shipmentDisplay = shipmentTypes[changedValuesToUse.shipment_type];
    shipmentDisplay += ' shipment';
    delete changedValuesToUse.shipment_type;
  }

  const dbFieldsToDisplay = Object.keys(dbFieldToDisplayName).filter((dbField) => {
    return changedValuesToUse[dbField];
  });

  return (
    <>
      <span className={styles.shipmentType}>{shipmentDisplay}</span>
      {dbFieldsToDisplay.map((modelField) => {
        const { displayName, displayValue } = retrieveTextToDisplay(modelField, changedValuesToUse[modelField]);

        return (
          <div key={modelField} className={descriptionListStyles.row}>
            <b>{displayName}</b>: {displayValue}
          </div>
        );
      })}
    </>
  );
};

LabeledDetails.propTypes = {
  historyRecord: HistoryLogRecordShape,
  getDetailsLabeledDetails: PropTypes.func,
};

LabeledDetails.defaultProps = {
  historyRecord: {},
  getDetailsLabeledDetails: null,
};

export default LabeledDetails;

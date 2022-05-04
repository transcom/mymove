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
  let changeValuesToUse = historyRecord.changedValues;
  let shipmentDisplay = '';
  // run custom function to mutate changedValues to display if not null
  if (getDetailsLabeledDetails) {
    changeValuesToUse = getDetailsLabeledDetails(historyRecord);
  }

  // Check for shipment_type in values that need changing
  if (changeValuesToUse.shipment_type !== null) {
    shipmentDisplay = shipmentTypes[changeValuesToUse.shipment_type];
    shipmentDisplay += ' shipment';
    delete changeValuesToUse.shipment_type;
  }

  const dbFieldsToDisplay = Object.keys(dbFieldToDisplayName).filter((dbField) => {
    return changeValuesToUse[dbField];
  });

  return (
    <div>
      <span className={styles.shipmentType}>{shipmentDisplay}</span>
      {dbFieldsToDisplay.map((modelField) => {
        const { displayName, displayValue } = retrieveTextToDisplay(modelField, changeValuesToUse[modelField]);

        return (
          <div key={modelField} className={descriptionListStyles.row}>
            <b>{displayName}</b>: {displayValue}
          </div>
        );
      })}
    </div>
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

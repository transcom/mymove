import React from 'react';
import PropTypes from 'prop-types';

import {
  HistoryLogRecordShape,
  dbFieldToDisplayName,
  dbWeightFields,
  optionFields,
} from 'constants/historyLogUIDisplayName';
import descriptionListStyles from 'styles/descriptionList.module.scss';

const retrieveTextToDisplay = (fieldName, value) => {
  const displayName = dbFieldToDisplayName[fieldName];
  let displayValue = value;

  if (displayName === dbFieldToDisplayName.storage_in_transit) {
    displayValue = `${displayValue} days`;
  } else if (dbWeightFields.includes(fieldName)) {
    displayValue = `${displayValue} lbs`;
  } else if (optionFields[displayValue]) {
    displayValue = optionFields[displayValue];
  }

  return {
    displayName,
    displayValue,
  };
};

const LabeledDetails = ({ historyRecord, getDetailsLabeledDetails }) => {
  let changeValuesToUse = historyRecord.changedValues;
  // run custom function to mutate changedValues to display if not null
  if (getDetailsLabeledDetails) {
    changeValuesToUse = getDetailsLabeledDetails(historyRecord);
  }

  const dbFieldsToDisplay = Object.keys(dbFieldToDisplayName).filter((dbField) => {
    return changeValuesToUse[dbField];
  });

  return (
    <div>
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

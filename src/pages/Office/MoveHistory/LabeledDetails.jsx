import React from 'react';

import { HistoryLogValuesShape, dbFieldToDisplayName, dbWeightFields } from 'constants/historyLogUIDisplayName';
import descriptionListStyles from 'styles/descriptionList.module.scss';

const retrieveTextToDisplay = (fieldName, value) => {
  const displayName = dbFieldToDisplayName[fieldName];
  let displayValue = value;

  if (displayName === dbFieldToDisplayName.storage_in_transit) {
    displayValue = `${displayValue} days`;
  } else if (dbWeightFields.includes(fieldName)) {
    displayValue = `${displayValue} lbs`;
  }

  return {
    displayName,
    displayValue,
  };
};

const LabeledDetails = ({ changedValues }) => {
  const dbFieldsToDisplay = Object.keys(dbFieldToDisplayName).filter((dbField) => {
    return changedValues[dbField];
  });

  return (
    <div>
      {dbFieldsToDisplay.map((modelField) => {
        const { displayName, displayValue } = retrieveTextToDisplay(modelField, changedValues[modelField]);

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
  changedValues: HistoryLogValuesShape,
};

LabeledDetails.defaultProps = {
  changedValues: {},
};

export default LabeledDetails;

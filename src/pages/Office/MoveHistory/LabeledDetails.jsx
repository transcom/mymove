import React from 'react';
import PropTypes from 'prop-types';

import {
  HistoryLogValuesShape,
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

const LabeledDetails = ({ changedValues, context, getDetailsLabeledDetails }) => {
  let changeValuesToUse = changedValues;
  // run custom function to mutate changedValues to display if not null
  if (getDetailsLabeledDetails) {
    changeValuesToUse = getDetailsLabeledDetails({ changedValues, context });
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
  changedValues: HistoryLogValuesShape,
  context: PropTypes.arrayOf(PropTypes.object),
  getDetailsLabeledDetails: PropTypes.func,
};

LabeledDetails.defaultProps = {
  changedValues: {},
  context: null,
  getDetailsLabeledDetails: null,
};

export default LabeledDetails;

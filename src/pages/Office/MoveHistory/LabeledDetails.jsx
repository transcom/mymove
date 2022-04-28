import React from 'react';
import PropTypes from 'prop-types';

import {
  HistoryLogValuesShape,
  dbFieldToDisplayName,
  dbWeightFields,
  HistoryLogContextShape,
  optionFields,
} from 'constants/historyLogUIDisplayName';
import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatMoveHistoryAgent } from 'utils/formatters';

const LabeledDetails = ({ changedValues, oldValues, context, getDetailsLabeledDetails }) => {
  const backfilledChangedValues = {
    email: oldValues.email,
    first_name: oldValues.first_name,
    last_name: oldValues.last_name,
    phone: oldValues.phone,
    ...changedValues,
  };

  let changedValuesWithFormattedItems = {
    ...changedValues,
    agent: formatMoveHistoryAgent(backfilledChangedValues),
  };

  // run custom function to mutate changedValues to display if not null
  if (getDetailsLabeledDetails) {
    changedValuesWithFormattedItems = getDetailsLabeledDetails({ changedValuesWithFormattedItems, oldValues, context });
  }

  const dbFieldsToDisplay = Object.keys(dbFieldToDisplayName).filter((dbField) => {
    return changedValuesWithFormattedItems[dbField];
  });

  const retrieveTextToDisplay = (fieldName, value) => {
    let displayName = dbFieldToDisplayName[fieldName];
    let displayValue = value;

    if (displayName === dbFieldToDisplayName.storage_in_transit) {
      displayValue = `${displayValue} days`;
    } else if (dbWeightFields.includes(fieldName)) {
      displayValue = `${displayValue} lbs`;
    } else if (optionFields[displayValue]) {
      displayValue = optionFields[displayValue];
    } else if (displayName === dbFieldToDisplayName.agent) {
      const agentType = changedValues.agent_type ?? oldValues.agent_type;

      if (agentType === 'RECEIVING_AGENT') {
        displayName = 'Receiving agent';
      } else if (agentType === 'RELEASING_AGENT') {
        displayName = 'Releasing agent';
      }
    }

    return {
      displayName,
      displayValue,
    };
  };

  return (
    <div>
      {dbFieldsToDisplay.map((modelField) => {
        const { displayName, displayValue } = retrieveTextToDisplay(
          modelField,
          changedValuesWithFormattedItems[modelField],
        );

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
  oldValues: HistoryLogValuesShape,
  context: HistoryLogContextShape,
  getDetailsLabeledDetails: PropTypes.func,
};

LabeledDetails.defaultProps = {
  changedValues: {},
  oldValues: {},
  context: [],
  getDetailsLabeledDetails: null,
};

export default LabeledDetails;

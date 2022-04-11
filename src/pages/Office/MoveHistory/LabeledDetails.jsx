import React from 'react';

import { HistoryLogValuesShape, dbFieldsToModel, dbFieldToDisplayName } from 'constants/historyLogUIDisplayName';
import descriptionListStyles from 'styles/descriptionList.module.scss';

const LabeledDetails = ({ changedValues }) => {
  const dbFieldsToDisplay = Object.keys(dbFieldsToModel).filter((dbField) => {
    return changedValues[dbField];
  });

  return (
    <div>
      {dbFieldsToDisplay.map((modelField) => (
        <div key={modelField} className={descriptionListStyles.row}>
          <b>{dbFieldToDisplayName[modelField]}</b>: {changedValues[modelField]}
        </div>
      ))}
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

import React from 'react';
import { string, bool } from 'prop-types';

import serviceInfoTableStyles from './ServiceInfoTable.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';

const ServiceInfoTable = ({
  affiliation,
  currentDutyStationName,
  currentDutyStationPhone,
  edipi,
  firstName,
  isEditable,
  lastName,
  editURL,
  rank,
}) => {
  return (
    <div className={serviceInfoTableStyles.serviceInfoContainer}>
      <div className={serviceInfoTableStyles.reviewHeader}>
        <h2>Service info</h2>
        {isEditable && <a href={editURL}>Edit</a>}
      </div>
      {!isEditable && (
        <div>
          To change information in this section, contact the {currentDutyStationName} transportation office{' '}
          {currentDutyStationPhone ? ` at ${currentDutyStationPhone}.` : '.'}
        </div>
      )}
      <div className={serviceInfoTableStyles.serviceInfoSection}>
        <dl className={descriptionListStyles.descriptionList}>
          <div className={descriptionListStyles.row}>
            <dt>Name</dt>
            <dd>
              {firstName} {lastName}
            </dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Branch</dt>
            <dd>{affiliation}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Rank</dt>
            <dd>{rank}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>DOD ID#</dt>
            <dd>{edipi}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Current duty station</dt>
            <dd>{currentDutyStationName}</dd>
          </div>
        </dl>
      </div>
    </div>
  );
};

ServiceInfoTable.propTypes = {
  affiliation: string.isRequired,
  currentDutyStationName: string.isRequired,
  currentDutyStationPhone: string,
  edipi: string.isRequired,
  firstName: string.isRequired,
  isEditable: bool,
  lastName: string.isRequired,
  editURL: string,
  rank: string.isRequired,
};

ServiceInfoTable.defaultProps = {
  currentDutyStationPhone: '',
  editURL: '',
  isEditable: true,
};

export default ServiceInfoTable;

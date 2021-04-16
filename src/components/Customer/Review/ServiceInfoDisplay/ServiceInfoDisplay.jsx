import React from 'react';
import { string, bool } from 'prop-types';
import { Link } from 'react-router-dom';

import serviceInfoDisplayStyles from './ServiceInfoDisplay.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';

const ServiceInfoDisplay = ({
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
    <div className={serviceInfoDisplayStyles.serviceInfoContainer}>
      <div className={serviceInfoDisplayStyles.header}>
        <h2>Service info</h2>
        {isEditable && <Link to={editURL}>Edit</Link>}
      </div>
      {!isEditable && (
        <div className={serviceInfoDisplayStyles.whoToContactContainer}>
          To change information in this section, contact the {currentDutyStationName} transportation office{' '}
          {currentDutyStationPhone ? ` at ${currentDutyStationPhone}.` : '.'}
        </div>
      )}
      <div className={serviceInfoDisplayStyles.serviceInfoSection}>
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
            <dt>DoD ID#</dt>
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

ServiceInfoDisplay.propTypes = {
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

ServiceInfoDisplay.defaultProps = {
  currentDutyStationPhone: '',
  editURL: '',
  isEditable: true,
};

export default ServiceInfoDisplay;

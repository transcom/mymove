import React from 'react';
import { string, bool } from 'prop-types';
import { Link } from 'react-router-dom';

import serviceInfoDisplayStyles from './ServiceInfoDisplay.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';

const ServiceInfoDisplay = ({
  affiliation,
  originDutyStationName,
  originTransportationOfficeName,
  originTransportationOfficePhone,
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
          To change information in this section, contact the {originTransportationOfficeName} transportation office
          {originTransportationOfficePhone ? ` at ${originTransportationOfficePhone}.` : '.'}
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
            <dd>{originDutyStationName}</dd>
          </div>
        </dl>
      </div>
    </div>
  );
};

ServiceInfoDisplay.propTypes = {
  affiliation: string.isRequired,
  originDutyStationName: string.isRequired,
  originTransportationOfficeName: string.isRequired,
  originTransportationOfficePhone: string,
  edipi: string.isRequired,
  firstName: string.isRequired,
  isEditable: bool,
  lastName: string.isRequired,
  editURL: string,
  rank: string.isRequired,
};

ServiceInfoDisplay.defaultProps = {
  originTransportationOfficePhone: '',
  editURL: '',
  isEditable: true,
};

export default ServiceInfoDisplay;

import React from 'react';
import { string, bool } from 'prop-types';
import { Link, useLocation } from 'react-router-dom';

import serviceInfoDisplayStyles from './ServiceInfoDisplay.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';

const editButtonStyle = serviceInfoDisplayStyles['edit-btn'];

const ServiceInfoDisplay = ({
  affiliation,
  originTransportationOfficeName,
  originTransportationOfficePhone,
  edipi,
  emplid,
  firstName,
  isEditable,
  showMessage,
  lastName,
  editURL,
}) => {
  const { state } = useLocation();

  return (
    <div className={serviceInfoDisplayStyles.serviceInfoContainer}>
      <div className={serviceInfoDisplayStyles.header}>
        <h2>Service info</h2>
        {isEditable && (
          <Link className={editButtonStyle} to={editURL} state={state}>
            Edit
          </Link>
        )}
      </div>
      {!isEditable && showMessage && (
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
            <dt>DoD ID#</dt>
            <dd>{edipi}</dd>
          </div>

          {affiliation === 'Coast Guard' && (
            <div className={descriptionListStyles.row}>
              <dt>EMPLID</dt>
              <dd>{emplid}</dd>
            </div>
          )}
        </dl>
      </div>
    </div>
  );
};

ServiceInfoDisplay.propTypes = {
  affiliation: string.isRequired,
  originTransportationOfficeName: string.isRequired,
  originTransportationOfficePhone: string,
  edipi: string.isRequired,
  firstName: string.isRequired,
  isEditable: bool,
  showMessage: bool,
  lastName: string.isRequired,
  editURL: string,
};

ServiceInfoDisplay.defaultProps = {
  originTransportationOfficePhone: '',
  editURL: '',
  isEditable: true,
  showMessage: false,
};

export default ServiceInfoDisplay;

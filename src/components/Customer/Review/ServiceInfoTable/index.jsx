import React from 'react';
import classnames from 'classnames';
import { string, func, bool } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import reviewStyles from '../Review.module.scss';

import serviceInfoTableStyles from './ServiceInfoTable.module.scss';

const ServiceInfoTable = ({
  affiliation,
  currentDutyStationName,
  currentDutyStationPhone,
  edipi,
  firstName,
  isEditable,
  lastName,
  onEditClick,
  rank,
}) => {
  const containerClassNames = classnames(
    reviewStyles['review-container'],
    reviewStyles['profile-container'],
    serviceInfoTableStyles.ServiceInfoTable,
  );
  const tableClassNames = classnames('table--stacked', reviewStyles['review-table']);
  return (
    <div className={containerClassNames}>
      <div className={classnames(reviewStyles['review-header'], serviceInfoTableStyles.ReviewHeader)}>
        <h2>Service info</h2>
        {isEditable && (
          <Button unstyled className={reviewStyles['edit-btn']} data-testid="edit-profile-table" onClick={onEditClick}>
            Edit
          </Button>
        )}
      </div>
      {!isEditable && (
        <div>
          To change information in this section, contact the {currentDutyStationName} transportation office{' '}
          {currentDutyStationPhone ? ` at ${currentDutyStationPhone}.` : '.'}
        </div>
      )}
      <table className={tableClassNames}>
        <colgroup>
          <col />
          <col />
        </colgroup>
        <tbody>
          <tr>
            <th scope="row">Name</th>
            <td>
              {firstName} {lastName}
            </td>
          </tr>
          <tr>
            <th scope="row">Branch</th>
            <td>{affiliation}</td>
          </tr>
          <tr>
            <th scope="row">Rank</th>
            <td>{rank}</td>
          </tr>
          <tr>
            <th scope="row">DOD ID#</th>
            <td>{edipi}</td>
          </tr>
          <tr>
            <th className={reviewStyles['table-divider-top']} scope="row" style={{ borderBottom: 'none' }}>
              Current duty station
            </th>
            <td className={reviewStyles['table-divider-top']} style={{ borderBottom: 'none' }}>
              {currentDutyStationName}
            </td>
          </tr>
        </tbody>
      </table>
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
  onEditClick: func,
  rank: string.isRequired,
};

ServiceInfoTable.defaultProps = {
  currentDutyStationPhone: '',
  onEditClick: () => {},
  isEditable: true,
};

export default ServiceInfoTable;

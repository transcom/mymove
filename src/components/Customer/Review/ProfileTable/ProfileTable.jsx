import React from 'react';
import classnames from 'classnames';
import { string, func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import reviewStyles from '../Review.module.scss';

import { customerRoutes } from 'constants/routes';

const ProfileTable = ({
  affiliation,
  city,
  currentDutyStationName,
  edipi,
  email,
  firstName,
  lastName,
  onEditClick,
  postalCode,
  rank,
  state,
  streetAddress1,
  streetAddress2,
  telephone,
}) => {
  const containerClassNames = classnames(reviewStyles['review-container'], reviewStyles['profile-container']);
  const tableClassNames = classnames('table--stacked', reviewStyles['review-table']);
  const editProfilePath = customerRoutes.PROFILE_PATH;
  return (
    <div className={containerClassNames}>
      <div className={reviewStyles['review-header']}>
        <h2>Profile</h2>
        <Button
          unstyled
          className={reviewStyles['edit-btn']}
          data-testid="edit-profile-table"
          onClick={() => onEditClick(editProfilePath)}
        >
          Edit
        </Button>
      </div>
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
              Current duty location
            </th>
            <td className={reviewStyles['table-divider-top']} style={{ borderBottom: 'none' }}>
              {currentDutyStationName}
            </td>
          </tr>
          <tr>
            <th scope="row">Best contact phone</th>
            <td>{telephone}</td>
          </tr>
          <tr>
            <th scope="row">Personal email</th>
            <td>{email}</td>
          </tr>
          <tr>
            <th scope="row">Current mailing address</th>
            <td>
              {streetAddress1} {streetAddress2}
              <br />
              {city}, {state} {postalCode}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );
};

ProfileTable.propTypes = {
  affiliation: string.isRequired,
  city: string.isRequired,
  currentDutyStationName: string.isRequired,
  edipi: string.isRequired,
  email: string.isRequired,
  firstName: string.isRequired,
  lastName: string.isRequired,
  onEditClick: func.isRequired,
  postalCode: string.isRequired,
  rank: string.isRequired,
  state: string.isRequired,
  streetAddress1: string.isRequired,
  streetAddress2: string,
  telephone: string.isRequired,
};

ProfileTable.defaultProps = {
  streetAddress2: '',
};

export default ProfileTable;

/* eslint-ignore */
import React from 'react';
import classnames from 'classnames';
import { string } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import TableDivider from '../TableDivider';
import reviewStyles from '../Review.module.scss';

const ProfileTable = ({
  firstName,
  lastName,
  affiliation,
  rank,
  edipi,
  currentDutyStationName,
  telephone,
  email,
  streetAddress1,
  streetAddress2,
  city,
  state,
  postalCode,
}) => {
  const containerClassNames = classnames(reviewStyles['review-container'], reviewStyles['profile-container']);
  const tableClassNames = classnames('table--stacked', reviewStyles['review-table']);
  return (
    <div className={containerClassNames}>
      <div className={reviewStyles['review-header']}>
        <h3>Profile</h3>
        <Button unstyled className={reviewStyles['edit-btn']}>
          Edit
        </Button>
      </div>
      <table className={tableClassNames}>
        <colgroup>
          <col style={{ width: '40%' }} />
          <col style={{ width: '60%' }} />
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
          <TableDivider />
          <tr>
            <th scope="row" style={{ borderTop: 'none' }}>
              Contact info
            </th>
            <td style={{ borderTop: 'none' }} />
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
  firstName: string.isRequired,
  lastName: string.isRequired,
  affiliation: string.isRequired,
  rank: string.isRequired,
  edipi: string.isRequired,
  currentDutyStationName: string.isRequired,
  telephone: string.isRequired,
  email: string.isRequired,
  streetAddress1: string.isRequired,
  streetAddress2: string,
  city: string.isRequired,
  state: string.isRequired,
  postalCode: string.isRequired,
};

ProfileTable.defaultProps = {
  streetAddress2: '',
};

export default ProfileTable;

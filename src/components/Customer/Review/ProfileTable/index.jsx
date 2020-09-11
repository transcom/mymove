/* eslint-ignore */
import React from 'react';
import { shape, string } from 'prop-types';
import { Link } from 'react-router-dom';

import TableDivider from '../TableDivider';
import reviewStyles from '../Review.module.scss';

const ProfileTable = ({ serviceMember }) => (
  <div className="review-container">
    <div className="stackedtable-header">
      <h3>Profile</h3>
      <Link>Edit</Link>
    </div>
    <table className={`table--stacked ${reviewStyles['review-table']}`}>
      <colgroup>
        <col style={{ width: '40%' }} />
        <col style={{ width: '60%' }} />
      </colgroup>
      <tbody>
        <tr>
          <th scope="row">Name</th>
          <td>
            {serviceMember.first_name} {serviceMember.last_name}
          </td>
        </tr>
        <tr>
          <th scope="row">Branch</th>
          <td>{serviceMember.affiliation}</td>
        </tr>
        <tr>
          <th scope="row">Rank</th>
          <td>{serviceMember.rank}</td>
        </tr>
        <tr>
          <th scope="row">DOD ID#</th>
          <td>{serviceMember.edipi}</td>
        </tr>
        <tr>
          <th scope="row" style={{ borderBottom: 'none' }}>
            Current duty station
          </th>
          <td style={{ borderBottom: 'none' }}>{serviceMember.current_station.name}</td>
        </tr>
        <TableDivider className="" />
        <tr>
          <th scope="row" style={{ borderTop: 'none' }}>
            Contact info
          </th>
          <td style={{ borderTop: 'none' }} />
        </tr>
        <tr>
          <th scope="row">Best contact phone</th>
          <td>{serviceMember.telephone}</td>
        </tr>
        <tr>
          <th scope="row">Personal email</th>
          <td>{serviceMember.telephone}</td>
        </tr>
        <tr>
          <th scope="row">Current mailing address</th>
          <td>
            {serviceMember.residential_address.street_address_1} {serviceMember.residential_address.street_address_2}
            <br />
            {serviceMember.residential_address.city}, {serviceMember.residential_address.state}{' '}
            {serviceMember.residential_address.postal_code}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
);

ProfileTable.propTypes = {
  serviceMember: shape({
    first_name: string.isRequired,
    last_name: string.isRequired,
    affiliation: string.isRequired,
    rank: string.isRequired,
    edipi: string.isRequired,
    current_station: shape({
      name: string.isRequired,
    }).isRequired,
    telephone: string.isRequired,
    personal_email: string.isRequired,
    residential_address: shape({
      street_address_1: string.isRequired,
      street_address_2: string.isRequired,
      city: string.isRequired,
      state: string.isRequired,
      postal_code: string.isRequired,
    }).isRequired,
  }).isRequired,
};

export default ProfileTable;

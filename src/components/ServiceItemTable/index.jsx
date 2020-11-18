import React from 'react';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheck } from '@fortawesome/free-solid-svg-icons';

const ServiceItemTable = () => (
  <div className="table--service-item">
    <table>
      <col />
      <col style={{ width: '50%' }} />
      <col />
      <thead className="table--small">
        <tr>
          <th>Date approved</th>
          <th>Service item</th>
          <th>Code</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>
            <span className="gray-out">
              <FontAwesomeIcon icon={faCheck} />
            </span>
            05 Feb 2020
            <span className="gray-out">RJB</span>
          </td>
          <td>Domestic line haul</td>
          <td>DLH</td>
        </tr>
        <tr>
          <td>
            <span className="gray-out">
              <FontAwesomeIcon icon={faCheck} />
            </span>
            05 Feb 2020
            <span className="gray-out">RJB</span>
          </td>
          <td>Domestic line haul</td>
          <td>DLH</td>
        </tr>
      </tbody>
    </table>
  </div>
);

export default ServiceItemTable;

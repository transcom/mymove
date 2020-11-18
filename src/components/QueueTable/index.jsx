import React from 'react';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faChevronLeft, faChevronRight } from '@fortawesome/free-solid-svg-icons';

const QueueTable = () => (
  <div className="table--queue">
    <table>
      <thead style={{ height: '55px' }}>
        <tr>
          <th className="sortAscending">Table Cell Content</th>
          <th>Status</th>
          <th>Confirmation</th>
          <th>Branch</th>
          <th>Original duty station</th>
          <th>Last modified by</th>
        </tr>
      </thead>
      <tbody>
        <tr className="filter">
          <td>
            <select className="usa-select">
              <option>All</option>
            </select>
          </td>
          <td>
            <select className="usa-select">
              <option>All</option>
            </select>
          </td>
          <td>
            <input className="usa-input" id="input-type-text" name="input-type-text" type="text" />
          </td>
          <td>
            <select className="usa-select">
              <option>All</option>
            </select>
          </td>
          <td>
            <input className="usa-input" id="input-type-text" name="input-type-text" type="text" />
          </td>
          <td>
            <input className="usa-input" id="input-type-text" name="input-type-text" type="text" />
          </td>
        </tr>
        <tr>
          <td>Payment requested</td>
          <td>
            <a href="#">Clark-Nuñez, Sofía</a>
          </td>
          <td>GIW13</td>
          <td>Marines</td>
          <td>Camp Pendleton</td>
          <td>SW - J. Garimundi</td>
        </tr>
        <tr>
          <td>Payment requested</td>
          <td>
            <a href="#">Clark-Nuñez, Sofía</a>
          </td>
          <td>GIW13</td>
          <td>Marines</td>
          <td>Camp Pendleton</td>
          <td>SW - J. Garimundi</td>
        </tr>
      </tbody>
    </table>
    <div className="display-flex">
      <div className="tcontrol--rows-per-page">
        <select className="usa-select" name="table-rows-per-page">
          <option value="1">1</option>
          <option value="2">2</option>
          <option value="3">3</option>
        </select>
        <p>rows per page</p>
      </div>
      <div className="tcontrol--pagination">
        <Button disabled className="usa-button--unstyled">
          <span className="icon">
            <FontAwesomeIcon icon={faChevronLeft} />
          </span>
          <span>Prev</span>
        </Button>
        <select className="usa-select" name="table-pagination">
          <option value="1">1</option>
          <option value="2">2</option>
          <option value="3">3</option>
        </select>
        <Button className="usa-button--unstyled">
          <span>Next</span>
          <span className="icon">
            <FontAwesomeIcon icon={faChevronRight} />
          </span>
        </Button>
      </div>
    </div>
  </div>
);

export default QueueTable;

import React from 'react';
import { Link } from 'react-router-dom';
import { get } from 'lodash';
import PropTypes from 'prop-types';

import ppmBlack from 'shared/icon/ppm-black.svg';
import { formatCentsRange, formatCents } from 'shared/formatters';
import { formatDateSM } from 'shared/formatters';

import './Review.css';

export default function PPMShipmentSummary(props) {
  const { movePath, ppm } = props;

  const editDateAndLocationAddress = movePath + '/edit-date-and-location';
  const editWeightAddress = movePath + '/edit-weight';
  const privateStorageString = get(ppm, 'estimated_storage_reimbursement')
    ? `(spend up to ${ppm.estimated_storage_reimbursement.toLocaleString()} on private storage)`
    : '';
  const sitDisplay = get(ppm, 'has_sit', false)
    ? `${ppm.days_in_storage} days ${privateStorageString}`
    : 'Not requested';

  return (
    <div className="usa-grid-full ppm-container">
      <h3>
        <img src={ppmBlack} alt="PPM shipment" /> Shipment - You move your stuff
        (PPM)
      </h3>
      <div className="usa-width-one-half review-section ppm-review-section">
        <table>
          <tbody>
            <tr>
              <th>
                Dates &amp; Locations
                <span className="align-right">
                  <Link to={editDateAndLocationAddress}>Edit</Link>
                </span>
              </th>
            </tr>
            <tr>
              <td> Move Date: </td>
              <td>{formatDateSM(get(ppm, 'planned_move_date'))}</td>
            </tr>
            <tr>
              <td> Pickup ZIP Code: </td>
              <td> {ppm.pickup_postal_code}</td>
            </tr>
            {ppm.has_additional_postal_code && (
              <tr>
                <td> Additional Pickup: </td>
                <td> {ppm.additional_pickup_postal_code}</td>
              </tr>
            )}
            <tr>
              <td> Delivery ZIP Code: </td>
              <td> {ppm.destination_postal_code}</td>
            </tr>
            <tr>
              <td> Storage: </td>
              <td>{sitDisplay}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div className="usa-width-one-half review-section ppm-review-section">
        <table>
          <tbody>
            <tr>
              <th>
                Weight
                <span className="align-right">
                  <Link to={editWeightAddress}>Edit</Link>
                </span>
              </th>
            </tr>
            <tr>
              <td> Estimated Weight: </td>
              <td> {ppm.weight_estimate.toLocaleString()} lbs</td>
            </tr>
            <tr>
              <td> Estimated PPM Incentive: </td>
              <td>
                {' '}
                {formatCentsRange(
                  ppm.incentive_estimate_min,
                  ppm.incentive_estimate_max,
                )}
              </td>
            </tr>
            {ppm.has_requested_advance && (
              <tr>
                <td> Advance: </td>
                <td> ${formatCents(ppm.advance.requested_amount)}</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

PPMShipmentSummary.propTypes = {
  ppm: PropTypes.object.required,
  movePath: PropTypes.string.required,
};

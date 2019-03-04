import React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import { get } from 'lodash';
import PropTypes from 'prop-types';

import { selectReimbursement } from 'shared/Entities/modules/ppms';
import ppmBlack from 'shared/icon/ppm-black.svg';
import { formatCentsRange, formatCents } from 'shared/formatters';
import { formatDateSM } from 'shared/formatters';

import './Review.css';

function PPMShipmentSummary(props) {
  const { advance, movePath, ppm, isHHGPPMComboMove } = props;

  const editDateAndLocationAddress = movePath + '/edit-date-and-location';
  const editWeightAddress = movePath + '/edit-weight';

  const privateStorageString = get(ppm, 'estimated_storage_reimbursement')
    ? `(spend up to ${formatCents(ppm.estimated_storage_reimbursement)} on private storage)`
    : '';
  const sitDisplay = get(ppm, 'has_sit', false)
    ? `${ppm.days_in_storage} days ${privateStorageString}`
    : 'Not requested';

  return (
    <div className="usa-grid-full ppm-container">
      <h3>
        <img src={ppmBlack} alt="PPM shipment" /> Shipment - You move your stuff (PPM)
      </h3>
      <div className="usa-width-one-half review-section ppm-review-section">
        <p className="heading">
          Dates & Locations
          <span className="edit-section-link">
            <Link data-cy="edit-ppm-dates" to={editDateAndLocationAddress}>
              Edit
            </Link>
          </span>
        </p>

        <table>
          <tbody>
            <tr>
              <td> Move Date: </td>
              <td>{formatDateSM(get(ppm, 'original_move_date'))}</td>
            </tr>
            <tr>
              <td> Pickup ZIP Code: </td>
              <td> {ppm && ppm.pickup_postal_code}</td>
            </tr>
            {ppm.has_additional_postal_code && (
              <tr>
                <td> Additional Pickup: </td>
                <td> {ppm.additional_pickup_postal_code}</td>
              </tr>
            )}
            <tr>
              <td> Delivery ZIP Code: </td>
              <td> {ppm && ppm.destination_postal_code}</td>
            </tr>
            {!isHHGPPMComboMove && (
              <tr>
                <td> Storage: </td>
                <td>{sitDisplay}</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <div className="usa-width-one-half review-section ppm-review-section">
        <p className="heading">
          Weight
          <span className="edit-section-link">
            <Link data-cy="edit-ppm-weight" to={editWeightAddress}>
              Edit
            </Link>
          </span>
        </p>

        <table>
          <tbody>
            <tr>
              <td> Estimated Weight: </td>
              <td> {ppm && ppm.weight_estimate.toLocaleString()} lbs</td>
            </tr>
            <tr>
              <td> Estimated PPM Incentive: </td>
              <td> {ppm && formatCentsRange(ppm.incentive_estimate_min, ppm.incentive_estimate_max)}</td>
            </tr>
            {ppm.has_requested_advance && (
              <tr>
                <td> Advance: </td>
                <td> ${formatCents(advance.requested_amount)}</td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}

PPMShipmentSummary.propTypes = {
  ppm: PropTypes.object.isRequired,
  movePath: PropTypes.string.isRequired,
};

function mapStateToProps(state, ownProps) {
  const { ppm } = ownProps;
  const advance = selectReimbursement(state, ppm.advance);
  return { ...ownProps, advance };
}

export default connect(mapStateToProps)(PPMShipmentSummary);

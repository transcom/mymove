import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import { get } from 'lodash';
import { object, string, shape, bool, number } from 'prop-types';
import { Grid } from '@trussworks/react-uswds';

import IconWithTooltip from 'shared/ToolTip/IconWithTooltip';
import { formatCentsRange, formatCents } from 'shared/formatters';
import { formatDateSM } from 'shared/formatters';
import { hasShortHaulError } from 'utils/incentives';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentPPM,
  selectReimbursementById,
} from 'store/entities/selectors';
import { selectPPMEstimateError } from 'store/onboarding/selectors';

import './Review.css';

export class PPMShipmentSummary extends Component {
  chooseEstimateText(ppmEstimate) {
    if (hasShortHaulError(ppmEstimate.rateEngineError)) {
      return (
        <td data-testid="estimateError">
          MilMove does not presently support short-haul PPM moves. Please contact your PPPO.
        </td>
      );
    }

    if (ppmEstimate.hasEstimateError || ppmEstimate.rateEngineError) {
      return (
        <td data-testid="estimateError">
          Not ready yet{' '}
          <IconWithTooltip toolTipText="We expect to receive rate data covering your move dates by the end of this month. Check back then to see your estimated incentive." />
        </td>
      );
    }

    return (
      <td data-testid="estimate">
        {' '}
        {ppmEstimate && formatCentsRange(ppmEstimate.incentive_estimate_min, ppmEstimate.incentive_estimate_max)}
      </td>
    );
  }

  render() {
    const { advance, movePath, ppm, ppmEstimate, estimated_storage_reimbursement } = this.props;
    const editDateAndLocationAddress = movePath + '/edit-date-and-location';
    const editWeightAddress = movePath + '/edit-weight';

    const privateStorageString = get(ppm, 'estimated_storage_reimbursement')
      ? `= ${estimated_storage_reimbursement} estimated reimbursement`
      : '';
    const sitDisplay = get(ppm, 'has_sit', false)
      ? `${ppm.weight_estimate && ppm.weight_estimate.toLocaleString()} lbs for ${
          ppm.days_in_storage
        } days ${privateStorageString}`
      : 'Not requested';

    return (
      <div data-testid="ppm-summary" className="review-content">
        <h4>Shipment - You move your stuff (PPM)</h4>
        <Grid row>
          <Grid tablet={{ col: true }}>
            <div className="review-section">
              <p className="heading">
                Dates & Locations
                <span className="edit-section-link">
                  <Link data-testid="edit-ppm-dates" to={editDateAndLocationAddress} className="usa-link">
                    Edit
                  </Link>
                </span>
              </p>

              <table>
                <tbody>
                  <tr>
                    <td> Scheduled move date: </td>
                    <td>{formatDateSM(get(ppm, 'original_move_date'))}</td>
                  </tr>
                  <tr>
                    <td> Pickup ZIP: </td>
                    <td> {ppm && ppm.pickup_postal_code}</td>
                  </tr>
                  {ppm.has_additional_postal_code && (
                    <tr>
                      <td> Additional pickup: </td>
                      <td> {ppm.additional_pickup_postal_code}</td>
                    </tr>
                  )}
                  <tr>
                    <td> Delivery ZIP: </td>
                    <td> {ppm && ppm.destination_postal_code}</td>
                  </tr>
                  <tr>
                    <td> Storage: </td>
                    <td data-testid="sit-display">{sitDisplay}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </Grid>
          <Grid tablet={{ col: true }}>
            <div className="review-section">
              <p className="heading">
                Pre-move Estimated Weight
                <span className="edit-section-link">
                  <Link data-testid="edit-ppm-weight" to={editWeightAddress} className="usa-link">
                    Edit
                  </Link>
                </span>
              </p>

              <table>
                <tbody>
                  <tr>
                    <td> Estimated weight: </td>
                    <td> {ppm.weight_estimate && ppm.weight_estimate.toLocaleString()} lbs</td>
                  </tr>
                  <tr>
                    <td> Estimated PPM incentive: </td>
                    {this.chooseEstimateText(ppmEstimate)}
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
          </Grid>
        </Grid>
      </div>
    );
  }
}

PPMShipmentSummary.propTypes = {
  ppm: object.isRequired,
  orders: object.isRequired,
  movePath: string.isRequired,
  ppmEstimate: shape({
    hasEstimateError: bool.isRequired,
    rateEngineError: Error.isRequired,
    originDutyLocationZip: string.isRequired,
    incentive_estimate_min: number,
    incentive_estimate_max: number,
  }).isRequired,
};

function mapStateToProps(state, ownProps) {
  const { ppm } = ownProps;
  const advance = selectReimbursementById(state, ppm.advance) || {};
  const { incentive_estimate_min, incentive_estimate_max, estimated_storage_reimbursement } =
    selectCurrentPPM(state) || {};

  const ppmEstimateError = selectPPMEstimateError(state);
  let hasError = !!ppmEstimateError;
  const serviceMember = selectServiceMemberFromLoggedInUser(state);

  return {
    ...ownProps,
    advance,
    ppmEstimate: {
      hasEstimateError: hasError,
      rateEngineError: ppmEstimateError,
      originDutyLocationZip: serviceMember?.current_location?.address?.postalCode,
      incentive_estimate_min,
      incentive_estimate_max,
    },
    estimated_storage_reimbursement,
  };
}

export default connect(mapStateToProps)(PPMShipmentSummary);

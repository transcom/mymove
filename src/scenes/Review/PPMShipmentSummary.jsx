import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import { get } from 'lodash';
import { object, string, shape, bool, number } from 'prop-types';
import IconWithTooltip from 'shared/ToolTip/IconWithTooltip';
import { selectPPMForMove, selectReimbursement } from 'shared/Entities/modules/ppms';
import { formatCentsRange, formatCents } from 'shared/formatters';
import { formatDateSM } from 'shared/formatters';
import { getPpmWeightEstimate } from 'scenes/Moves/Ppm/ducks';

import './Review.css';

class PPMShipmentSummary extends Component {
  componentDidUpdate() {
    if (
      !this.props.ppmEstimate.hasEstimateInProgress &&
      !this.props.ppmEstimate.hasEstimateSuccess &&
      !this.props.ppmEstimate.hasEstimateError
    ) {
      this.props.getPpmWeightEstimate(
        this.props.original_move_date,
        this.props.pickup_postal_code,
        this.props.originDutyStationZip,
        this.props.destination_postal_code,
        this.props.weight_estimate,
      );
    }
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
      <div data-cy="ppm-summary">
        <h3>Shipment - You move your stuff (PPM)</h3>
        <div className="grid-row grid-gap review-content">
          <div className="grid-col-6 review-section ppm-review-section">
            <p className="heading">
              Dates & Locations
              <span className="edit-section-link">
                <Link data-cy="edit-ppm-dates" to={editDateAndLocationAddress} className="usa-link">
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
                <tr>
                  <td> Storage: </td>
                  <td data-cy="sit-display">{sitDisplay}</td>
                </tr>
              </tbody>
            </table>
          </div>

          <div className="grid-col-6 review-section ppm-review-section">
            <p className="heading">
              Weight
              <span className="edit-section-link">
                <Link data-cy="edit-ppm-weight" to={editWeightAddress} className="usa-link">
                  Edit
                </Link>
              </span>
            </p>

            <table>
              <tbody>
                <tr>
                  <td> Estimated Weight: </td>
                  <td> {ppm.weight_estimate && ppm.weight_estimate.toLocaleString()} lbs</td>
                </tr>
                <tr>
                  <td> Estimated PPM Incentive: </td>
                  {ppmEstimate.hasEstimateError ? (
                    <td>
                      Not ready yet{' '}
                      <IconWithTooltip toolTipText="We expect to receive rate data covering your move dates by the end of this month. Check back then to see your estimated incentive." />
                    </td>
                  ) : (
                    <td>
                      {' '}
                      {ppmEstimate &&
                        formatCentsRange(ppmEstimate.incentive_estimate_min, ppmEstimate.incentive_estimate_max)}
                    </td>
                  )}
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
      </div>
    );
  }
}

PPMShipmentSummary.propTypes = {
  ppm: object.isRequired,
  movePath: string.isRequired,
  ppmEstimate: shape({
    hasEstimateError: bool.isRequired,
    hasEstimateSuccess: bool.isRequired,
    hasEstimateInProgress: bool.isRequired,
    originDutyStationZip: string.isRequired,
    incentive_estimate_min: number,
    incentive_estimate_max: number,
  }).isRequired,
};

function mapStateToProps(state, ownProps) {
  const { ppm } = ownProps;
  const advance = selectReimbursement(state, ppm.advance);
  const { incentive_estimate_min, incentive_estimate_max, estimated_storage_reimbursement } = selectPPMForMove(
    state,
    ppm.move_id,
  );
  return {
    ...ownProps,
    advance,
    ppmEstimate: {
      hasEstimateError: state.ppm.hasEstimateError,
      hasEstimateSuccess: state.ppm.hasEstimateSuccess,
      hasEstimateInProgress: state.ppm.hasEstimateInProgress,
      originDutyStationZip: state.serviceMember.currentServiceMember.current_station.address.postal_code,
      incentive_estimate_min,
      incentive_estimate_max,
    },
    estimated_storage_reimbursement,
  };
}

const mapDispatchToProps = {
  getPpmWeightEstimate,
};

export default connect(mapStateToProps, mapDispatchToProps)(PPMShipmentSummary);

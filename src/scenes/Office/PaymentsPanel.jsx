import { get, isEmpty } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import { approveReimbursement } from './ducks';
import { no_op } from 'shared/utils';
import { formatDate } from './helpers';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';

class PaymentsTable extends Component {
  approveReimbursement = () => {
    this.props.approveReimbursement(this.props.advance.id);
  };

  renderAdvanceAction = () => {
    if (this.props.ppm.status === 'APPROVED') {
      if (this.props.advance.status === 'APPROVED') {
        return <div>{/* Further actions to come*/}</div>;
      } else {
        return (
          <React.Fragment>
            <div onClick={this.approveReimbursement}>
              <FontAwesomeIcon
                aria-hidden
                className="icon approval-ready"
                icon={faCheck}
                title="Approve"
              />
              <span className="tooltiptext">Approve</span>
            </div>
          </React.Fragment>
        );
      }
    } else {
      return (
        <React.Fragment>
          <FontAwesomeIcon
            aria-hidden
            className="icon approval-blocked"
            icon={faCheck}
            title="Can't approve payment until shipment is approved."
          />
          <span className="tooltiptext">
            Can't approve payment until shipment is approved.
          </span>
        </React.Fragment>
      );
    }
  };

  render() {
    const advance = this.props.advance;

    return (
      <div className="payment-panel">
        <div className="payment-panel-title">Payments</div>
        <table className="payment-table">
          <tbody>
            <tr>
              <th className="payment-table-column-title" />
              <th className="payment-table-column-title">Amount</th>
              <th className="payment-table-column-title">Disbursement</th>
              <th className="payment-table-column-title">Requested on</th>
              <th className="payment-table-column-title">Status</th>
              <th className="payment-table-column-title">Actions</th>
            </tr>
            {!isEmpty(advance) ? (
              <React.Fragment>
                <tr>
                  <th className="payment-table-subheader" colSpan="6">
                    Payments against PPM Incentive
                  </th>
                </tr>
                <tr>
                  <td className="payment-table-column-content">Advance </td>
                  <td className="payment-table-column-content">
                    ${get(advance, 'requested_amount', '').toLocaleString()}.00
                  </td>
                  <td className="payment-table-column-content">
                    {advance.method_of_receipt}
                  </td>
                  <td className="payment-table-column-content">
                    {formatDate(advance.requested_date)}
                  </td>
                  <td className="payment-table-column-content">
                    {advance.status === 'APPROVED' ? (
                      <div>
                        <FontAwesomeIcon
                          aria-hidden
                          className="icon approval-ready"
                          icon={faCheck}
                          title="Approved"
                        />{' '}
                        Approved
                      </div>
                    ) : (
                      <div>
                        <FontAwesomeIcon
                          aria-hidden
                          className="icon approval-waiting"
                          icon={faClock}
                          title="Awaiting Review"
                        />{' '}
                        Awaiting review
                      </div>
                    )}
                  </td>
                  <td className="payment-table-column-content">
                    <span className="tooltip">
                      {this.renderAdvanceAction()}
                    </span>
                    {/* Disabling unimplemented feature for now.
                    <span className="tooltip">
                      <FontAwesomeIcon
                        aria-hidden
                        className="icon payment-action"
                        title="Delete"
                        icon={faTimes}
                      />
                      <span className="tooltiptext">Delete</span>
                    </span>
                    <span className="tooltip">
                      <FontAwesomeIcon
                        aria-hidden
                        title="Edit"
                        className="icon payment-action"
                        icon={faPencil}
                      />
                      <span className="tooltiptext">Edit</span>
                    </span>
                    */}
                  </td>
                </tr>
              </React.Fragment>
            ) : (
              <tr>
                <th className="payment-table-subheader">
                  No payments requested
                </th>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    );
  }
}

const mapStateToProps = state => ({
  ppm: get(state, 'office.officePPMs[0]', {}),
  advance: get(state, 'office.officePPMs[0].advance', {}),
  hasError: false,
  errorMessage: state.office.error,
});

const mapDispatchToProps = dispatch =>
  bindActionCreators({ approveReimbursement, update: no_op }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(PaymentsTable);

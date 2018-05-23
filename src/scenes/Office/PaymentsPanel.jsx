import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
// import { Link } from 'react-router-dom';
import { bindActionCreators } from 'redux';

// import { updatePaymentInfo } from './ducks';
import { no_op } from 'shared/utils';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import faPencil from '@fortawesome/fontawesome-free-solid/faPencilAlt';
import faTimes from '@fortawesome/fontawesome-free-solid/faTimes';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';

const PaymentsTable = props => {
  const ppm = props.ppm;
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
            <th className="payment-table-column-title" />
          </tr>
          {ppm ? (
            <React.Fragment>
              <tr>
                <th className="payment-table-subheader" colSpan="6">
                  Payments against PPM Incentive
                </th>
              </tr>
              <tr>
                <td className="payment-table-column-content">Advance </td>
                <td className="payment-table-column-content Todo">{ppm.id}</td>
                <td className="payment-table-column-content Todo">
                  {ppm.planned_move_date}
                </td>
                <td className="payment-table-column-content Todo">
                  {ppm.planned_move_date}
                </td>
                <td className="payment-table-column-content Todo">
                  {ppm.status === 'APPROVED' ? (
                    <span>
                      <FontAwesomeIcon
                        className="icon approval-ready"
                        icon={faCheck}
                      />{' '}
                      Approved
                    </span>
                  ) : (
                    <FontAwesomeIcon
                      className="icon approval-waiting"
                      icon={faClock}
                    />
                  )}
                </td>
                <td className="payment-table-column-content">
                  <span className="tooltip">
                    {ppm.status === 'APPROVED' ? (
                      <React.Fragment>
                        <FontAwesomeIcon
                          className="icon approval-ready"
                          icon={faCheck}
                        />
                        <span className="tooltiptext">Approve</span>
                      </React.Fragment>
                    ) : (
                      <React.Fragment>
                        <FontAwesomeIcon
                          className="icon approval-blocked"
                          icon={faCheck}
                        />
                        <span
                          className="tooltiptext"
                          aria-label="Can't approve payment until shipment is approved."
                        >
                          Can't approve payment until shipment is approved.
                        </span>
                      </React.Fragment>
                    )}
                  </span>
                  <span className="tooltip" aria-label="Delete">
                    <FontAwesomeIcon
                      className="icon payment-action"
                      title="Delete"
                      icon={faTimes}
                    />
                    <span className="tooltiptext">Delete</span>
                  </span>
                  <span className="tooltip" aria-label="Edit">
                    <FontAwesomeIcon
                      className="icon payment-action"
                      icon={faPencil}
                    />
                    <span className="tooltiptext">Edit</span>
                  </span>
                </td>
              </tr>
            </React.Fragment>
          ) : (
            <tr>
              <th className="payment-table-subheader">No payments requested</th>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};

function mapStateToProps(state) {
  return {
    initialValues: {},
    ppm: get(state, 'office.officePPMs[0]', {}),
    hasError: false,
    errorMessage: state.office.error,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: no_op,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(PaymentsTable);

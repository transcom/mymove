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

// import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';

const PaymentsTable = props => {
  const ppm = props.ppm;
  return (
    <table className="payment-table">
      <tbody>
        <tr>
          <th className="payment-table-title" colSpan="5">
            Payments
          </th>
        </tr>
        <tr>
          <td className="payment-table-column-title" />
          <th className="payment-table-column-title">Amount</th>
          <th className="payment-table-column-title">Disbursement</th>
          <th className="payment-table-column-title">Requested on</th>
          <th className="payment-table-column-title">Approved</th>
        </tr>
        {ppm ? (
          <React.Fragment>
            <tr>
              <th className="payment-table-subheader" colSpan="5">
                Payments against PPM Incentive
              </th>
            </tr>
            <tr>
              <td className="payment-table-column-content" />
              <td className="payment-table-column-content">{ppm.id}</td>
              <td className="payment-table-column-content">{ppm.status}</td>
              <td className="payment-table-column-content">
                {ppm.planned_move_date}
              </td>
              <td className="payment-table-column-content">
                <span className="tooltip">
                  <FontAwesomeIcon className="icon" icon={faCheck} />
                  {ppm.status === 'APPROVED' ? (
                    <span className="tooltiptext">Approve</span>
                  ) : (
                    <span className="tooltiptext">
                      Can't approve payment until shipment is approved.
                    </span>
                  )}
                </span>
                <span className="tooltip">
                  <FontAwesomeIcon
                    className="icon"
                    title="Delete"
                    icon={faTimes}
                  />
                  <span className="tooltiptext">Delete</span>
                </span>
                <span className="tooltip">
                  <FontAwesomeIcon className="icon" icon={faPencil} />
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

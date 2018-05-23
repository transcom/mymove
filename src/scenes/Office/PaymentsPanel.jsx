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
    <div className="usa-grid">
      <table className="payment-table">
        <tbody>
          <tr>
            <th className="payment-table-title">Payments</th>
          </tr>
          <tr>
            <th className="payment-table-column-title" />
            <th className="payment-table-column-title">Amount</th>
            <th className="payment-table-column-title">Disbursement</th>
            <th className="payment-table-column-title">Requested on</th>
            <th className="payment-table-column-title">Approved</th>
          </tr>
          {ppm ? (
            [
              <tr>
                <th className="payment-table-subheader">
                  Payments against PPM Incentive
                </th>
              </tr>,
              <tr>
                <td className="payment-table-column-content" />
                <td className="payment-table-column-content">I</td>
                <td className="payment-table-column-content">Like</td>
                <td className="payment-table-column-content">Dogs</td>
                <td className="payment-table-column-content">
                  <FontAwesomeIcon className="icon" icon={faCheck} />
                  <FontAwesomeIcon className="icon" icon={faTimes} />
                  <FontAwesomeIcon className="icon" icon={faPencil} />
                </td>
              </tr>,
            ]
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

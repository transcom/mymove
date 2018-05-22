import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';
import { bindActionCreators } from 'redux';

// import { updatePaymentInfo } from './ducks';
import { no_op } from 'shared/utils';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';

const PaymentsTable = props => {
  const ppm = props.ppm;
  return (
    <div className="usa-grid">
      <div className="payment-table">
        <div className="payment-table payment-table-header">Payments</div>
        {ppm.has_requested_advance ? (
          [
            <div className="payment-table-content">Amount</div>,
            <div className="payment-table-cell payment-table-cell-header">
              Disbursement
            </div>,
            <div className="payment-table-cell payment-table-cell-header">
              Requested on
            </div>,
            <div className="payment-table-cell payment-table-cell-header">
              Approved
            </div>,
            <div className="payment-table-cell payment-table-cell-header">
              Payments against PPM incentive
            </div>,
            <div className="payment-table-cell">ppm.requested_amount</div>,
            <div className="payment-table-cell">ppm.method_of_receipt</div>,
          ]
        ) : (
          <div className="payment-table payment-table-content">
            <Link to="blank">
              <FontAwesomeIcon
                className="icon"
                icon={faPlusCircle}
                flip="horizontal"
              />{' '}
              Add a payment
            </Link>
          </div>
        )}
      </div>
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

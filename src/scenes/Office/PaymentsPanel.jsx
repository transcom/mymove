import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

// import { updatePaymentInfo } from './ducks';
import { no_op } from 'shared/utils';

// import FontAwesomeIcon from '@fortawesome/react-fontawesome';

const PaymentsTable = props => {
  const ppm = props.ppm;
  return (
    <div className="usa-grid">
      <table>
        <tbody>
          <tr>
            <th />
            <th>Amount</th>
            <th>Disbursement</th>
            <th>Requested on</th>
            <th>Approved</th>
          </tr>
          {ppm.has_requested_advance ? (
            [
              <tr>
                <th>Payments against PPM Incentive</th>
              </tr>,
              <tr>
                <td />
                <td>{ppm.requested_amount}</td>
                <td>{ppm.method_of_receipt}</td>
                <td>Dogs</td>
                <td>Dogs</td>
              </tr>,
            ]
          ) : (
            <tr>
              <th>No payments requested</th>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};

function mapStateToProps(state) {
  // let serviceMember = get(state, 'office.officeServiceMember', {});
  // let backupContact = get(state, 'office.officeBackupContacts.0', {}); // there can be only one

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

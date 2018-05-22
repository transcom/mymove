// import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

// import { updatePaymentInfo } from './ducks';
import { no_op } from 'shared/utils';

// import FontAwesomeIcon from '@fortawesome/react-fontawesome';

class PaymentsTable extends Component {
  // const backupAddress = props.backupMailingAddress;
  // const backupContact = props.backupContact;
  render() {
    return (
      <React.Fragment>
        <table>
          <tbody>
            <tr>
              <th />
              <th>Amount</th>
              <th>Disbursement</th>
              <th>Requested on</th>
              <th>Approved</th>
            </tr>
            <tr>
              <th>Payments against PPM Incentive</th>
            </tr>
            <tr>
              <td />
              <td>Dogs</td>
              <td>Dogs</td>
              <td>Dogs</td>
              <td>Dogs</td>
            </tr>
          </tbody>
        </table>
      </React.Fragment>
    );
  }
}

function mapStateToProps(state) {
  // let serviceMember = get(state, 'office.officeServiceMember', {});
  // let backupContact = get(state, 'office.officeBackupContacts.0', {}); // there can be only one

  return {
    // reduxForm
    initialValues: {},

    // addressSchema: get(state, 'swagger.spec.definitions.Address', {}),
    // backupContactSchema: get(
    //   state,
    //   'swagger.spec.definitions.ServiceMemberBackupContactPayload',
    //   {},
    // ),
    // backupMailingAddress: serviceMember.backup_mailing_address,
    // backupContact: backupContact,

    // getUpdateArgs: function() {
    //   let values = getFormValues(formName)(state);
    //   return [
    //     serviceMember.id,
    //     { backup_mailing_address: values.backupMailingAddress },
    //     backupContact.id,
    //     values.backupContact,
    //   ];
    // },

    hasError: false,
    // errorMessage: state.office.error,
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

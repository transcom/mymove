import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getPaymentRequestList, selectPaymentRequests } from 'shared/Entities/modules/paymentRequests';

class PaymentRequestIndex extends React.Component {
  componentDidMount() {
    this.props.getPaymentRequestList();
  }

  render() {
    console.log(this.props.paymentRequests);
    return (
      <>
        <h1>Payment Requests</h1>
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Final</th>
              <th>Rejection Reason</th>
              <th>Service Item IDs</th>
            </tr>
          </thead>
          {this.props.paymentRequests.map(pr => (
            <tr>
              <td>{pr.id}</td>
              <td>{pr.isFinal}</td>
              <td>{pr.rejectionReason}</td>
              <td>{pr.serviceItemIDs}</td>
            </tr>
          ))}
        </table>
      </>
    );
  }
}

const mapStateToProps = state => ({
  paymentRequests: selectPaymentRequests(state),
});

const mapDispatchToProps = dispatch => bindActionCreators({ getPaymentRequestList }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(PaymentRequestIndex);

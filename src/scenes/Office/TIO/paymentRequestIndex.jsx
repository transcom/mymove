import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getPaymentRequestList, selectPaymentRequests } from 'shared/Entities/modules/paymentRequests';
import { Link } from 'react-router-dom';

class PaymentRequestIndex extends React.Component {
  componentDidMount() {
    this.props.getPaymentRequestList();
  }

  render() {
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
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {this.props.paymentRequests.map(pr => (
              <tr key={pr.id}>
                <td>
                  <Link to={`/payment_requests/${pr.id}`}>{pr.id}</Link>
                </td>
                <td>{`${pr.isFinal}`}</td>
                <td>{pr.rejectionReason}</td>
                <td>{pr.serviceItemIDs}</td>
                <td>{pr.status}</td>
              </tr>
            ))}
          </tbody>
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

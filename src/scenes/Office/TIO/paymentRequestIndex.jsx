import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getPaymentRequestList, selectPaymentRequests } from 'shared/Entities/modules/paymentRequests';
import { Link, withRouter } from 'react-router-dom';

class PaymentRequestIndex extends React.Component {
  componentDidMount() {
    // TODO - only get payment requests associated with the given move order ID
    this.props.getPaymentRequestList();
  }

  render() {
    const {
      params: { moveOrderId },
    } = this.props.match;

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
            {this.props.paymentRequests.map((pr) => (
              <tr key={pr.id}>
                <td>
                  <Link to={`/moves/${moveOrderId}/payment-requests/${pr.id}`}>{pr.id}</Link>
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

const mapStateToProps = (state) => ({
  paymentRequests: selectPaymentRequests(state),
});

const mapDispatchToProps = (dispatch) => bindActionCreators({ getPaymentRequestList }, dispatch);

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(PaymentRequestIndex));

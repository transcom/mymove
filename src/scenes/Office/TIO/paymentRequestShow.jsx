import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { selectPaymentRequest, getPaymentRequest, updatePaymentRequest } from 'shared/Entities/modules/paymentRequests';

class PaymentRequestShow extends React.Component {
  componentDidMount() {
    this.props.getPaymentRequest(this.props.id);
  }

  approvePaymentRequest() {
    this.props.updatePaymentRequest(this.props.id);
  }

  denyPaymentRequest() {
    this.props.updatePaymentRequest(this.props.id);
  }

  render() {
    const {
      id,
      paymentRequest: { isFinal, rejectionReason, serviceItemIDs, status },
    } = this.props;
    return (
      <>
        <div>
          <h1>Payment Request Id {id}</h1>
          <ul>
            <li>isFinal: {`${isFinal}`}</li>
            <li>rejectionReason: {rejectionReason}</li>
            <li>serviceItemIds: {serviceItemIDs}</li>
            <li>status: {status}</li>
          </ul>

          <button className="usa-button usa-button--outline" onClick={this.approvePaymentRequest}>
            Approve
          </button>
          <button className="usa-button usa-button--outline" onClick={this.denyPaymentRequest}>
            Deny
          </button>
        </div>
      </>
    );
  }
}
const mapStateToProps = (state, props) => {
  const id = props.match.params.id;
  return {
    id,
    paymentRequest: selectPaymentRequest(state, id),
  };
};

const mapDispatchToProps = dispatch => bindActionCreators({ getPaymentRequest, updatePaymentRequest }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(PaymentRequestShow);

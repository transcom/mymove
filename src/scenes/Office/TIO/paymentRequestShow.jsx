import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { selectPaymentRequest, getPaymentRequest } from 'shared/Entities/modules/paymentRequests';

class PaymentRequestShow extends React.Component {
  componentDidMount() {
    this.props.getPaymentRequest(this.props.id);
  }

  render() {
    const {
      id,
      paymentRequest: { isFinal, rejectionReason, serviceItemIDs },
    } = this.props;
    return (
      <div>
        <h1>Payment Request Id {id}</h1>
        <ul>
          <li>isFinal: {`${isFinal}`}</li>
          <li>rejectionReason: {rejectionReason}</li>
          <li>serviceItemIds: {serviceItemIDs}</li>
          <li></li>
        </ul>
      </div>
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

const mapDispatchToProps = dispatch => bindActionCreators({ getPaymentRequest }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(PaymentRequestShow);

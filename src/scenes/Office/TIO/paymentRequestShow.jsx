import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getPaymentRequest } from 'shared/Entities/modules/paymentRequests';

class PaymentRequestShow extends React.Component {
  componentDidMount() {
    this.props.getPaymentRequest(this.props.id);
  }

  render() {
    const { id } = this.props;
    return <h1>Payment Request Id {id}</h1>;
  }
}
const mapStateToProps = (_state, props) => ({
  id: props.match.params.id,
});

const mapDispatchToProps = dispatch => bindActionCreators({ getPaymentRequest }, dispatch);

export default connect(mapStateToProps, mapDispatchToProps)(PaymentRequestShow);

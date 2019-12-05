import React from 'react';
import { connect } from 'react-redux';

class PaymentRequestShow extends React.Component {
  render() {
    const { id } = this.props;
    return <h1>Payment Request Id {id}</h1>;
  }
}
const mapStateToProps = (_state, props) => ({
  id: props.match.params.id,
});

export default connect(mapStateToProps)(PaymentRequestShow);

import React from 'react';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { getEntitlements, getAllCustomerMoves, getCustomerInfo } from 'shared/Entities/modules/moveTaskOrders';

class TOO extends React.Component {
  componentDidMount() {
    this.props.getAllCustomerMoves();
  }

  render() {
    return (
      <div>
        <h2>All Customer Moves</h2>
      </div>
    );
  }
}
const mapStateToProps = state => {
  const entitlements = get(state, 'entities.entitlements');
  const customer = get(state, 'entities.customer', {});
  return {};
};
const mapDispatchToProps = {
  getEntitlements,
  getCustomerInfo,
  getAllCustomerMoves,
};

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(TOO);

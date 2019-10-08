import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { getEntitlements } from 'shared/Entities/modules/moveTaskOrders';

class TOO extends React.Component {
  componentDidMount() {
    this.props.getEntitlements('fake_move_task_order_id');
  }

  render() {
    return <h1>TOO Placeholder Page</h1>;
  }
}
const mapStateToProps = state => {
  return {};
};
const mapDispatchToProps = dispatch => bindActionCreators({ getEntitlements }, dispatch);

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(TOO);

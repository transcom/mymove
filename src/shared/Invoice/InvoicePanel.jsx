import React, { Component } from 'react';
import PropTypes from 'prop-types';
import './InvoicePanel.css';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import BasicPanel from 'shared/BasicPanel';
import {
  makeGetUnbilledShipmentLineItems,
  makeTotalFromUnbilledLineItems,
} from 'shared/Entities/modules/shipmentLineItems';

export class InvoicePanel extends Component {
  constructor() {
    super();
    this.state = {
      canApprove: null,
      shipmentState: null,
    };
  }

  componentDidUpdate(prevProps) {
    if (!this.state.shipmentState || prevProps.shipmentState !== this.props.shipmentState) {
      this.setState({ shipmentState: this.props.shipmentState });
    }
    if (!this.state.canApprove || prevProps.canApprove !== this.props.canApprove) {
      this.setState({ canApprove: this.props.canApprove });
    }
  }

  approvePayment = () => {
    this.props.onApprovePayment(this.props.shipmentId);
  };

  render() {
    let invoicingContent = (
      <tr className="empty-content">
        <td>No line items</td>
      </tr>
    );
    if (this.props.shipmentLineItems.length > 0) {
      //stand up a table
    }

    return (
      <div className="invoice-panel">
        <BasicPanel title={'Invoicing'}>
          {this.state.canApprove &&
            this.state.shipmentState === 'DELIVERED' && (
              <div className="usa-width-one-whole align-right">
                <button className="button button-secondary" disabled={!this.canApprove} onClick={this.approvePayment}>
                  Approve Payment
                  {this.canApprove && this.shipmentState === 'DELIVERED'}
                </button>
              </div>
            )}
          <table cellSpacing={0}>
            <tbody>{invoicingContent}</tbody>
          </table>
        </BasicPanel>
      </div>
    );
  }
}

InvoicePanel.propTypes = {
  shipmentLineItems: PropTypes.array,
  shipmentId: PropTypes.string,
  shipmentState: PropTypes.string,
  lineItemsTotal: PropTypes.number,
  onApprovePayment: PropTypes.func,
  canApprove: PropTypes.bool,
};

//https://github.com/reduxjs/reselect#sharing-selectors-with-props-across-multiple-component-instances
const makeMapStateToProps = () => {
  //using a memoized selector
  const mapStateToProps = (state, ownProps) => {
    const getLineItems = makeGetUnbilledShipmentLineItems();
    const getLineItemSum = makeTotalFromUnbilledLineItems();
    return {
      shipmentLineItems: getLineItems(state, ownProps.shipmentId),
      lineItemsTotal: getLineItemSum(state, ownProps.shipmentId),
    };
  };
  return mapStateToProps;
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}
export default connect(makeMapStateToProps, mapDispatchToProps)(InvoicePanel);
